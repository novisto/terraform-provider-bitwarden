package bitwarden

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
	"github.com/samber/lo"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ItemLoginURI struct {
	Match int    `json:"match"`
	URI   string `json:"uri"`
}

type ItemLogin struct {
	URIs                 []ItemLoginURI `json:"uris"`
	Username             string         `json:"username"`
	Password             string         `json:"password"`
	TOTP                 string         `json:"totp"`
	PasswordRevisionDate string         `json:"passwordRevisionDate"`
}

type ItemSecureNote struct {
	Type int `json:"type"`
}

// Item Generic model for an item in the bw CLI
type Item struct {
	Object         string         `json:"object"`
	ID             string         `json:"id"`
	OrganizationId string         `json:"organizationId"`
	FolderID       string         `json:"folderId"`
	Type           int            `json:"type"`
	Reprompt       int            `json:"reprompt"`
	Name           string         `json:"name"`
	Notes          string         `json:"notes"`
	Favorite       bool           `json:"favorite"`
	Login          ItemLogin      `json:"login"`
	SecureNote     ItemSecureNote `json:"secureNote"`
	CollectionIDs  []string       `json:"collectionIds"`
	RevisionDate   string         `json:"revisionDate"`
}

type ItemResponse struct {
	Data Item `json:"data"`
}

// ItemCreate Represents the create and update payload for an item in the bw CLI
type ItemCreate struct {
	OrganizationId string          `json:"organizationId"`
	CollectionIDs  []string        `json:"collectionIds"`
	FolderID       *string         `json:"folderId"`
	Type           int             `json:"type"`
	Name           string          `json:"name"`
	Notes          string          `json:"notes"`
	Favorite       bool            `json:"favorite,omitempty"`
	Fields         []string        `json:"fields"`
	Login          *struct{}       `json:"login"` // Format unknown
	SecureNote     *ItemSecureNote `json:"secureNote"`
	Card           *struct{}       `json:"card"`     // Format unknown
	Identity       *struct{}       `json:"identity"` // Format unknown
	Reprompt       int             `json:"reprompt"`
}

func PrepareSecureNoteCreate(secureNote SecureNote) ItemCreate {
	var folderId *string = nil
	if !secureNote.FolderID.Null {
		folderId = &secureNote.FolderID.Value
	}

	var favorite = false
	if !secureNote.Favorite.Null {
		favorite = secureNote.Favorite.Value
	}

	var reprompt = 0
	if !secureNote.Reprompt.Null && secureNote.Reprompt.Value {
		reprompt = 1
	}

	return ItemCreate{
		OrganizationId: secureNote.OrganizationId.Value,
		CollectionIDs:  secureNote.CollectionIDs,
		FolderID:       folderId,
		Type:           2,
		Name:           secureNote.Name.Value,
		Notes:          secureNote.Notes.Value,
		Favorite:       favorite,
		Fields:         nil,
		Login:          nil,
		SecureNote:     &ItemSecureNote{Type: 0},
		Card:           nil,
		Identity:       nil,
		Reprompt:       reprompt,
	}
}

type Client struct {
	Password string
	Session  string
}

type bwServeClient struct {
	Command    *exec.Cmd
	restClient *resty.Client
}

func bitwardenServe(password string) (*bwServeClient, error) {
	bwClient := bwServeClient{}

	bwPort := strconv.Itoa(rand.Intn(65000-10000) + 10000)
	ln, err := net.Listen("tcp", ":"+bwPort)

	for err != nil {
		bwPort = strconv.Itoa(rand.Intn(65000-10000) + 10000)
		ln, err = net.Listen("tcp", ":"+bwPort)
	}
	err = ln.Close()
	if err != nil {
		return nil, err
	}

	bwClient.Command = exec.Command("bw", "serve", "--port", bwPort)
	if err := bwClient.Command.Start(); err != nil {
		return nil, err
	}

	bwClient.restClient = resty.New()
	bwClient.restClient.SetBaseURL("http://localhost:" + bwPort)

	bwTimedout := true
	var errorResp *resty.Response
	var bwErr error
	start := time.Now()
	for time.Since(start) < time.Second*time.Duration(10) {
		errorResp, bwErr = bwClient.restClient.R().Get("/status")

		if bwErr == nil && errorResp.StatusCode() == 200 {
			bwTimedout = false
			break
		}
		time.Sleep(time.Second)

	}
	if bwTimedout {
		bwClient.Close()
		if bwErr != nil {
			return nil, fmt.Errorf(
				"bitwarden serve did not start in a reasonable time with error %s",
				bwErr,
			)
		} else if errorResp != nil {
			return nil, fmt.Errorf(
				"bitwarden serve did not start in a reasonable time http error [%s] %s",
				errorResp.StatusCode(),
				errorResp.Body(),
			)
		} else {
			return nil, fmt.Errorf("bitwarden serve did not start in a reasonable time")
		}
	}

	resp, err := bwClient.restClient.R().SetBody(map[string]string{"password": password}).Post("/unlock")
	if err != nil {
		bwClient.Close()
		return nil, err
	}

	if resp.StatusCode() != 200 {
		bwClient.Close()
		return nil, fmt.Errorf("error unlocking bitwarden\n%s", resp.Body())
	}

	err = bwClient.Sync()
	if err != nil {
		bwClient.Close()
		return nil, err
	}

	return &bwClient, nil
}

func (bwClient *bwServeClient) Close() {
	_ = bwClient.Command.Process.Kill()
}

func (bwClient *bwServeClient) Sync() error {
	resp, err := bwClient.restClient.R().Post("/sync")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error syncing\n%s", resp.Body())
	}

	return nil
}

func NewClient(password string) (*Client, error) {
	c := Client{Password: password}

	out, err := RunCommand("bw", "--version")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", out, err))
	}

	requiredVersion, err := version.NewVersion("1.22.0")
	if err != nil {
		return nil, err
	}

	clientVersion, err := version.NewVersion(strings.TrimSpace(strings.ReplaceAll(out, "\n", "")))
	if err != nil {
		return nil, err
	}

	if clientVersion.LessThan(requiredVersion) {
		return nil, fmt.Errorf("bitwarden client version(%s) must be equal or greater than 1.22.0", out)
	}

	return &c, nil
}

func (c *Client) CreateSecureNote(secureNote SecureNote) (*Item, error) {
	createPayload := PrepareSecureNoteCreate(secureNote)

	bwClient, err := bitwardenServe(c.Password)
	if err != nil {
		return nil, err
	}

	resp, err := bwClient.restClient.R().SetBody(createPayload).Post("/object/item")
	if err != nil {
		return nil, err
	}
	bwClient.Close()

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("bitwarden error when creating secure note\n%s", resp.Body())
	}

	var decoded ItemResponse
	err = json.Unmarshal(resp.Body(), &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded.Data, nil
}

func (c *Client) UpdateSecureNote(id string, secureNote SecureNote) (*Item, error) {
	updatePayload := PrepareSecureNoteCreate(secureNote)

	bwClient, err := bitwardenServe(c.Password)
	if err != nil {
		return nil, err
	}

	resp, err := bwClient.restClient.R().SetBody(updatePayload).Put(fmt.Sprintf("/object/item/%s", id))
	if err != nil {
		return nil, err
	}
	bwClient.Close()

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("bitwarden error when updating secure note\n%s", resp.Body())
	}

	var decoded ItemResponse
	err = json.Unmarshal(resp.Body(), &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded.Data, nil
}

func (c *Client) GetItem(id string) (*Item, error) {
	bwClient, err := bitwardenServe(c.Password)
	if err != nil {
		return nil, err
	}

	resp, err := bwClient.restClient.R().Get(fmt.Sprintf("/object/item/%s", id))
	if err != nil {
		return nil, err
	}
	bwClient.Close()

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("bitwarden error when fetching secure note\n%s", resp.Body())
	}

	var decoded ItemResponse
	err = json.Unmarshal(resp.Body(), &decoded)
	if err != nil {
		return nil, err
	}

	// This is a fix for BW cli that returns duplicated values for collectionIDs
	decoded.Data.CollectionIDs = lo.Uniq[string](decoded.Data.CollectionIDs)

	return &decoded.Data, nil
}

func (c *Client) MoveItem(id string, newOrgId string) error {
	bwClient, err := bitwardenServe(c.Password)
	if err != nil {
		return err
	}

	resp, err := bwClient.restClient.R().Post(fmt.Sprintf("/move/%s/%s", id, newOrgId))
	if err != nil {
		return err
	}
	bwClient.Close()

	if resp.StatusCode() != 200 {
		return fmt.Errorf("bitwarden error when moving secure note\n%s", resp.Body())
	}

	return nil
}

func (c *Client) DeleteItem(id string) error {
	bwClient, err := bitwardenServe(c.Password)
	if err != nil {
		return err
	}

	resp, err := bwClient.restClient.R().Delete(fmt.Sprintf("/object/item/%s", id))
	if err != nil {
		return err
	}
	bwClient.Close()

	if resp.StatusCode() != 200 {
		return fmt.Errorf("bitwarden error when deleting secure note\n%s", resp.Body())
	}

	return nil
}

package bitwarden

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	var collectionIDs []string
	secureNote.CollectionIDs.ElementsAs(context.TODO(), collectionIDs, false)

	return ItemCreate{
		OrganizationId: secureNote.OrganizationId.Value,
		CollectionIDs:  collectionIDs,
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

func NewClient(password string) (*Client, error) {
	c := Client{Password: password}

	out, err := RunCommand("bw", "unlock", c.Password, "--raw")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", out, err))
	}
	c.Session = out

	err = c.Sync()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) Sync() error {
	out, err := RunCommand("bw", "sync", "-f", "--session", c.Session)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\n%s", out, err))
	}
	return nil
}

func (c *Client) CreateSecureNote(secureNote SecureNote) (*Item, error) {
	createPayload := PrepareSecureNoteCreate(secureNote)

	marshal, err := json.Marshal(createPayload)
	if err != nil {
		return nil, err
	}
	b64payload := base64.StdEncoding.EncodeToString(marshal)

	RandSleep(5)

	out, err := RunCommand(
		"bw", "create", "item", "--organizationid", secureNote.OrganizationId.Value, b64payload, "--session", c.Session,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", out, err))
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

func (c *Client) UpdateSecureNote(id string, secureNote SecureNote) (*Item, error) {
	updatePayload := PrepareSecureNoteCreate(secureNote)

	marshal, err := json.Marshal(updatePayload)
	if err != nil {
		return nil, err
	}
	b64payload := base64.StdEncoding.EncodeToString(marshal)

	RandSleep(5)

	out, err := RunCommand(
		"bw",
		"edit",
		"item",
		id,
		"--organizationid",
		secureNote.OrganizationId.Value,
		b64payload,
		"--session",
		c.Session,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", out, err))
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

func (c *Client) GetItem(id string) (*Item, error) {
	RandSleep(5)

	out, err := RunCommand("bw", "get", "item", id, "--session", c.Session)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", out, err))
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, err
	}

	// This is a fix for BW cli that returns duplicated values for collectionIDs
	decoded.CollectionIDs = Unique(decoded.CollectionIDs)

	return &decoded, nil
}

func (c *Client) MoveItem(id string, newOrgId string) error {
	RandSleep(5)

	out, err := RunCommand("bw", "move", id, newOrgId, "--session", c.Session)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\n%s", out, err))
	}

	return nil
}

func (c *Client) DeleteItem(id string) error {
	RandSleep(5)

	out, err := RunCommand("bw", "delete", "item", id, "--session", c.Session)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\n%s", out, err))
	}

	return nil
}

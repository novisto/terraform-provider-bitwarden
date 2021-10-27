package bitwarden

import (
	"encoding/base64"
	"encoding/json"
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

func NewClient(password string) (*Client, string, error) {
	c := Client{Password: password}

	out, err := RunCommand("bw", "unlock", c.Password, "--raw")
	if err != nil {
		return nil, out, err
	}
	c.Session = out

	out, err = c.Sync()
	if err != nil {
		return nil, out, err
	}

	return &c, "", nil
}

func (c *Client) Sync() (string, error) {
	out, err := RunCommand("bw", "sync", "-f", "--session", c.Session)
	if err != nil {
		return out, err
	}
	return "", nil
}

func (c *Client) CreateSecureNote(secureNote SecureNote) (*Item, string, error) {
	out, err := c.Sync()
	if err != nil {
		return nil, out, err
	}

	createPayload := PrepareSecureNoteCreate(secureNote)

	marshal, err := json.Marshal(createPayload)
	if err != nil {
		return nil, "", err
	}
	b64payload := base64.StdEncoding.EncodeToString(marshal)

	out, err = RunCommand(
		"bw", "create", "item", "--organizationid", secureNote.OrganizationId.Value, b64payload, "--session", c.Session,
	)
	if err != nil {
		return nil, out, err
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, "", err
	}

	return &decoded, out, nil
}

func (c *Client) UpdateSecureNote(id string, secureNote SecureNote) (*Item, string, error) {
	out, err := c.Sync()
	if err != nil {
		return nil, out, err
	}

	updatePayload := PrepareSecureNoteCreate(secureNote)

	marshal, err := json.Marshal(updatePayload)
	if err != nil {
		return nil, "", err
	}
	b64payload := base64.StdEncoding.EncodeToString(marshal)

	out, err = RunCommand(
		"bw", "edit", "item", id, "--organizationid", secureNote.OrganizationId.Value, b64payload, "--session", c.Session,
	)
	if err != nil {
		return nil, out, err
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, "", err
	}

	return &decoded, out, nil
}

func (c *Client) GetItem(id string) (*Item, string, error) {
	out, err := c.Sync()
	if err != nil {
		return nil, out, err
	}

	out, err = RunCommand("bw", "get", "item", id, "--session", c.Session)
	if err != nil {
		return nil, out, err
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, out, err
	}

	return &decoded, out, nil
}

func (c *Client) MoveItem(id string, newOrgId string) (string, error) {
	out, err := c.Sync()
	if err != nil {
		return out, err
	}

	out, err = RunCommand("bw", "move", id, newOrgId, "--session", c.Session)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (c *Client) DeleteItem(id string) (string, error) {
	out, err := c.Sync()
	if err != nil {
		return out, err
	}

	out, err = RunCommand("bw", "delete", "item", id, "--session", c.Session)
	if err != nil {
		return out, err
	}

	return out, nil
}

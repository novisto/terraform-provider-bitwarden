package bitwarden

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"time"
)

type Client struct {
	Password string
	Session  string
}

func NewClient(password string) (*Client, error) {
	c := Client{Password: password}

	log.Printf("Created client")

	session, err := RunCommand("bw", "unlock", c.Password, "--raw")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	c.Session = session

	c.Sync()

	return &c, nil
}

func (c *Client) Sync() {
	log.Printf("Running sync...")
	out, err := RunCommand("bw", "sync", "-f", "--session", c.Session)
	if err != nil {
		log.Printf(out)
		log.Fatal(err)
	}

	log.Printf("Sync done...")
}

type ItemLoginURI struct {
	Match int    `json:"match"`
	URI   string `json:"uri"`
}

type ItemLogin struct {
	URIs                 []ItemLoginURI `json:"uris"`
	Username             string         `json:"username"`
	Password             string         `json:"password"`
	TOTP                 string         `json:"totp"`
	PasswordRevisionDate time.Time      `json:"passwordRevisionDate"`
}

type ItemSecureNote struct {
	Type int `json:"type"`
}

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
	RevisionDate   time.Time      `json:"revisionDate"`
}

type SecureNoteCreate struct {
	OrganizationId string         `json:"organizationId"`
	FolderID       string         `json:"folderId"`
	Name           string         `json:"name"`
	Notes          string         `json:"notes"`
	Favorite       bool           `json:"favorite"`
	CollectionIDs  []string       `json:"collectionIds"`
	Type           int            `json:"type"`
	Reprompt       int            `json:"reprompt"`
	SecureNote     ItemSecureNote `json:"secureNote"`
}

func (c *Client) CreateSecureNote(secureNote SecureNote) (*Item, error) {
	c.Sync()

	var reprompt int
	if !secureNote.Reprompt.Null {
		if secureNote.Reprompt.Value {
			reprompt = 1
		} else {
			reprompt = 0
		}
	} else {
		reprompt = 0
	}

	createPayload := SecureNoteCreate{
		OrganizationId: secureNote.OrganizationId.Value,
		FolderID:       secureNote.FolderID.Value,
		Name:           secureNote.Name.Value,
		Notes:          secureNote.Notes.Value,
		Favorite:       secureNote.Favorite.Value,
		CollectionIDs:  secureNote.CollectionIDs,
		Type:           2,
		Reprompt:       reprompt,
		SecureNote:     ItemSecureNote{Type: 0},
	}

	log.Printf("To JSON...")

	marshal, err := json.Marshal(createPayload)
	if err != nil {
		return nil, err
	}

	log.Printf("Got JSON: " + string(marshal[:]))
	log.Printf("To b64...")

	b64payload := base64.StdEncoding.EncodeToString(marshal)
	log.Printf(string(marshal))
	log.Printf(b64payload)

	out, err := RunCommand(
		"bw", "create", "item", "--organizationid", secureNote.OrganizationId.Value, b64payload, "--session", c.Session,
	)
	log.Printf(out)
	if err != nil {
		return nil, err
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

func (c *Client) GetItem(id string) (*Item, error) {
	_, err := RunCommand("bw", "sync", "-f")
	if err != nil {
		return nil, err
	}

	out, err := RunCommand("bw", "get", "item", id)
	if err != nil {
		return nil, err
	}

	var decoded Item
	err = json.Unmarshal([]byte(out), &decoded)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}

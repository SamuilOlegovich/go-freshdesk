package freshdesk

import (
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Contact struct {
	// ID of the contact
	ID uint64 `json:"id,omitempty"`

	// External ID of the contact
	UniqueExternalID string `json:"unique_external_id,omitempty"`

	// Set to true if the contact has been verified
	Active bool `json:"active"`

	// Set to true if the contact has been deleted. Note that this attribute will only be present for deleted contacts
	Deleted bool `json:"deleted,omitempty"`

	// Contact creation timestamp
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Contact updated timestamp
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Set to true if the contact can see all tickets that are associated with the company to which he belongs
	ViewAllTickets bool `json:"view_all_tickets"`

	// Name of the contact
	Name string `json:"name"`

	// Job title of the contact
	JobTitle *string `json:"job_title,omitempty"`

	// A short description of the contact
	Description *string `json:"description,omitempty"`

	// Address of the contact
	Address *string `json:"address,omitempty"`

	// Primary email address of the contact.
	// If you want to associate additional email(s) with this contact, use the other_emails attribute
	Email string `json:"email,omitempty"`

	// Additional emails associated with the contact
	OtherEmails []string `json:"other_emails,omitempty"`

	// Telephone number of the contact
	Phone *string `json:"phone,omitempty"`

	// Mobile number of the contact
	Mobile *string `json:"mobile,omitempty"`

	// ["Alpa","Thor","Jackpot"]
	OtherPhones []map[string]string `json:"other_phone_numbers,omitempty"`

	// Twitter handle of the contact
	TwitterID *string `json:"twitter_id,omitempty"`

	// Time zone in which the contact resides
	TimeZone *string `json:"time_zone,omitempty"`

	// Language of the contact
	Language *string `json:"language,omitempty"`

	// ID of the primary company to which this contact belongs
	CompanyID *string `json:"company_id,omitempty"`

	// Additional companies associated with the contact
	OtherCompanies []map[string]string `json:"other_companies,omitempty"`

	// Tags associated with this contact
	Tags []string `json:"tags,omitempty"`

	// Key value pair containing the name and value of the custom fields.
	// See https://support.freshdesk.com/support/solutions/articles/216553
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

func (c *Contact) Hash() (string, error) {
	in, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	x := md5.Sum(in) //nolint:gosec
	return hex.EncodeToString(x[:]), nil
}

type ContactsClient interface {
	Create(t *Contact) (*Contact, error)
	Update(id uint64, t *Contact) (*Contact, error)
	View(id uint64) (*Contact, error)
	ListAll() ([]*Contact, error)
	Delete(id uint64) error
	Restore(id uint64) error
}

type contactsClient struct {
	*client
}

// Create creates a new contact
func (c *contactsClient) Create(t *Contact) (*Contact, error) {
	req, err := c.client.newRequest(http.MethodPost, "contacts", t)
	if err != nil {
		return nil, err
	}

	res := new(Contact)
	if err = c.client.do(req, res, http.StatusCreated); err != nil {
		return nil, err
	}

	return res, nil
}

// Update updates an existing contact
func (c *contactsClient) Update(id uint64, t *Contact) (*Contact, error) {
	req, err := c.client.newRequest(http.MethodPut, fmt.Sprintf("contacts/%d", id), t)
	if err != nil {
		return nil, err
	}

	res := new(Contact)
	if err := c.client.do(req, res, http.StatusOK); err != nil {
		return nil, err
	}

	return res, nil
}

// View gets an existing contact by id
func (c *contactsClient) View(id uint64) (*Contact, error) {
	req, err := c.client.newRequest(http.MethodGet, fmt.Sprintf("contacts/%d", id), nil)
	if err != nil {
		return nil, err
	}

	res := new(Contact)
	if err := c.client.do(req, res, http.StatusOK); err != nil {
		return nil, err
	}

	return res, nil
}

// ListAll lists all existing contacts
func (c *contactsClient) ListAll() ([]*Contact, error) {
	req, err := c.client.newRequest(http.MethodGet, "contacts", nil)
	if err != nil {
		return nil, err
	}

	var res []*Contact
	if err := c.client.do(req, &res, http.StatusOK); err != nil {
		return nil, err
	}

	return res, nil
}

// Delete deletes an existing contact
func (c *contactsClient) Delete(id uint64) error {
	req, err := c.client.newRequest(http.MethodDelete, fmt.Sprintf("contacts/%d", id), nil)
	if err != nil {
		return err
	}

	return c.client.do(req, nil, http.StatusOK)
}

// Restore restores previously deleted contact
func (c *contactsClient) Restore(id uint64) error {
	req, err := c.client.newRequest(http.MethodPut, fmt.Sprintf("contacts/%d/restore", id), nil)
	if err != nil {
		return err
	}

	return c.client.do(req, nil, http.StatusOK)
}

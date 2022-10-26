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
	// Additional phones associated with the contact
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

type ContactField struct {
	// Set to true if the field can be updated by customers during signup
	EditableInSignup bool `json:"editable_in_signup"`
	// ID of the contact field
	ID uint64 `json:"id"`
	// Display name for the field (as seen by agents)
	Label string `json:"label"`
	// Name of the contact field
	Name string `json:"name"`
	// Position of the contact field
	Position uint64 `json:"position"`
	// Set to true if the field is not a custom field
	Default bool `json:"default"`
	// For custom contact fields, type of value associated with the field will be given (Examples custom_date, custom_text...)
	Type string `json:"type"`
	// Customers can edit the field in the customer portal
	CustomersCanEdit bool `json:"customers_can_edit"`
	// Display name for the field (as seen in the customer portal)
	LabelForCustomers string `json:"label_for_customers"`
	// Set to true if the field is mandatory in the customer portal
	RequiredForCustomers bool `json:"required_for_customers"`
	// Customers can see the field in the customer portal
	DisplayedForCustomers bool `json:"displayed_for_customers"`
	// Set to true if the field is mandatory for agents
	RequiredForAgents bool `json:"required_for_agents"`
	// List of values supported by the field
	Choices map[string]string `json:"choices,omitempty"`
}

type ContactsClient interface {
	Create(t *Contact) (*Contact, error)
	Update(id uint64, t *Contact) (*Contact, error)
	View(id uint64) (*Contact, error)
	ListAll() ([]*Contact, error)
	Delete(id uint64) error
	HardDelete(id uint64, force bool) error
	Restore(id uint64) error
	ListAllContactFields() ([]*ContactField, error)
	SearchContacts(keyword string) (contacts []*Contact, err error)
	SendInvite(id uint64) error
	Merge(primaryID uint64, secondaryIDs []uint64, attrs *Contact) error
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
	err = c.client.do(req, res, http.StatusCreated)

	return res, err
}

// Update updates an existing contact
func (c *contactsClient) Update(id uint64, t *Contact) (*Contact, error) {
	req, err := c.client.newRequest(http.MethodPut, fmt.Sprintf("contacts/%d", id), t)
	if err != nil {
		return nil, err
	}

	res := new(Contact)
	err = c.client.do(req, res, http.StatusOK)

	return res, err
}

// View gets an existing contact by id
func (c *contactsClient) View(id uint64) (*Contact, error) {
	req, err := c.client.newRequest(http.MethodGet, fmt.Sprintf("contacts/%d", id), nil)
	if err != nil {
		return nil, err
	}

	res := new(Contact)
	err = c.client.do(req, res, http.StatusOK)

	return res, err
}

// ListAll lists all existing contacts
func (c *contactsClient) ListAll() ([]*Contact, error) {
	req, err := c.client.newRequest(http.MethodGet, "contacts", nil)
	if err != nil {
		return nil, err
	}

	var res []*Contact
	err = c.client.do(req, &res, http.StatusOK)

	return res, err
}

// Delete soft-deletes an existing contact
func (c *contactsClient) Delete(id uint64) error {
	req, err := c.client.newRequest(http.MethodDelete, fmt.Sprintf("contacts/%d", id), nil)
	if err != nil {
		return err
	}

	return c.client.do(req, nil, http.StatusOK)
}

// HardDelete deletes an existing contact permanently
func (c *contactsClient) HardDelete(id uint64, force bool) error {
	type params struct {
		// Contact ID
		ID uint64 `json:"id"`
		// Send as true to force hard delete of a contact that is not already soft deleted
		Force bool `json:"force"`
	}

	in, err := json.Marshal(&params{
		ID:    id,
		Force: force,
	})
	if err != nil {
		return err
	}

	req, err := c.client.newRequest(http.MethodDelete, fmt.Sprintf("contacts/%d/hard_delete", id), in)
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

// ListAllContactFields lists all contact fields
func (c *contactsClient) ListAllContactFields() ([]*ContactField, error) {
	req, err := c.client.newRequest(http.MethodGet, "contact_fields", nil)
	if err != nil {
		return nil, err
	}

	var res []*ContactField
	err = c.client.do(req, &res, http.StatusOK)

	return res, err
}

// SearchContacts searches for a contact using their name
func (c *contactsClient) SearchContacts(keyword string) (contacts []*Contact, err error) {
	req, err := c.client.newRequest(http.MethodPost, fmt.Sprintf("contacts/autocomplete?term=%s", keyword), nil)
	if err != nil {
		return
	}

	err = c.client.do(req, &contacts, http.StatusNoContent)
	return
}

// SendInvite used to send an activation email to an existing contact for email verification
func (c *contactsClient) SendInvite(id uint64) error {
	req, err := c.client.newRequest(http.MethodPost, fmt.Sprintf("contacts/%d/send_invite", id), nil)
	if err != nil {
		return err
	}

	return c.client.do(req, nil, http.StatusNoContent)
}

// Merge used to merge two or more duplicate contacts together.
func (c *contactsClient) Merge(primaryID uint64, secondaryIDs []uint64, attrs *Contact) error {
	type params struct {
		// ID of the primary contact
		PrimaryContactID uint64 `json:"primary_contact_id"`
		// IDs of the secondary contacts
		SecondaryContactIDs []uint64 `json:"secondary_contact_ids"`
		// Contains attributes that need to be updated in the primary contact during merge (optional)
		Contact *Contact `json:"contact,omitempty"`
	}

	in, err := json.Marshal(&params{
		PrimaryContactID:    primaryID,
		SecondaryContactIDs: secondaryIDs,
		Contact:             attrs,
	})
	if err != nil {
		return err
	}

	req, err := c.client.newRequest(http.MethodPost, "contacts/merge", in)
	if err != nil {
		return err
	}

	return c.client.do(req, nil, http.StatusNoContent)
}

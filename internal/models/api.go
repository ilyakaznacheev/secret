package models

// this file was generated with go-swagger but changed to support xml.
// Do not re-generate the model with go-swagger until it doesn't support xml.

import (
	"encoding/xml"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SecretResponse secret
// swagger:model Secret
type SecretResponse struct {
	XMLName xml.Name `json:"-" xml:"Secret"`

	// The date and time of the creation
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty" xml:"createdAt,omitempty"`

	// The secret cannot be reached after this time
	// Format: date-time
	ExpiresAt *strfmt.DateTime `json:"expiresAt,omitempty" xml:"expiresAt,omitempty"`

	// Unique hash to identify the secrets
	Hash string `json:"hash,omitempty" xml:"hash,omitempty"`

	// How many times the secret can be viewed
	RemainingViews int32 `json:"remainingViews,omitempty" xml:"remainingViews,omitempty"`

	// The secret itself
	SecretText string `json:"secretText,omitempty" xml:"secretText,omitempty"`
}

// Validate validates this secret
func (m *Secret) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExpiresAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Secret) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Secret) validateExpiresAt(formats strfmt.Registry) error {

	if swag.IsZero(m.ExpiresAt) { // not required
		return nil
	}

	if err := validate.FormatOf("expiresAt", "body", "date-time", m.ExpiresAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Secret) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Secret) UnmarshalBinary(b []byte) error {
	var res Secret
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

package domain

import (

	"github.com/google/uuid"

)

// AuthRequest struct
type AuthRequest struct {
	ID               uuid.UUID       `json:"-"`
	Identifier       *string         `json:"identifier"`
}

// FromAuthRequester TODO
func FromAuthRequester() (*AuthRequest, error) {
	res := AuthRequest{
		
	}

	return &res, nil
}


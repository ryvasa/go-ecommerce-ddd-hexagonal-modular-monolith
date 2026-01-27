package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrEmptyUserID   = errors.New("user ID is required")
	ErrEmptyTitle    = errors.New("title is required")
	ErrAlreadyDone   = errors.New("cart already done")
	ErrAlreadyUndone = errors.New("cart is not done")
)

type Cart struct {
	id     string
	userID string
	title  string
	done   bool
}

// NewCart creates a new Cart with validation
func NewCart(userID, title string) (*Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrEmptyUserID
	}
	if strings.TrimSpace(title) == "" {
		return nil, ErrEmptyTitle
	}
	return &Cart{
		id:     uuid.NewString(),
		userID: userID,
		title:  title,
		done:   false,
	}, nil
}

// RehydrateCart reconstructs a Cart from persistence
func RehydrateCart(id, userID, title string, done bool) *Cart {
	return &Cart{
		id:     id,
		userID: userID,
		title:  title,
		done:   done,
	}
}

// Getters
func (c *Cart) ID() string     { return c.id }
func (c *Cart) UserID() string { return c.userID }
func (c *Cart) Title() string  { return c.title }
func (c *Cart) IsDone() bool   { return c.done }

// Domain behaviors
func (c *Cart) MarkDone() error {
	if c.done {
		return ErrAlreadyDone
	}
	c.done = true
	return nil
}

func (c *Cart) MarkUndone() error {
	if !c.done {
		return ErrAlreadyUndone
	}
	c.done = false
	return nil
}

func (c *Cart) UpdateTitle(newTitle string) error {
	if strings.TrimSpace(newTitle) == "" {
		return ErrEmptyTitle
	}
	c.title = newTitle
	return nil
}

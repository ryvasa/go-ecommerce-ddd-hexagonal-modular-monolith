package entity

import (
	"testing"
)

func TestNewCart(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		title   string
		wantErr error
	}{
		{
			name:    "valid cart",
			userID:  "user-123",
			title:   "Shopping List",
			wantErr: nil,
		},
		{
			name:    "empty userID",
			userID:  "",
			title:   "Shopping List",
			wantErr: ErrEmptyUserID,
		},
		{
			name:    "whitespace userID",
			userID:  "   ",
			title:   "Shopping List",
			wantErr: ErrEmptyUserID,
		},
		{
			name:    "empty title",
			userID:  "user-123",
			title:   "",
			wantErr: ErrEmptyTitle,
		},
		{
			name:    "whitespace title",
			userID:  "user-123",
			title:   "   ",
			wantErr: ErrEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := NewCart(tt.userID, tt.title)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewCart() error = %v, wantErr %v", err, tt.wantErr)
				}
				if cart != nil {
					t.Error("NewCart() should return nil cart on error")
				}
			} else {
				if err != nil {
					t.Errorf("NewCart() unexpected error: %v", err)
				}
				if cart == nil {
					t.Fatal("NewCart() returned nil cart")
				}
				if cart.ID() == "" {
					t.Error("NewCart() should generate ID")
				}
				if cart.UserID() != tt.userID {
					t.Errorf("Cart.UserID() = %q, want %q", cart.UserID(), tt.userID)
				}
				if cart.Title() != tt.title {
					t.Errorf("Cart.Title() = %q, want %q", cart.Title(), tt.title)
				}
				if cart.IsDone() {
					t.Error("NewCart() should create cart with done=false")
				}
			}
		})
	}
}

func TestCart_MarkDone(t *testing.T) {
	tests := []struct {
		name     string
		cartDone bool
		wantErr  error
		wantDone bool
	}{
		{
			name:     "undone cart can be marked done",
			cartDone: false,
			wantErr:  nil,
			wantDone: true,
		},
		{
			name:     "already done cart returns error",
			cartDone: true,
			wantErr:  ErrAlreadyDone,
			wantDone: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := RehydrateCart("id", "user-123", "title", tt.cartDone)

			err := cart.MarkDone()

			if err != tt.wantErr {
				t.Errorf("Cart.MarkDone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if cart.IsDone() != tt.wantDone {
				t.Errorf("Cart.IsDone() = %v, want %v", cart.IsDone(), tt.wantDone)
			}
		})
	}
}

func TestCart_MarkUndone(t *testing.T) {
	tests := []struct {
		name     string
		cartDone bool
		wantErr  error
		wantDone bool
	}{
		{
			name:     "done cart can be marked undone",
			cartDone: true,
			wantErr:  nil,
			wantDone: false,
		},
		{
			name:     "already undone cart returns error",
			cartDone: false,
			wantErr:  ErrAlreadyUndone,
			wantDone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := RehydrateCart("id", "user-123", "title", tt.cartDone)

			err := cart.MarkUndone()

			if err != tt.wantErr {
				t.Errorf("Cart.MarkUndone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if cart.IsDone() != tt.wantDone {
				t.Errorf("Cart.IsDone() = %v, want %v", cart.IsDone(), tt.wantDone)
			}
		})
	}
}

func TestCart_UpdateTitle(t *testing.T) {
	tests := []struct {
		name     string
		newTitle string
		wantErr  error
	}{
		{
			name:     "valid new title",
			newTitle: "New Shopping List",
			wantErr:  nil,
		},
		{
			name:     "empty title",
			newTitle: "",
			wantErr:  ErrEmptyTitle,
		},
		{
			name:     "whitespace title",
			newTitle: "   ",
			wantErr:  ErrEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := RehydrateCart("id", "user-123", "Old Title", false)

			err := cart.UpdateTitle(tt.newTitle)

			if err != tt.wantErr {
				t.Errorf("Cart.UpdateTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && cart.Title() != tt.newTitle {
				t.Errorf("Cart.Title() = %q, want %q", cart.Title(), tt.newTitle)
			}
		})
	}
}

func TestRehydrateCart(t *testing.T) {
	cart := RehydrateCart("cart-id", "user-id", "My Cart", true)

	if cart.ID() != "cart-id" {
		t.Errorf("Cart.ID() = %q, want %q", cart.ID(), "cart-id")
	}
	if cart.UserID() != "user-id" {
		t.Errorf("Cart.UserID() = %q, want %q", cart.UserID(), "user-id")
	}
	if cart.Title() != "My Cart" {
		t.Errorf("Cart.Title() = %q, want %q", cart.Title(), "My Cart")
	}
	if !cart.IsDone() {
		t.Error("Cart.IsDone() = false, want true")
	}
}

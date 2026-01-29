package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
)

// Mock implementations for testing

type MockCartRepository struct {
	saveFunc    func(ctx context.Context, cart *entity.Cart) error
	findAllFunc func(ctx context.Context) ([]*entity.Cart, error)
}

func (m *MockCartRepository) Save(ctx context.Context, cart *entity.Cart) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, cart)
	}
	return nil
}

func (m *MockCartRepository) FindAll(ctx context.Context) ([]*entity.Cart, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return nil, nil
}

type MockUserReader struct {
	existsFunc func(ctx context.Context, userID string) (bool, error)
}

func (m *MockUserReader) Exists(ctx context.Context, userID string) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, userID)
	}
	return true, nil
}

type MockTxManager struct{}

func (m *MockTxManager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}

type MockEventPublisher struct{}

func (m *MockEventPublisher) Publish(ctx context.Context, events ...sharedEvent.Event) error {
	return nil
}

func TestCartService_Create(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		title       string
		repo        *MockCartRepository
		userReader  *MockUserReader
		wantErr     bool
		errContains string
	}{
		{
			name:   "successful cart creation",
			userID: "user-123",
			title:  "Shopping List",
			repo: &MockCartRepository{
				saveFunc: func(ctx context.Context, cart *entity.Cart) error {
					return nil
				},
			},
			userReader: &MockUserReader{
				existsFunc: func(ctx context.Context, userID string) (bool, error) {
					return true, nil
				},
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: "nonexistent-user",
			title:  "Shopping List",
			repo:   &MockCartRepository{},
			userReader: &MockUserReader{
				existsFunc: func(ctx context.Context, userID string) (bool, error) {
					return false, nil
				},
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "empty title validation fails",
			userID: "user-123",
			title:  "",
			repo:   &MockCartRepository{},
			userReader: &MockUserReader{
				existsFunc: func(ctx context.Context, userID string) (bool, error) {
					return true, nil
				},
			},
			wantErr:     true,
			errContains: "title",
		},
		{
			name:   "empty userID validation fails",
			userID: "",
			title:  "Shopping List",
			repo:   &MockCartRepository{},
			userReader: &MockUserReader{
				existsFunc: func(ctx context.Context, userID string) (bool, error) {
					return true, nil
				},
			},
			wantErr:     true,
			errContains: "user ID",
		},
		{
			name:   "repository save error",
			userID: "user-123",
			title:  "Shopping List",
			repo: &MockCartRepository{
				saveFunc: func(ctx context.Context, cart *entity.Cart) error {
					return errors.New("database error")
				},
			},
			userReader: &MockUserReader{
				existsFunc: func(ctx context.Context, userID string) (bool, error) {
					return true, nil
				},
			},
			wantErr:     true,
			errContains: "database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := &MockTxManager{}
			publisher := &MockEventPublisher{}
			service := NewCartService(tt.repo, tt.userReader, txManager, publisher)

			err := service.Create(context.Background(), tt.userID, tt.title)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Create() expected error, got nil")
				} else if tt.errContains != "" {
					errStr := err.Error()
					if !contains(errStr, tt.errContains) {
						t.Errorf("Create() error = %q, want containing %q", errStr, tt.errContains)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Create() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCartService_List(t *testing.T) {
	tests := []struct {
		name      string
		repo      *MockCartRepository
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty list",
			repo: &MockCartRepository{
				findAllFunc: func(ctx context.Context) ([]*entity.Cart, error) {
					return []*entity.Cart{}, nil
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "multiple carts",
			repo: &MockCartRepository{
				findAllFunc: func(ctx context.Context) ([]*entity.Cart, error) {
					cart1 := entity.RehydrateCart("1", "user-1", "Cart 1", false)
					cart2 := entity.RehydrateCart("2", "user-2", "Cart 2", true)
					return []*entity.Cart{cart1, cart2}, nil
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "repository error",
			repo: &MockCartRepository{
				findAllFunc: func(ctx context.Context) ([]*entity.Cart, error) {
					return nil, errors.New("database error")
				},
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := &MockTxManager{}
			userReader := &MockUserReader{}
			publisher := &MockEventPublisher{}
			service := NewCartService(tt.repo, userReader, txManager, publisher)

			carts, err := service.List(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Errorf("List() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("List() unexpected error: %v", err)
				}
				if len(carts) != tt.wantCount {
					t.Errorf("List() returned %d carts, want %d", len(carts), tt.wantCount)
				}
			}
		})
	}
}

// Helper function for simple substring check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

//go:build integration

package persistence

import (
	"context"
	"testing"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCartTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto migrate the cart model
	if err := db.AutoMigrate(&CartModel{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestGormCartRepository_Save(t *testing.T) {
	db := setupCartTestDB(t)
	repo := NewGormCartRepository(db)

	tests := []struct {
		name    string
		userID  string
		title   string
		wantErr bool
	}{
		{
			name:    "save new cart successfully",
			userID:  "user-123",
			title:   "My Shopping List",
			wantErr: false,
		},
		{
			name:    "save cart with different user",
			userID:  "user-456",
			title:   "Another List",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := entity.NewCart(tt.userID, tt.title)
			if err != nil {
				t.Fatalf("failed to create cart entity: %v", err)
			}

			err = repo.Save(context.Background(), cart)

			if tt.wantErr {
				if err == nil {
					t.Error("Save() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Save() unexpected error: %v", err)
				}

				// Verify cart was saved
				var model CartModel
				if err := db.First(&model, "id = ?", cart.ID()).Error; err != nil {
					t.Errorf("failed to find saved cart: %v", err)
				}
				if model.Title != tt.title {
					t.Errorf("saved title = %q, want %q", model.Title, tt.title)
				}
				if model.UserID != tt.userID {
					t.Errorf("saved userID = %q, want %q", model.UserID, tt.userID)
				}
				if model.Done != false {
					t.Error("saved Done should be false for new cart")
				}
			}
		})
	}
}

func TestGormCartRepository_FindAll(t *testing.T) {
	db := setupCartTestDB(t)
	repo := NewGormCartRepository(db)

	tests := []struct {
		name       string
		setupCarts int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "empty database returns empty slice",
			setupCarts: 0,
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear database for this test
			db.Exec("DELETE FROM carts")

			// Setup carts if needed
			for i := 0; i < tt.setupCarts; i++ {
				cart, _ := entity.NewCart("user-123", "Cart Title")
				repo.Save(context.Background(), cart)
			}

			carts, err := repo.FindAll(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("FindAll() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("FindAll() unexpected error: %v", err)
				}
				if len(carts) != tt.wantCount {
					t.Errorf("FindAll() returned %d carts, want %d", len(carts), tt.wantCount)
				}
			}
		})
	}
}

func TestGormCartRepository_FindAll_MultipleCarts(t *testing.T) {
	db := setupCartTestDB(t)
	repo := NewGormCartRepository(db)

	// Setup multiple carts
	cart1, _ := entity.NewCart("user-1", "Cart 1")
	cart2, _ := entity.NewCart("user-2", "Cart 2")
	cart3, _ := entity.NewCart("user-1", "Cart 3")

	_ = repo.Save(context.Background(), cart1)
	_ = repo.Save(context.Background(), cart2)
	_ = repo.Save(context.Background(), cart3)

	carts, err := repo.FindAll(context.Background())

	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}

	if len(carts) != 3 {
		t.Errorf("FindAll() returned %d carts, want 3", len(carts))
	}
}

func TestGormCartRepository_FindAll_RehydratesCartCorrectly(t *testing.T) {
	db := setupCartTestDB(t)
	repo := NewGormCartRepository(db)

	// Setup a cart
	originalCart, _ := entity.NewCart("user-123", "Test Cart")
	_ = repo.Save(context.Background(), originalCart)

	// Find all and verify rehydration
	carts, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}

	if len(carts) != 1 {
		t.Fatalf("expected 1 cart, got %d", len(carts))
	}

	cart := carts[0]

	if cart.ID() != originalCart.ID() {
		t.Errorf("cart.ID() = %q, want %q", cart.ID(), originalCart.ID())
	}
	if cart.UserID() != originalCart.UserID() {
		t.Errorf("cart.UserID() = %q, want %q", cart.UserID(), originalCart.UserID())
	}
	if cart.Title() != originalCart.Title() {
		t.Errorf("cart.Title() = %q, want %q", cart.Title(), originalCart.Title())
	}
	if cart.IsDone() != originalCart.IsDone() {
		t.Errorf("cart.IsDone() = %v, want %v", cart.IsDone(), originalCart.IsDone())
	}
}

package cart

import (
	"github.com/google/wire"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/adapter/in/http"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/adapter/out/persistence"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/service"
)

var Module = wire.NewSet(
	persistence.NewGormCartRepository,
	service.NewCartService,
	http.NewHandler,
)

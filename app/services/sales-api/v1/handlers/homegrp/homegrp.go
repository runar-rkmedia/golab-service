// Package homegrp maintains the group of handlers for home access.
package homegrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ardanlabs/service/business/core/home"
	"github.com/ardanlabs/service/business/data/page"
	"github.com/ardanlabs/service/business/web/v1/mid"
	"github.com/ardanlabs/service/business/web/v1/response"
	"github.com/ardanlabs/service/foundation/web"
)

// Set of error variables for handling home group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Handlers manages the set of home enpoints.
type Handlers struct {
	home *home.Core
}

// New constructs a handlers for route access.
func New(home *home.Core) *Handlers {
	return &Handlers{
		home: home,
	}
}

// Create adds a new home to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewHome
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nh, err := toCoreNewHome(ctx, app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	hme, err := h.home.Create(ctx, nh)
	if err != nil {
		return fmt.Errorf("create: hme[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppHome(hme), http.StatusCreated)
}

// Update updates a home in the system.
func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateHome
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	hme := mid.GetHome(ctx)

	updHme, err := toCoreUpdateHome(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	hme, err = h.home.Update(ctx, hme, updHme)
	if err != nil {
		return fmt.Errorf("update: homeID[%s] app[%+v]: %w", hme.ID, app, err)
	}

	return web.Respond(ctx, w, toAppHome(hme), http.StatusOK)
}

// Delete deletes a home from the system.
func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	hme := mid.GetHome(ctx)

	if err := h.home.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: homeID[%s]: %w", hme.ID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Query returns a list of homes with paging.
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := page.Parse(r)
	if err != nil {
		return err
	}

	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	homes, err := h.home.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.home.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, response.NewPageDocument(toAppHomes(homes), total, page.Number, page.RowsPerPage), http.StatusOK)
}

// QueryByID returns a home by its ID.
func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, toAppHome(mid.GetHome(ctx)), http.StatusOK)
}

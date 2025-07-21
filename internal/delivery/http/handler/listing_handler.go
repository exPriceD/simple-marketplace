package handler

import (
	"encoding/json"
	"net/http"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/mapping"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/middleware"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/response"
)

type ListingHandler struct {
	create usecase.CreateListingUseCase
	list   usecase.ListListingsUseCase
}

func NewListingHandler(create usecase.CreateListingUseCase, list usecase.ListListingsUseCase) *ListingHandler {
	return &ListingHandler{create: create, list: list}
}

var badJSONListErr = apperror.New("http.decode", "bad_json", apperror.KindValidation, nil)

func (h *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req mapping.CreateListingRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), badJSONListErr)
		return
	}
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()),
			apperror.New("listing.create", "unauthorized", apperror.KindAuth, nil))
		return
	}
	userLogin, _ := middleware.UserLogin(r.Context())
	out, err := h.create.Execute(r.Context(), usecase.CreateListingInput{
		AuthorID:    userID,
		AuthorLogin: userLogin,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
	})
	if err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), err)
		return
	}
	out.IsOwner = true
	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"listing": out,
	})
}

func (h *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	q := mapping.ParseListQuery(r)
	currentUserID, _ := middleware.UserID(r.Context())
	out, err := h.list.Execute(r.Context(), usecase.ListListingsInput{
		Limit:         q.Limit,
		Offset:        q.Offset,
		SortBy:        q.SortBy,
		SortDir:       q.SortDir,
		PriceMin:      q.PriceMin,
		PriceMax:      q.PriceMax,
		CurrentUserID: currentUserID,
	})
	if err != nil {
		response.WriteAppError(w, middleware.RequestIDValue(r.Context()), err)
		return
	}

	if currentUserID != "" {
		for i := range out.Items {
			if out.Items[i].AuthorID == currentUserID {
				out.Items[i].IsOwner = true
			}
		}
	}
	response.WriteJSON(w, http.StatusOK, out)
}

package hcategory

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/888NiKiToS888/catalog-service/internal/app/entity"
	rhandler "github.com/888NiKiToS888/catalog-service/internal/app/handler/http"
	"github.com/888NiKiToS888/catalog-service/internal/app/service"
	"github.com/888NiKiToS888/catalog-service/internal/pkg/http/binding"
	"github.com/888NiKiToS888/catalog-service/internal/pkg/http/httph"
)

type handler struct {
	svcCategory service.Category
}

func NewHandler(svcCategory service.Category) rhandler.Category {
	return &handler{svcCategory: svcCategory}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestCategoryCreate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.svcCategory.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := entity.ResponseCategory{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	httph.SendEncoded(w, r, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, entity.ErrIncorrectParameters.Error())
		return
	}

	category, err := h.svcCategory.GetByGUID(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, err.Error())
		} else {
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := entity.ResponseCategory{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, entity.ErrIncorrectParameters.Error())
		return
	}

	var req entity.RequestCategoryUpdate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.svcCategory.Update(r.Context(), guid, req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.ErrorApply(w, http.StatusNotFound, err.Error())
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := entity.ResponseCategory{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, entity.ErrIncorrectParameters.Error())
		return
	}

	err = h.svcCategory.Delete(r.Context(), guid)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.ErrorApply(w, http.StatusNotFound, err.Error())
		case errors.Is(err, entity.ErrCategoryHasProducts):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	httph.SendEmpty(w, http.StatusNoContent)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svcCategory.List(r.Context())
	if err != nil {
		httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]entity.ResponseCategory, len(categories))
	for i, cat := range categories {
		resp[i] = entity.ResponseCategory{
			GUID:      cat.GUID,
			Name:      cat.Name,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		}
	}
	httph.SendEncoded(w, r, http.StatusOK, resp)
}

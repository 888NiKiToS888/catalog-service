package sproduct

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/888NiKiToS888/catalog-service/internal/app/entity"
	"github.com/888NiKiToS888/catalog-service/internal/app/repository"
	"github.com/888NiKiToS888/catalog-service/internal/app/service"
)

type svc struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &svc{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *svc) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	existing, err := s.repoProduct.List(ctx, &req.Name, nil)
	if err != nil {
		return entity.Product{}, err
	}
	if len(existing) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}

	_, err = s.repoCategory.GetByGUID(ctx, req.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}

	now := time.Now()
	product := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repoProduct.Create(ctx, product); err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *svc) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	return s.repoProduct.GetByGUID(ctx, guid)
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	if req.Name != nil && *req.Name != product.Name {
		existing, err := s.repoProduct.List(ctx, req.Name, nil)
		if err != nil {
			return entity.Product{}, err
		}
		if len(existing) > 0 {
			return entity.Product{}, entity.ErrAlreadyExists
		}
		product.Name = *req.Name
	}

	if req.CategoryGUID != nil && *req.CategoryGUID != product.CategoryGUID {
		_, err := s.repoCategory.GetByGUID(ctx, *req.CategoryGUID)
		if err != nil {
			return entity.Product{}, err
		}
		product.CategoryGUID = *req.CategoryGUID
	}

	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}

	product.UpdatedAt = time.Now()

	if err := s.repoProduct.Update(ctx, product); err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}
	return s.repoProduct.Delete(ctx, guid)
}

func (s *svc) List(ctx context.Context) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, nil, nil)
}

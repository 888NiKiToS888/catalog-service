package pproduct

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/888NiKiToS888/catalog-service/internal/app/entity"
	"github.com/888NiKiToS888/catalog-service/internal/app/repository"
	rcpostgres "github.com/888NiKiToS888/catalog-service/internal/app/repository/conn/postgres"
	"github.com/888NiKiToS888/catalog-service/internal/app/util"
)

type (
	repoPg struct {
		*_DB
	}
	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Product {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) error {
	_, err := r.NewInsert().
		Model(&product).
		Exec(ctx)
	return err
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var product entity.Product
	err := r.NewSelect().
		Model(&product).
		Where("guid = ?", guid).
		Scan(ctx)
	if err != nil {
		return product, util.ReplaceErr1(err, sql.ErrNoRows, entity.ErrNotFound)
	}
	return product, nil
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) error {
	res, err := r.NewUpdate().
		Model(&product).
		WherePK().
		ExcludeColumn("id", "created_at").
		Exec(ctx)
	if err != nil {
		return err
	}
	return rcpostgres.UpdateErr(res, err)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.NewDelete().
		Model((*entity.Product)(nil)).
		Where("guid = ?", guid).
		Exec(ctx)
	return rcpostgres.DeleteErr(err)
}

func (r *repoPg) List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	query := r.NewSelect().Model(&products)
	if name != nil {
		query = query.Where("name = ?", *name)
	}
	if categoryGUID != nil {
		query = query.Where("category_guid = ?", *categoryGUID)
	}
	err := query.Scan(ctx)
	return products, err
}

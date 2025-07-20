package listingrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/domain/listing"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/database/postgres"
	"github.com/jackc/pgx/v5"
)

// ListingRepositoryPG — реализация порта ListingRepository для Postgres.
type ListingRepositoryPG struct {
	db *postgres.DB
}

func NewListingRepositoryPG(db *postgres.DB) *ListingRepositoryPG {
	return &ListingRepositoryPG{db: db}
}

// Create вставляет запись и возвращает *Listing с назначенным ID (и фактическим created_at из БД).
func (r *ListingRepositoryPG) Create(ctx context.Context, l *listing.Listing) (*listing.Listing, error) {
	row := ListingToRow(l)
	var id int64
	var createdAtPg = row.CreatedAt
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO listings (title,description,image_url,price,author_id,created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at
	`, row.Title, row.Description, row.ImageURL, row.Price, row.AuthorID, row.CreatedAt).Scan(&id, &createdAtPg)
	if err != nil {
		return nil, err
	}

	titleVO, _ := listing.NewTitle(row.Title)
	descVO, _ := listing.NewDescription(row.Description)
	imgVO, _ := listing.NewImageURL(row.ImageURL)
	priceVO, _ := listing.NewPrice(row.Price)

	return listing.RehydrateListing(id, titleVO, descVO, imgVO, priceVO, row.AuthorID, createdAtPg), nil
}

func (r *ListingRepositoryPG) List(ctx context.Context, f port.ListingFilter) ([]*listing.Listing, error) {
	var (
		args       []any
		clauses    []string
		argIndex   = 1
		orderField = "created_at"
	)

	if f.PriceMin != nil {
		clauses = append(clauses, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *f.PriceMin)
		argIndex++
	}
	if f.PriceMax != nil {
		clauses = append(clauses, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *f.PriceMax)
		argIndex++
	}

	if f.SortBy == "price" {
		orderField = "price"
	}
	orderDir := "DESC"
	if strings.ToLower(f.SortDir) == "asc" {
		orderDir = "ASC"
	}

	query := `
		SELECT id, title, description, image_url, price, author_id, created_at
		FROM listings
	`
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY %s %s, id ASC", orderField, orderDir)

	limit := f.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*listing.Listing
	for rows.Next() {
		var lr ListingRow
		if err := rows.Scan(&lr.ID, &lr.Title, &lr.Description, &lr.ImageURL, &lr.Price, &lr.AuthorID, &lr.CreatedAt); err != nil {
			return nil, err
		}
		l, err := RowToListing(lr)
		if err != nil {
			return nil, err
		}
		results = append(results, l)
	}
	return results, rows.Err()
}

func (r *ListingRepositoryPG) Count(ctx context.Context, f port.ListingFilter) (int64, error) {
	var (
		args     []any
		clauses  []string
		argIndex = 1
	)
	if f.PriceMin != nil {
		clauses = append(clauses, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *f.PriceMin)
		argIndex++
	}
	if f.PriceMax != nil {
		clauses = append(clauses, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *f.PriceMax)
		argIndex++
	}

	query := `SELECT COUNT(*) FROM listings`
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	var cnt int64
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}

var _ port.ListingRepository = (*ListingRepositoryPG)(nil)

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/adesokandijam/go-rest-api-boilerplate/internal/validator"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Item struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     int       `json:"version"`
}
type ItemModel struct{ DB *sql.DB }

func ValidateItem(v *validator.Validator, item *Item) {
	v.Check(item.Name != "", "name", "must be provided")
	v.Check(len(item.Name) <= 100, "name", "must not exceed 100 characters")
	v.Check(len(item.Description) <= 1000, "description", "must not exceed 1000 characters")
}
func (m ItemModel) Insert(item *Item) error {
	query := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at, version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, item.Name, item.Description).Scan(&item.ID, &item.CreatedAt, &item.Version)
}
func (m ItemModel) Get(id int64) (*Item, error) {
	query := `SELECT id, created_at, name, description, version FROM items WHERE id = $1`
	var item Item
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.CreatedAt, &item.Name, &item.Description, &item.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRecordNotFound
	}
	return &item, err
}
func (m ItemModel) GetAll() ([]*Item, error) {
	query := `SELECT id, created_at, name, description, version FROM items ORDER BY id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Item{}
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.CreatedAt, &item.Name, &item.Description, &item.Version); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}
func (m ItemModel) Update(item *Item) error {
	query := `UPDATE items SET name=$1, description=$2, version=version+1 WHERE id=$3 AND version=$4 RETURNING version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, item.Name, item.Description, item.ID, item.Version).Scan(&item.Version)
}
func (m ItemModel) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, `DELETE FROM items WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return ErrRecordNotFound
	}
	return err
}

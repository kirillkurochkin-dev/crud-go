package psql

import (
	"context"
	"crud-go/internal/entity"
	"database/sql"
)

type Phones struct {
	db *sql.DB
}

func NewPhone(db *sql.DB) *Phones {
	return &Phones{db: db}
}

func (p *Phones) GetPhoneById(ctx context.Context, id int64) (entity.Phone, error) {
	var ph entity.Phone
	err := p.db.QueryRow("SELECT * FROM phones WHERE id = $1", id).
		Scan(&ph.Id, &ph.Brand, &ph.Model, &ph.Year, &ph.OS, &ph.Processor)
	if err != nil {
		return ph, err
	}

	return ph, nil
}

func (p *Phones) GetAllPhones(ctx context.Context) ([]entity.Phone, error) {
	var phones []entity.Phone

	rows, err := p.db.Query("SELECT * FROM phones")
	if err != nil {
		return phones, err
	}

	for rows.Next() {
		var ph entity.Phone
		err := rows.Scan(&ph.Id, &ph.Brand, &ph.Model, &ph.Year, &ph.OS, &ph.Processor)
		if err != nil {
			return phones, err
		}
		phones = append(phones, ph)
	}

	return phones, nil
}

func (p *Phones) CreatePhone(ctx context.Context, ph entity.PhoneInputDto) error {
	_, err := p.db.Exec("INSERT INTO phones (brand, model, year, os, processor) VALUES ($1, $2, $3, $4, $5)",
		ph.Brand, ph.Model, ph.Year, ph.OS, ph.Processor)
	return err
}

func (p *Phones) UpdatePhoneById(ctx context.Context, id int64, ph entity.PhoneInputDto) error {
	_, err := p.db.Exec("UPDATE phones SET brand=$1, model=$2, year=$3, os=$4, processor=$5 WHERE id=$6",
		ph.Brand, ph.Model, ph.Year, ph.OS, ph.Processor, id)
	return err
}

func (p *Phones) DeletePhoneById(ctx context.Context, id int64) error {
	_, err := p.db.Exec("DELETE FROM phones WHERE id = $1", id)
	return err
}

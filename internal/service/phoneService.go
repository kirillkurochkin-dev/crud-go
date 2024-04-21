package service

import (
	"context"
	"crud-go/internal/entity"
)

type PhonesRepository interface {
	GetPhoneById(ctx context.Context, id int) (entity.Phone, error)
	GetAllPhones(ctx context.Context) ([]entity.Phone, error)
	CreatePhone(ctx context.Context, ph entity.PhoneInputDto) error
	UpdatePhoneById(ctx context.Context, id int, ph entity.PhoneInputDto) error
	DeletePhoneById(ctx context.Context, id int) error
}

type Phones struct {
	repository PhonesRepository
}

func NewPhones(repository PhonesRepository) *Phones {
	return &Phones{
		repository: repository,
	}
}

func (p *Phones) GetPhoneById(ctx context.Context, id int) (entity.Phone, error) {
	return p.repository.GetPhoneById(ctx, id)
}

func (p *Phones) GetAllPhones(ctx context.Context) ([]entity.Phone, error) {
	return p.repository.GetAllPhones(ctx)
}

func (p *Phones) CreatePhone(ctx context.Context, ph entity.PhoneInputDto) error {
	return p.repository.CreatePhone(ctx, ph)
}

func (p *Phones) UpdatePhoneById(ctx context.Context, id int, ph entity.PhoneInputDto) error {
	return p.repository.UpdatePhoneById(ctx, id, ph)
}

func (p *Phones) DeletePhoneById(ctx context.Context, id int) error {
	return p.repository.DeletePhoneById(ctx, id)
}

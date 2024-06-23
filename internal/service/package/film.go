package _package

import (
	"homework/internal/models"
)

const (
	FilmPrice models.Price       = 1
	FilmType  models.PackageType = "film"
)

// FilmPackage implements ValidatePackage
type FilmPackage struct {
}

func NewFilmPackage() *FilmPackage {
	return &FilmPackage{}
}

// ValidatePackage provides validation
func (p *FilmPackage) Validate(weight models.Weight) error {
	return nil
}

func (c *FilmPackage) Apply(order *models.Order) {
	order.PackageType = FilmType
	order.PackagePrice = FilmPrice
	order.OrderPrice += FilmPrice
}

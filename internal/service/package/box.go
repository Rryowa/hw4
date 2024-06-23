package _package

import (
	"homework/internal/models"
	"homework/internal/util"
)

const (
	BoxPrice     models.Price       = 20
	MaxBoxWeight models.Weight      = 30
	BoxType      models.PackageType = "box"
)

type BoxPackage struct {
}

func NewBoxPackage() *BoxPackage {
	return &BoxPackage{}
}

func (p *BoxPackage) Validate(weight models.Weight) error {
	if weight < MaxBoxWeight {
		return nil
	}
	return util.ErrWeightExceeds
}

func (c *BoxPackage) Apply(order *models.Order) {
	order.PackageType = BoxType
	order.PackagePrice = BoxPrice
	order.OrderPrice += BoxPrice
}

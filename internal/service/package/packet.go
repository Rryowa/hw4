package _package

import (
	"homework/internal/models"
	"homework/internal/util"
)

const (
	PacketPrice     models.Price       = 5
	MaxPacketWeight models.Weight      = 10
	PacketType      models.PackageType = "packet"
)

type PacketPackage struct {
}

func NewPacketPackage() *PacketPackage {
	return &PacketPackage{}
}

func (p *PacketPackage) Validate(weight models.Weight) error {
	if weight < MaxPacketWeight {
		return nil
	}
	return util.ErrWeightExceeds
}

func (c *PacketPackage) Apply(order *models.Order) {
	order.PackageType = PacketType
	order.PackagePrice = PacketPrice
	order.OrderPrice += PacketPrice
}

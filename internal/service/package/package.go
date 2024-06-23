package _package

import (
	"homework/internal/models"
	"homework/internal/util"
)

// TODO: add LRU cache using Strategy pattern
type packageContext struct {
	strategies map[models.PackageType]PackageStrategy
}

type PackageService interface {
	ValidatePackage(weight models.Weight, packageType models.PackageType) error
	ApplyPackage(order *models.Order, packageType models.PackageType)
}

type PackageStrategy interface {
	Validate(weight models.Weight) error
	Apply(order *models.Order)
}

func NewPackageService() PackageService {
	return &packageContext{
		strategies: map[models.PackageType]PackageStrategy{
			FilmType:   NewFilmPackage(),
			PacketType: NewPacketPackage(),
			BoxType:    NewBoxPackage(),
		},
	}
}

func (pc *packageContext) ValidatePackage(weight models.Weight, packageType models.PackageType) error {
	if strategy, ok := pc.strategies[packageType]; ok {
		return strategy.Validate(weight)
	}
	return util.ErrPackageTypeInvalid
}

func (pc *packageContext) ApplyPackage(order *models.Order, packageType models.PackageType) {
	if strategy, ok := pc.strategies[packageType]; ok {
		strategy.Apply(order)
		return
	}
	//Assume that film has no weight limit
	pc.strategies[FilmType].Apply(order)
}

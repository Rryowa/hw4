package _package

import (
	"homework/internal/models"
	"homework/internal/util"
)

type PackageService interface {
	ValidatePackage(weight models.Weight, packageType models.PackageType) error
	ApplyPackage(order *models.Order, packageType models.PackageType)
}

// TODO: refactor from switch to map
// TODO: add LRU cache using Strategy pattern
// TODO: refactor PackageService to strategy pattern
type packageService struct {
	filmPackage   *FilmPackage
	packetPackage *PacketPackage
	boxPackage    *BoxPackage
}

func NewPackageService() PackageService {
	return &packageService{
		filmPackage:   NewFilmPackage(),
		packetPackage: NewPacketPackage(),
		boxPackage:    NewBoxPackage(),
	}
}

func (ps *packageService) ValidatePackage(weight models.Weight, packageType models.PackageType) error {
	switch packageType {
	case FilmType:
		if err := ps.filmPackage.Validate(weight); err != nil {
			return err
		}
	case PacketType:
		if err := ps.packetPackage.Validate(weight); err != nil {
			return err
		}
	case BoxType:
		if err := ps.boxPackage.Validate(weight); err != nil {
			return err
		}
	case "":
		return nil
	default:
		return util.ErrPackageTypeInvalid
	}
	return nil
}

func (ps *packageService) ApplyPackage(order *models.Order, packageType models.PackageType) {
	switch packageType {
	case FilmType:
		ps.filmPackage.Apply(order)
	case PacketType:
		ps.packetPackage.Apply(order)
	case BoxType:
		ps.boxPackage.Apply(order)
	case "":
		ps.filmPackage.Apply(order)
	}
}

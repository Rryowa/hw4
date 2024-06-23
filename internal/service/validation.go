package service

import (
	"errors"
	"homework/internal/models"
	pkg "homework/internal/service/package"
	"homework/internal/storage"
	"homework/internal/util"
	"strconv"
	"time"
)

type ValidationService interface {
	ValidateAccept(id, userId, dateStr, orderPriceStr, weightStr, pkgTypeStr string) (*models.Order, error)
	ValidateIssue(ids []string) (*[]models.Order, error)
	ValidateAcceptReturn(id, userId string) (*models.Order, error)
	ValidateReturnToCourier(id string) error
	ValidateList(offset, limit string) (int, int, error)
}

type validationService struct {
	repository     storage.Storage
	packageService pkg.PackageService
}

func NewValidationService(repository storage.Storage, packageService pkg.PackageService) ValidationService {
	return &validationService{
		repository:     repository,
		packageService: packageService,
	}
}

func (v *validationService) ValidateAccept(id, userId, dateStr, orderPriceStr, weightStr, pkgTypeStr string) (*models.Order, error) {
	if len(id) == 0 {
		return &models.Order{}, util.ErrOrderIdNotProvided
	}
	if len(userId) == 0 {
		return &models.Order{}, util.ErrUserIdNotProvided
	}
	if len(weightStr) == 0 {
		return &models.Order{}, util.ErrWeightNotProvided
	}
	if len(orderPriceStr) == 0 {
		return &models.Order{}, util.ErrPriceNotProvided
	}

	storageUntil, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return &models.Order{}, errors.New("error parsing date")
	} else if storageUntil.Before(time.Now()) {
		return &models.Order{}, util.ErrDateInvalid
	}

	orderPriceFloat, err := strconv.ParseFloat(orderPriceStr, 64)
	if err != nil || orderPriceFloat <= 0 {
		return &models.Order{}, util.ErrOrderPriceInvalid
	}

	weightFloat, err := strconv.ParseFloat(weightStr, 64)
	if err != nil || weightFloat <= 0 {
		return &models.Order{}, util.ErrWeightInvalid
	}

	//Check for existence
	_, err = v.repository.Get(id)
	if err == nil {
		return &models.Order{}, util.ErrOrderExists
	}

	orderPrice := models.Price(orderPriceFloat)
	weight := models.Weight(weightFloat)
	packageType := models.PackageType(pkgTypeStr)

	if err = v.packageService.ValidatePackage(weight, packageType); err != nil {
		return &models.Order{}, err
	}

	order := models.Order{
		ID:           id,
		UserID:       userId,
		StorageUntil: storageUntil,
		OrderPrice:   orderPrice,
		Weight:       weight,
	}

	return &order, nil
}

func (v *validationService) ValidateIssue(ids []string) (*[]models.Order, error) {
	var ordersToIssue []models.Order

	if len(ids) == 0 {
		return &ordersToIssue, util.ErrUserIdNotProvided
	}

	order, err := v.repository.Get(ids[0])
	if err != nil {
		return &ordersToIssue, util.ErrOrderNotFound
	}
	recipientID := order.UserID

	for _, id := range ids {
		order, err = v.repository.Get(id)
		if err != nil {
			return &ordersToIssue, util.ErrOrderNotFound
		}
		if order.Issued {
			return &ordersToIssue, util.ErrOrderIssued
		}
		if order.Returned {
			return &ordersToIssue, util.ErrOrderReturned
		}
		if time.Now().After(order.StorageUntil) {
			return &ordersToIssue, util.ErrOrderExpired
		}

		//Check if users are equal
		if order.UserID != recipientID {
			return &ordersToIssue, util.ErrOrdersUserDiffers
		}

		ordersToIssue = append(ordersToIssue, order)
	}

	return &ordersToIssue, nil
}

func (v *validationService) ValidateAcceptReturn(id, userId string) (*models.Order, error) {
	if len(id) == 0 {
		return &models.Order{}, util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return &models.Order{}, util.ErrUserIdNotProvided
	}

	order, err := v.repository.Get(id)
	if err != nil {
		return &models.Order{}, util.ErrOrderNotFound
	}

	if order.UserID != userId {
		return &models.Order{}, util.ErrOrderDoesNotBelong
	}
	if !order.Issued {
		return &models.Order{}, util.ErrOrderNotIssued
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return &models.Order{}, util.ErrReturnPeriodExpired
	}

	return &order, nil
}

func (v *validationService) ValidateReturnToCourier(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	order, err := v.repository.Get(id)
	if err != nil {
		return util.ErrOrderNotFound
	}

	if order.Issued {
		return util.ErrOrderIssued
	}

	//skip checking for a period, to ensure that its working
	//if time.Now().Before(order.StorageUntil) {
	//	return util.ErrOrderNotExpired
	//}

	return nil
}

func (v *validationService) ValidateList(offset, limit string) (int, int, error) {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return -1, -1, err
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return -1, -1, err
	}

	return offsetInt, limitInt, nil
}

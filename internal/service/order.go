package service

import (
	"fmt"
	"homework/internal/models"
	pkg "homework/internal/service/package"
	"homework/internal/storage"
	"homework/pkg/hash"
	"strings"
	"time"
)

type OrderService interface {
	Accept(order *models.Order, pkgTypeStr string) error
	Issue(ordersToIssue *[]models.Order) error
	Return(orders *models.Order) error
	ReturnToCourier(id string) error
	ListReturns(offset, limit int) ([]models.Order, error)
	ListOrders(userId string, offset, limit int) ([]models.Order, error)
	PrintList(orders []models.Order)
}

type orderService struct {
	repository     storage.Storage
	packageService pkg.PackageService
}

func NewOrderService(repository storage.Storage, packageService pkg.PackageService) OrderService {
	return &orderService{
		repository:     repository,
		packageService: packageService,
	}
}

func (os *orderService) Accept(order *models.Order, pkgTypeStr string) error {
	os.packageService.ApplyPackage(order, models.PackageType(pkgTypeStr))

	fmt.Print("Calculating hash.")

	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Print(" .")
			}
		}
	}()

	go func(order *models.Order, ticker *time.Ticker, done chan struct{}) {
		order.Hash = hash.GenerateHash()
		ticker.Stop()
		done <- struct{}{}
	}(order, ticker, done)

	<-done

	return os.repository.Insert(*order)
}

func (os *orderService) Issue(orders *[]models.Order) error {
	for _, order := range *orders {
		order.Issued = true
		order.IssuedAt = time.Now()
	}

	return os.repository.IssueUpdate(*orders)
}

func (os *orderService) Return(order *models.Order) error {
	order.Returned = true

	return os.repository.Update(*order)
}

func (os *orderService) ReturnToCourier(id string) error {
	return os.repository.Delete(id)
}

func (os *orderService) ListReturns(offset, limit int) ([]models.Order, error) {
	return os.repository.GetReturns(offset, limit)
}

func (os *orderService) ListOrders(userId string, offset, limit int) ([]models.Order, error) {
	return os.repository.GetOrders(userId, offset, limit)
}

func (os *orderService) PrintList(orders []models.Order) {
	if len(orders) == 0 {
		defer fmt.Printf("\n\n")
	}
	fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n", "id", "user_id", "storage_until", "issued_at", "returned", "order_price", "weight", "package_type", "package_price")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range orders {
		fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n",
			order.ID,
			order.UserID,
			order.StorageUntil.Format("2006-01-02"),
			order.IssuedAt.Format("2006-01-02"),
			order.Returned,
			order.OrderPrice,
			order.Weight,
			order.PackageType,
			order.PackagePrice)
	}
	fmt.Printf("\n")
}

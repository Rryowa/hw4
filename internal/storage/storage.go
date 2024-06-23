package storage

import "homework/internal/models"

//go:generate mockery --name Storage
type Storage interface {
	Insert(order models.Order) error
	Update(order models.Order) error
	IssueUpdate(orders []models.Order) error
	Delete(id string) error
	Get(id string) (models.Order, error)
	GetReturns(offset, limit int) ([]models.Order, error)
	GetOrders(userId string, offset, limit int) ([]models.Order, error)
}

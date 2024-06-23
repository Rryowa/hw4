package models

import (
	"time"
)

type Price float64
type Weight float64
type PackageType string

type Order struct {
	ID           string      `db:"id"`
	UserID       string      `db:"user_id"`
	StorageUntil time.Time   `db:"storage_until"`
	Issued       bool        `db:"issued"`
	IssuedAt     time.Time   `db:"issued_at"`
	Returned     bool        `db:"returned"`
	OrderPrice   Price       `db:"order_price"`
	Weight       Weight      `db:"weight"`
	PackageType  PackageType `db:"package_type"`
	PackagePrice Price       `db:"package_price"`
	Hash         string      `db:"hash"`
}

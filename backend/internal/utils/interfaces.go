package utils

import "gorm.io/gorm"

type Repository[T any] interface {
	createFilter(qr *gorm.DB, query *T)
}

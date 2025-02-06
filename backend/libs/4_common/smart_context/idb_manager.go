package smart_context

import "gorm.io/gorm"

type IDbManager interface {
	GetGORM() *gorm.DB
}

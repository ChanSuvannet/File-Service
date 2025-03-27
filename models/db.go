// models/db.go
package models

import (
	"gorm.io/gorm"
)

// DB will hold the global database connection instance
var DB *gorm.DB

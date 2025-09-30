package db

import (
	"log"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for all models
// Note: This will create tables if they don't exist, but won't modify existing columns
func AutoMigrate(db *gorm.DB) error {
	log.Println("🔄 Running database auto-migration...")
	log.Println("⚠️  Note: Tables from schema.sql already exist, skipping migration")
	log.Println("   To rebuild tables, run: mysql -u root -p go_react_admin < docs/schema.sql")

	// Skip migration if you already imported schema.sql
	// If you want to use GORM auto-migration, drop all tables first:
	// mysql -u root -p -e "DROP DATABASE go_react_admin; CREATE DATABASE go_react_admin;"

	// // Core tables (Phase 1)
	// err := db.AutoMigrate(
	// 	&model.User{},
	// 	&model.Role{},
	// 	&model.Dept{},
	// 	&model.Menu{},
	// 	&model.RolePermission{},
	// 	&model.UserRole{},
	// 	&model.RoleMenu{},
	// 	&model.LoginLog{},
	// )
	// if err != nil {
	// 	return err
	// }

	// log.Println("✅ Core tables migrated successfully")

	// // Phase 2 tables (optional for now)
	// err = db.AutoMigrate(
	// 	&model.OperationLog{},
	// 	&model.AuditLog{},
	// )
	// if err != nil {
	// 	log.Printf("⚠️  Phase 2 tables migration warning: %v", err)
	// } else {
	// 	log.Println("✅ Phase 2 tables migrated successfully")
	// }

	log.Println("✅ Database auto-migration skipped (using schema.sql)")

	// Verify models are defined (no actual migration)
	_ = model.User{}
	_ = model.Role{}
	_ = model.Dept{}
	_ = model.Menu{}
	_ = model.RolePermission{}

	return nil
}

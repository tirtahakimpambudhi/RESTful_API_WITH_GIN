package helper

import (
	"gorm.io/gorm"
)

func CommitOrRollback(tx *gorm.DB) {
	err := recover()
	defer tx.Rollback()
	if err != nil {
		tx.Commit()
		return
	}
}

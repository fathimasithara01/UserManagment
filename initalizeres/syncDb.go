package initalizeres

import "github.com/fathima-sithara/UserManagment/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}

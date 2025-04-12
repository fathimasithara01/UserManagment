package initalizeres

import (
	"log"
	"time"
)

func Pooling() {
	sqldb, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get sql db")
		return
	}

	sqldb.SetMaxOpenConns(100)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(30 * time.Minute)

}

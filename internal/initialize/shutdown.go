package initialize

import "dyc/internal/db"

func shutdown()  {
	db.Close()
}
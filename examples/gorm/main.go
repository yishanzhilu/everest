package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Product .
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Read
	var product Product

	err = db.First(&product, "code = ?", "L12133").Error // find product with code l1212
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Print("sql no")
		} else {
			fmt.Print(err)
		}
	}
}

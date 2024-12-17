package main

import (
	"log"

	_ "github.com/lib/pq"
	wire_excel "github.com/ryvasa/go-super-farmer/cmd/excel/pkg/wire"
)

func main() {
	_, err := wire_excel.InitializeExcelApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
}

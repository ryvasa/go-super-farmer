package main

import (
	_ "github.com/lib/pq"
	wire_excel "github.com/ryvasa/go-super-farmer/cmd/excel/pkg/wire"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
)

func main() {
	_, err := wire_excel.InitializeExcelApp()
	if err != nil {
		logrus.Log.Fatalf("failed to initialize app: %v", err)
	}
}

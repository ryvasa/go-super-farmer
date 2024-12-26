package main

import (
	_ "github.com/lib/pq"
	wire_excel "github.com/ryvasa/go-super-farmer/cmd/mail/pkg/wire"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
)

func main() {
	app, err := wire_excel.InitializeMailApp()
	defer app.RabbitMQ.Close()

	if err != nil {
		logrus.Log.Fatalf("failed to initialize app: %v", err)
	}
}

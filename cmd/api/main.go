package main

import (
	"github.com/ryvasa/go-super-farmer/cmd/api/pkg/wire"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
)

func main() {
	app, err := wire.InitializeApp()
	if err != nil {
		logrus.Log.Fatal("failed to initialize app", err)
	}

	if err := app.Start(); err != nil {
		logrus.Log.Fatal("failed to initialize app", err)
	}
}

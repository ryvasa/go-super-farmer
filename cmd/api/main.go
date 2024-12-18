package main

import (
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/wire"
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

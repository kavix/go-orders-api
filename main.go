package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/kavix/orders-api-golang/application"
)

func main() {
	app := application.New()

	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Fail to Start %s", err)
	}

}

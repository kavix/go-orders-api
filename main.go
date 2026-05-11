package main

import (
	"context"
	"fmt"

	"github.com/kavix/orders-api-golang/application"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Fail to Start %s", err)
	}

}

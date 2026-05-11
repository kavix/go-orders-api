package handeler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an Order")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Print("List all Orders")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Get Order By Id")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Order Update")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Order Deleted")
}

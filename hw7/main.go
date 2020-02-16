package main

import (
	"fmt"
	"github.com/sabitrayan/hw7"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	bookStore, err := book_store.NewBookStore("books.json")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Methods("GET").Path("/{id}").HandlerFunc(.Endpoints.GetBook(bookStore, "id"))
	router.Methods("POST").Path("/").HandlerFunc(bookCreateBook(bookStore))
	router.Methods("GET").Path("/exit").HandlerFunc(ExitWithSave(bookStore))
	router.Methods("LIST").Path("/list").HandlerFunc(ListBooks(bookStore))
	router.Methods("UPDATE").Path("/up").HandlerFunc(UpdateBook(bookStore,"id"))
	router.Methods("DELETE").Path("/del").HandlerFunc(DeleteBook(bookStore,"id"))
	fmt.Println("Server started")
	http.ListenAndServe("0.0.0.0:8080", router)
	time.Sleep(2 * time.Second)
}


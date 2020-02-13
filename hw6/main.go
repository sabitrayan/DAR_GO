package main

import (
	"encoding/json"
	"fmt"
	"github.com/fullacc/darintern/day6hw/book_store"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	bookStore, err := book_store.NewBookStore("books.json")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Methods("GET").Path("/{id}").HandlerFunc(GetBook(bookStore, "id"))
	router.Methods("POST").Path("/").HandlerFunc(CreateBook(bookStore))
	router.Methods("GET").Path("/exit").HandlerFunc(ExitWithSave(bookStore))
	router.Methods("LIST").Path("/list").HandlerFunc(ListBooks(bookStore))
	router.Methods("UPDATE").Path("/up").HandlerFunc(UpdateBook(bookStore,"id"))
	router.Methods("DELETE").Path("/del").HandlerFunc(DeleteBook(bookStore,"id"))
	fmt.Println("Server started")
	http.ListenAndServe("0.0.0.0:8080", router)
	time.Sleep(2 * time.Second)
}

func DeleteBook(store book_store.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		err := store.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
	}
}
func ListBooks(store book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		books, err := store.ListBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("I'm sorry: " + err.Error()))
			return
		}
		data, err := json.Marshal(books)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}
func GetBook(store book_store.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		book, err := store.GetBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("I'm sorry: " + err.Error()))
			return
		}
		data, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}
func UpdateBook(store book_store.BookStore, idParam string) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) { // eto kak
		vars := mux.Vars(r)//what is it
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		book := &book_store.Book{}
		result, err := store.UpdateBook(id,book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		response, err := json.Marshal(result)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusOK)

	}
}
func CreateBook(store book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))            //why byte's
			return
		}
		book := &book_store.Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		result, err := store.Create(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func ExitWithSave(book book_store.BookStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := book.SaveBooks("books.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
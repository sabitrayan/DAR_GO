package book_store

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)
type Endpoints interface {

	GetBook(idParam string) func(w http.ResponseWriter, r *http.Request)

	CreateBook() func(w http.ResponseWriter, r *http.Request)

	ListBooks() func(w http.ResponseWriter, r *http.Request)

	UpdateBook(idParam string) func (w http.ResponseWriter,r *http.Request)

	DeleteBook(idParam string) func(w http.ResponseWriter,r *http.Request)
}

type endpointsFactory struct {
	bookStore BookStore
}

func (ef *endpointsFactory) DeleteBook(idParam string) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		err := ef.bookStore.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
	}
}
func (ef *endpointsFactory) ListBooks() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		books, err := ef.bookStore.ListBooks()
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
func (ef endpointsFactory)GetBook(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		book, err := ef.bookStore.GetBook(id)
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
func (ef endpointsFactory) UpdateBook(idParam string) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) { // eto kak
		vars := mux.Vars(r)//what is it
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		book := &Book{}
		result, err := ef.bookStore.UpdateBook(id,book)
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
func (ef endpointsFactory)CreateBook() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))            //why byte's
			return
		}
		book := &Book{}
		if err := json.Unmarshal(data, book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		result, err := ef.bookStore.Create(book)
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

func (ef endpointsFactory)ExitWithSave() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ef.bookStore.SaveBooks("books.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
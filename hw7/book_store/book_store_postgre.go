package book_store

import "github.com/go-pg/pg"


type postgreStoreBook struct {
	db *pg.DB
}
func (ps *postgreStoreBook) CreateBook(book Book) (Book, error) {
	return book, ps.db.Insert(&book)
}

func (ps *postgreStoreBook) GetListBook() ([]Book, error) {
	var books []Book
	err := ps.db.Model(&books).Select()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (ps *postgreStoreBook) GetBook(id int) (*Book, error) {
	book := &Book{ID: id}
	err := ps.db.Select(book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (ps *postgreStoreBook) UpdateBook(book Book, id int) (*Book, error) {
	new_book := &Book{
		ID:          id,
		Title:        book.Title,
		Description: book.Description,
		Author:      book.Description,
	}
	err := ps.db.Update(new_book)
	if err != nil {
		return nil, err
	}
	return new_book, nil
}

func (ps *postgreStoreBook) DeleteBook(id int) error {
	book := &Book{ID: id}
	err := ps.db.Delete(book)
	if err != nil {
		return err
	}
	return nil
}

func (ps *postgreStore) CloseDB() error {
	err := ps.db.Close()
	if err != nil {
		return err
	}
	return nil
}


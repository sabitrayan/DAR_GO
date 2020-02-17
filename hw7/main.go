package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dairovolzhas/go-intern/lesson7/book_store"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	postgreConfigPath = "postgreconfig.json"
	configPath = "config.json"
	flags      = []cli.Flag{

	}
)

func main() {
	app := &cli.App{
		Flags:  flags,
		Name:   "dar",
		Usage:  "bar",
		Action: run,
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := startServer(); err != nil {
		return err
	}
	return nil
}

func startServer() error {

	config := book_store.Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	postgreConfig := book_store.PostgreConfig{}
	file, err = os.Open(postgreConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader = bufio.NewReader(file)
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &postgreConfig)
	if err != nil {
		return err
	}


	router := mux.NewRouter()

	bookStore, err := book_store.CreateBookStore(config.PathToBookStore)
	if err != nil {
		return err
	}

	endpoints, err := book_store.CreateEndPointFactory(bookStore)
	if err != nil {
		return err
	}

	router.Methods("GET").Path("/").HandlerFunc(endpoints.BooksListHandler())
	router.Methods("POST").Path("/").HandlerFunc(endpoints.BooksCreateHandler())
	router.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(endpoints.BookGetHandler("id"))
	router.Methods("PUT").Path("/{id:[0-9]+}").HandlerFunc(endpoints.BookUpdateHandler("id"))
	router.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(endpoints.BookDeleteHandler("id"))
	router.Methods("POST").Path("/save").HandlerFunc(book_store.SaveBookStoreHandler(bookStore, config.PathToBookStore))


	endpointsPostgre, err := book_store.NewPostgreBookStore(postgreConfig)
	if err != nil {
		return err
	}


	router.Methods("GET").Path("/postgre").HandlerFunc(endpointsPostgre.BooksListHandler())
	router.Methods("POST").Path("/postgre").HandlerFunc(endpointsPostgre.BooksCreateHandler())
	router.Methods("GET").Path("/postgre/{id:[0-9]+}").HandlerFunc(endpointsPostgre.BookGetHandler("id"))
	router.Methods("PUT").Path("/postgre/{id:[0-9]+}").HandlerFunc(endpointsPostgre.BookUpdateHandler("id"))
	router.Methods("DELETE").Path("/postgre/{id:[0-9]+}").HandlerFunc(endpointsPostgre.BookDeleteHandler("id"))

	fmt.Println("Server Started")
	fmt.Println(config.Port)

	go func() {
		http.ListenAndServe("localhost:"+config.Port, router)
	}()

	c := make(chan os.Signal)
	d := make(chan bool)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		d <- true
	}()
	<-d
	err = bookStore.SaveBookStore(config.PathToBookStore)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("\r- Book Store saved!")
	}

	os.Exit(1)

	return nil
}

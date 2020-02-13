package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"github.com/fatih/color"
	"os"

	//"os"

)

func rec(dir string, lvl int, last bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for i, f := range files {
		if last {
			for i:=0; i < lvl; i++ {
				fmt.Printf("    ")
			}
		} else {
			for i:=0; i < lvl; i++ {
				fmt.Printf("║   ")
			}
		}
		if i != len(files)-1 {
			fmt.Printf("╠═══ ")
		}else{
			fmt.Printf("╚═══ ")
		}
		if f.IsDir() {
			color.Set(color.FgMagenta)
			fmt.Println(f.Name())
			color.Unset()
		} else {
			fmt.Println(f.Name())
		}

		if f.IsDir() {
			rec(dir+"/"+f.Name(), lvl+1, i == len(files)-1)
		}
	}
}

func main(){
	var path string
	arg := os.Args
	if len(arg) <= 1 {
		path = "."
	}
	if len(arg) > 1 {
		path = arg[1]
	}
	//path = "."
	rec(path, 0, true)
}
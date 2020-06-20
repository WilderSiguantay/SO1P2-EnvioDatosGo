package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	fmt.Println("Server Running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

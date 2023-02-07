package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Home)
	http.ListenAndServe(":3000", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wellcome to home section")
}

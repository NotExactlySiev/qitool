package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Handler struct {
}

func (Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.Copy(os.Stdout, r.Body)
	fmt.Println()
	fmt.Println()
}

func main() {

	// /report/v1/index/bdlog"
	h := Handler{}
	http.ListenAndServe("127.0.0.1:44556", &h)
}

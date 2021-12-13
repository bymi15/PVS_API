package main

import (
	"fmt"
	"net/http"

	"github.com/bymi15/PVS_API/functions/src/utils"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	utils.ServeFunction("/api/helloworld", handler)
}

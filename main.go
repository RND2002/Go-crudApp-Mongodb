// mongodb+srv://codesite81:<password>@cluster0.upbvcsv.mongodb.net/?retryWrites=true&w=majority
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rnd2002/mongodb/router"
)

func main() {
	fmt.Println(("App started...."))

	r := router.Router()

	fmt.Println("Server is getting started.....")

	log.Fatal(http.ListenAndServe(":3000", r))
}

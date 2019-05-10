package Api

import (
	"log"
	"net/http"
)

func StartApi(port string) {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port, router))
}

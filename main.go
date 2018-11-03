package main
import (
	"net/http"
	"github.com/saresend/calhacks5server/handlers"	
)



func main() {


	http.HandleFunc("/votes", handlers.SocketHandler);
	http.ListenAndServe(":8080", nil);
}


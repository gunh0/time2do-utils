package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Orders API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3333
// @BasePath /
func main() {
	log.Println("Start main.go")

	// gorilla/mux router
	router := mux.NewRouter()

	// use the logger middleware on complete router
	router.Use(logMW)

	// serving routes
	router.HandleFunc("/hello", helloHandler)
	router.HandleFunc("/secret", authMW(secretHandler)) // this route using the authentication middleware

	router.PathPrefix("/").HandlerFunc(notFoundHandler)

	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)
	router.PathPrefix("/swagger").Handler(corsHandler(httpSwagger.WrapHandler))

	// start server on port 3333
	log.Println("Server started on port 3333")
	http.ListenAndServe("localhost:3333", router)

}

// (open) hello route
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello World!</h1>")
}

// restricted secret route
func secretHandler(w http.ResponseWriter, r *http.Request) {
	// getting isAdmin from context and convert to bool
	adm := context.Get(r, "isAdmin").(bool)

	// creating response, depending on isAdmin status
	body := "<h1>Hello on secret route.</h1>\n<p>%s</p>"
	var response string
	if adm {
		response = fmt.Sprintf(body, "You are admin.")
	} else {
		response = fmt.Sprintf(body, "You are user.")
	}

	fmt.Fprintln(w, response)
}

// compare the two forms of writing the middleware

// for global use (using a http.Handler!)
func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}

// for use on route (using a http.HandlerFunc)
func authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read basic auth information
		usr, _, ok := r.BasicAuth()

		// if there is no basic auth (no matter which credentials)
		if !ok {
			errMsg := "Authentication error!"
			// return a 403 forbidden
			http.Error(w, errMsg, http.StatusForbidden)
			log.Println(errMsg)

			// stop processing route
			return
		}

		// let's assume we check the user against a database to get
		// his admin-right and put this to the request context
		context.Set(r, "isAdmin", true)

		// else continue processing
		log.Printf("User %s logged in.", usr)
		next(w, r)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>404 Page Not Found</h1>")
}

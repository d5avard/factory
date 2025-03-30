package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/d5avard/factory/internal"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tmpl, err := template.ParseFiles("./web/templates/fixer.html")
			if err != nil {
				internal.HttpError(w, r, err.Error(), http.StatusInternalServerError)
				return
			}

			data := map[string]interface{}{
				"fixed_code":  "Fixed Code",
				"explanation": "Explanation",
			}

			if err = tmpl.Execute(w, data); err != nil {
				internal.HttpError(w, r, err.Error(), http.StatusInternalServerError)
				return
			}

			internal.LogRequest(r, http.StatusText(http.StatusOK), http.StatusOK)
			return

		case http.MethodPost:
			// Parse form data in case of POST requests
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}
			// Process POST request logic
			w.Write([]byte("Handling POST request"))
		default:
			// Send a response for unsupported methods
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

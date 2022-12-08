package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	tpl, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.",
			http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.",
			http.StatusInternalServerError)
		return
	}
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ Page</h1>
  <ul>
	<li>
	  <b>Is there a free version?</b>
	  Yes! We offer a free trial for 30 days on any paid plans.
	</li>
	<li>
	  <b>What are your support hours?</b>
	  We have support staff answering emails 24/7, though response
	  times may be a bit slower on weekends.
	</li>
	<li>
	  <b>How do I contact support?</b>
	  Email us - <a href="mailto:jhashem@hashem">jhashem@hashem.com</a>
	</li>
  </ul>
  `)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:jh@gmail.com\">jh@gmail.com</a>.")
}

func MyRequestHandler(w http.ResponseWriter, r *http.Request) {
	// fetch the url parameter `"userID"` from the request of a matching
	// routing pattern. An example routing pattern could be: /users/{userID}
	userID := chi.URLParam(r, "userID")

	// respond to the client
	w.Write([]byte(fmt.Sprintf("Here is the userID: %v", userID)))
}

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// }

func main() {

	logger := httplog.NewLogger("httplog", httplog.Options{
		// JSON: true,
		Concise: true,
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
	})

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/articles/{userID}", MyRequestHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

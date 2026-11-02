package main
import (

	"net/http"
	"log"
)

	// Initialize the form store
	// Used to create and store forms in memory
	//Is a map of the forms
	var Store = FormStore{
		Forms: make(map[int]*Form),
	}

func main() {
	
	// Create the server mux
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// Register for home page, create form page, and create form post handler
	mux.HandleFunc("/", home)
	mux.HandleFunc("/create", createForm)
	mux.HandleFunc("/createPost", createFormPost)

	mux.HandleFunc("/displayForm", displayForm)
	mux.HandleFunc("/respondToForm", respondToForm)
	mux.HandleFunc("/respondToFormPost", respondToFormPost)
	mux.HandleFunc("/viewResponses", viewResponses)

	mux.HandleFunc("/startFromTemplate", startFromTemplate)
	mux.HandleFunc("/deleteQuestion", deleteQuestion)
	
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

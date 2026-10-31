package main


import (
	"net/http"
	"html/template"
	"log"
)


// Handler for the home page
// Generates the home page template
func home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Handler for the create form page
// Generates the create form template
func createForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	
}

//Creates a new form and adds it to the store
//Called on post requests

func createFormPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	question := r.FormValue("question")

	// Check if the form with ID 1 already exists in the store
	//If not create a new one. 
	form := Store.Forms[1]
	if form == nil {
    	form = &Form{ID: 1, Title: title}
    	Store.Forms[1] = form
		log.Println("Created new form with ID:", form.ID)
	}

	//Append the new question to the form's questions slice
	form.Questions = append(form.Questions, Questions{Question: question})

	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

}

func displayForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/displayform.html")
	
	printForms(Store.Forms[1])

	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, Store.Forms[1])
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}
}

// Used to debug the forms
func printForms( form *Form) {
	log.Println("Form ID:", form.ID)
	log.Println("Form Title:", form.Title)
	log.Println("Form Questions:")
	for _, q := range form.Questions {
		log.Println("  - Question:", q.Question)
		log.Println("    Response:", q.Response)
	}
}

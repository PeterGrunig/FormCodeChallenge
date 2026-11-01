package main


import (
	"net/http"
	"html/template"
	"log"
	"strconv"
	"math/rand"

)


// Handler for the home page
// Generates the home page template
func home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	err = ts.Execute(w, Store.Forms)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Handler for the create form page
// Generates the create form template
func createForm(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	if id == 0 {
		Newid := rand.Intn(1_000_000) + 1
		id = Newid
		form := &Form{ID: id}
		Store.Forms[id] = form
		log.Println("Created new form with ID:", form.ID)
	}
	

	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, Store.Forms[id])
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
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}
	// Check if the form with the given ID already exists in the store
	// If not, create a new one.
	form := Store.Forms[id]
	if form == nil {
    	form = &Form{ID: id, Title: title}
    	Store.Forms[id] = form
		log.Println("Created new form with ID:", form.ID)
	}

	form.Title = title
	//Append the new question to the form's questions slice
	form.Questions = append(form.Questions, Questions{Question: question})

	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, form)
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

}

// Handler for displaying a form
func displayForm(w http.ResponseWriter, r *http.Request) {
	
	
	id, err := strconv.Atoi(r.FormValue("id"))
	printForms(Store.Forms[id])
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}
	
	ts, err := template.ParseFiles("templates/displayform.html")

	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, Store.Forms[id])
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}
}

// Shows the page that responds to a form
func respondToForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}
	ts, err := template.ParseFiles("templates/respondToForm.html")
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, Store.Forms[id])
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}
}

// Save function for the form responses
func respondToFormPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("formID"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	form, ok := Store.Forms[id]
	if !ok {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	for i := range form.Questions {
		key := "responses[" + strconv.Itoa(i) + "]"
		response := r.FormValue(key)
		form.Questions[i].Response = append(form.Questions[i].Response, response)
    }

	http.Redirect(w, r, "/home", http.StatusSeeOther)

}


// Handler for viewing responses to a form
func viewResponses(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	form, ok := Store.Forms[id]
	if !ok {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	ts, err := template.ParseFiles("templates/viewresponses.html")
	if err != nil {
		http.Error(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, form)
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

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed data/local_templates.json
var localTemplateJSON []byte

type localTemplateBundle struct {
	Questions []templateQuestion `json:"questions"`
}

type createFormPageData struct {
	Form               *Form
	Notice             string
	NextQuestionNumber int
}

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
// If the form ID is 0, it creates a new form with a random ID and adds it to the store
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
	sortQuestionsByOrder(Store.Forms[id])

	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, createFormPageData{
		Form:               Store.Forms[id],
		NextQuestionNumber: nextQuestionNumber(Store.Forms[id]),
	})
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// Creates a new form and adds it to the store
// Called on post requests
func createFormPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	question := r.FormValue("question")
	orderStr := r.FormValue("order")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}
	order, err := strconv.Atoi(orderStr)
	if err != nil || order < 1 {
		http.Error(w, "Question number must be a positive integer", http.StatusBadRequest)
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
	form.Questions = append(form.Questions, Questions{Question: question, Order: order})
	sortQuestionsByOrder(form)

	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, createFormPageData{
		Form:               form,
		NextQuestionNumber: nextQuestionNumber(form),
	})
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// Handler for displaying a form
func displayForm(w http.ResponseWriter, r *http.Request) {
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
	sortQuestionsByOrder(form)
	printForms(form)

	ts, err := template.ParseFiles("templates/displayform.html")

	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, form)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
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
	form, ok := Store.Forms[id]
	if !ok {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}
	sortQuestionsByOrder(form)
	ts, err := template.ParseFiles("templates/respondToForm.html")
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, form)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
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
	sortQuestionsByOrder(form)

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
	sortQuestionsByOrder(form)

	ts, err := template.ParseFiles("templates/viewresponses.html")
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, form)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func startFromTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	form := Store.Forms[id]
	if form == nil {
		form = &Form{ID: id, Title: "New Form from Template"}
		Store.Forms[id] = form
		log.Println("Created new form with ID:", form.ID)
	}

	form.Title = "New Form from Template"
	sortQuestionsByOrder(form)

	questions, apiErr := fetchTemplateQuestions()
	notice := ""
	if apiErr != nil {
		log.Printf("template API failed; using local template bundle: %v", apiErr)
		questions, err = loadLocalTemplateQuestions()
		if err != nil {
			http.Error(w, "Failed to load local templates: "+err.Error(), http.StatusInternalServerError)
			return
		}
		notice = "The template API is unavailable. Local template questions are being used instead."
	}

	firstNewQuestionNumber := nextQuestionNumber(form)
	for i := range questions {
		questions[i].Order = firstNewQuestionNumber + i
	}
	form.Questions = append(form.Questions, questions...)
	sortQuestionsByOrder(form)

	ts, err := template.ParseFiles("templates/createform.html")
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, createFormPageData{
		Form:               form,
		Notice:             notice,
		NextQuestionNumber: nextQuestionNumber(form),
	})
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func fetchTemplateQuestions() ([]Questions, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	questions := make([]Questions, 0, 3)

	for i := 0; i < 3; i++ {
		resp, err := client.Get("https://uselessfacts.jsph.pl/api/v2/facts/random?language=en")
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("template API returned %s", resp.Status)
		}

		var fact templateQuestion
		decodeErr := json.NewDecoder(resp.Body).Decode(&fact)
		resp.Body.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("decode template API response: %w", decodeErr)
		}

		question := strings.TrimSpace(fact.Question)
		if question == "" {
			return nil, fmt.Errorf("template API returned an empty question")
		}
		questions = append(questions, Questions{Question: question})
	}

	return questions, nil
}

func loadLocalTemplateQuestions() ([]Questions, error) {
	var bundle localTemplateBundle
	if err := json.Unmarshal(localTemplateJSON, &bundle); err != nil {
		return nil, fmt.Errorf("decode local template bundle: %w", err)
	}
	if len(bundle.Questions) == 0 {
		return nil, fmt.Errorf("local template bundle contains no questions")
	}

	questions := make([]Questions, 0, len(bundle.Questions))
	for _, item := range bundle.Questions {
		question := strings.TrimSpace(item.Question)
		if question != "" {
			questions = append(questions, Questions{Question: question})
		}
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("local template bundle contains no usable questions")
	}

	return questions, nil
}

func deleteQuestion(w http.ResponseWriter, r *http.Request) {

	log.Println("Delete question handler called")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(r.FormValue("index"))

	if err != nil {
		http.Error(w, "Invalid question index", http.StatusBadRequest)
		return
	}

	form := Store.Forms[id]
	if form == nil {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}
	sortQuestionsByOrder(form)

	if index < 0 || index >= len(form.Questions) {
		http.Error(w, "Question index out of range", http.StatusBadRequest)
		return
	}

	// Remove the question at the specified index
	form.Questions = append(form.Questions[:index], form.Questions[index+1:]...)
	sortQuestionsByOrder(form)

	// Redirect back to the create form page
	http.Redirect(w, r, fmt.Sprintf("/create?id=%d", id), http.StatusSeeOther)

}

func sortQuestionsByOrder(form *Form) {
	if form == nil {
		return
	}

	sort.SliceStable(form.Questions, func(i, j int) bool {
		return form.Questions[i].Order < form.Questions[j].Order
	})
}

func nextQuestionNumber(form *Form) int {
	next := 1
	if form == nil {
		return next
	}
	for _, question := range form.Questions {
		if question.Order >= next {
			next = question.Order + 1
		}
	}
	return next
}

// Used to debug the forms
func printForms(form *Form) {
	log.Println("Form ID:", form.ID)
	log.Println("Form Title:", form.Title)
	log.Println("Form Questions:")
	for _, q := range form.Questions {
		log.Println("  - Question:", q.Question)
		log.Println("    Response:", q.Response)
	}
}

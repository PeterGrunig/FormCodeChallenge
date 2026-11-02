package main

// Structs for the form and questions
// Has a question and response form
type Questions struct {
	Question string
	Response []string
	Order int 
	MultiplyChoise bool
	Required bool
}

// Main struct for the form
// Has an ID, title, and a slice of questions
type Form struct {
	ID          int
	Title       string
	Questions   []Questions
}

// Stores the forms in memory
type FormStore struct {
	Forms map[int]*Form
}

type templateQuestion struct {
	Question string `json:"text"`
}
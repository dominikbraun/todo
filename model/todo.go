// Package model provides domain entities for the application. These entities
// are expected and returned by the API and will be stored in the database.
package model

// ToDo represents a ToDo item, typically consisting of multiple sub-tasks.
type ToDo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Tasks       []Task `json:"tasks,omitempty"`
}

// Task represents a sub-task that is part of a ToDo item.
type Task struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

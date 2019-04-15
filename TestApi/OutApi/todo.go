package OutApi

import "time"

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

type Todos []Todo

func NewTodo() Todo {
	todo := Todo{
		0,
		"empty",
		false,
		time.Now(),
	}
	return todo
}

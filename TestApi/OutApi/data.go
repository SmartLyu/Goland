package main

import (
	"errors"
	"time"
)

var currentId int

var todos Todos

//type Mysql struct {
//	user   string
//	host   string
//	port   string
//	passwd string
//}

// 数据库信息
//var mysql = Mysql{
//	"root",
//	"127.0.0.1",
//	"3306",
//	"passwd",
//}

// Give us some seed data
func init() {
	first := NewTodo()
	RepoCreateTodo(first)

	first = Todo{
		Id:        2,
		Name:      "Write presentation",
		Completed: true,
		Due:       time.Now(),
	}
	RepoCreateTodo(first)

	RepoCreateTodo(Todo{
		3,
		"Host meetup",
		true,
		time.Now(),
	})
}

func RepoFindTodo(id int) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	return NewTodo()
}

func RepoCreateTodo(t Todo) Todo {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)
	return t
}

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return errors.New("Could not find Todo with id of " + string(id) + " to delete")
}

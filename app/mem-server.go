package main

import (
	"net/http"
	"sync"

	"github.com/google/uuid"

	api "doubleslash.de/coding-dojo-api/app/api"
	"github.com/labstack/echo/v4"
)

type InMemoryServer struct {
	lock  *sync.Mutex
	tasks *map[string]api.PostTask
}

func NewInMemoryServer() *InMemoryServer {
	return &InMemoryServer{
		lock:  &sync.Mutex{},
		tasks: &map[string]api.PostTask{},
	}
}

func (s InMemoryServer) GetTasks(ctx echo.Context) error {
	s.lock.Lock()
	taskList := make([]api.GetTask, 0, len(*s.tasks))
	for id, postTask := range *s.tasks {
		getTask := api.GetTask{
			Description: postTask.Description,
			Title:       postTask.Title,
			Uuid:        id,
		}
		taskList = append(taskList, getTask)
	}
	s.lock.Unlock()
	return ctx.JSON(http.StatusOK, taskList)
}

func (s InMemoryServer) PostTask(ctx echo.Context) error {
	var postTask api.PostTask
	err := ctx.Bind(&postTask)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Task JSON invalid")
	}
	id := uuid.New().String()
	s.lock.Lock()
	(*s.tasks)[id] = postTask
	s.lock.Unlock()
	getTask := api.GetTask{
		Description: postTask.Description,
		Title:       postTask.Title,
		Uuid:        id,
	}
	return ctx.JSON(http.StatusCreated, getTask)
}

func (s InMemoryServer) DeleteTask(ctx echo.Context, uuid string) error {
	s.lock.Lock()
	_, exists := (*s.tasks)[uuid]
	if !exists {
		s.lock.Unlock()
		return ctx.String(http.StatusNotFound, "Task not found")
	}
	delete(*s.tasks, uuid)
	s.lock.Unlock()
	return ctx.String(http.StatusOK, "Task deleted")

}

func (s InMemoryServer) GetTask(ctx echo.Context, uuid string) error {
	s.lock.Lock()
	postTask, exists := (*s.tasks)[uuid]
	if !exists {
		s.lock.Unlock()
		return ctx.String(http.StatusNotFound, "Task not found")
	}
	s.lock.Unlock()
	getTask := api.GetTask{
		Description: postTask.Description,
		Title:       postTask.Title,
		Uuid:        uuid,
	}
	return ctx.JSON(http.StatusOK, getTask)
}

func (s InMemoryServer) ReplaceTask(ctx echo.Context, uuid string) error {
	var postTask api.PostTask
	err := ctx.Bind(&postTask)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Task JSON invalid")
	}

	s.lock.Lock()
	_, exists := (*s.tasks)[uuid]
	if !exists {
		s.lock.Unlock()
		return ctx.String(http.StatusNotFound, "Task not found")
	}
	(*s.tasks)[uuid] = postTask
	s.lock.Unlock()

	return ctx.String(http.StatusOK, "Task replaced")
}

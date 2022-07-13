package main

import (
	"net/http"

	"github.com/google/uuid"

	api "doubleslash.de/coding-dojo-api/app/api"
	"github.com/labstack/echo/v4"
)

type Server struct {
	TaskList *[]api.GetTask
}

func NewServer() *Server {
	return &Server{
		TaskList: &[]api.GetTask{},
	}
}

func sendError(ctx echo.Context, code int, message string) error {
	return ctx.String(code, message)
}

func (s Server) GetTasks(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, s.TaskList)
}

func (s Server) PostTask(ctx echo.Context) error {
	var postTask api.PostTask
	err := ctx.Bind(&postTask)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Task JSON invalid")
	}
	id := uuid.New().String()
	getTask := api.GetTask{
		Description: postTask.Description,
		Title:       postTask.Title,
		Uuid:        id,
	}
	*s.TaskList = append(*s.TaskList, getTask)
	return ctx.JSON(http.StatusCreated, getTask)
}

func (s Server) DeleteTask(ctx echo.Context, uuid string) error {
	index := s.getIndexByUuid(uuid)
	if index == nil {
		return ctx.String(http.StatusNotFound, "Task not found")
	}
	s.removeTaskByIndex(*index)
	return ctx.String(http.StatusOK, "Task deleted")
}

func (s Server) GetTask(ctx echo.Context, uuid string) error {
	index := s.getIndexByUuid(uuid)
	if index == nil {
		return ctx.String(http.StatusNotFound, "Task not found")
	}
	return ctx.JSON(http.StatusOK, (*s.TaskList)[*index])
}

func (s Server) ReplaceTask(ctx echo.Context, uuid string) error {
	index := s.getIndexByUuid(uuid)
	if index == nil {
		return ctx.String(http.StatusNotFound, "Task not found")
	}

	var postTask api.PostTask
	err := ctx.Bind(&postTask)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Task JSON invalid")
	}
	getTask := api.GetTask{
		Description: postTask.Description,
		Title:       postTask.Title,
		Uuid:        uuid,
	}
	(*s.TaskList)[*index] = getTask
	return ctx.String(http.StatusOK, "Task replaced")
}

func (s Server) getIndexByUuid(uuid string) *int {
	for index, task := range *s.TaskList {
		if task.Uuid == uuid {
			return &index
		}
	}
	return nil
}

func (s Server) removeTaskByIndex(index int) {
	length := len(*s.TaskList)
	(*s.TaskList)[index] = (*s.TaskList)[length-1]
	*s.TaskList = (*s.TaskList)[:length-1]
}

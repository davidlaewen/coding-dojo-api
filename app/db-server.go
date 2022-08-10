package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"

	api "doubleslash.de/coding-dojo-api/app/api"
	"doubleslash.de/coding-dojo-api/app/database"
	"github.com/labstack/echo/v4"
)

type DatabaseServer struct {
	db database.Database
}

func NewDatabaseServer(db database.Database) *DatabaseServer {
	return &DatabaseServer{
		db: db,
	}
}

func (s DatabaseServer) GetTasks(ctx echo.Context) error {
	tasks, err := s.db.GetTasks()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, tasks)
}

func (s DatabaseServer) PostTask(ctx echo.Context) error {
	var postTask api.PostTask
	err := ctx.Bind(&postTask)
	if err != nil {
		log.Println("Task JSON invalid: ", err.Error())
		return ctx.String(http.StatusBadRequest, "Task JSON invalid")
	}
	id := uuid.New().String()
	getTask := api.GetTask{
		Description: postTask.Description,
		Title:       postTask.Title,
		Uuid:        id,
	}
	_, err = s.db.InsertTask(getTask)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, getTask)
}

func (s DatabaseServer) DeleteTask(ctx echo.Context, uuid string) error {
	exists, err := s.db.ExistsTask(uuid)
	if err != nil {
		return err
	}
	if !exists {
		return ctx.String(http.StatusNotFound, "Task not found")
	}

	err = s.db.DeleteTask(uuid)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s DatabaseServer) GetTask(ctx echo.Context, uuid string) error {
	task, err := s.db.GetTask(uuid)
	if err == sql.ErrNoRows {
		return ctx.String(http.StatusNotFound, "Task not found")
	}
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, task)
}

func (s DatabaseServer) ReplaceTask(ctx echo.Context, uuid string) error {
	var postTask api.PostTask
	err := ctx.Bind(&postTask)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Task JSON invalid")
	}

	exists, err := s.db.ExistsTask(uuid)
	if err != nil {
		return err
	}

	if !exists {
		return ctx.String(http.StatusNotFound, "Task not found")
	}

	err = s.db.UpdateTask(uuid, postTask)
	if err != nil {
		return err
	}
	// TODO: Update row in table

	return ctx.String(http.StatusOK, "Task replaced")
}

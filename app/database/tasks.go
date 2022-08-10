package database

import (
	"log"

	"doubleslash.de/coding-dojo-api/app/api"
	_ "github.com/lib/pq"
)

func (d Database) InitTableTasks() error {
	initDynStmt :=
		`CREATE TABLE IF NOT EXISTS tasks (
			uuid text PRIMARY KEY,
			title text NOT NULL,
			description text
		);
		CREATE INDEX IF NOT EXISTS tasks_uuid ON tasks(uuid);`
	_, err := d.DB.Exec(initDynStmt)
	return err
}

func (d Database) EmptyTableTasks() error {
	emptyDynStmt := `DELETE FROM tasks;`
	_, err := d.DB.Exec(emptyDynStmt)
	return err
}

func (d Database) InsertTask(task api.GetTask) (*string, error) {
	insertDynStmt := `
		INSERT INTO tasks(uuid, title, description)
		VALUES ($1, $2, $3)
		RETURNING uuid;`
	var id string
	err := d.DB.QueryRow(insertDynStmt, task.Uuid, task.Title, task.Description).Scan(&id)
	if err != nil {
		return nil, err
	}
	log.Println("Inserted task: ", task.Title, task.Description)
	return &id, nil
}

func (d Database) GetTasks() ([]api.GetTask, error) {
	query := `SELECT uuid, title, description FROM tasks;`
	rows, err := d.DB.Query(query)
	if err != nil {
		return []api.GetTask{}, err
	}
	tasks := []api.GetTask{}
	var task api.GetTask
	for rows.Next() {
		err := rows.Scan(&task.Uuid, &task.Title, &task.Description)
		if err != nil {
			return []api.GetTask{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (d Database) GetTask(uuid string) (api.GetTask, error) {
	queryDynStmt := `
		SELECT uuid, title, description
		FROM tasks
		WHERE uuid = $1;`
	var task api.GetTask
	err := d.DB.QueryRow(queryDynStmt, uuid).Scan(
		&task.Uuid, &task.Title, &task.Description)
	if err != nil {
		return api.GetTask{}, err
	}
	return task, nil
}

func (d Database) DeleteTask(uuid string) error {
	deleteDynStmt := `
		DELETE FROM tasks
		WHERE uuid = $1;`
	_, err := d.DB.Exec(deleteDynStmt, uuid)
	if err != nil {
		return err
	}
	log.Println("Deleted task with UUID ", uuid)
	return nil
}

func (d Database) UpdateTask(uuid string, task api.PostTask) error {
	updateDynStmt := `
		UPDATE tasks
		SET title = $2,
			description = $3
		WHERE uuid = $1;`
	_, err := d.DB.Exec(updateDynStmt, uuid, task.Title, task.Description)
	return err
}

func (d Database) FindTask(task api.PostTask) (string, error) {
	queryDynStmt := `
		SELECT uuid
		FROM tasks
		WHERE title = $2 AND description = $3;`
	var uuid string
	err := d.DB.QueryRow(queryDynStmt, task.Title, task.Description).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (d Database) ExistsTask(uuid string) (bool, error) {
	queryDynStmt := `
		SELECT EXISTS (
			SELECT 1
			FROM tasks
			WHERE uuid = $1);`
	var exists bool
	err := d.DB.QueryRow(queryDynStmt, uuid).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

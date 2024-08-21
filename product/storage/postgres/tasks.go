package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "gitlab.com/bahodirova/product/genproto/product"
)

type TaskManager struct {
	db *sql.DB
}

func NewTaskManager(db *sql.DB) *TaskManager {
	return &TaskManager{db: db}
}

func (t *TaskManager) CreateTask(req *pb.CreatetaskReq) (*pb.Empty, error) {
	id := uuid.NewString()
	query := `
	INSERT INTO
		tasks
			(id,

			title,
			description,
			user_id_assigned_to,
			status,
			date)
	VALUES
		($1,$2,$3,$4,$5,$6)`

	_, err := t.db.Exec(
		query,
		id,
		req.Title,
		req.Description,
		req.UserIdAssignedTo,
		req.Status,
		req.Date,
	)

	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (t *TaskManager) GetTask(req *pb.GetById) (*pb.Task, error) {
	query := `
		SELECT
			id,
			title,
			description,
			user_id_assigned_to,
			status,
			date,
			created_at
		FROM
			tasks
		WHERE
			id = $1
		AND
			deleted_at = 0
		`

	row := t.db.QueryRow(query, req.Id)
	var task pb.Task
	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.UserIdAssignedTo,
		&task.Status,
		&task.Date,
		&task.CreatedAt,
	)
	if err != nil {
		log.Println("no rows result set")
		return nil, err
	}
	return &task, nil
}

func (t *TaskManager) GetAllTasks(req *pb.GetAllTasksReq) (*pb.GetAllTasksRes, error) {

	query := `SELECT
			id,
			title,
			description,
			user_id_assigned_to,
			status,
			date,
			created_at
		FROM
			tasks`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.Status != "" {
		filters = append(filters, fmt.Sprintf("status = $%d", argCount))
		args = append(args, req.Status)
		argCount++
	}

	if req.Date != "" {
		filters = append(filters, fmt.Sprintf("date = $%d", argCount))
		args = append(args, req.Date)
		argCount++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", req.Filter.Limit)
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", req.Filter.Offset)
		}
	}
	rows, err := t.db.Query(query, args...)
	if err != nil {
		log.Println("no rows result set")
		return nil, err
	}
	defer rows.Close()
	tasks := []*pb.Task{}
	for rows.Next() {
		var task pb.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.UserIdAssignedTo,
			&task.Status,
			&task.Date,
			&task.CreatedAt,
		)
		if err != nil {
			log.Println("no rows result set")
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return &pb.GetAllTasksRes{Tasks: tasks}, nil
}

func (t *TaskManager) UpdateTask(req *pb.UpdateTaskReq) (*pb.UpdateTaskRes, error) {

	var args []interface{}
	var conditions []string

	if req.Title != "" && req.Title != "string" {
		args = append(args, req.Title)
		conditions = append(conditions, fmt.Sprintf("title = $%d", len(args)))
	}

	if req.Description != "" && req.Description != "string" {
		args = append(args, req.Description)
		conditions = append(conditions, fmt.Sprintf("description = $%d", len(args)))
	}

	if req.Status != "" && req.Status != "string" {
		args = append(args, req.Status)
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}

	if req.Date != "" && req.Date != "string" {
		args = append(args, req.Date)
		conditions = append(conditions, fmt.Sprintf("date = $%d", len(args)))
	}

	if req.UserIdAssignedTo != 0 {
		args = append(args, req.UserIdAssignedTo)
		conditions = append(conditions, fmt.Sprintf("user_id_assigned_to = $%d", len(args)))
	}

	args = append(args, time.Now())
	conditions = append(conditions, fmt.Sprintf("updated_at = $%d", len(args)))

	query := `UPDATE tasks SET ` + strings.Join(conditions, ", ") + ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, req.Id)

	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pb.UpdateTaskRes{Success: true, Message: "Task updated successfully"}, nil

}
func (t *TaskManager) DeleteTask(req *pb.GetById) (*pb.DeleteTaskRes, error) {
	query := `
	UPDATE
		tasks
	SET
		deleted_at = extract(epoch from now())
	WHERE
		id = $1
	`

	_, err := t.db.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskRes{Success: true, Message: "Task deleted successfully"}, nil
}

func (t *TaskManager) GetTasksByUser(req *pb.GetByUserReq) (*pb.GetAllTasksRes, error) {

	query := `
		SELECT
			id,
			title,
			description,
			user_id_assigned_to,
			status,
			date,
			created_at
		FROM
			tasks
		WHERE
			user_id_assigned_to = $1
	`
	
	var args []interface{}
	args = append(args, req.UserIdAssignedTo)
	argCount := 2

	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}

	rows, err := t.db.Query(query, args...)
	if err != nil {
		log.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	tasks := []*pb.Task{}
	for rows.Next() {
		task := pb.Task{}
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.UserIdAssignedTo,
			&task.Status,
			&task.Date,
			&task.CreatedAt,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return &pb.GetAllTasksRes{Tasks: tasks}, nil
}



func (t *TaskManager) SearchTasks(req *pb.SearchTasksReq) (*pb.GetAllTasksRes, error) {

	query := `
		SELECT
			id,
			title,
			description,
			user_id_assigned_to,
			status,
			date,
			created_at
		FROM
			tasks

	`
	
	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.Title != "" {
		filters = append(filters, fmt.Sprintf("title = $%d", argCount))
		args = append(args, req.Title)
		argCount++
	}
	if req.Description != "" {
		filters = append(filters, fmt.Sprintf("description = $%d", argCount))
		args = append(args, req.Description)
		argCount++
	}
	if req.Status != "" {
		filters = append(filters, fmt.Sprintf("status = $%d", argCount))
		args = append(args, req.Status)
		argCount++
	}
	if req.Date != "" {
		filters = append(filters, fmt.Sprintf("date = $%d", argCount))
		args = append(args, req.Date)
		argCount++
	}
	
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", req.Filter.Limit)
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", req.Filter.Offset)
		}
	}

	rows, err := t.db.Query(query, args...)
	if err != nil {
		log.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	tasks := []*pb.Task{}
	for rows.Next() {
		task := pb.Task{}
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.UserIdAssignedTo,
			&task.Status,
			&task.Date,
			&task.CreatedAt,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return &pb.GetAllTasksRes{Tasks: tasks}, nil
}

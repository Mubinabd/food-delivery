package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage/postgres"
)

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewTaskManager(db)

	req := &pb.CreatetaskReq{
		Title:            "Task Title",
		Description:      "Task Description",
		UserIdAssignedTo: 123,
		Status:           "Open",
		Date:             "2023-07-28",
	}

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(sqlmock.AnyArg(), req.Title, req.Description, req.UserIdAssignedTo, req.Status, req.Date).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = manager.CreateTask(req)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewTaskManager(db)

	taskID := uuid.NewString()
	req := &pb.GetById{Id: taskID}

	mock.ExpectQuery("SELECT id, title, description, user_id_assigned_to, status, date, created_at FROM tasks").
		WithArgs(taskID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "user_id_assigned_to", "status", "date", "created_at"}).
			AddRow(taskID, "Task Title", "Task Description", 123, "Open", "2023-07-28", "2023-07-01"))

	task, err := manager.GetTask(req)
	require.NoError(t, err)

	assert.Equal(t, taskID, task.Id)
	assert.Equal(t, "Task Title", task.Title)
	assert.Equal(t, "Task Description", task.Description)
	assert.Equal(t, int32(123), task.UserIdAssignedTo)
	assert.Equal(t, "Open", task.Status)
	assert.Equal(t, "2023-07-28", task.Date)
	assert.Equal(t, "2023-07-01", task.CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewTaskManager(db)

	req := &pb.UpdateTaskReq{
		Id:               "id1",
		Title:            "Updated Task Title",
		Description:      "Updated Task Description",
		UserIdAssignedTo: 123,
		Status:           "InProgress",
		Date:             "2023-07-30",
	}

	query := `UPDATE tasks SET title = \$1, description = \$2, status = \$3, date = \$4, user_id_assigned_to = \$5, updated_at = \$6 WHERE id = \$7`
	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(req.Title, req.Description, req.Status, req.Date, req.UserIdAssignedTo, sqlmock.AnyArg(), req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := manager.UpdateTask(req)
	require.NoError(t, err)

	assert.True(t, res.Success)
	assert.Equal(t, "Task updated successfully", res.Message)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetTasksByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewTaskManager(db)

	req := &pb.GetByUserReq{UserIdAssignedTo: 123}

	mock.ExpectQuery("SELECT id, title, description, user_id_assigned_to, status, date, created_at FROM tasks").
		WithArgs(req.UserIdAssignedTo).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "user_id_assigned_to", "status", "date", "created_at"}).
			AddRow("id1", "Task Title 1", "Task Description 1", 123, "Open", "2023-07-28", "2023-07-01"))

	res, err := manager.GetTasksByUser(req)
	require.NoError(t, err)

	assert.Len(t, res.Tasks, 1)
	assert.Equal(t, "id1", res.Tasks[0].Id)
	assert.Equal(t, "Task Title 1", res.Tasks[0].Title)
	assert.Equal(t, "Task Description 1", res.Tasks[0].Description)
	assert.Equal(t, int32(123), res.Tasks[0].UserIdAssignedTo)
	assert.Equal(t, "Open", res.Tasks[0].Status)
	assert.Equal(t, "2023-07-28", res.Tasks[0].Date)
	assert.Equal(t, "2023-07-01", res.Tasks[0].CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSearchTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewTaskManager(db)

	req := &pb.SearchTasksReq{
		Title:       "Task Title 1",
		Description: "Task Description 1",
		Status:      "Open",
		Date:        "2023-07-28",
	}

	mock.ExpectQuery("SELECT id, title, description, user_id_assigned_to, status, date, created_at FROM tasks").
		WithArgs(req.Title, req.Description, req.Status, req.Date).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "user_id_assigned_to", "status", "date", "created_at"}).
			AddRow("id1", "Task Title 1", "Task Description 1", 123, "Open", "2023-07-28", "2023-07-01"))

	res, err := manager.SearchTasks(req)
	require.NoError(t, err)

	assert.Len(t, res.Tasks, 1)
	assert.Equal(t, "id1", res.Tasks[0].Id)
	assert.Equal(t, "Task Title 1", res.Tasks[0].Title)
	assert.Equal(t, "Task Description 1", res.Tasks[0].Description)
	assert.Equal(t, int32(123), res.Tasks[0].UserIdAssignedTo)
	assert.Equal(t, "Open", res.Tasks[0].Status)
	assert.Equal(t, "2023-07-28", res.Tasks[0].Date)
	assert.Equal(t, "2023-07-01", res.Tasks[0].CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage/postgres"
)

func TestCreateNotification(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewNotificationManager(db)

	req := &pb.CreateNotificationReq{
		UserId:  "user-123",
		Message: "Test message",
		IsRead:  false,
	}

	mock.ExpectExec("INSERT INTO notifications").
		WithArgs(sqlmock.AnyArg(), req.UserId, req.Message, req.IsRead).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = manager.CreateNotification(req)
	assert.NoError(t, err)
}

func TestGetNotification(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewNotificationManager(db)

	id := uuid.NewString()
	rows := sqlmock.NewRows([]string{"id", "user_id", "message", "is_read", "created_at"}).
		AddRow(id, "user-123", "Test message", false, "2024-08-01T00:00:00Z")

	mock.ExpectQuery("SELECT id, user_id, message, is_read, created_at FROM notifications WHERE id = ?").
		WithArgs(id).
		WillReturnRows(rows)

	req := &pb.GetById{Id: id}
	notification, err := manager.GetNotification(req)
	assert.NoError(t, err)
	assert.NotNil(t, notification)
	assert.Equal(t, id, notification.Id)
	assert.Equal(t, "user-123", notification.UserId)
	assert.Equal(t, "Test message", notification.Message)
	assert.Equal(t, false, notification.IsRead)
}

func TestGetAllNotifications(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewNotificationManager(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "message", "is_read", "created_at"}).
		AddRow(uuid.NewString(), "user-123", "Test message", false, "2024-08-01T00:00:00Z")

	mock.ExpectQuery("SELECT id, user_id, message, is_read, created_at FROM notifications").
		WillReturnRows(rows)

	req := &pb.GetAllNotificationsReq{UserId: "user-123"}
	notifications, err := manager.GetAllNotifications(req)
	assert.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.Greater(t, len(notifications.Notifications), 0)
}

func TestMarkNotificationAsRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewNotificationManager(db)

	id := uuid.NewString()

	mock.ExpectExec("UPDATE notifications SET is_read = true WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.MarkNotificationAsReadReq{Id: id}
	resp, err := manager.MarkNotificationAsRead(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Notification marked as read successfully", resp.Message)
}

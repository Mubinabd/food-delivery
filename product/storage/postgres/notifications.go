package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	pb "gitlab.com/bahodirova/product/genproto/product"
)

type NotificationManager struct {
	db *sql.DB
}

func NewNotificationManager(db *sql.DB) *NotificationManager {
	return &NotificationManager{db: db}
}

func (n *NotificationManager) CreateNotification(req *pb.CreateNotificationReq) (*pb.Empty, error) {
	id := uuid.NewString()
	query := `
		INSERT INTO
			notifications
				(id, 
				user_id, 
				message,  
				is_read)
		VALUES
			($1, $2, $3, $4)`

	_, err := n.db.Exec(
		query,
		id,
		req.UserId,
		req.Message,
		req.IsRead)

	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (n *NotificationManager) GetNotification(req *pb.GetById) (*pb.Notification, error) {

	query := `
		SELECT
			id,
			user_id,
			message,
			is_read,
			created_at
		FROM
			notifications
		WHERE
			id = $1`

	row := n.db.QueryRow(query, req.Id)
	var notification pb.Notification

	err := row.Scan(
		&notification.Id,
		&notification.UserId,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
	)

	if err != nil {
		log.Println("no rows result set")
		return nil, err
	}
	return &notification, nil
}

func (n *NotificationManager) GetAllNotifications(req *pb.GetAllNotificationsReq) (*pb.GetAllNotificationsRes, error) {
	query := `
		SELECT
			id,
			user_id,
			message,
			is_read,
			created_at
		FROM
			notifications`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.UserId != "" {
		filters = append(filters, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, req.UserId)
		argCount++
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
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

	rows, err := n.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*pb.Notification

	for rows.Next() {
		var notification pb.Notification
		err := rows.Scan(
			&notification.Id,
			&notification.UserId,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}
	return &pb.GetAllNotificationsRes{Notifications: notifications}, nil
}

func (n *NotificationManager) MarkNotificationAsRead(req *pb.MarkNotificationAsReadReq) (*pb.MarkNotificationAsReadResp, error) {
	query := `
		UPDATE
			notifications
		SET
			is_read = true
		WHERE
			id = $1`

	_, err := n.db.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.MarkNotificationAsReadResp{Success: true, Message: "Notification marked as read successfully"}, nil
}


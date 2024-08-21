package postgres

import (
	"database/sql"
	"fmt"

	u "gitlab.com/bahodirova/product/storage"
	"googlemaps.github.io/maps"

	_ "github.com/lib/pq"
	"gitlab.com/bahodirova/product/config"
)

type Storage struct {
	Db               *sql.DB
	TaskS            u.Task
	CartS            u.Cart
	OrderS           u.Order
	ProductS         u.Product
	OrderItemS       u.OrderItem
	NotificationS    u.Notification
	CourierLocationS u.CourierLocation
}

func NewPostgresStorage() (*Storage, error) {

	conn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		config.Load().PostgresHost,
		config.Load().PostgresUser,
		config.Load().PostgresDatabase,
		config.Load().PostgresPassword,
		config.Load().PostgresPort,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	mapsAPIKey := "AIzaSyD10zkY_fPmM0iQNhEQQPgmtxwyRdIBHxM"

	return &Storage{
		Db:               db,
		CartS:            NewCartManager(db),
		TaskS:            NewTaskManager(db),
		OrderS:           NewOrderManager(db),
		ProductS:         NewProductManager(db),
		OrderItemS:       NewOrderItemManager(db),
		NotificationS:    NewNotificationManager(db),
		CourierLocationS: NewCourierLocationManager(db, mapsAPIKey),
	}, nil

}

func (s *Storage) Product() u.Product {
	if s.ProductS == nil {
		s.ProductS = &ProductManager{s.Db}
	}

	return s.ProductS
}
func (s *Storage) Task() u.Task {
	if s.TaskS == nil {
		s.TaskS = &TaskManager{s.Db}
	}

	return s.TaskS
}

func (s *Storage) Order() u.Order {
	if s.OrderS == nil {
		s.OrderS = &OrderManager{s.Db}
	}

	return s.OrderS
}
func (s *Storage) OrderItem() u.OrderItem {
	if s.OrderItemS == nil {
		s.OrderItemS = &OrderItemManager{s.Db}
	}

	return s.OrderItemS
}
func (s *Storage) Notification() u.Notification {
	if s.NotificationS == nil {
		s.NotificationS = &NotificationManager{s.Db}
	}

	return s.NotificationS
}

func (s *Storage) CourierLocation() u.CourierLocation {
	if s.CourierLocationS == nil {
		s.CourierLocationS = &CourierLocationManager{s.Db, &maps.Client{}}
	}

	return s.CourierLocationS
}
func (s *Storage) Cart() u.Cart {
	if s.CartS == nil {
		s.CartS = &CartManager{s.Db}
	}

	return s.CartS
}

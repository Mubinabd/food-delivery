package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"gitlab.com/bahodirova/product/storage/postgres"
)

func TestCreateCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	manager := postgres.NewCartManager(db)

	req := &pb.CreateCartReq{
		UserId:    "user-id",
		ProductId: "product-id",
		Quantity:  1,
		Options:   "some options",
		Name:      "cart name",
		Number:    1,
	}

	mock.ExpectExec("INSERT INTO cart").
		WithArgs(sqlmock.AnyArg(), req.UserId, req.ProductId, req.Quantity, sqlmock.AnyArg(), req.Name, req.Number).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = manager.CreateCart(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

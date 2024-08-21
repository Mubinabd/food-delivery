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

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewProductManager(db)

	req := &pb.CreateProductRequest{
		Name:        "Product Name",
		Description: "Product Description",
		Price:       1000,
		ImageUrl:    "http://image.url",
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(sqlmock.AnyArg(), req.Name, req.Description, req.Price, req.ImageUrl).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = manager.CreateProduct(req)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewProductManager(db)

	productID := uuid.NewString()
	req := &pb.GetById{Id: productID}

	mock.ExpectQuery("SELECT id, name, description, price, image_url, created_at FROM products").
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "image_url", "created_at"}).
			AddRow(productID, "Product Name", "Product Description", 1000, "http://image.url", "2023-07-01"))

	product, err := manager.GetProduct(req)
	require.NoError(t, err)

	assert.Equal(t, productID, product.Id)
	assert.Equal(t, "Product Name", product.Name)
	assert.Equal(t, "Product Description", product.Description)
	assert.Equal(t, int32(1000), product.Price)
	assert.Equal(t, "http://image.url", product.ImageUrl)
	assert.Equal(t, "2023-07-01", product.CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewProductManager(db)

	req := &pb.GetAllProductsReq{}

	mock.ExpectQuery("SELECT id, name, description, price, image_url, created_at FROM products").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "image_url", "created_at"}).
			AddRow("id1", "Product Name 1", "Product Description 1", 1000, "http://image1.url", "2023-07-01").
			AddRow("id2", "Product Name 2", "Product Description 2", 2000, "http://image2.url", "2023-07-02"))

	res, err := manager.GetAllProducts(req)
	require.NoError(t, err)

	assert.Len(t, res.Products, 2)
	assert.Equal(t, "id1", res.Products[0].Id)
	assert.Equal(t, "Product Name 1", res.Products[0].Name)
	assert.Equal(t, "Product Description 1", res.Products[0].Description)
	assert.Equal(t, int32(1000), res.Products[0].Price)
	assert.Equal(t, "http://image1.url", res.Products[0].ImageUrl)
	assert.Equal(t, "2023-07-01", res.Products[0].CreatedAt)

	assert.Equal(t, "id2", res.Products[1].Id)
	assert.Equal(t, "Product Name 2", res.Products[1].Name)
	assert.Equal(t, "Product Description 2", res.Products[1].Description)
	assert.Equal(t, int32(2000), res.Products[1].Price)
	assert.Equal(t, "http://image2.url", res.Products[1].ImageUrl)
	assert.Equal(t, "2023-07-02", res.Products[1].CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewProductManager(db)

	req := &pb.UpdateProductRequest{
		Id:          "id1",
		Name:        "Updated Product Name",
		Description: "Updated Product Description",
		Price:       1500,
		ImageUrl:    "http://updatedimage.url",
	}

	query := `UPDATE products SET name = \$1, description = \$2, price = \$3, image_url = \$4, updated_at = \$5 WHERE id = \$6`
	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(req.Name, req.Description, req.Price, req.ImageUrl, sqlmock.AnyArg(), req.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := manager.UpdateProduct(req)
	require.NoError(t, err)

	assert.True(t, res.Success)
	assert.Equal(t, "Product updated successfully", res.Message)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}


func TestSearchProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	manager := postgres.NewProductManager(db)

	req := &pb.SearchProductsReq{Name: "Product Name 1"}

	mock.ExpectQuery("SELECT id, name, description, price, image_url, created_at FROM products").
		WithArgs(req.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "image_url", "created_at"}).
			AddRow("id1", "Product Name 1", "Product Description 1", 1000, "http://image1.url", "2023-07-01"))

	res, err := manager.SearchProducts(req)
	require.NoError(t, err)

	assert.Len(t, res.Products, 1)
	assert.Equal(t, "id1", res.Products[0].Id)
	assert.Equal(t, "Product Name 1", res.Products[0].Name)
	assert.Equal(t, "Product Description 1", res.Products[0].Description)
	assert.Equal(t, int32(1000), res.Products[0].Price)
	assert.Equal(t, "http://image1.url", res.Products[0].ImageUrl)
	assert.Equal(t, "2023-07-01", res.Products[0].CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

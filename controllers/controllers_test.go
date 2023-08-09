package controllers

import (
	"bytes"
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject1/models" // Replace with the actual package import
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductHandler(t *testing.T) {
	// Create a Gin router
	router := gin.Default()

	dsn := "host=localhost user=postgres password=docker dbname=CRUD-db port=5432 sslmode=disable"
	mockDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Unable to connect to the database")
	}

	// Create a ProductController instance with the mock DB
	productCtrl := &ProductController{DB: mockDB}

	// Set up a test route with the CreateProductHandler
	router.POST("/products/", productCtrl.CreateProductHandler)

	// Create a mock product for testing
	mockProduct := models.Product{
		ID:       100,
		Name:     "Test Product",
		Type:     "Type A",
		Quantity: 10,
	}

	// Convert the mock product to JSON
	productJSON, _ := json.Marshal(mockProduct)

	// Create a request for the test route
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	assert.NoError(t, err)

	// Create a response recorder to capture the response
	rec := httptest.NewRecorder()

	// Serve the request through the router
	router.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Deserialize the response JSON
	var responseProduct models.Product
	err = json.Unmarshal(rec.Body.Bytes(), &responseProduct)
	assert.NoError(t, err)

	// Check if the response product matches the expected product
	assert.Equal(t, mockProduct, responseProduct)
}

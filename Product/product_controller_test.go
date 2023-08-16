package product_test

import (
	"awesomeProject1/Product"
	"awesomeProject1/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProductHandler(t *testing.T) {
	models.InitLogger()

	t.Run("ExistingID", func(t *testing.T) {
		mockRepo := &product.MockProductRepository{
			GetProductByIDFn: func(productID uint) (*product.Product, error) {
				// Mock implementation to simulate an existing product with the same ID
				return &product.Product{}, nil
			},
		}

		ctrl := product.NewProductController(mockRepo)

		r := gin.Default()
		r.POST("/products", ctrl.CreateProductHandler)

		// Create a mock request with a product that has an existing ID
		productJSON := `{"id": 123, "name": "testproduct", "type": "test", "quantity": 5}`
		w := performRequest(r, "POST", "/products", productJSON)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Product with same ID exists")
	})

	// Add more test cases as needed...
}

func TestDeleteProductHandler(t *testing.T) {
	models.InitLogger()

	t.Run("InvalidID", func(t *testing.T) {
		mockRepo := &product.MockProductRepository{}
		ctrl := product.NewProductController(mockRepo)

		r := gin.Default()
		r.DELETE("/products/:id", ctrl.DeleteProductHandler)

		// Create a mock request with an invalid product ID
		w := performRequest(r, "DELETE", "/products/invalid-id", "")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	// Add more test cases for product not found and successful deletion...
}

func TestUpdateProductHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &product.MockProductRepository{
		GetProductByIDFn: func(productID uint) (*product.Product, error) {
			if productID == 1 {
				return &product.Product{ID: 1}, nil
			}
			return nil, nil
		},
		UpdateProductFn: func(product *product.Product, existingPID uint) error {
			if existingPID == 1 {
				return nil
			}
			return errors.New("product not found")
		},
	}

	ctrl := product.NewProductController(mockRepo)

	r := gin.Default()
	r.PUT("/products/:id", ctrl.UpdateProductHandler)

	// Test case 1: Invalid ID
	w := performRequest(r, "PUT", "/products/invalid-id", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 2: Product not found
	w = performRequest(r, "PUT", "/products/2", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 3: Invalid JSON
	w = performRequest(r, "PUT", "/products/1", "application/json")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 4: Successful update
	validJSON := `{"id": 1, "name": "new_product", "type": "updated", "quantity": 10}`
	w = performRequest(r, "PUT", "/products/1", validJSON)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllProductsHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &product.MockProductRepository{
		GetAllProductsFn: func() ([]product.Product, error) {
			products := []product.Product{
				{ID: 1, Name: "product1", Type: "type1", Quantity: 5},
				{ID: 2, Name: "product2", Type: "type2", Quantity: 10},
			}
			return products, nil
		},
	}

	ctrl := product.NewProductController(mockRepo)

	r := gin.Default()
	r.GET("/products", ctrl.GetAllProductsHandler)

	w := performRequest(r, "GET", "/products", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var responseProducts []product.Product
	err := json.Unmarshal(w.Body.Bytes(), &responseProducts)
	assert.NoError(t, err)

	expectedProducts := []product.Product{
		{ID: 1, Name: "product1", Type: "type1", Quantity: 5},
		{ID: 2, Name: "product2", Type: "type2", Quantity: 10},
	}

	assert.Equal(t, expectedProducts, responseProducts)
}

func TestGetSpecificProductHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &product.MockProductRepository{
		GetProductByIDFn: func(productID uint) (*product.Product, error) {
			if productID == 1 {
				return &product.Product{ID: 1, Name: "product1", Type: "type1", Quantity: 5}, nil
			} else if productID == 2 {
				return &product.Product{ID: 2, Name: "product2", Type: "type2", Quantity: 10}, nil
			}
			return nil, nil // Simulate product not found
		},
	}

	ctrl := product.NewProductController(mockRepo)

	r := gin.Default()
	r.GET("/products/:id", ctrl.GetSpecificProductHandler)

	// Test case: Product with ID 1 exists
	w1 := performRequest(r, "GET", "/products/1", "")
	assert.Equal(t, http.StatusOK, w1.Code)

	var product1 product.Product
	err := json.Unmarshal(w1.Body.Bytes(), &product1)
	assert.NoError(t, err)

	expectedProduct1 := product.Product{ID: 1, Name: "product1", Type: "type1", Quantity: 5}
	assert.Equal(t, expectedProduct1, product1)

	// Test case: Product with ID 2 exists
	w2 := performRequest(r, "GET", "/products/2", "")
	assert.Equal(t, http.StatusOK, w2.Code)

	var product2 product.Product
	err = json.Unmarshal(w2.Body.Bytes(), &product2)
	assert.NoError(t, err)

	expectedProduct2 := product.Product{ID: 2, Name: "product2", Type: "type2", Quantity: 10}
	assert.Equal(t, expectedProduct2, product2)

	// Test case: Product with non-existing ID
	w3 := performRequest(r, "GET", "/products/3", "")
	assert.Equal(t, http.StatusNotFound, w3.Code)

	var response3 map[string]string
	err = json.Unmarshal(w3.Body.Bytes(), &response3)
	assert.NoError(t, err)

	expectedResponse3 := map[string]string{"error": "product not found"}
	assert.Equal(t, expectedResponse3, response3)

	// Test case: Invalid ID
	w4 := performRequest(r, "GET", "/products/invalid-id", "")
	assert.Equal(t, http.StatusBadRequest, w4.Code)

	var response4 map[string]string
	err = json.Unmarshal(w4.Body.Bytes(), &response4)
	assert.NoError(t, err)

	expectedResponse4 := map[string]string{"error": "Invalid ID"}
	assert.Equal(t, expectedResponse4, response4)
}
func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

package controllers_test

import (
	"awesomeProject1/controllers"
	"awesomeProject1/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	models.InitLogger()

	t.Run("ExistingID", func(t *testing.T) {
		mockRepo := &controllers.MockUserRepository{
			GetUserByIDFn: func(userID uint) (*models.User, error) {
				// Mock implementation to simulate an existing user with the same ID
				return &models.User{}, nil
			},
		}

		ctrl := controllers.NewUserController(mockRepo)

		r := gin.Default()
		r.POST("/users", ctrl.CreateUserHandler)

		// Create a mock request with a user that has an existing ID
		userJSON := `{"id": 123, "username": "testuser", "email": "test@example.com", "role": "user"}`
		w := performRequest(r, "POST", "/users", userJSON)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "User with same ID exists")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockRepo := &controllers.MockUserRepository{}

		ctrl := controllers.NewUserController(mockRepo)

		r := gin.Default()
		r.POST("/users", ctrl.CreateUserHandler)

		// Create a mock request with invalid JSON
		w := performRequest(r, "POST", "/users", "invalid-json")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// Add more assertions if needed...
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &controllers.MockUserRepository{
			CreateUserFn: func(user *models.User) error {
				// Mock implementation for successful user creation
				return nil
			},
			GetUserByIDFn: func(userID uint) (*models.User, error) {
				// Mock implementation to simulate a user not found (for new ID)
				return nil, nil
			},
		}

		ctrl := controllers.NewUserController(mockRepo)

		r := gin.Default()
		r.POST("/users", ctrl.CreateUserHandler)

		// Create a mock request with a valid user JSON
		userJSON := `{"id": 123, "username": "testuser", "email": "test@example.com", "role": "user"}`
		w := performRequest(r, "POST", "/users", userJSON)

		assert.Equal(t, http.StatusCreated, w.Code)
		// Add more assertions if needed...
	})
}
func TestDeleteUserHandler(t *testing.T) {
	models.InitLogger()

	t.Run("InvalidID", func(t *testing.T) {
		mockRepo := &controllers.MockUserRepository{}
		ctrl := controllers.NewUserController(mockRepo)

		r := gin.Default()
		r.DELETE("/users/:id", ctrl.DeleteUserHandler)

		// Create a mock request with an invalid user ID
		w := performRequest(r, "DELETE", "/users/invalid-id", "")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo := &controllers.MockUserRepository{
			GetUserByIDFn: func(userID uint) (*models.User, error) {
				// Mock implementation to simulate user not found
				return nil, nil
			},
		}

		ctrl := controllers.NewUserController(mockRepo)

		r := gin.Default()
		r.DELETE("/users/:id", ctrl.DeleteUserHandler)

		// Create a mock request with a user ID that doesn't exist
		w := performRequest(r, "DELETE", "/users/123", "")

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "user not found")
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &controllers.MockUserRepository{
			GetUserByIDFn: func(userID uint) (*models.User, error) {
				// Mock implementation to simulate existing user
				return &models.User{}, nil
			},
			DeleteUserFn: func(userID uint) error {
				// Mock implementation for successful user deletion
				return nil
			},
		}

		ctrl := controllers.NewUserController(mockRepo)

		r := gin.Default()
		r.DELETE("/users/:id", ctrl.DeleteUserHandler)

		// Create a mock request with an existing user ID
		w := performRequest(r, "DELETE", "/users/123", "")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User deleted successfully")
	})
}

func TestUpdateUserHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &controllers.MockUserRepository{
		GetUserByIDFn: func(userID uint) (*models.User, error) {
			if userID == 1 {
				return &models.User{ID: 1}, nil
			}
			return nil, nil
		},
		UpdateUserFn: func(user *models.User, existingUID uint) error {
			if existingUID == 1 {
				return nil
			}
			return errors.New("user not found")
		},
	}

	ctrl := controllers.NewUserController(mockRepo)

	r := gin.Default()
	r.PUT("/users/:id", ctrl.UpdateUserHandler)

	// Test case 1: Invalid ID
	w := performRequest(r, "PUT", "/users/invalid-id", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 2: User not found
	w = performRequest(r, "PUT", "/users/2", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 3: Invalid JSON
	w = performRequest(r, "PUT", "/users/1", "application/json")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 4: Successful update
	validJSON := `{"id": 1, "username": "new_username", "email": "new_email@example.com", "role": "user"}`
	w = performRequest(r, "PUT", "/users/1", validJSON)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetAllUsersHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &controllers.MockUserRepository{
		GetAllUsersFn: func() ([]models.User, error) {
			users := []models.User{
				{ID: 1, Username: "user1", Email: "user1@example.com", Role: "user"},
				{ID: 2, Username: "user2", Email: "user2@example.com", Role: "admin"},
			}
			return users, nil
		},
	}

	ctrl := controllers.NewUserController(mockRepo)

	r := gin.Default()
	r.GET("/users", ctrl.GetAllUsersHandler)

	w := performRequest(r, "GET", "/users", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUsers []models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)

	expectedUsers := []models.User{
		{ID: 1, Username: "user1", Email: "user1@example.com", Role: "user"},
		{ID: 2, Username: "user2", Email: "user2@example.com", Role: "admin"},
	}

	assert.Equal(t, expectedUsers, responseUsers)
}
func TestGetSpecificUserHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &controllers.MockUserRepository{
		GetUserByIDFn: func(userID uint) (*models.User, error) {
			if userID == 1 {
				return &models.User{ID: 1, Username: "user1", Email: "user1@example.com", Role: "user"}, nil
			} else if userID == 2 {
				return &models.User{ID: 2, Username: "user2", Email: "user2@example.com", Role: "admin"}, nil
			}
			return nil, nil // Simulate user not found
		},
	}

	ctrl := controllers.NewUserController(mockRepo)

	r := gin.Default()
	r.GET("/users/:id", ctrl.GetSpecificUserHandler)

	// Test case: User with ID 1 exists
	w1 := performRequest(r, "GET", "/users/1", "")
	assert.Equal(t, http.StatusOK, w1.Code)

	var user1 models.User
	err := json.Unmarshal(w1.Body.Bytes(), &user1)
	assert.NoError(t, err)

	expectedUser1 := models.User{ID: 1, Username: "user1", Email: "user1@example.com", Role: "user"}
	assert.Equal(t, expectedUser1, user1)

	// Test case: User with ID 2 exists
	w2 := performRequest(r, "GET", "/users/2", "")
	assert.Equal(t, http.StatusOK, w2.Code)

	var user2 models.User
	err = json.Unmarshal(w2.Body.Bytes(), &user2)
	assert.NoError(t, err)

	expectedUser2 := models.User{ID: 2, Username: "user2", Email: "user2@example.com", Role: "admin"}
	assert.Equal(t, expectedUser2, user2)

	// Test case: User with non-existing ID
	w3 := performRequest(r, "GET", "/users/3", "")
	assert.Equal(t, http.StatusNotFound, w3.Code)

	var response3 map[string]string
	err = json.Unmarshal(w3.Body.Bytes(), &response3)
	assert.NoError(t, err)

	expectedResponse3 := map[string]string{"error": "user not found"}
	assert.Equal(t, expectedResponse3, response3)

	// Test case: Invalid ID
	w4 := performRequest(r, "GET", "/users/invalid-id", "")
	assert.Equal(t, http.StatusBadRequest, w4.Code)

	var response4 map[string]string
	err = json.Unmarshal(w4.Body.Bytes(), &response4)
	assert.NoError(t, err)

	expectedResponse4 := map[string]string{"error": "Invalid ID"}
	assert.Equal(t, expectedResponse4, response4)
}

func TestCreateProductHandler(t *testing.T) {
	models.InitLogger()

	t.Run("ExistingID", func(t *testing.T) {
		mockRepo := &controllers.MockProductRepository{
			GetProductByIDFn: func(productID uint) (*models.Product, error) {
				// Mock implementation to simulate an existing product with the same ID
				return &models.Product{}, nil
			},
		}

		ctrl := controllers.NewProductController(mockRepo)

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
		mockRepo := &controllers.MockProductRepository{}
		ctrl := controllers.NewProductController(mockRepo)

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
	mockRepo := &controllers.MockProductRepository{
		GetProductByIDFn: func(productID uint) (*models.Product, error) {
			if productID == 1 {
				return &models.Product{ID: 1}, nil
			}
			return nil, nil
		},
		UpdateProductFn: func(product *models.Product, existingPID uint) error {
			if existingPID == 1 {
				return nil
			}
			return errors.New("product not found")
		},
	}

	ctrl := controllers.NewProductController(mockRepo)

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
	mockRepo := &controllers.MockProductRepository{
		GetAllProductsFn: func() ([]models.Product, error) {
			products := []models.Product{
				{ID: 1, Name: "product1", Type: "type1", Quantity: 5},
				{ID: 2, Name: "product2", Type: "type2", Quantity: 10},
			}
			return products, nil
		},
	}

	ctrl := controllers.NewProductController(mockRepo)

	r := gin.Default()
	r.GET("/products", ctrl.GetAllProductsHandler)

	w := performRequest(r, "GET", "/products", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var responseProducts []models.Product
	err := json.Unmarshal(w.Body.Bytes(), &responseProducts)
	assert.NoError(t, err)

	expectedProducts := []models.Product{
		{ID: 1, Name: "product1", Type: "type1", Quantity: 5},
		{ID: 2, Name: "product2", Type: "type2", Quantity: 10},
	}

	assert.Equal(t, expectedProducts, responseProducts)
}

func TestGetSpecificProductHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &controllers.MockProductRepository{
		GetProductByIDFn: func(productID uint) (*models.Product, error) {
			if productID == 1 {
				return &models.Product{ID: 1, Name: "product1", Type: "type1", Quantity: 5}, nil
			} else if productID == 2 {
				return &models.Product{ID: 2, Name: "product2", Type: "type2", Quantity: 10}, nil
			}
			return nil, nil // Simulate product not found
		},
	}

	ctrl := controllers.NewProductController(mockRepo)

	r := gin.Default()
	r.GET("/products/:id", ctrl.GetSpecificProductHandler)

	// Test case: Product with ID 1 exists
	w1 := performRequest(r, "GET", "/products/1", "")
	assert.Equal(t, http.StatusOK, w1.Code)

	var product1 models.Product
	err := json.Unmarshal(w1.Body.Bytes(), &product1)
	assert.NoError(t, err)

	expectedProduct1 := models.Product{ID: 1, Name: "product1", Type: "type1", Quantity: 5}
	assert.Equal(t, expectedProduct1, product1)

	// Test case: Product with ID 2 exists
	w2 := performRequest(r, "GET", "/products/2", "")
	assert.Equal(t, http.StatusOK, w2.Code)

	var product2 models.Product
	err = json.Unmarshal(w2.Body.Bytes(), &product2)
	assert.NoError(t, err)

	expectedProduct2 := models.Product{ID: 2, Name: "product2", Type: "type2", Quantity: 10}
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

// Helper function to perform a request
func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

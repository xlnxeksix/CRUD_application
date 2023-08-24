package user_test

import (
	"awesomeProject1/Models"
	"awesomeProject1/User"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	models.InitLogger()

	t.Run("ExistingID", func(t *testing.T) {
		mockRepo := &user.MockUserRepository{
			GetUserByIDFn: func(userID uint) (*user.User, error) {
				// Mock implementation to simulate an existing user with the same ID
				return &user.User{}, nil
			},
		}

		ctrl := user.NewUserController(mockRepo)

		r := gin.Default()
		r.POST("/users", ctrl.CreateUserHandler)

		// Create a mock request with a user that has an existing ID
		userJSON := `{"username": "testuser", "email": "test@example.com", "role": "user"}`
		w := performRequest(r, "POST", "/users", userJSON)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "User with same ID exists")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockRepo := &user.MockUserRepository{}

		ctrl := user.NewUserController(mockRepo)

		r := gin.Default()
		r.POST("/users", ctrl.CreateUserHandler)

		// Create a mock request with invalid JSON
		w := performRequest(r, "POST", "/users", "invalid-json")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// Add more assertions if needed...
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &user.MockUserRepository{
			CreateUserFn: func(user *user.User) error {
				// Mock implementation for successful user creation
				return nil
			},
			GetUserByIDFn: func(userID uint) (*user.User, error) {
				// Mock implementation to simulate a user not found (for new ID)
				return nil, nil
			},
		}

		ctrl := user.NewUserController(mockRepo)

		r := gin.Default()
		r.POST("/users", ctrl.CreateUserHandler)

		// Create a mock request with a valid user JSON
		userJSON := `{"username": "testuser", "email": "test@example.com", "role": "user"}`
		w := performRequest(r, "POST", "/users", userJSON)

		assert.Equal(t, http.StatusCreated, w.Code)
		// Add more assertions if needed...
	})
}
func TestDeleteUserHandler(t *testing.T) {
	models.InitLogger()

	t.Run("InvalidID", func(t *testing.T) {
		mockRepo := &user.MockUserRepository{}
		ctrl := user.NewUserController(mockRepo)

		r := gin.Default()
		r.DELETE("/users/:id", ctrl.DeleteUserHandler)

		// Create a mock request with an invalid user ID
		w := performRequest(r, "DELETE", "/users/invalid-id", "")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid ID")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo := &user.MockUserRepository{
			GetUserByIDFn: func(userID uint) (*user.User, error) {
				// Mock implementation to simulate user not found
				return nil, nil
			},
		}

		ctrl := user.NewUserController(mockRepo)

		r := gin.Default()
		r.DELETE("/users/:id", ctrl.DeleteUserHandler)

		// Create a mock request with a user ID that doesn't exist
		w := performRequest(r, "DELETE", "/user/123", "")

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "user not found")
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &user.MockUserRepository{
			GetUserByIDFn: func(userID uint) (*user.User, error) {
				// Mock implementation to simulate existing user
				return &user.User{}, nil
			},
			DeleteUserFn: func(userID uint) error {
				// Mock implementation for successful user deletion
				return nil
			},
		}

		ctrl := user.NewUserController(mockRepo)

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
	mockRepo := &user.MockUserRepository{
		GetUserByIDFn: func(userID uint) (*user.User, error) {
			if userID == 1 {
				return &user.User{Username: "test_user_updated", Email: "updated@gmail.com", Role: "test_user"}, nil
			}
			return nil, nil
		},
		UpdateUserFn: func(user *user.User, existingUID uint) error {
			if existingUID == 1 {
				return nil
			}
			return errors.New("user not found")
		},
	}

	ctrl := user.NewUserController(mockRepo)

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
	validJSON := `{"username": "new_username", "email": "new_email@example.com", "role": "user"}`
	w = performRequest(r, "PUT", "/users/1", validJSON)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetAllUsersHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &user.MockUserRepository{
		GetAllUsersFn: func() ([]user.User, error) {
			users := []user.User{
				{Username: "user1", Email: "user1@example.com", Role: "user"},
				{Username: "user2", Email: "user2@example.com", Role: "admin"},
			}
			return users, nil
		},
	}

	ctrl := user.NewUserController(mockRepo)

	r := gin.Default()
	r.GET("/users", ctrl.GetAllUsersHandler)

	w := performRequest(r, "GET", "/users", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUsers []user.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)

	expectedUsers := []user.User{
		{Username: "user1", Email: "user1@example.com", Role: "user"},
		{Username: "user2", Email: "user2@example.com", Role: "admin"},
	}

	assert.Equal(t, expectedUsers, responseUsers)
}
func TestGetSpecificUserHandler(t *testing.T) {
	models.InitLogger()
	mockRepo := &user.MockUserRepository{
		GetUserByIDFn: func(userID uint) (*user.User, error) {
			if userID == 1 {
				return &user.User{Username: "user1", Email: "user1@example.com", Role: "user"}, nil
			} else if userID == 2 {
				return &user.User{Username: "user2", Email: "user2@example.com", Role: "admin"}, nil
			}
			return nil, nil // Simulate user not found
		},
	}

	ctrl := user.NewUserController(mockRepo)

	r := gin.Default()
	r.GET("/users/:id", ctrl.GetSpecificUserHandler)

	// Test case: User with ID 1 exists
	w1 := performRequest(r, "GET", "/users/1", "")
	assert.Equal(t, http.StatusOK, w1.Code)

	var user1 user.User
	err := json.Unmarshal(w1.Body.Bytes(), &user1)
	assert.NoError(t, err)

	expectedUser1 := user.User{Username: "user1", Email: "user1@example.com", Role: "user"}
	assert.Equal(t, expectedUser1, user1)

	// Test case: User with ID 2 exists
	w2 := performRequest(r, "GET", "/users/2", "")
	assert.Equal(t, http.StatusOK, w2.Code)

	var user2 user.User
	err = json.Unmarshal(w2.Body.Bytes(), &user2)
	assert.NoError(t, err)

	expectedUser2 := user.User{Username: "user2", Email: "user2@example.com", Role: "admin"}
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

func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

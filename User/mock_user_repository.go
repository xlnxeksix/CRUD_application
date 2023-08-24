package user

type MockUserRepository struct {
	CreateUserFn       func(user *User) error
	GetUserByIDFn      func(userID uint) (*User, error)
	DeleteUserFn       func(userID uint) error
	UpdateUserFn       func(user *User, existingUID uint) error
	GetAllUsersFn      func() ([]User, error)
	AuthenticateUserFn func(username, passwd string) (*User, error)
	GetUserRoleFn      func(username string) (string, error)
}

func (m *MockUserRepository) CreateUser(user *User) error {
	return m.CreateUserFn(user)
}

func (m *MockUserRepository) GetUserByID(userID uint) (*User, error) {
	return m.GetUserByIDFn(userID)
}

func (m *MockUserRepository) DeleteUser(userID uint) error {
	return m.DeleteUserFn(userID)
}

func (m *MockUserRepository) UpdateUser(user *User, existingUID uint) error {
	return m.UpdateUserFn(user, existingUID)
}

func (m *MockUserRepository) GetAllUsers() ([]User, error) {
	return m.GetAllUsersFn()
}

func (m *MockUserRepository) AuthenticateUser(username, passwd string) (*User, error) {
	return m.AuthenticateUserFn(username, passwd)
}

func (m *MockUserRepository) GetUserRole(username string) (string, error) {
	return m.GetUserRoleFn(username)
}

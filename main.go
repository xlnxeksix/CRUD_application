package main

/*
import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createUser(db *gorm.DB, id int, userName string, email string, role string) {
	user := User{id: id, userName: userName, email: email, role: role}
	db.Create(&user)
}
func getUsers(db *gorm.DB) []User {
	var users []User
	db.Find(&users)
	return users
}
func main() {
	// Connect to the PostgreSQL database
	dsn := "host=localhost user=postgres password=docker dbname=user_database port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate the model
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to auto migrate")
	}

	// Start the CRUD operations
	// ...

	createUser(db, 0, "firstUser", "first@example.com", "founder")
	var users []User

	users = getUsers(db)
	for _, v := range users {
		fmt.Println(v.id)
	}
	// Retrieve users
	//var users []User
	//db.Find(&users)

	// Update a user
	//user.userName = "Jane Doe"
	//db.Save(&user)

	// Delete a user
	//db.Delete(&user)
}
*/

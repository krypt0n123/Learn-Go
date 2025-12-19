package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Todo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func Init() {

	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//迁移schema
	DB.AutoMigrate(&Todo{})
}

func CreateTodo(todo Todo) error {
	result := DB.Create(&todo)
	return result.Error
}

func GetAllTodos() ([]Todo, error) {
	var todos []Todo
	result := DB.Find(&todos)
	return todos, result.Error
}

func UpdateTodo(todo Todo) error {
	result := DB.Save(&todo)
	return result.Error
}

func DeleteTodo(id string) error {
	result := DB.Where("id=?", id).Delete(&Todo{})
	return result.Error
}

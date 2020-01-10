package main

import (
	"fmt"
)

// Model pointer value
type Model struct {
	ID int
}

// Print implementation.
func (m *Model) Print() {
	fmt.Print(m.ID)
}

// User pointer value
type User struct {
	Model
	Name string
	Age  *int `gorm:"default:18"`
}

// Print implementation.
func (u *User) Print() {
	fmt.Print(u.Name)
}

func main() {
	u := &User{}
	u.ID = 1
	u.Name = "123"
	u.Print()
}

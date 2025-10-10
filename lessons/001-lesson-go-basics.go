// Lesson 001: Go Basics for Building a Chat App
// ================================================
// Welcome Konrad! Since you're coming from Python and JavaScript,
// let's explore Go's key differences and similarities.

package main

import (
	"fmt"
)

// 1. FUNCTIONS IN GO
// Unlike Python's 'def' or JavaScript's 'function', Go uses 'func'
// Go is statically typed - you must declare types for parameters and return values
func greetUser(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to Go.", name)
}

// 2. VARIABLES AND TYPES
// Go has several ways to declare variables:
func variableExamples() {
	// Method 1: var keyword with type
	var message string = "Starting our chat app journey"

	// Method 2: var with type inference
	var userCount = 10

	// Method 3: Short declaration (only inside functions!)
	// This is like Python's simple assignment, but with :=
	isOnline := true

	fmt.Println(message, userCount, isOnline)
}

// 3. STRUCTS - Go's way of creating custom types
// Think of these like Python classes but without methods attached directly
// Or like JavaScript objects but with a fixed structure
type User struct {
	Name     string
	ID       int
	IsOnline bool
}

// 4. METHODS - Functions attached to structs
// The (u User) part is called a "receiver" - it's like 'self' in Python
func (u User) GetStatus() string {
	if u.IsOnline {
		return fmt.Sprintf("%s is online", u.Name)
	}
	return fmt.Sprintf("%s is offline", u.Name)
}

// 5. SLICES - Go's dynamic arrays (like Python lists or JS arrays)
func sliceExample() {
	// Creating a slice of strings for chat messages
	messages := []string{
		"Hey there!",
		"How's the Go learning going?",
		"Building cool stuff!",
	}

	// Appending to a slice (like Python's append or JS push)
	messages = append(messages, "New message!")

	// Iterating with range (like Python's enumerate)
	for index, msg := range messages {
		fmt.Printf("Message %d: %s\n", index, msg)
	}
}

// MAIN FUNCTION - Entry point of every Go program
func main() {
	fmt.Println("=== Go Basics Lesson ===")

	// Testing our greeting function
	greeting := greetUser("Konrad")
	fmt.Println(greeting)
	fmt.Println()

	// Testing variables
	fmt.Println("Variable examples:")
	variableExamples()
	fmt.Println()

	// Testing structs and methods
	fmt.Println("Struct example:")
	entrepreneur := User{
		Name:     "Konrad",
		ID:       1,
		IsOnline: true,
	}
	fmt.Println(entrepreneur.GetStatus())
	fmt.Println()

	// Testing slices
	fmt.Println("Slice example:")
	sliceExample()
}

// To run this lesson:
// 1. Save this file
// 2. In your terminal, navigate to the app directory
// 3. Run: go run 001-lesson-go-basics.go
//
// Try it out and let me know what you see!

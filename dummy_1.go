package main

import "fmt"

// Buggy method with multiple violations
func ProcessData(data []string) string {
	result := data[0] + data[1]                               
	db.Exec("DELETE FROM users WHERE name = '" + result + "'") 
	fmt.Println("Processing:", result)                         
	return result                                              
}

package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Book(c *gin.Context) {
	if c.Request.Method == "POST" {
		type UserData struct {
			BookName        string `json:"book_name"`
			Author          string `json:"author"`
			PublicationYear string `json:"publication_year"`
		}
		var publicationYear string
		var bookName string
		var author string
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found"})
			return
		}
		if userRole != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You dont have access to add books"})
			return
		}
		var user_data UserData
		if err := c.Bind(&user_data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}
		if user_data.BookName != "" {
			bookName = user_data.BookName
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "book name should not be empty"})
			return
		}
		if user_data.Author != "" {
			author = user_data.Author
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "author should not be empty"})
			return
		}
		year, _ := strconv.Atoi(user_data.PublicationYear)
		if year >= 1000 && year <= time.Now().Year() {
			publicationYear = user_data.PublicationYear
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication year"})
			return
		}
		newRecord := []string{bookName, author, string(publicationYear)}
		dataToWrite := [][]string{newRecord}
		file, err := os.OpenFile("./csv/regularUser.csv", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open CSV file"})
			return
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		if err := writer.WriteAll(dataToWrite); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data to CSV file"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})
	} else if c.Request.Method == "GET" {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found"})
			return
		}

		var filePaths []string
		switch userRole {
		case "admin":
			filePaths = []string{"./csv/adminUser.csv", "./csv/regularUser.csv"}
		case "user":
			filePaths = []string{"./csv/regularUser.csv"}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role"})
			return
		}

		var allData []map[string]string
		for _, filePath := range filePaths {
			file, err := os.Open(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error opening file: %s", err.Error())})
				return
			}
			defer file.Close()

			reader := csv.NewReader(file)
			data, err := reader.ReadAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading CSV file: %s", err.Error())})
				return
			}
			columnNames := data[0]
			dataWithoutHeader := data[1:]
			for _, row := range dataWithoutHeader {
				record := make(map[string]string)
				for i, column := range columnNames {
					record[column] = row[i]
				}
				allData = append(allData, record)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   allData,
			"length": len(allData),
		})
	}
}
func RemoveBook(c *gin.Context) {
	type UserData struct {
		BookName string `json:"book_name"`
	}
	var bookName string
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found"})
		return
	}
	if userRole != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You dont have access to remove books"})
		return
	}
	var user_data UserData
	if err := c.Bind(&user_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Extract the book name provided by the user
	if user_data.BookName != "" {
		bookName = strings.ToLower(user_data.BookName)
	}

	// Open the CSV file
	file, err := os.OpenFile("./csv/regularUser.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open CSV file"})
		return
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read CSV file"})
		return
	}

	// Filter out rows that match the given book name
	filteredRows := make([][]string, 0)
	for _, row := range rows {
		if strings.ToLower(row[0]) != bookName { // Assuming the book name is in the first column
			filteredRows = append(filteredRows, row)
		}
	}

	// Write the filtered data back to the CSV file
	file.Truncate(0) // Clear the file contents
	file.Seek(0, 0)  // Rewind to the beginning of the file

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(filteredRows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data to CSV file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record removed from CSV file successfully"})
}

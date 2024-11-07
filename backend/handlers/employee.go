// handlers/employee.go
package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"radhaLLC/db"
	"radhaLLC/models"
	"radhaLLC/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper function to convert byte slice to ReadCloser
func toReadCloser(body []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewBuffer(body))
}

// AddEmployee handles employee creation
func AddEmployee(c *gin.Context) {
	var employee models.Employee

	// Log the raw request for debugging
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}
	log.Printf("Received request body: %s", string(body))

	// Restore the request body
	c.Request.Body = toReadCloser(body)

	// Bind JSON to employee struct
	if err := c.ShouldBindJSON(&employee); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid employee data",
			"details": err.Error(),
		})
		return
	}

	// Generate EmployeeID
	employee.EmployeeID = utils.GenerateEmployeeID(employee.FirstName, employee.LastName)

	// Set default values
	employee.Status = "Active"
	employee.LeaveBalance = 20
	employee.LastLeaveDate = nil

	// Get MongoDB collection with error handling
	collection, err := db.GetCollection("employees")
	if err != nil {
		log.Printf("Error getting MongoDB collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database connection error",
			"details": err.Error(),
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert employee
	result, err := collection.InsertOne(ctx, employee)
	if err != nil {
		log.Printf("Error inserting employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add employee",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully inserted employee with ID: %v", result.InsertedID)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Employee added successfully",
		"employeeID": employee.EmployeeID,
	})
}

// GetEmployees handles fetching all employees with pagination
func GetEmployees(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	skip := (page - 1) * limit

	collection, err := db.GetCollection("employees")
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get total count
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting employees"})
		return
	}

	// Set find options
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "firstName", Value: 1}})

	// Execute find
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Printf("Error finding employees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching employees"})
		return
	}
	defer cursor.Close(ctx)

	// Decode results
	var employees []models.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		log.Printf("Error decoding employees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding employees"})
		return
	}

	// Calculate total pages
	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
		"pagination": gin.H{
			"currentPage": page,
			"totalPages":  totalPages,
			"totalItems":  total,
			"limit":       limit,
		},
	})
}

// GetEmployeeByID handles fetching a single employee by ID
func GetEmployeeByID(c *gin.Context) {
	id := c.Param("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection, err := db.GetCollection("employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var employee models.Employee
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&employee)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// UpdateEmployee handles updating an employee
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Read and log the raw request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}
	log.Printf("Update request body: %s", string(body))

	// Restore the request body
	c.Request.Body = toReadCloser(body)

	var updateData models.Employee
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	collection, err := db.GetCollection("employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create update document
	update := bson.M{
		"$set": bson.M{
			"firstName":    updateData.FirstName,
			"lastName":     updateData.LastName,
			"department":   updateData.Department,
			"position":     updateData.Position,
			"status":       updateData.Status,
			"salary":       updateData.Salary,
			"hireDate":     updateData.HireDate,
			"leaveBalance": updateData.LeaveBalance,
			"contactInfo":  updateData.ContactInfo,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Printf("Error updating employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating employee"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

// DeleteEmployee handles employee deletion
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection, err := db.GetCollection("employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error deleting employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting employee"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

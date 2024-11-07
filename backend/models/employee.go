// models/employee.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Employee structure matching MongoDB fields
type Employee struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID    string             `json:"employeeID" bson:"employeeID"` // Changed to string for custom ID format
	FirstName     string             `json:"firstName" bson:"firstName"`
	LastName      string             `json:"lastName" bson:"lastName"`
	Department    string             `json:"department" bson:"department"`
	Position      string             `json:"position" bson:"position"`
	Status        string             `json:"status" bson:"status"`
	Salary        int                `json:"salary" bson:"salary"`
	HireDate      string             `json:"hireDate" bson:"hireDate"` // Format as date string in frontend
	LastLeaveDate *string            `json:"lastLeaveDate" bson:"lastLeaveDate"`
	LeaveBalance  int                `json:"leaveBalance" bson:"leaveBalance"`
	ContactInfo   struct {
		Email string `json:"email" bson:"email"`
		Phone string `json:"phone" bson:"phone"`
	} `json:"contactInfo" bson:"contactInfo"`
}

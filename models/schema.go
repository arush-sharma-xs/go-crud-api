package models

import "time"

// Library: (ID, Name)
type Library struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `form:"name" json:"name" validate:"required" binding:"required"`
}

// Users: (ID, Name, Email, ContactNumber, Role, LibID)
type Users struct {
	ID            uint   `gorm:"primaryKey;"`
	Name          string `gorm:"omitempty" json:"name" binding:"required"`
	Email         string `json:"email" binding:"required, email"`
	ContactNumber string `json:"contactNumber" binding:"required"`
	Role          string `json:"role" `
	LibID         uint   `json:"libId" binding:"required"`
}

// BookInventory : (ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies)
type BookInventory struct {
	ISBN            uint   `json:"isbn" gorm:"primaryKey;autoIncrement:true"`
	LibID           uint   `json:"libId" validate:"required"`
	Title           string `json:"title" validate:"required"`
	Authors         string `json:"authors" validate:"required"`
	Publisher       string `json:"publisher" validate:"required"`
	Version         uint   `json:"version" validate:"required"`
	TotalCopies     uint   `json:"totalCopies" validate:"required"`
	AvailableCopies uint   `json:"availableCopies"`
}

// RequestEvents: (ReqID, BookID, ReaderID, RequestDate, ApprovalDate, ApproverID, RequestType)
type RequestEvents struct {
	ReqID        uint      `gorm:"primaryKey;autoIncrement:true"`
	BookID       uint      `json:"bookId" validate:"required"`
	ReaderID     uint      `json:"readerId" validate:"required"`
	RequestDate  time.Time `gorm:"default:current_timestamp"`
	ApprovalDate time.Time `gorm:"type:timestamp;default:null"`
	ApproverID   uint      `gorm:"default:null"`
	RequestType  string    `json:"requestType" validate:"required"`
}

// IssueRegistery: (IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate
type IssueRegistery struct {
	IssueID            uint      `gorm:"primaryKey;autoIncrement:true"`
	ISBN               uint      `json:"isbn" validate:"required"`
	ReaderID           uint      `json:"readerId" validate:"required"`
	IssueApproverID    uint      `json:"issueApproverId" validate:"required"`
	IssueStatus        string    `json:"issueStatus" validate:"required"`
	IssueDate          time.Time `gorm:"default:current_timestamp"`
	ExpectedReturnDate time.Time `gorm:"type:timestamp;null;default:null"`
	ReturnDate         time.Time `gorm:"type:timestamp;null;default:null"`
	ReturnApproverID   uint      `gorm:"default:null"`
}

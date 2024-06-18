package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xs-arush-0856/lms/models"
)

type AuthRequest struct {
	LibID uint   `json:"libId"`
	Email string `json:"email" validate:"required"`
}

func auth(c *gin.Context) {
	var authRequest AuthRequest
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Now Check if user found or not;
	var user models.Users
	res := models.DB.First(&user, "email = ? AND lib_id = ?", authRequest.Email, authRequest.LibID)

	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found with this email in Library"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userData": user})
}

func handleCreateLibrary(c *gin.Context) {
	var newLib models.Library
	if err := c.ShouldBindJSON(&newLib); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	LibNameSize := len(newLib.Name)

	if LibNameSize < 3 || LibNameSize > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library Length should between 3 and 10."})
		return
	}

	// checking if library already exist
	res := models.DB.First(&newLib, "name = ?", newLib.Name)

	// if record found, res badrequest
	if res.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library already exist, Enter new name"})
		return
	}

	models.DB.Create(&newLib)
	c.JSON(http.StatusOK, gin.H{"status": "Library Created."})
}

func addUser(c *gin.Context) {
	
	var user models.Users
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if useralreadyFound
	checkUserRecord := models.DB.First(&user, "email = ?", user.Email, user.LibID)
	if checkUserRecord.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist."})
		return
	}

	//check if library Exist or not
	var lib models.Library
	checkLibraryRecord := models.DB.Model(&models.Library{}).First(&lib, "id = ?", user.LibID)
	log.Println("Rows", checkLibraryRecord.RowsAffected)

	if checkLibraryRecord.RowsAffected <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library not exist."})
		return
	}

	models.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"status": "User Created..."})
}

func getBook(c *gin.Context) {
	var books []models.BookInventory
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"books": books})
}

func addBook(c *gin.Context) {
	var book models.BookInventory
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fetchedBook models.BookInventory
	if err := models.DB.Where(models.BookInventory{Title: book.Title, Authors: book.Authors, LibID: book.LibID}).First(&fetchedBook).Error; err != nil {
		// create new Book if not exist.
		models.DB.Create(&book)
		c.JSON(http.StatusOK, gin.H{"status": "Book Added!"})
		return
	}

	// Increse the book count if exist.
	models.DB.Model(&book).Where("title=?", book.Title).Where("authors=?", book.Authors).Updates(models.BookInventory{TotalCopies: fetchedBook.TotalCopies + book.TotalCopies, AvailableCopies: fetchedBook.AvailableCopies + book.TotalCopies})
	c.JSON(http.StatusOK, gin.H{"status": "Book Updated!"})
}

func removeBook(c *gin.Context) {
	// first check if book exist or not;
	bookId, _ := strconv.Atoi(c.Param("isbn"))

	var book models.BookInventory
	if err := models.DB.Where("isbn = ?", bookId).First(&book).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "Book not found !"})
		return
	}

	// checkissueed Copy
	issuedCopy := book.TotalCopies - book.AvailableCopies

	if book.TotalCopies > issuedCopy && book.AvailableCopies > 0 {
		// descrease count
		models.DB.Model(&book).Where("isbn = ?", bookId).Updates(models.BookInventory{AvailableCopies: book.AvailableCopies - 1})
		models.DB.Model(&book).Where("isbn = ?", bookId).Updates(models.BookInventory{TotalCopies: book.AvailableCopies + issuedCopy})
		c.JSON(http.StatusOK, gin.H{"status": "Book removed"})

		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not available"})
		return
	}
}

func updateBook(c *gin.Context) {
	bookId, _ := strconv.Atoi(c.Param("isbn"))
	var book models.BookInventory
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if book exist or Not
	var fetchBook models.BookInventory
	if err := models.DB.Where("isbn = ?", bookId).First(&fetchBook).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "Book not found !"})
		return
	}

	models.DB.Model(&book).Where("isbn = ?", bookId).Updates(book)
	c.JSON(http.StatusOK, gin.H{"status": "Book Updated!"})
}

func getIssueRequest(c *gin.Context) {
	var issueList []models.RequestEvents
	models.DB.Find(&issueList).Where("approver_id = ?", 1)
	c.JSON(http.StatusOK, gin.H{"issues": issueList})
}

func searchBook(c *gin.Context) {
	queryKey := c.Param("key")
	queryValue := c.Param("value")

	// Check if book exist or Not
	var fetchBook []models.BookInventory
	if err := models.DB.Where(queryKey+" LIKE ?", "%"+queryValue+"%").Find(&fetchBook).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "Book not found ! Can't Update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": fetchBook})
}

func raiseIssue(c *gin.Context) {

	var raisedIssue models.RequestEvents
	if err := c.ShouldBindJSON(&raisedIssue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if books is found
	var fetchBook models.BookInventory
	res := models.DB.Where(models.BookInventory{ISBN: raisedIssue.BookID}).Find(&fetchBook)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	//  check if issue request already in go or not
	var checkRequestExist models.RequestEvents
	checkRequestRes := models.DB.Where(models.RequestEvents{BookID: raisedIssue.BookID, ReaderID: raisedIssue.ReaderID, RequestType: raisedIssue.RequestType}).Find(&checkRequestExist)

	if checkRequestRes.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book " + raisedIssue.RequestType + " already raised.."})
		return
	}

	if raisedIssue.RequestType == "issue" {
		// check if book copies is available or not
		if fetchBook.AvailableCopies <= 0 {
			// return when book will be available.
			c.JSON(http.StatusBadRequest, gin.H{"status": "Books copies are not available,"})
			return
		}
		models.DB.Create(&raisedIssue)
		c.JSON(http.StatusOK, gin.H{"status": "Issue Raised"})
	}

	if raisedIssue.RequestType == "return" {
		// check if book issued or not
		var checkIssuedBook models.IssueRegistery
		checkIssuedBookRes := models.DB.Where(models.IssueRegistery{ISBN: raisedIssue.BookID, ReaderID: raisedIssue.ReaderID}).Find(&checkIssuedBook)

		if checkIssuedBookRes.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Book never issued"})
			return
		}

		models.DB.Create(&raisedIssue)
		c.JSON(http.StatusOK, gin.H{"status": "Issue Raised"})
		return
	}
}


func handleApproveRequest(c *gin.Context) {
	// 	(IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate,
	// ExpectedReturnDate, ReturnDate, ReturnApproverID)

	reqId,_ := strconv.Atoi(c.Param("reqId"))

	var issueReg models.IssueRegistery
	if err := c.ShouldBindJSON(&issueReg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	expectedDate := time.Now().AddDate(0, 0, 14)
	models.DB.Create(&issueReg)
	models.DB.Model(&issueReg).Updates(models.IssueRegistery{ExpectedReturnDate: expectedDate})

	// update request Event dataes
	var reqEvent models.RequestEvents
	models.DB.Model(&reqEvent).Where("req_id = ?", reqId).Updates(models.RequestEvents{ApproverID: issueReg.IssueApproverID, ApprovalDate: time.Now()})
	//descrese available book count
	
	var book models.BookInventory
	models.DB.First(&book).Where("isbn = ?", issueReg.ISBN)
	models.DB.Model(&book).Where("isbn = ?", issueReg.ISBN).Updates(models.BookInventory{AvailableCopies: book.AvailableCopies - 1 })

}

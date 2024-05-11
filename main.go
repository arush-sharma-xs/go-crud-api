package main

// Library: (ID, Name)
type Library struct {
	ID		uint
	Name	string
}
	
// Users: (ID, Name, Email, ContactNumber, Role, LibID)
type Users struct {
	ID							uint	
	Name						string
	Email 					string
	ContactNumber 	uint
	Role 						string
	LibID 					uint
}

// BookInventory : (ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies)
type BookInventory struct {
	ISBN						uint
	LibID						uint
	Title						string
	Authors 				string
	Publisher				string
	Version					string
	TotalCopies			uint
	AvailableCopies	uint
}

// RequestEvents: (ReqID, BookID, ReaderID, RequestDate, ApprovalDate, ApproverID, RequestType)
type RequestEvents struct {
	ReqID						uint
	BookID					uint
	ReaderID				string
	RequestDate			string
	ApprovalDate		string
	ApproverID			string
	RequestType			uint
}

// IssueRegistery: (IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate,
type IssueRegistery struct {
	IssueID						uint
	ISBN							uint
	ReaderID					uint
	IssueApproverID		uint
	IssueStatus				string
	IssueDate					string
}

	// ExpectedReturnDate, ReturnDate, ReturnApproverID)

type ExpectedReturnDate struct {
	ReturnDate				uint
	ReturnApproverID	uint
	
}

func main() {
		
}
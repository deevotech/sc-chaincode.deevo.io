package main

const (
	CK_AUDIT_OBJ     = "auditedObject~audit"
	CK_AUDITOR_AUDIT = "auditor~audit"
)

// Org model
type Org struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
}

// Party model
type Party struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	ORG        string `json:"org"`
}

// Location model
type Location struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Party      string `json:"party"`
}

// Product model
type Product struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}

// Log model
type Log struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Content    string `json:"content"`
	Time       int64  `json:"name"`
	Location   string `json:"location"`
	ObjectID   string `json:"objectID"`
}

// InitData model
type InitData struct {
	ORG       Org        `json:"org"`
	Parties   []Party    `json:"parties"`
	Locations []Location `json:"locations"`
	Products  []Product  `json:"products"`
	Auditors  []Auditor  `json:"auditors"`
}

// Auditor model
type Auditor struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
}

// AuditAction model
type AuditAction struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Time       int64  `json:"time"`
	Auditor    string `json:"auditor"`
	Location   string `json:"location"`
	ObjectID   string `json:"objectID"`
}

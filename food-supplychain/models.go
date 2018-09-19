package main

import "sort"

const (
	CK_AUDIT_OBJ     = "auditedObject~audit"
	CK_AUDITOR_AUDIT = "auditor~audit"

	TYPE_LOG         = "log"
	TYPE_AUDITACTION = "auditAction"
	TYPE_AUDITOR     = "auditor"
)

// InitData model
type InitData struct {
	Traceable []Traceable `json:"traceable"`
	Auditors  []Auditor   `json:"auditors"`
}

// LiteModel model
type LiteModel struct {
	ObjectType string `json:"objectType"`
	ID         string `json:"id"`
}

// Traceable model
type Traceable struct {
	ObjectType string `json:"objectType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Parent     string `json:"parent"`
}

// Log model
type Log struct {
	ObjectType string   `json:"objectType"`
	ID         string   `json:"id"`
	Time       int64    `json:"time"`
	Ref        []string `json:"ref"`
	CTE        string   `json:"cte"`
	Content    string   `json:"content"`
	Asset      string   `json:"asset"`
	Product    string   `json:"product"`
	Location   string   `json:"location"`
}

// Equals compare 2 logs
func (l *Log) Equals(other Log) bool {
	if l.ObjectType != other.ObjectType {
		return false
	}
	if l.ID != other.ID {
		return false
	}
	if l.Time != other.Time {
		return false
	}
	if l.CTE != other.CTE {
		return false
	}
	if l.Content != other.Content {
		return false
	}
	if l.Asset != other.Asset {
		return false
	}
	if l.Product != other.Product {
		return false
	}
	if l.Location != other.Location {
		return false
	}
	if len(l.Ref) != len(other.Ref) {
		return false
	}
	ref1 := make([]string, len(l.Ref))
	ref2 := make([]string, len(other.Ref))
	copy(l.Ref, ref1)
	copy(other.Ref, ref2)
	sort.Strings(ref1)
	sort.Strings(ref2)
	for index, item := range ref1 {
		if item != ref2[index] {
			return false
		}
	}
	return true
}

// Auditor model
type Auditor struct {
	ObjectType string `json:"objectType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
}

// AuditAction model
type AuditAction struct {
	ObjectType string `json:"objectType"`
	ID         string `json:"id"`
	Time       int64  `json:"time"`
	Auditor    string `json:"auditor"`
	Location   string `json:"location"`
	ObjectID   string `json:"objectID"`
	Content    string `json:"content"`
}

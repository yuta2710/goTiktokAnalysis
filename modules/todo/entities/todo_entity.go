package entities

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type MemberRole string

// Owner, Collaborators, Viewer, Assignee, Reviewers
const (
	Owner        MemberRole = "Owner"
	Collaborator MemberRole = "Collaborator"
	Viewer       MemberRole = "Viewer"
	Assignee     MemberRole = "Assignee"
	Reviewer     MemberRole = "Reviewer"
	Member       MemberRole = "Member"
)

type Priority string

const (
	Low    Priority = "Low"
	Medium Priority = "Medium"
	High   Priority = "High"
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}

//	func (p Priority) Value() (driver.Value, error) {
//		return int(p), nil
//	}
func (p *Priority) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan Priority: %v is not a string", value)
	}

	switch strValue {
	case "Low":
		*p = Low
	case "Medium":
		*p = Medium
	case "High":
		*p = High
	default:
		return fmt.Errorf("invalid priority value: %v", strValue)
	}

	return nil
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	var priorityStr string
	// Unmarshal the data into a string
	if err := json.Unmarshal(data, &priorityStr); err != nil {
		return fmt.Errorf("priority should be a string, got %s", data)
	}

	// Convert string to Priority
	switch priorityStr {
	case "Low":
		*p = Low
	case "Medium":
		*p = Medium
	case "High":
		*p = High
	default:
		return fmt.Errorf("invalid priority value: %s", priorityStr)
	}

	return nil
}

func ConvertPriorityToEnum(priority Priority) string {
	switch priority {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "" // Or handle invalid case appropriately
	}
}

type Todo struct {
	shared.BaseSQLModel
	UserId      int          `gorm:"column:user_id;not null" json:"userId"`                             // Matches the user_id foreign key in the database
	Title       string       `gorm:"column:title;type:VARCHAR(255);not null" json:"title"`              // Matches the title column in the database
	Description string       `gorm:"column:description;type:TEXT" json:"description"`                   // Matches the description column in the database
	IsCompleted bool         `gorm:"column:is_completed;type:BOOLEAN;default:false" json:"isCompleted"` // Matches the is_completed column in the database
	DueDate     time.Time    `gorm:"column:due_date;type:TIMESTAMP" json:"dueDate"`                     // Matches the due_date column in the database
	Priority    Priority     `gorm:"column:priority;type:priority;not null" json:"priority"`            // Matches the priority column as ENUM
	CreatedAt   time.Time    `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`      // Tracks the creation time
	UpdatedAt   time.Time    `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updatedAt"`      // Tracks the update time
	Members     []TodoMember `gorm:"foreignKey:TodoId;references:Id" json:"members"`                    // For TodoMember association
}

// title, description, isCompleted, dueDate, priority

type TodoMember struct {
	// ID     uint       `gorm:"primaryKey"`
	TodoId int        `gorm:"column:todo_id;not null" json:"-"`
	UserId int        `gorm:"column:user_id; not null"`
	Role   MemberRole `gorm:"column:role"`
}

func (td *Todo) Mask(dbType shared.DbType) {
	td.BaseSQLModel.Mask(dbType)
}

// The package models contains the data structures used in the application.
// It defines the User, Question, Answer, and Room types.
package models

type UserRole string

const (
	AdminRole   UserRole = "admin"
	RegularRole UserRole = "user"
	GuestRole   UserRole = "guest"
)

type User struct {
	ID        int64    `json:"id" db:"id"`
	Email     string   `json:"email" db:"email"`
	Name      string   `json:"name" db:"name"`
	Role      UserRole `json:"role" db:"role"`
	CreatedAt string   `json:"created_at" db:"created_at"`
}

type contextKey string

const (
	UserContextKey contextKey = "user_id"
)

func (u User) IsAdmin() bool {
	return u.Role == AdminRole
}

func (r UserRole) IsValid() bool {
	switch r {
	case AdminRole, RegularRole, GuestRole:
		return true
	}
	return false
}

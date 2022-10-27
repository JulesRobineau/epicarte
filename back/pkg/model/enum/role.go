package enum

import "database/sql/driver"

type Role string

const (
	SUPERADMIN Role = "superadmin"
	ADMIN      Role = "admin"
	STUDENT    Role = "student"
)

func (r *Role) Scan(value interface{}) error {
	*r = Role(value.(string))
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case SUPERADMIN, ADMIN, STUDENT:
		return true
	default:
		return false
	}
}

// HasPermission checks if the role has the permission
func (r Role) HasPermission(permission Role) bool {
	perm := map[Role][]Role{
		SUPERADMIN: {SUPERADMIN, ADMIN, STUDENT},
		ADMIN:      {ADMIN, STUDENT},
		STUDENT:    {STUDENT},
	}[r]

	for _, p := range perm {
		if p == permission {
			return true
		}
	}
	return false
}

package enum

import "database/sql/driver"

type Role string

const (
	SUPERADMIN Role = "superadmin"
	ADMIN      Role = "admin"
	USER       Role = "user"
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
	case SUPERADMIN, ADMIN, USER:
		return true
	default:
		return false
	}
}

// HasPermission checks if the role has the permission
func (r Role) HasPermission(permission Role) bool {
	perm := map[Role][]Role{
		SUPERADMIN: {SUPERADMIN, ADMIN, USER},
		ADMIN:      {ADMIN, USER},
		USER:       {USER},
	}[r]

	for _, p := range perm {
		if p == permission {
			return true
		}
	}
	return false
}

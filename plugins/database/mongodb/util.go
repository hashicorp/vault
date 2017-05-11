package mongodb

type createUserCommand struct {
	Username string        `bson:"createUser"`
	Password string        `bson:"pwd"`
	Roles    []interface{} `bson:"roles"`
}
type mongodbRole struct {
	Role string `json:"role" bson:"role"`
	DB   string `json:"db"   bson:"db"`
}

type mongodbRoles []mongodbRole

type mongoDBStatement struct {
	DB    string       `json:"db"`
	Roles mongodbRoles `json:"roles"`
}

// Convert array of role documents like:
//
// [ { "role": "readWrite" }, { "role": "readWrite", "db": "test" } ]
//
// into a "standard" MongoDB roles array containing both strings and role documents:
//
// [ "readWrite", { "role": "readWrite", "db": "test" } ]
//
// MongoDB's createUser command accepts the latter.
func (roles mongodbRoles) toStandardRolesArray() []interface{} {
	var standardRolesArray []interface{}
	for _, role := range roles {
		if role.DB == "" {
			standardRolesArray = append(standardRolesArray, role.Role)
		} else {
			standardRolesArray = append(standardRolesArray, role)
		}
	}
	return standardRolesArray
}

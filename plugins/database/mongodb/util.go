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

// writeConcern is the mgo.Safe struct with JSON tags
// More info: https://godoc.org/gopkg.in/mgo.v2#Safe
type writeConcern struct {
	W        int    `json:"w"`        // Min # of servers to ack before success
	WMode    string `json:"w_mode"`   // Write mode for MongoDB 2.0+ (e.g. "majority")
	WTimeout int    `json:"wtimeout"` // Milliseconds to wait for W before timing out
	FSync    bool   `json:"fsync"`    // Sync via the journal if present, or via data files sync otherwise
	J        bool   `json:"j"`        // Sync via the journal if present
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

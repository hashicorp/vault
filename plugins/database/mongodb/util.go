package mongodb

import (
	"encoding/json"
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

type createUserCommand struct {
	Username string        `bson:"createUser"`
	Password string        `bson:"pwd"`
	Roles    []interface{} `bson:"roles"`
}

type upsertUserCommand struct {
	Username string        `bson:"updateUser"`
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

// mgo.Role is a named string type
func (roles mongodbRoles) toStandardRolesStringArray() []mgo.Role {
	var standardRolesArray []mgo.Role
	for _, role := range roles {
		if role.DB == "" {
			standardRolesArray = append(standardRolesArray, mgo.Role(role.Role))
		} else {
			b, err := json.Marshal(role)
			if err != nil {
				fmt.Println("error:", err)
			}
			s := string(b)
			standardRolesArray = append(standardRolesArray, mgo.Role(s))
		}
	}
	return standardRolesArray
}

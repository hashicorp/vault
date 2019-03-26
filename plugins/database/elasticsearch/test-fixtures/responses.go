package fixtures

const (
	CreateRoleResponse = `{
	  "role": {
	    "created": true 
	  }
	}`

	GetRoleResponse = `{
	  "role-name": {
	    "cluster" : [ "all" ],
	    "indices" : [
	      {
	        "names" : [ "index1", "index2" ],
	        "privileges" : [ "all" ],
	        "field_security" : {
	          "grant" : [ "title", "body" ]
	        }
	      }
	    ],
	    "applications" : [ ],
	    "run_as" : [ "other_user" ],
	    "metadata" : {
	      "version" : 1
	    },
	    "transient_metadata": {
	      "enabled": true
	    }
	  }
	}`

	DeleteRoleResponse = `{
	  "found" : true
	}`

	CreateUserResponse = `{
	  "user": {
	    "created" : true
	  },
	  "created": true 
	}`

	ChangePasswordResponse = `{}`

	DeleteUserResponse = `{
	  "found" : true
	}`
)

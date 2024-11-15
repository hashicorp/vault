package gocb

type userManagerProvider interface {
	GetAllUsers(opts *GetAllUsersOptions) ([]UserAndMetadata, error)
	GetUser(name string, opts *GetUserOptions) (*UserAndMetadata, error)
	UpsertUser(user User, opts *UpsertUserOptions) error
	DropUser(name string, opts *DropUserOptions) error
	GetRoles(opts *GetRolesOptions) ([]RoleAndDescription, error)
	GetGroup(groupName string, opts *GetGroupOptions) (*Group, error)
	GetAllGroups(opts *GetAllGroupsOptions) ([]Group, error)
	UpsertGroup(group Group, opts *UpsertGroupOptions) error
	DropGroup(groupName string, opts *DropGroupOptions) error
	ChangePassword(newPassword string, opts *ChangePasswordOptions) error
}

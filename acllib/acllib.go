package acllib

import (
	"database/sql"
)

type ACL struct {
	db *sql.DB
}

func Open(db *sql.DB) ACL {
	var acl ACL
	acl.db = db

	return acl
}

func (acl *ACL) Close() {
	acl.db.Close()
}

func (acl *ACL) Setup() error {
	return acl.createTables()
}

func (acl *ACL) createTables() error {
	sqlStmt := `
	create table users (username text not null primary key, password text, isAdmin bool);
	create table objects (uuid text not null primary key, owner text, allowedViewers text, isPublic bool);
	`
	_, err := acl.db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = acl.db.Exec("insert into users(username, password, isAdmin) values('root','root',true)")
	return err
}

// CreateUser adds a new user to the database.
// It takes as input the identity (username) of the user currently logged in,
// a new username, password, and boolean that indicates if the new user is an admin.
// Only admin users can create a new user. Otherwise, CreateUser should return an error.
func (acl *ACL) CreateUser(identity string, newUserName string, newUserPass string, newUserIsAdmin bool) error {
	return nil
}

// Login checks the database to see if the supplied password matches the database entry for a given user.
// It takes as input the username and password to be checked.
// If the correct password is provided, Login should return true.
// If the password is incorrect, it should return false.
// If the username is not found in the database, it should return an error.
func (acl *ACL) Login(username string, password string) (error, bool) {
	return nil, false
}

// ChangePassword changes the password stored in the database for a user.
// It takes as input the identity (username) of the user currently logged in, 
// the username of the user whose password is to be changed, and the new password.
// A user is allowed to change their own password, and an admin can change any password. 
// Otherwise, ChangePassword should return an error.
func (acl *ACL) ChangePassword(identity string, username string, newPass string) error {
	return nil
}

// CreateObjectPrivate creates a new row in the objects database for a private object.
// It takes as input the identity (username) of the user currently logged in, 
// the id of the object, and a list of users who are allowed to view the object.
// Any user is allowed to create a private object.
func (acl *ACL) CreateObjectPrivate(identity string, objUuid string, allowedViewers []string) error {
	return nil

}

// CreateObjectPublic creates a new row in the objects database for a public object.
// It takes as input the identity (username) of the user currently logged in and the id of the object.
// Any user is allowed to create a public object.
func (acl *ACL) CreateObjectPublic(identity string, objUuid string) error {
	return nil
}

// AllowUserOnObj adds a user to the list of users allowed to view a particular object.
// It takes as input the identity (username) of the user currently logged in, 
// the id of the object, and the user to be added.
// If the user is already in the list of allowed users, there is no side effect.
// If the object is public, the function should return an error.
// Only the object owner or an admin can add a user to an object's list. 
// Otherwise, the function should return an error.
func (acl *ACL) AllowUserOnObj(identity string, objUuid string, user string) error {
	return nil
}

// DisallowUserOnObj removes a user from the list of users allowed to view an object.
// It takes as input the identity (username) of the user currently logged in,
// the id of the object, and the user to be removed.
// If the user is not in the list, there is no side effect.
// If the object is public, the function should return an error.
// Only the object owner or an admin can add a user to an object's list.
// Otherwise, the function should return an error.
// An owner cannot be removed from the list; the function should return an error in this case.
func (acl *ACL) DisallowUserOnObj(identity string, objUuid string, user string) error {
	return nil
}

// DeleteObject removes an object from the objects database.
// It takes as input the identity (username) of the user currently logged in
// and the id of the object to be deleted.
// Only the owner or an admin can delete an object.
// Otherwise, the function should return an error.
// If the object does not exist, the function should return an error.
func (acl *ACL) DeleteObject(identity string, objUuid string) error {
	return nil
}

// Check checks if a user can view an object.
// It takes as input the identity (username) of the user to check and the id of the object.
// If the object does not exist, the function should return an error.
// An admin can view any object. 
func (acl *ACL) Check(identity string, objUuid string) (error, bool) {
	return nil, false
}

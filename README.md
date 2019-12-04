# COS316, Assignment 6: Access Control Lists

## Due: December 13 at 17:00

## Background
An access control list is a list of permissions attached to an object. Oftentimes, websites need to maintain an access control list for their users that keeps track of what objects a user is allowed to view or edit. 

One example of an object could be a photo posted to Facebook. Each photo has an associated list of allowed users, who can see the photo. All other users not on the list should not be able to see it. Google Docs is another example of an access control list: a user can create a document and decide which users can view and/or edit the document. 

For this assignment, you will be implementing an access control list for a website where users can create objects and allow certain users to view them. The library you implement is responsible for checking if a user is allowed to visit the object.


## Objectives
The objective of this assignment is to implement an access control list library for a website. 
It should support the following APIs:

```go
// CreateUser adds a new user to the database. 
// It takes as input the identity (username) of the user currently logged in,
// a new username, password, and boolean that indicates if the new user is an admin.
// Only admin users can create a new user. Otherwise, CreateUser should return an error.
func (acl *ACL) CreateUser(identity string, newUserName string, newUserPass string, newUserIsAdmin bool) error

// Login checks the database to see if the supplied password matches the database entry for a given user.
// It takes as input the username and password to be checked.
// If the correct password is provided, Login should return true.
// If the password is incorrect, it should return false.
// If the username is not found in the database, it should return an error.
func (acl *ACL) Login(username string, password string) (error, bool)

// ChangePassword changes the password stored in the database for a user.
// It takes as input the identity (username) of the user currently logged in,
// the username of the user whose password is to be changed, and the new password.
// A user is allowed to change their own password, and an admin can change any password.
// Otherwise, ChangePassword should return an error. 
func (acl *ACL) ChangePassword(identity string, username string, newPass string) error

// CreateObjectPrivate creates a new row in the objects database for a private object.
// It takes as input the identity (username) of the user currently logged in,
// the id of the object, and a list of users who are allowed to view the object.
// Any user is allowed to create a private object.
func (acl *ACL) CreateObjectPrivate(identity string, objUuid string, allowedViewers []string) error

// CreateObjectPublic creates a new row in the objects database for a public object.
// It takes as input the identity (username) of the user currently logged in and the id of the object.
// Any user is allowed to create a public object.
func (acl *ACL) CreateObjectPublic(identity string, objUuid string) error

// AllowUserOnObj adds a user to the list of users allowed to view a particular object.
// It takes as input the identity (username) of the user currently logged in,
// the id of the object, and the user to be added.
// If the user is already in the list of allowed users, there is no side effect.
// If the object is public, the function should return an error.
// Only the object owner or an admin can add a user to an object's list.
// Otherwise, the function should return an error.
func (acl *ACL) AllowUserOnObj(identity string, objUuid string, user string) error

// DisallowUserOnObj removes a user from the list of users allowed to view an object.
// It takes as input the identity (username) of the user currently logged in,
// the id of the object, and the user to be removed.
// If the user is not in the list, there is no side effect.
// If the object is public, the function should return an error.
// Only the object owner or an admin can add a user to an object's list.
// Otherwise, the function should return an error.
// An owner cannot be removed from the list; the function should return an error in this case.
func (acl *ACL) DisallowUserOnObj(identity string, objUuid string, user string) error

// DeleteObject removes an object from the objects database.
// It takes as input the identity (username) of the user currently logged in
// and the id of the object to be deleted.
// Only the owner or an admin can delete an object.
// Otherwise, the function should return an error.
// If the object does not exist, the function should return an error.
func (acl *ACL) DeleteObject(identity string, objUuid string) error

// Check checks if a user can view an object.
// It takes as input the identity (username) of the user to check and the id of the object.
// If the object does not exist, the function should return an error.
// An admin can view any object. 
func (acl *ACL) Check(identity string, objUuid string) (error, bool)

```

### Users
There are two types of users: normal and admin. The admin has certain privileges (as described in the API). Each user has a unique username and a password. The code provided for you ensures that the user table is set up with one admin user (username: "root"). 

### Objects
Each object is represented by a unique UUID, which is a string.

### The Database
Each ACL object uses a sqlite database as backend.

There are two tables used in the database: the users table and the objects table. The code to create these tables is provided for you; you do not need to create them.

The users table has three columns: username (string, primary key), password (string), and isAdmin (bool).

The objects table has three columns:  uuid (string, primary key), isPublic (bool), and allowedUsers (string, comma separated list of usernames).

## Resources
As part of this assignment, you will need to write code that composes SQL
queries. We recommend you consult the [SQL precept slides](https://docs.google.com/presentation/d/18cojpvYedtSQEwi3Nq9jd7cAWSRjBrb_iB_jZGVmpVk/edit) or [SQLite documentation](https://www.sqlite.org/index.html) for a refresher on how you might accomplish this.

[Go sql package][go_sql]: SQL queries in Go

[go_sql]: https://golang.org/pkg/database/sql/

## Submission & Grading
Your assignment will be automatically submitted every time you push your changes
to your GitHub Classroom repo. Within a couple minutes of your submission, the
autograder will make a comment on your commit listing the output of our testing
suite when run against your code. **Note that you will be graded only on your
changes to the `acllib` package**, and not on your changes to any other files,
though you may modify any files you wish.

You may submit and receive feedback in this way as many times as you like,
whenever you like, but a substantial lateness penalty will be applied to
submissions past the deadline.




# COS316, Assignment 6/Final Project: Secure Object Relational Mapper (Secure DORM)

# Extensions for Final Project

For the final project, we remove a few of the original expectations from Assignment 6.
Please select some portion of the project to deviate from. For example, instead of creating
capabilities, explore using join tables for permissions, or examine whether it's possible to
create more fine-grained capabilities, such as update/create/delete permissions.
You are free to deviate from the original specification however you would like and we expect
some modification to the original design, as well as a clear explanation of why you modified the design.

Since the design may deviate from the original, the graderbot tests will not work, so its
expected that you create your own tests and demonstrate that your implementation is correct.
The existing tests can be used as a starting point for your own tests.

# Secure Dopey Object Relational Mapper (Secure DORM)

This project asks you to extend the Dopey Object Relational Mapper (DORM) that
you built in assignment 4, with a focus on security. As a reminder, the ORM you built translates language-level
objects---Go structs---to and from a "relational mapping"---a SQLite database. In this
assignment you will use capabilities to allow applications to enforce a security policy
about who can read and write objects in the database.

For example, assume that we have an application that keeps track of users and their posts.
In the application code, we would have a Go struct modeling each user:

```go
type User struct {
    ID       int64 `dorm:"primary_key"`
    Username string
}
```

This would translate into the database schema below:

```sql
create table user (
    id integer primary key autoincrement,
    username text
)
```

where the `id` is a primary key and `username` is some user-specified string identifier.

Your DORM implementation from assignment 4 allowed a programmer to do something like the following:

```go
// Create two new users
user1 := &User{
    Username: "alevy",
}

user2 := &User{
    Username: "wlloyd",
}

// Insert the users into the database
dorm.Create(user1)
dorm.Create(user2)
...

// Get all users
allUsers := []User{}
dorm.Find(&allUsers)

for _, user := range allUsers {
    fmt.Printf("Found user %s\n", user.Username)
}
```

The problem is that by using the DORM interface, code executing on behalf of any user will
have access to all of the user data in the user table when calling `Find()`. For example,
if code executing in a web server in response to a request by user `alevy` calls `Find()`,
it will return user structs for `alevy` and `wlloyd`. This creates creates a potential
vulnerabilty where user `alevy` may be able to extract `wlloyd`'s sensitive information
(e.g., encrypted password) from the application. In this assignment, you'll use capabilities to
design a set of policies that protect against this, by specifying what actions a user can
perform.

For example, assuming the `wlloyd` user is already created in the database, we can modify the
code above to use capabilities to enforce this security policy as follows:

```go
dorm := NewSecureDB(...)
...

// Create user and set root capability
// allowing them to read and write their
// own user object.
user := &User{
    Username: "alevy",
}
dorm.Create(user)
...

// Get user's capability by username
cap := dorm.GetCapability("alevy")

...

// Get users from's Alevy's capability
users := []User{}
dorm.Find(cap, &users)

for _, user := range users {
    fmt.Printf("Found user %s\n", user.Username)
}
```

In the code above, we expect that `allUsers` only contains `alevy`'s user struct,
even though both `alevy` and `wlloyd` exist the database.

## Expressing Security Policies

How you choose to represent security policies is up to you.
As an example, you can extend Go's struct tags:
consider the following user struct, which is the same as the one above:

```go
type User struct {
    ID       int64 `dorm:"primary_key" cap:""`
    Username string
}
```

The `cap` tag will be used by your capabilities (more details on them below) to
control how read and write permissions are evaluated for each object type. Here,
the empty value (as denoted by the `""`) indicates that a user's ID field
is used to determine whether a capability allows the caller to read or
write the object. Consider the following example:

```go
user := &User{
    ID:       1,
    Username: "alevy",
}

cap := dorm.GetCapability(user)

cap.CanRead(user) // Should return true

user.ID = 7
cap.CanRead(user) // Should return false
```

Because the user object's ID field no longer matches `alevy`'s, it returns false.

Unfortunately, this simple case does not map well to all security policies that
an application may want to express. For instance, consider a user's posts. A
microblogging application may want a user to be able to read all posts from the
users that they follow. It is not clear how to express this policy in terms
of post IDs.

Fortunately, we can use the flexibility of Go's struct tags to extend the idea
above and express this more complicated policy. For example, consider the
following post struct:

```go
type Post struct {
	ID     int64 `dorm:"primary_key"`
	UserID int64 `cap:"read=User.ID"`
	Text   string
}
```

Here, the `cap` tag has an associated value `read=User.ID`. Since
the `cap` tag is on the `UserID` field, this tag indicates that
a post should be readable if the capability permits the caller to read the
user ID that is stored in the post's `UserID` field. Again, to make this
concrete, consider the following example code:

```go
user := &User{
    ID:       1,
    Username: "alevy",
}

post := &Post{
    ID:       100,
    UserID:   1,
    Text:     "Hello world!",
}

cap := dorm.GetCapability("alevy")

cap.CanRead(post) // Should return true

post.UserID = 7
cap.CanRead(post) // Should return false
```

Note that the first call to `CanRead()` is true even though user `alevy`'s
root capability (more info on root capabilities below), only included permissions
to read and write `alevy`'s **user** struct. In this way, the application allows that
once we have the capability to read a user's struct, we can also read all of their
posts.

More generally, the tag value can include one or both of `read=<struct name>.<field name>`
and `write=<struct name>.<field name>`. If it includes both, then the two will be separated
by a semicolon, for example, `cap:"read=User.ID;write=User.ID"`.

## Capabilities API

As an example of a possible implementable API,
we provide a set of template interfaces for creating and manipulating capabilities.

The `secure_dorm` package (in `cap_manager.go`) includes two interfaces,
`Capability` and `CapabilityManager`, which are defined as follows:

```go
/*
 * Capability is a structure that encodes information about which
 * objects a user can read or write.
 */
type Capability struct {
}

/*
 * Given a capability and an object, calling cap.CanRead(object) returns
 * true if the capability permits the user to read the object. CanRead
 * expects that its argument is a pointer to a struct.
 *
 * To be explicit, `object` will have type: *MyStruct,
 * where MyStruct is any arbitrary struct subject to the restrictions
 * discussed later in this document.
 *
 * Example usage to test if the caller can read a post object:
 *    type Post struct = { ... }
 *    cap := ...
 *    post := &Post{}
 *    ok := cap.CanRead(post)
 */
func (cap *Capability) CanRead(object interface{}) bool

/*
 * Given a capability and an object, calling cap.CanWrite(object) returns
 * true if the capability permits the user to write to the object.
 *
 * As mentioned in the description of CanRead, the argument `object`
 * will be a pointer to a model.
 */
func (cap *Capability) CanWrite(object interface{}) bool

/*
 * The capability manager allows users to create and modify capabilities.
 */
type CapabilityManager struct {
}

/*
 * Creates a new instance of a capability manager.
 */
func NewCapabilityManager() *CapabilityManager

/*
 * A root capability sets up a user's initial permissions. Given a unique username and
 * two slices of pointers to objects, cm.SetRootCapability(username, readSet, writeSet)
 * associates a root capability with the username. The root capability is expected to allow
 * reading and writing all objects in readSet and writeSet, respectively. For instance,
 * a newly created user's root capability might just include the ability to read and write
 * their own user object. Thus, after creating the new object, the user's root capability
 * would be set by `cm.SetRootCapability(user.Username, []interface{user}, []interface{user}).`
 */
func (cm *CapabilityManager) SetRootCapability(username string, readSet []interface{}, writeSet []interface{})

/*
 * Given a unique username, cm.GetRootCapability(username) returns the user's
 * root capability (or nil if one has not yet been set). A root capability is the user's initial
 * set of permissions. For instance, a newly created user's root capability might
 * just include the ability to read and write their own user object.
 */
func (cm *CapabilityManager) GetRootCapability(username string) *Capability

/*
 * Given a capability and an object, cm.AddReadCapability(cap, object) returns a new capability
 * that includes all capabilities of cap plus the ability to read object. That is, if newCap is
 * the new capability, then calling newCap.CanRead(object) should return true. Note, however, that
 * the original capability should not be modified, so calling cap.CanRead(object) should still
 * return false. Similarly, root capabilities should not change.
 */
func (cm *CapabilityManager) AddReadCapability(cap *Capability, object interface{}) *Capability

/*
 * Given a capability and an object, cm.AddWriteapCability(cap, object) returns a new capability
 * that includes all capabilities of cap plus the ability to write object. Like mentioned above
 * for `AddReadCapability()`, the original capability and all root capabilities should not be modified.
 */
func (cm *CapabilityManager) AddWriteCapability(cap *Capability, object interface{}) *Capability

/*
 * Given a capability and an object, cm.RemoveReadCapability(cap, object) returns a new capability
 * that includes all capabilities of cap minus the ability to read object. Like mentioned above
 * for `AddReadCapability()`, the original capability and all root capabilities should not be modified.
 */
func (cm *CapabilityManager) RemoveReadCapability(cap *Capability, object interface{}) *Capability

/*
 * Given a capability and an object, cm.RemoveWriteCapability(cap, object) returns a new capability
 * that includes all capabilities of cap minus the ability to write object. Like mentioned above
 * for `AddReadCapability()`, the original capability and all root capabilities should not be modified.
 */
func (cm *CapabilityManager) RemoveWriteCapability(cap *Capability, object interface{}) *Capability
```

For the `Capability`, you will be required to implement the following functions:
`CanRead()` and `CanWrite()`.

For the `CapabiltyManager`, you will be required to implement the following functions:
`NewCapabilityManager()`, `GetRootCapability()`, `SetRootCapability()`, `AddReadCapability()`,
`AddWriteCapability()`, `RemoveReadCapability()`, and `RemoveWriteCapability()`.

## Secure DORM API

Once you've implemented the code to create and manipulate capabilities, you can then use it to
implement your secure DORM interface. The `secure_dorm` package exposes the following API:

```go
// SecureDB handle
type SecureDB struct {
	inner DB
	cm    *CapabilityManager
}

// NewSecureDB returns a new SecureDB. It wraps a db object that
// implements the DB interface (e.g., the implementation in dorm.go).
// It also takes a capability manager to enforce secure access to
// the SQL database.
func NewSecureDB(db DB, cm *CapabilityManager) *SecureDB

// Close closes db's database connection.
func (db *SecureDB) Close() error

// Find returns all rows in a given table that the capability `cap`
// allows the caller to read. It stores all matching rows in the
// slice provided as an argument.
//
// The argument `result` will be a pointer to an empty slice of models.
// To be explicit, it will have type: *[]MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.
// You may assume the slice referenced by `result` is empty.
//
// Example usage to find all UserComment entries in the database:
//    type UserComment struct = { ... }
//    cap := ...
//    result := []UserComment{}
//    db.Find(cap, &result)
func (db *SecureDB) Find(cap *Capability, result interface{})

// First queries a database for the first row in a table that the
// capability `cap` allows the caller to read. Note that this may not
// be the first row in the table. It stores the matching row in the
// struct provided as an argument. If no such entry exists, First
// returns false; else it returns true.
//
// The argument `result` will be a pointer to a model.
// To be explicit, it will have type: *MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.
//
// Example usage to find the first UserComment entry in the database:
//    type UserComment struct = { ... }
//    cap := ...
//    result := &UserComment{}
//    ok := db.First(cap, result)
func (db *SecureDB) First(cap *Capability, result interface{})

// Create adds the specified model to the appropriate database table
// if the capability `cap` allows the caller to write the object.
// The table for the model *must* already exist, and Create() should
// panic if it does not.
//
// Create returns true if the model/object was successfully created;
// otherwise it returns false.
func (db *SecureDB) Create(cap *Capability, model interface{}) bool
```

You will be required to implement the following functions:
`Find()`, `First()`, and `Create()`.

The functions `NewDB` and `Close` are provided for you.
There is no need to modify these functions for your implementation,
although you are welcome to if it will help your implementation.

You'll notice that `SecureDB` wraps a `DB` object. This `DB` object
exposes the insecure DORM interface, and an implementation of the
interface is provided to you in `dorm.go`. This way, you should not
have to re-implement most of the functionality from assignment 4.

In particular, the `DB` interface is defined as follows:

```go
/*
 * The DB interface implemented by DORM.
 */
type DB interface {
	// DORM Close
	Close() error

	// DORM ToUnderscoreCase
	ToUnderscoreCase(n string) string

	// (Insecure) DORM Find
	Find(result interface{})

	// (Insecure) DORM First
	First(result interface{}) bool

	// (Insecure) DORM Create
	Create(model interface{})
}
```

### Reflection Tips

To further remove some of the pains of using Go's `reflect`
package, we have provided you with two utility functions, which
can be found in `utils.go`. You are not required to use either of
these functions but may find them useful when implementing `SecureDB`.
Their definitions and descriptions are below:

```go
/*
 * Given a pointer to a slice of structs, returns a pointer to a new slice
 * of the same type.
 */
func NewSliceFromSlice(result interface{}) interface{}

/*
 * Given a pointer to a struct, returns a pointer to a new slice of
 * structs of the same type.
 */
func NewSliceFromStruct(result interface{}) interface{}
```

Separately, when interacting with struct tags through the reflection API,
we recommend that you carefully read about the different between
`StructTag.Get()` and `StructTag.Lookup()`.

Finally, to simplify the presentation of some of the code examples above,
we used the syntax `[]interface{object}`. Unforunately, this is not actually
valid Go code and thus will not compile. The reason why is unimportant,
but you can achieve the desired result using two lines of code. For instance,

```go
var list []interface{}
list = append(list, object)
```

### Restrictions on Structs

Similar to the DORM assignment, we will make several simplifying
assumptions about the sorts of structs that make valid DB models.

In particular:
* You may assume that the fields of structs will all be primitive types
  (e.g. `string`, `int`, `int64`, `bool`, ...). This means you need not
  worry about `map`, `slice`, or `struct` types being included as
  fields.
  Primitive types are handled natively by the `sql` library we are
  using, so you should *not* have to do any special work to support
  these types. In contrast, `map`, `slice`, or nested `struct` types
  add complexity to the ORM implementation, so you are not responsible
  for supporting them.
* You may assume there will be no nested structs in your models.
* You may assume that all field names will be in a valid `camelCase` or
  `CamelCase` format.
* If any of the fields of the `model` provided to `Create()` are
  tagged with `dorm:"primary_key"`, you should assume that the type
  of that field will be `int64`.
* You can assume that only one field of any model argument will have a
  `cap` tag.


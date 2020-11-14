package secure_dorm

// SecureDB handle
type SecureDB struct {
	inner DB
	cm    *CapabilityManager
}

// NewSecureDB returns a new SecureDB. It wraps a db object that
// implements the DB interface (e.g., the implementation in dorm.go).
// It also takes a capability manager to enforce secure access to
// the SQL database.
// This function is provided for you. You DO NOT need to modify it.
func NewSecureDB(db DB, cm *CapabilityManager) *SecureDB {
	return &SecureDB{inner: db, cm: cm}
}

// Close closes db's database connection.
// This function is provided for you. You DO NOT need to modify it.
func (db *SecureDB) Close() error {
	return db.inner.Close()
}

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
func (db *SecureDB) Find(cap *Capability, result interface{}) {
}

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
func (db *SecureDB) First(cap *Capability, result interface{}) bool {
	return false
}

// Create adds the specified model to the appropriate database table
// if the capability `cap` allows the caller to write the object.
// The table for the model *must* already exist, and Create() should
// panic if it does not.
//
// Create returns true if the model/object was successfully created;
// otherwise it returns false.
func (db *SecureDB) Create(cap *Capability, model interface{}) bool {
	return false
}

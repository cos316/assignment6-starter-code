package secure_dorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// Assume public field names (or struct names) are CamelCase
//
// toUnderscoreCase converts a string of the form NameOfField (CamelCase) to a string of the
// form name_of_field, lower-case-words separated by underscores, and returns that string.
//
// A word in CamelCase is any number of capital letters _not_ followed by a
// lower-case letter or number, or a capitalized word (one capital letter, followed by
// any number of lower case letters).
//
// Numbers never cause a word split.
//
// A capital letter is a new word, unless it does not preceded a non-upper-case-letter.
func toUnderscoreCase(n string) string {
	words := []string{}
	l := 0
	for s := n; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], func(c rune) bool { return !unicode.IsUpper(c) })
		if l < 0 {
			l = len(s)
		} else if l == 0 {
			l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
			if l <= 0 {
				l = len(s)
			}
		}
		words = append(words, strings.ToLower(s[:l]))
	}
	return strings.Join(words, "_")
}

func tableName(t reflect.Type) string {
	return toUnderscoreCase(t.Name())
}

// Return the table name in underscore-case based on the struct's CamelCase name (rules above)
func TableName(result interface{}) string {
	return tableName(reflect.TypeOf(result).Elem())
}

type column struct {
	name  string
	index []int
}

func columnNames(cls []column) []string {
	names := []string{}
	for _, c := range cls {
		names = append(names, c.name)
	}
	return names
}

func columnValues(cls []column) [][]int {
	vals := [][]int{}
	for _, c := range cls {
		vals = append(vals, c.index)
	}
	return vals
}

func _columns(t reflect.Value, prefix []int) []column {
	numFields := t.NumField()
	cls := []column{}
	for i := 0; i < numFields; i++ {
		ft := t.Type().Field(i)
		f := t.Field(i)
		if ft.Type.Kind() == reflect.Struct {
			cls = append(cls, _columns(f, append(prefix, i))...)
		} else if f.CanInterface() {
			cls = append(cls, column{
				name:  toUnderscoreCase(ft.Name),
				index: append(prefix, i),
			})
		}
	}
	return cls
}

func columns(t reflect.Value) []column {
	return _columns(t, []int{})
}

// Return a list of column names based on public field names
//
// Public fields only, uses CamelCase to underscore_case rules above
func ColumnNames(v interface{}) []string {
	return columnNames(columns(reflect.ValueOf(v).Elem()))
}

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

/*
 * DBImpl is the struct that implements the DB interface.
 * This should look very similar to the DB struct from the
 * DORM assignment.
 */
type DBImpl struct {
	inner *sql.DB
}

// Creates a new DB type
func NewDB(conn *sql.DB) *DBImpl {
	return &DBImpl{inner: conn}
}

// Closes the underlying DB connection
func (db *DBImpl) Close() error {
	return db.inner.Close()
}

// Closes the underlying DB connection
func (db *DBImpl) ToUnderscoreCase(n string) string {
	return toUnderscoreCase(n)
}

// Finds all rows in a table
//
// Arguments:
//  - result: a pointer to an empty slice of models (e.g. *[]MyModel).
//
// Result:
//
// Populates the slice with all rows from  the table associated with the model's struct type.
func (db *DBImpl) Find(result interface{}) {
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		panic("Find's argument must be of a slice of model structs")
	}
	if reflect.TypeOf(result).Elem().Kind() != reflect.Slice {
		panic("Find's argument must be of a slice of model structs")
	}
	rslice := reflect.ValueOf(result).Elem()

	if reflect.TypeOf(result).Elem().Elem().Kind() != reflect.Struct {
		panic("Find's argument must be of a slice of model structs")
	}
	modelType := reflect.TypeOf(result).Elem().Elem()
	table := tableName(modelType)
	cols := columns(reflect.New(modelType).Elem())
	rows, err := db.inner.Query("select " + strings.Join(columnNames(cols), ",") + " from " + table)
	if err != nil {
		panic(err)
	}

	colValues := columnValues(cols)
	fields := make([]interface{}, len(colValues))
	for i, idx := range colValues {
		field := reflect.New(modelType.FieldByIndex(idx).Type).Interface()
		fields[i] = field
	}

	for rows.Next() {
		err = rows.Scan(fields...)
		if err != nil {
			panic(err)
		}
		f := reflect.New(modelType)
		for i, field := range fields {
			f.Elem().FieldByIndex(colValues[i]).Set(reflect.ValueOf(field).Elem())
		}
		rslice.Set(reflect.Append(rslice, f.Elem()))
	}
}

// Finds the first row in a table
//
// Arguments:
//  - result: a pointer to a model (e.g. *MyModel).
//
// Result:
//
// Same as Find, but only the first row in the table (using natural table
// order). If there are no rows, return false (and don't have to do anything
// with the argument), otherwise return true.
func (db *DBImpl) First(result interface{}) bool {
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		panic("Find's argument must be of a pointer to a model struct")
	}

	if reflect.TypeOf(result).Elem().Kind() != reflect.Struct {
		panic("Find's argument must be of a pointer to a model struct")
	}

	modelType := reflect.TypeOf(result).Elem()

	table := tableName(modelType)
	cols := columns(reflect.New(modelType).Elem())
	row := db.inner.QueryRow("select " + strings.Join(columnNames(cols), ",") + " from " + table + " limit 1")

	colValues := columnValues(cols)
	fields := make([]interface{}, len(colValues))
	for i, idx := range colValues {
		field := reflect.New(modelType.FieldByIndex(idx).Type).Interface()
		fields[i] = field
	}

	err := row.Scan(fields...)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		panic(err)
	}

	f := reflect.ValueOf(result)
	for i, field := range fields {
		f.Elem().FieldByIndex(colValues[i]).Set(reflect.ValueOf(field).Elem())
	}

	return true
}

// Inserts a new row into the database, table based on struct type.
//
// Argument:
//   - model: the new value to insert--a pointer to a struct (e.g. *MyModel).
//   If the model type has a field annotated with `dorm:"primary_key"`, Create
//   should allow the database to populate that column automatically, and
//   Create should set the field to the database-generated value (hint: LastInsertId()).
//
func (db *DBImpl) Create(model interface{}) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		panic("Find's argument must be of a pointer to a model struct")
	}

	if reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		panic("Find's argument must be of a pointer to a model struct")
	}

	modelType := reflect.TypeOf(model).Elem()

	table := tableName(modelType)
	cols := columns(reflect.New(modelType).Elem())
	value := reflect.ValueOf(model).Elem()

	fields := []interface{}{}
	placeholders := []string{}
	includedColumns := []string{}
	var primaryKey reflect.Value
	isPrimaryKey := func(idx []int) bool {
		field := modelType.FieldByIndex(idx)
		dorm, ok := field.Tag.Lookup("dorm")
		return ok && dorm == "primary_key"
	}
	for _, c := range cols {
		if isPrimaryKey(c.index) {
			primaryKey = value.FieldByIndex(c.index)
		} else {
			field := value.FieldByIndex(c.index).Interface()
			fields = append(fields, field)
			includedColumns = append(includedColumns, c.name)
			placeholders = append(placeholders, "?")
		}
	}

	query := fmt.Sprintf("insert into %v (%v) values (%v)", table, strings.Join(includedColumns, ","), strings.Join(placeholders, ","))
	result, err := db.inner.Exec(query, fields...)

	if err != nil {
		panic(err)
	}

	if primaryKey.IsValid() {
		mid, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		primaryKey.Set(reflect.ValueOf(mid))
	}
}

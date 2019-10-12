# SCIM Filter Transpiler

Transpiles Tokenized SCIM v2 Filter into a SQL command. It utilizes https://github.com/di-wu/scim-filter-parser to do the work of parsing a raw filter query parameter. It uses https://github.com/Masterminds/squirrel to do the SQL building.

## Install

```
go get github.com/articulate/scim-filter-transpiler
```

## Usage

```go
import "github.com/articulate/scim-filter-transpiler"

parser := NewParser(
  // Attribute map that tells us how to map our attribute names.
  // Any missing path will be returned as is.
  map[string]string{
    "id":           "users.id",
    "username":     "users.username",
    "emails.value": "emails.value",
    "emails":       "emails.value",
    "emails.type":  "emails.type",
  },
  // Resource table name
  "users",
  // Include any necessary joins
  []string{"LEFT JOIN emails ON emails.user_id = users.id"},
)

// Use ToSql if you already have a parsed filter.
sql, _ := parser.ToSqlFromString(`userName eq "andy@example.com"`, "users.id")

// We can even use Squirrel to query our DB, or use ToSql to get the raw query and params.
// Builds the following query:
// SELECT users.id FROM users LEFT JOIN emails ON emails.user_id = users.id WHERE users.username = ?
// With andy@example.com as the only parameter
rows, _ := sql.Limit(10).Offset(10).RunWith(db).Query()
```

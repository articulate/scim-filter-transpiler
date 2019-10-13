# SCIM Filter Transpiler

Transpiles Tokenized SCIM v2 Filter into a SQL command. It utilizes https://github.com/di-wu/scim-filter-parser to do the work of parsing a raw filter query parameter.

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
)

// Use ToSql if you already have a parsed filter.
query, params, _ := parser.ToSqlFromString(`userName eq "andy@example.com"`)
```

package scim

import (
	"reflect"
	"testing"

	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var (
	attrMap = map[string]string{
		"username":  "users.username",
		"emails":    "emails.value",
		"ims.value": "imsValue",
		"ims":       "imsValue",
		"ims.type":  "imsType",
	}
	joins = []string{
		"LEFT JOIN emails ON emails.user_id = users.id",
	}
)

func TestFilterParser(t *testing.T) {
	var tests = []struct {
		attributeMap    map[string]string
		expected        string
		filter          string
		joins           []string
		params          []interface{}
		expectedResults []string
	}{
		{
			filter:          `userName eq "andy@example.com"`,
			expected:        `SELECT distinct on (users.id) users.id FROM users LEFT JOIN emails ON emails.user_id = users.id WHERE users.username = ?`,
			params:          []interface{}{"andy@example.com"},
			joins:           joins,
			attributeMap:    attrMap,
			expectedResults: []string{"0001"},
		},
		{
			filter:          `emails co "example.org" and (emails.type eq "work" and emails.value co "example.org")`,
			expected:        "SELECT distinct on (users.id) users.id FROM users LEFT JOIN emails ON emails.user_id = users.id WHERE (emails.value LIKE ? AND (emails.type = ? AND emails.value LIKE ?))",
			params:          []interface{}{"%example.org%", "work", "%example.org%"},
			joins:           joins,
			attributeMap:    attrMap,
			expectedResults: []string{"0002"},
		},
		{
			filter:          `emails[type eq "work" and value co "@example.com"] or ims[type eq "xmpp" and value co "@foo.com"]`,
			expected:        "SELECT distinct on (users.id) users.id FROM users LEFT JOIN emails ON emails.user_id = users.id WHERE ((emails.type = ? AND emails.value LIKE ?) OR (imsType = ? AND imsValue LIKE ?))",
			params:          []interface{}{"work", "%@example.com%", "xmpp", "%@foo.com%"},
			joins:           joins,
			attributeMap:    attrMap,
			expectedResults: []string{"0001"},
		},
		{
			filter:          `not emails co "example.com"`,
			expected:        `SELECT distinct on (users.id) users.id FROM users LEFT JOIN emails ON emails.user_id = users.id WHERE NOT (emails.value LIKE ?)`,
			params:          []interface{}{"%example.com%"},
			joins:           joins,
			attributeMap:    attrMap,
			expectedResults: []string{"0001", "0002", "0003", "0004"},
		},
	}

	db, err := sql.Open("postgres", "host=db port=5432 dbname=root user=root password=root sslmode=disable")

	if err != nil {
		t.Errorf("failed to connect to test database, error: %v", err)
		return
	}

	for _, test := range tests {
		parser := NewParser(test.attributeMap, "users", test.joins)

		sqlQuery, err := parser.ToSqlFromString(test.filter, "distinct on (users.id) users.id")

		if err != nil {
			t.Errorf("Expected to create a filter parser without an error but received an error %v", err)
		}

		query, params, err := sqlQuery.ToSql()

		if err != nil {
			t.Errorf("failed to parse sql query, error %v", err)
		}

		if query != test.expected {
			t.Errorf("Malformed SQL query, expected:\n%s\ngot:\n%s", test.expected, query)
		}

		if len(params) != len(test.params) || !reflect.DeepEqual(params, test.params) {
			t.Errorf(`Malformed parameters, expected %v, received %v`, test.params, params)
		}

		rows, err := sqlQuery.PlaceholderFormat(sq.Dollar).RunWith(db).Query()

		if err != nil {
			t.Errorf("failed to query test database, error: %v", err)
			return
		}

		var ids []string

		for rows.Next() {
			var id string
			err := rows.Scan(&id)
			ids = append(ids, id)

			if err != nil {
				t.Errorf("could not read test data from database, error: %v", err)
				return
			}
		}

		err = rows.Err()

		if err != nil {
			t.Errorf("could not read test data from database, error: %v", err)
		}

		if !reflect.DeepEqual(ids, test.expectedResults) {
			t.Errorf("expected results did not match the real results, got: %v, expected %v", ids, test.expectedResults)
		}

		rows.Close()
	}
}

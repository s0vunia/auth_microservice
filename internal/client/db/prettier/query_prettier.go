package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// PlaceholderDollar represents the placeholder used in SQL queries for dollar sign syntax
	PlaceholderDollar = "$"
	// PlaceholderQuestion represents the placeholder used in SQL queries for question mark syntax
	PlaceholderQuestion = "?"
)

// Pretty returns a prettier SQL query
func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, fmt.Sprintf("%s%s", placeholder, strconv.Itoa(i+1)), value, -1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}

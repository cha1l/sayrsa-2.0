package repository

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"strconv"
	"strings"
	"testing"
)

func TestGetPublicKeysQuery(t *testing.T) {
	type Case struct {
		name      string
		usernames []string
		result    string
	}

	cases := []Case{
		{
			name:      "classic",
			usernames: []string{"test1", "test2", "test3"},
			result:    fmt.Sprintf(`SELECT public_key FROM %s WHERE username=$1 OR username=$2 OR username=$3`, usersTable),
		},
		{
			name:      "just one",
			usernames: []string{"test"},
			result:    fmt.Sprintf(`SELECT public_key FROM %s WHERE username=$1`, usersTable),
		},
	}

	for _, val := range cases {
		assert.Equal(t, testQuery(val.usernames), val.result)
	}

}

func testQuery(usernames []string) string {
	var searchIndexes []string
	for i, _ := range usernames {
		searchIndexes = append(searchIndexes, "$"+strconv.Itoa(i+1))
	}
	indexes := strings.Join(searchIndexes, " OR username=")

	query := fmt.Sprintf(`SELECT public_key FROM %s WHERE username=%s`, usersTable, indexes)

	return query
}

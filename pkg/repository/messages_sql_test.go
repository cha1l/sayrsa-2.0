package repository

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"strconv"
	"strings"
	"testing"
)

func TestInsertMessageTextQuery(t *testing.T) {
	type Tests struct {
		Name   string
		Answer string
		Text   map[string]string
	}

	cases := []Tests{
		{
			Name:   "1 value",
			Answer: fmt.Sprintf("INSERT INTO %s (conv_message_id, text, for_user) VALUES ($1, $2, $3)", messageTextTable),
			Text: map[string]string{
				"vbnm251": "test",
			},
		},
		{
			Name:   "2 values",
			Answer: fmt.Sprintf("INSERT INTO %s (conv_message_id, text, for_user) VALUES ($1, $2, $3), ($4, $5, $6)", messageTextTable),
			Text: map[string]string{
				"vbnm251": "123",
				"test":    "user",
			},
		},
	}

	for _, val := range cases {
		assert.Equal(t, testMessagesQuery(val.Text), val.Answer)
	}
}

func testMessagesQuery(text map[string]string) string {
	valueList := make([]string, 0)
	cnt := 1

	for range text {
		object := fmt.Sprintf(`($%s, $%s, $%s)`, strconv.Itoa(cnt), strconv.Itoa(cnt+1), strconv.Itoa(cnt+2))
		valueList = append(valueList, object)
		cnt += 3
	}

	values := strings.Join(valueList, ", ")

	textQuery := fmt.Sprintf(`INSERT INTO %s (conv_message_id, text, for_user) VALUES %s`, messageTextTable, values)

	return textQuery
}

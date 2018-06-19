package slackbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEntity(t *testing.T) {
	tests := []struct {
		TestName string

		Entity   string
		ID, Name string
	}{
		{"Empty", "", "", ""},
		{"Username", "@foo", "", "foo"},
		{"Short", "<@1|a>", "1", "a"},
		{"Full", "<@123|foo>", "123", "foo"},
		{"JustID", "<@123>", "123", ""},
		{"Partial", "<@123|foo", "", "<@123|foo"},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			id, name := parseEntity(test.Entity)
			assert.Equal(t, test.ID, id, "id mismatch")
			assert.Equal(t, test.Name, name, "name mismatch")
		})
	}
}

package slackbot

import (
	"strings"
)

// Team contains the ID and domain of a Slack team.
type Team struct {
	ID     string
	Domain string
}

// Enterprise contains the ID and name of a Slack enterprise.
type Enterprise struct {
	ID   string
	Name string
}

// Channel contains the ID and name of a Slack channel.
type Channel struct {
	ID   string
	Name string
}

// ParseChannel parses an extended Channel entity.
func ParseChannel(c string) (channel Channel) {
	channel.ID, channel.Name = parseEntity(c)
	return
}

// User contains the ID and name of a Slack user.
type User struct {
	ID   string
	Name string
}

// ParseUser parses an extended User entity.
func ParseUser(u string) (user User) {
	user.ID, user.Name = parseEntity(u)
	return
}

func parseEntity(e string) (id, name string) {
	if len(e) < 1 {
		return
	}

	if len(e) >= 6 && e[0] == '<' && e[1] == '@' && e[len(e)-1] == '>' {
		e = e[1 : len(e)-1]

		i := strings.IndexByte(e, '|')
		if i > -1 {
			id = e[1:i]
			name = e[i+1:]
		} else {
			id = e[1:]
		}
	} else {
		if e[0] == '@' {
			name = e[1:]
		} else {
			name = e
		}
	}

	return
}

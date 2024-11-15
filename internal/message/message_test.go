package message

import (
	"testing"

	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusChangedMessage(t *testing.T) {
	type args struct {
		domain  string
		changes []diff.Change
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"create",
			args{
				"example.com",
				[]diff.Change{
					{Type: diff.CREATE, To: "a"},
				},
			},
			"The statuses on example.com have changed. Here are the changes:\n```\n + a\n```",
		},
		{
			"update",
			args{
				"example.com",
				[]diff.Change{
					{Type: diff.UPDATE, From: "a", To: "b"},
				},
			},
			"The statuses on example.com have changed. Here are the changes:\n```\n - a\n + b\n```",
		},
		{
			"delete",
			args{
				"example.com",
				[]diff.Change{
					{Type: diff.DELETE, From: "a"},
				},
			},
			"The statuses on example.com have changed. Here are the changes:\n```\n - a\n```",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewStatusChangedMessage(tt.args.domain, tt.args.changes))
		})
	}
}

func TestNewThresholdMessage(t *testing.T) {
	type args struct {
		domain   string
		timeLeft int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"example.com 7d",
			args{"example.com", 7},
			"example.com will expire in 7 days.",
		},
		{
			"example.com 1d",
			args{"example.com", 1},
			"example.com will expire in 1 days.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewThresholdMessage(tt.args.domain, tt.args.timeLeft))
		})
	}
}

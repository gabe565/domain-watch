package message

import (
	"fmt"

	"github.com/r3labs/diff/v3"
)

func NewStatusChangedMessage(domain string, changes []diff.Change) string {
	var added, removed string
	for _, change := range changes {
		switch change.Type {
		case diff.UPDATE:
			removed += fmt.Sprintf("\n - %s", change.From)
			added += fmt.Sprintf("\n + %s", change.To)
		case diff.CREATE:
			added += fmt.Sprintf("\n + %s", change.To)
		case diff.DELETE:
			removed += fmt.Sprintf("\n - %s", change.From)
		}
	}
	message := fmt.Sprintf(
		"The statuses on %s have changed. Here are the changes:\n```%s%s\n```",
		domain,
		removed,
		added,
	)
	return message
}

func NewThresholdMessage(domain string, timeLeft int) string {
	return fmt.Sprintf("%s will expire in %d days.", domain, timeLeft)
}

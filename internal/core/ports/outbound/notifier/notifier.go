package notifier

import "context"

type Notifier interface {
	SendMessage(ctx context.Context, content string) error
}

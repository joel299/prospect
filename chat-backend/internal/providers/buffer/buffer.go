package buffer

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Publisher interface {
	CreatePost(context.Context, Post) (string, error)
}
type Post struct{ Text, ProfileID, ScheduledAt string }
type Client struct {
	BaseURL, AccessToken string
	HTTP                 *http.Client
}

func (c Client) CreatePost(ctx context.Context, p Post) (string, error) {
	if c.BaseURL == "" || c.AccessToken == "" {
		return "", fmt.Errorf("buffer not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, nil)
	if err != nil {
		return "", err
	}
	_ = req
	_ = p
	_ = time.Now()
	return "", fmt.Errorf("buffer adapter endpoint mapping required")
}

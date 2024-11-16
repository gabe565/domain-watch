package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/util"
	"github.com/gotify/server/v2/model"
)

type Gotify struct {
	URL   *url.URL
	token string
}

func (g *Gotify) Setup(ctx context.Context, conf *config.Config) error {
	if conf.GotifyURL == "" {
		return fmt.Errorf("gotify %w: token", util.ErrNotConfigured)
	}

	var err error
	g.URL, err = url.Parse(conf.GotifyURL)
	if err != nil {
		return err
	}

	if g.token = conf.GotifyToken; g.token == "" {
		return fmt.Errorf("gotify %w: chat ID", util.ErrNotConfigured)
	}
	return g.Login(ctx)
}

func (g *Gotify) Login(ctx context.Context) error {
	u, err := g.URL.Parse("version")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", util.ErrUnexpectedStatus, resp.Status)
	}

	var version model.VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return err
	}

	slog.Info("Connected to Gotify", "version", version.Version)
	return nil
}

func (g *Gotify) Send(ctx context.Context, text string) error {
	if g.URL == nil {
		return nil
	}

	priority := 5
	payload := model.MessageExternal{
		Message:  text,
		Priority: &priority,
		Extras: map[string]any{
			"client::display": map[string]any{
				"contentType": "text/markdown",
			},
		},
	}

	u, err := g.URL.Parse("message")
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", g.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", util.ErrUnexpectedStatus, resp.Status)
	}

	return nil
}

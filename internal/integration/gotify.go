package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gabe565/domain-watch/internal/util"
	"github.com/gotify/server/v2/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Gotify struct {
	URL   *url.URL
	token string
}

func (g *Gotify) Flags(cmd *cobra.Command) error {
	cmd.Flags().String("gotify-url", "", "Gotify URL (include https:// and port if non-standard)")
	if err := viper.BindPFlag("gotify.url", cmd.Flags().Lookup("gotify-url")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("gotify-url", util.NoFileComp); err != nil {
		panic(err)
	}

	cmd.Flags().String("gotify-token", "", "Gotify app token")
	if err := viper.BindPFlag("gotify.token", cmd.Flags().Lookup("gotify-token")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("gotify-token", util.NoFileComp); err != nil {
		panic(err)
	}

	cmd.MarkFlagsRequiredTogether("gotify-url", "gotify-token")

	return nil
}

func (g *Gotify) Setup() (err error) {
	host := viper.GetString("gotify.url")
	if host == "" {
		return fmt.Errorf("gotify %w: token", util.ErrNotConfigured)
	}

	g.URL, err = url.Parse(host)
	if err != nil {
		return err
	}

	if g.token = viper.GetString("gotify.token"); g.token == "" {
		return fmt.Errorf("gotify %w: chat ID", util.ErrNotConfigured)
	}

	return g.Login()
}

func (g *Gotify) Login() error {
	u, err := g.URL.Parse("version")
	if err != nil {
		return err
	}

	resp, err := http.Get(u.String())
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

	log.WithFields(log.Fields{
		"version": version.Version,
	}).Info("connected to Gotify")

	return nil
}

func (g *Gotify) Send(text string) error {
	if g.URL == nil {
		return nil
	}

	payload := model.MessageExternal{
		Message:  text,
		Priority: 5,
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

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
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

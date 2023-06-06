package integration

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gabe565/domain-watch/internal/util"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Gotify struct {
	URL    *url.URL
	token  string
	client *client.GotifyREST
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

	return nil
}

func (g *Gotify) Setup() error {
	host := viper.GetString("gotify.url")
	if host == "" {
		return fmt.Errorf("gotify %w: token", util.ErrNotConfigured)
	}

	serverUrl, err := url.Parse(host)
	if err != nil {
		return err
	}

	if g.token = viper.GetString("gotify.token"); g.token == "" {
		return fmt.Errorf("gotify %w: chat ID", util.ErrNotConfigured)
	}

	return g.Login(serverUrl)
}

func (g *Gotify) Login(serverUrl *url.URL) error {
	g.client = gotify.NewClient(serverUrl, http.DefaultClient)

	version, err := g.client.Version.GetVersion(nil)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"version": version.Payload.Version,
	}).Info("connected to Gotify")

	return nil
}

func (g *Gotify) Send(text string) error {
	if g.client == nil {
		return nil
	}

	payload := message.NewCreateMessageParams()
	payload.Body = &models.MessageExternal{
		Message: text,
		Extras: map[string]any{
			"client::display": map[string]any{
				"contentType": "text/markdown",
			},
		},
	}

	_, err := g.client.Message.CreateMessage(payload, auth.TokenAuth(g.token))
	return err
}

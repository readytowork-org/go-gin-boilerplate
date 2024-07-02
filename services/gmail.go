package services

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"boilerplate-api/internal/config"
	"boilerplate-api/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"google.golang.org/api/gmail/v1"
)

type EmailParams struct {
	To              string
	SubjectData     string
	SubjectTemplate string
	BodyData        interface{}
	BodyTemplate    string
	Lang            string
}

type GmailService struct {
	*gmail.Service
	logger config.Logger
}

func NewGmailService(logger config.Logger, env config.Env) GmailService {
	ctx := context.Background()

	oauthConfig := oauth2.Config{
		ClientID:     env.MailClientID,
		ClientSecret: env.MailClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  env.HOST, // e.g: "http://localhost" or deployed API url
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}
	token := oauth2.Token{
		AccessToken:  env.MailAccesstoken,
		RefreshToken: env.MailRefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}
	var tokenSource = oauthConfig.TokenSource(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Fatal("failed to receive gmail client", err.Error())
	}

	return GmailService{
		Service: srv,
		logger:  logger,
	}
}

func (g GmailService) SendEmail(params EmailParams) (bool, error) {
	to := params.To
	emailBody, err := utils.ParseTemplate(params.BodyTemplate, params.BodyData)
	if err != nil {
		return false, errors.New("unable to parse email body template")
	}
	var msgString string
	emailTo := "To: " + to + "\r\n"
	msgString = emailTo
	subject := "Subject: " + params.SubjectData + "\n"
	msgString = msgString + subject
	msgString = msgString + "\n" + emailBody
	var msg []byte

	if params.Lang != "en" {
		msgStringJP, _ := utils.ToISO2022JP(msgString)
		msg = msgStringJP
	} else {
		msg = []byte(msgString)
	}
	message := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(msg),
	}
	_, err = g.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

package types

import (
	"context"
	"html"

	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"

	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

func generateEmail(ctx context.Context, user *view_model.NotifyUser, subject, content string, smtpConfig func(ctx context.Context) (*smtp.EmailConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), lastEmail bool) error {
	content = html.UnescapeString(content)
	message := &messages.Email{
		Recipients: []string{user.VerifiedEmail},
		Subject:    subject,
		Content:    content,
	}
	if lastEmail {
		message.Recipients = []string{user.LastEmail}
	}

	channelChain, err := senders.EmailChannels(ctx, smtpConfig, getFileSystemProvider, getLogProvider)
	if err != nil {
		return err
	}

	if channelChain.Len() == 0 {
		return caos_errors.ThrowPreconditionFailed(nil, "MAIL-83nof", "Errors.Notification.Channels.NotPresent")
	}
	return channelChain.HandleMessage(message)
}

func mapNotifyUserToArgs(user *view_model.NotifyUser) map[string]interface{} {
	return map[string]interface{}{
		"UserName":           user.UserName,
		"FirstName":          user.FirstName,
		"LastName":           user.LastName,
		"NickName":           user.NickName,
		"DisplayName":        user.DisplayName,
		"LastEmail":          user.LastEmail,
		"VerifiedEmail":      user.VerifiedEmail,
		"LastPhone":          user.LastPhone,
		"VerifiedPhone":      user.VerifiedPhone,
		"PreferredLoginName": user.PreferredLoginName,
		"LoginNames":         user.LoginNames,
		"ChangeDate":         user.ChangeDate,
		"CreationDate":       user.CreationDate,
	}
}

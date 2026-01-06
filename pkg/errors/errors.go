package apperr

import (
	"errors"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrTeamAlreadyExists     = errors.New("team already exists")
	ErrTeamNotFound          = errors.New("team not found")
	ErrPlayerNotFound        = errors.New("player not found")
	ErrTransferNotFound      = errors.New("transfer not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTooManyAttempts       = errors.New("too many login attempts")
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrPlayerAlreadyListed   = errors.New("player already listed for transfer")
	ErrCannotBuyOwnPlayer    = errors.New("cannot buy your own player")
	ErrTransferNotActive     = errors.New("transfer is not active")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidInput          = errors.New("invalid input")
	ErrInternal              = errors.New("internal server error")
)

var errorToMessageID = map[error]string{
	ErrUserAlreadyExists:   "errors.user_already_exists",
	ErrUserNotFound:        "errors.user_not_found",
	ErrTeamAlreadyExists:   "errors.team_already_exists",
	ErrTeamNotFound:        "errors.team_not_found",
	ErrPlayerNotFound:      "errors.player_not_found",
	ErrTransferNotFound:    "errors.transfer_not_found",
	ErrInvalidCredentials:  "errors.invalid_credentials",
	ErrTooManyAttempts:     "errors.too_many_attempts",
	ErrInsufficientFunds:   "errors.insufficient_funds",
	ErrPlayerAlreadyListed: "errors.player_already_listed",
	ErrCannotBuyOwnPlayer:  "errors.cannot_buy_own_player",
	ErrTransferNotActive:   "errors.transfer_not_active",
	ErrUnauthorized:        "errors.unauthorized",
	ErrForbidden:           "errors.forbidden",
	ErrInvalidInput:        "errors.invalid_input",
	ErrInternal:            "errors.internal_error",
}

func LocalizeError(err error, localizer *i18n.Localizer) string {
	for appErr, messageID := range errorToMessageID {
		if errors.Is(err, appErr) {
			msg, locErr := localizer.Localize(&i18n.LocalizeConfig{
				MessageID: messageID,
			})
			if locErr != nil {
				return err.Error()
			}

			return msg
		}
	}

	return err.Error()
}

func SQLError(op string, err error) error {
	return fmt.Errorf("%s: sql error: %w", op, err)
}

func SQLQueryError(op string, err error) error {
	return fmt.Errorf("%s: sql query error: %w", op, err)
}

func SQLExecError(op string, err error) error {
	return fmt.Errorf("%s: sql exec error: %w", op, err)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s': %s", e.Field, e.Message)
}

package handler

import (
	"github.com/google/wire"
	httpHandler "github.com/nori-plugins/authentication/internal/handler/http"
	"github.com/nori-plugins/authentication/internal/handler/http/authentication"
	"github.com/nori-plugins/authentication/internal/handler/http/mfa_recovery_code"
	"github.com/nori-plugins/authentication/internal/handler/http/mfa_secret"
)

var HandlerSet = wire.NewSet(
	wire.Struct(new(httpHandler.Handler), "R", "AuthenticationService",
		"MfaRecoveryCodeService", "Config", "AuthenticationHandler", "MfaRecoveryCodeHandler", "MfaSecretHandler"),
	authentication.New,
	mfa_recovery_code.New,
	mfa_secret.New,
)
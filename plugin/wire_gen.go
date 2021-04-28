// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nori-io/common/v4/pkg/domain/logger"
	"github.com/nori-io/common/v4/pkg/domain/registry"
	"github.com/nori-io/interfaces/database/orm/gorm"
	http2 "github.com/nori-io/interfaces/nori/http"
	"github.com/nori-plugins/authentication/internal/config"
	"github.com/nori-plugins/authentication/internal/handler/http"
	"github.com/nori-plugins/authentication/internal/handler/http/authentication"
	mfa_recovery_code4 "github.com/nori-plugins/authentication/internal/handler/http/mfa_recovery_code"
	mfa_totp3 "github.com/nori-plugins/authentication/internal/handler/http/mfa_totp"
	settings2 "github.com/nori-plugins/authentication/internal/handler/http/settings"
	social_provider3 "github.com/nori-plugins/authentication/internal/handler/http/social_provider"
	"github.com/nori-plugins/authentication/internal/helper/cookie"
	error2 "github.com/nori-plugins/authentication/internal/helper/error"
	"github.com/nori-plugins/authentication/internal/helper/goth_provider"
	mfa_recovery_code2 "github.com/nori-plugins/authentication/internal/helper/mfa_recovery_code"
	"github.com/nori-plugins/authentication/internal/helper/security"
	"github.com/nori-plugins/authentication/internal/repository/authentication_log"
	"github.com/nori-plugins/authentication/internal/repository/mfa_recovery_code"
	"github.com/nori-plugins/authentication/internal/repository/mfa_totp"
	"github.com/nori-plugins/authentication/internal/repository/session"
	"github.com/nori-plugins/authentication/internal/repository/social_provider"
	"github.com/nori-plugins/authentication/internal/repository/user"
	"github.com/nori-plugins/authentication/internal/service/auth"
	authentication_log2 "github.com/nori-plugins/authentication/internal/service/authentication_log"
	mfa_recovery_code3 "github.com/nori-plugins/authentication/internal/service/mfa_recovery_code"
	mfa_totp2 "github.com/nori-plugins/authentication/internal/service/mfa_totp"
	session2 "github.com/nori-plugins/authentication/internal/service/session"
	"github.com/nori-plugins/authentication/internal/service/settings"
	social_provider2 "github.com/nori-plugins/authentication/internal/service/social_provider"
	user2 "github.com/nori-plugins/authentication/internal/service/user"
	"github.com/nori-plugins/authentication/pkg/transactor"
)

// Injectors from wire.go:

func Initialize(registry2 registry.Registry, config2 config.Config, logger2 logger.FieldLogger) (*http.Handler, error) {
	httpHttp, err := http2.GetHttp(registry2)
	if err != nil {
		return nil, err
	}
	db, err := pg.GetGorm(registry2)
	if err != nil {
		return nil, err
	}
	params := transactor.Params{
		Db:     db,
		Logger: logger2,
	}
	transactorTransactor := transactor.New(params)
	userRepository := user.New(transactorTransactor)
	securityParams := security.Params{
		Config: config2,
	}
	securityHelper := security.New(securityParams)
	userParams := user2.Params{
		UserRepository: userRepository,
		Transactor:     transactorTransactor,
		Config:         config2,
		SecurityHelper: securityHelper,
	}
	userService := user2.New(userParams)
	authenticationLogRepository := authentication_log.New(transactorTransactor)
	authentication_logParams := authentication_log2.Params{
		AuthenticationLogRepository: authenticationLogRepository,
		Transactor:                  transactorTransactor,
	}
	authenticationLogService := authentication_log2.New(authentication_logParams)
	sessionRepository := session.New(transactorTransactor)
	sessionParams := session2.Params{
		SessionRepository: sessionRepository,
		Transactor:        transactorTransactor,
	}
	sessionService := session2.New(sessionParams)
	authParams := auth.Params{
		Config:                   config2,
		UserService:              userService,
		AuthenticationLogService: authenticationLogService,
		SessionService:           sessionService,
		Transactor:               transactorTransactor,
		SecurityHelper:           securityHelper,
	}
	authenticationService := auth.New(authParams)
	cookieParams := cookie.Params{
		Config: config2,
	}
	cookieHelper := cookie.New(cookieParams)
	errorParams := error2.Params{
		Logger: logger2,
	}
	errorHelper := error2.New(errorParams)
	authenticationParams := authentication.Params{
		AuthenticationService: authenticationService,
		SessionService:        sessionService,
		Logger:                logger2,
		Config:                config2,
		CookieHelper:          cookieHelper,
		ErrorHelper:           errorHelper,
	}
	authenticationHandler := authentication.New(authenticationParams)
	mfaRecoveryCodeRepository := mfa_recovery_code.New(transactorTransactor)
	mfa_recovery_codeParams := mfa_recovery_code2.Params{
		Config: config2,
	}
	mfaRecoveryCodesHelper := mfa_recovery_code2.New(mfa_recovery_codeParams)
	params2 := mfa_recovery_code3.Params{
		MfaRecoveryCodeRepository: mfaRecoveryCodeRepository,
		MfaRecoveryCodeHelper:     mfaRecoveryCodesHelper,
		Config:                    config2,
	}
	mfaRecoveryCodeService := mfa_recovery_code3.New(params2)
	params3 := mfa_recovery_code4.Params{
		MfaRecoveryCodeService: mfaRecoveryCodeService,
		Logger:                 logger2,
		CookieHelper:           cookieHelper,
		ErrorHelper:            errorHelper,
	}
	mfaRecoveryCodeHandler := mfa_recovery_code4.New(params3)
	mfaTotpRepository := mfa_totp.New(transactorTransactor)
	mfa_totpParams := mfa_totp2.Params{
		MfaTotpRepository: mfaTotpRepository,
		UserService:       userService,
		Config:            config2,
	}
	mfaTotpService := mfa_totp2.New(mfa_totpParams)
	params4 := mfa_totp3.Params{
		MfaTotpService: mfaTotpService,
		Logger:         logger2,
		CookieHelper:   cookieHelper,
		ErrorHelper:    errorHelper,
	}
	mfaTotpHandler := mfa_totp3.New(params4)
	settingsParams := settings.Params{
		SessionRepository: sessionRepository,
		UserService:       userService,
		SecurityHelper:    securityHelper,
	}
	settingsService := settings.New(settingsParams)
	params5 := settings2.Params{
		SettingsService: settingsService,
		Logger:          logger2,
		CookieHelper:    cookieHelper,
		ErrorHelper:     errorHelper,
	}
	settingsHandler := settings2.New(params5)
	socialProviderRepository := social_provider.New(transactorTransactor)
	social_providerParams := social_provider2.Params{
		SocialProviderRepository: socialProviderRepository,
	}
	socialProvider := social_provider2.New(social_providerParams)
	params6 := social_provider3.Params{
		SocialProviderService: socialProvider,
		Logger:                logger2,
		CookieHelper:          cookieHelper,
		ErrorHelper:           errorHelper,
	}
	socialProviderHandler := social_provider3.New(params6)
	gothProviderHelper := goth_provider.New()
	handler := &http.Handler{
		R:                      httpHttp,
		AuthenticationHandler:  authenticationHandler,
		MfaRecoveryCodeHandler: mfaRecoveryCodeHandler,
		MfaTotpHandler:         mfaTotpHandler,
		SettingsHandler:        settingsHandler,
		SocialProviderHandler:  socialProviderHandler,
		GothProviderHelper:     gothProviderHelper,
		SocialProviderService:  socialProvider,
	}
	return handler, nil
}

// wire.go:

var set = wire.NewSet(pg.GetGorm, http2.GetHttp)

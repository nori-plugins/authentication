package authentication

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/markbates/goth/gothic"
	"github.com/nori-plugins/authentication/internal/domain/errors"
	"github.com/nori-plugins/authentication/pkg/enum/social_provider_status"

	error2 "github.com/nori-plugins/authentication/internal/domain/helper/error"

	"github.com/nori-plugins/authentication/internal/domain/helper/cookie"
	"github.com/nori-plugins/authentication/pkg/enum/session_status"

	"github.com/nori-plugins/authentication/internal/config"

	"github.com/nori-plugins/authentication/internal/handler/http/response"

	"github.com/nori-io/common/v4/pkg/domain/logger"

	"github.com/nori-plugins/authentication/internal/domain/entity"

	"github.com/go-chi/chi"
	"github.com/nori-plugins/authentication/internal/domain/service"
)

type AuthenticationHandler struct {
	authenticationService service.AuthenticationService
	sessionService        service.SessionService
	logger                logger.FieldLogger
	config                config.Config
	cookieHelper          cookie.CookieHelper
	errorHelper           error2.ErrorHelper
	socialProviderService service.SocialProvider
}

type Params struct {
	AuthenticationService service.AuthenticationService
	SessionService        service.SessionService
	Logger                logger.FieldLogger
	Config                config.Config
	CookieHelper          cookie.CookieHelper
	ErrorHelper           error2.ErrorHelper
	SocialProviderService service.SocialProvider
}

func New(params Params) *AuthenticationHandler {
	return &AuthenticationHandler{
		authenticationService: params.AuthenticationService,
		sessionService:        params.SessionService,
		logger:                params.Logger,
		config:                params.Config,
		cookieHelper:          params.CookieHelper,
		errorHelper:           params.ErrorHelper,
	}
}

func (h *AuthenticationHandler) Session(w http.ResponseWriter, r *http.Request) {
	sessionId, err := h.cookieHelper.GetSessionID(r)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, http.ErrNoCookie.Error(), http.StatusUnauthorized)
	}

	sess, user, err := h.authenticationService.GetSessionInfo(r.Context(), sessionId)
	if err != nil {
		h.logger.Error("%s", err)
		h.errorHelper.Error(w, err)
	}

	response.JSON(w, r, SessionResponse{
		Success:  true,
		Message:  "session exists",
		Email:    user.Email,
		Phone:    user.PhoneCountryCode + user.PhoneNumber,
		OpenedAt: sess.OpenedAt,
	})
}

func (h *AuthenticationHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data, err := newSignUpData(r)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = h.authenticationService.SignUp(r.Context(), data)
	if err != nil {
		h.logger.Error("%s", err)
		h.errorHelper.Error(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthenticationHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	data, err := newSignInData(r)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sess, mfaType, err := h.authenticationService.SignIn(r.Context(), data)
	if err != nil {
		h.logger.Error("%s", err)
		h.errorHelper.Error(w, err)
	}

	h.cookieHelper.SetSession(w, sess)

	response.JSON(w, r, SignInResponse{
		Success: true,
		Message: "User sign in",
		MfaType: *mfaType,
	})
}

func (h *AuthenticationHandler) SignInMfa(w http.ResponseWriter, r *http.Request) {
	_, err := h.cookieHelper.GetSessionID(r)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, http.ErrNoCookie.Error(), http.StatusUnauthorized)
	}

	data, err := newSignInMfaData(r)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sess, err := h.authenticationService.SignInMfa(r.Context(), data)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.cookieHelper.SetSession(w, sess)

	response.JSON(w, r, SignInMfaResponse{
		Success: true,
		Message: "User sign in by mfa",
	})
}

func (h *AuthenticationHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	sessionId, err := h.cookieHelper.GetSessionID(r)
	if err != nil {
		h.logger.Error("%s", err)
		http.Error(w, http.ErrNoCookie.Error(), http.StatusUnauthorized)
	}

	data := &entity.Session{
		ID:         0,
		UserID:     0,
		SessionKey: nil,
		Status:     0,
		OpenedAt:   time.Time{},
		ClosedAt:   time.Time{},
		UpdatedAt:  time.Time{},
	}
	//@todo
	if data.Status != session_status.Active {
	}

	if err := h.authenticationService.SignOut(r.Context(), &entity.Session{SessionKey: []byte(sessionId)}); err != nil {
		h.logger.Error("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// todo: redirect

	h.cookieHelper.UnsetSession(w)
	http.Redirect(w, r, "/", 0)
}

func (h *AuthenticationHandler) HandleSocialProvider(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "social_provider")

	data := service.GetByNameData{Name: name}

	provider, err := h.socialProviderService.GetByName(r.Context(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if provider.Status != social_provider_status.Enabled {
		http.Error(w, errors.SocialProviderNotFound.Error(), http.StatusBadRequest)
	}

	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *AuthenticationHandler) HandleSocialProviderCallBack(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "social_provider")

	data := service.GetByNameData{Name: name}

	provider, err := h.socialProviderService.GetByName(r.Context(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if provider.Status != social_provider_status.Enabled {
		http.Error(w, errors.SocialProviderNotFound.Error(), http.StatusBadRequest)
	}

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(w, user)
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

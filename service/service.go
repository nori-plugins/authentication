package service

import (
	"context"

	"github.com/cheebo/gorest"
	"github.com/cheebo/rand"
	"github.com/nori-io/nori-common/interfaces"

	"github.com/nori-io/auth/service/database"
	//"github.com/cheebo/gorest"
	//	"github.com/cheebo/rand"
	"github.com/sirupsen/logrus"
)

type Service interface {
	SignUp(ctx context.Context, req SignUpRequest) (resp *SignUpResponse)
	SignIn(ctx context.Context, req SignInRequest) (resp *SignInResponse)
	SignOut(ctx context.Context, req SignOutRequest) (resp *SignOutResponse)
}

type Config struct {
	Sub func() string
	Iss func() string
}

type service struct {
	auth    interfaces.Auth
	db      database.Database
	session interfaces.Session
	cfg     *Config
	log     *logrus.Logger
}

func NewService(
	auth interfaces.Auth,
	session interfaces.Session,
	cfg *Config,
	log *logrus.Logger,
	db database.Database,
) Service {
	return &service{
		auth:    auth,
		db:      db,
		session: session,
		cfg:     cfg,
		log:     log,
	}
}

func (s *service) SignUp(ctx context.Context, req SignUpRequest) (resp *SignUpResponse) {
	var err error
		var model *database.AuthModel
		resp = &SignUpResponse{}
		errField := rest.ErrFieldResp{
			Meta: rest.ErrFieldRespMeta{
				ErrCode: 400,
			},
		}

		if model, err = s.db.Auth().FindByEmail(req.Email); err != nil {
			resp.Err = rest.ErrorInternal(err.Error())
			return resp
		}

		if model != nil && model.Id_Auth != 0 {
			errField.AddError("email", 400, "Email already exists.")
		}
		if errField.HasErrors() {
			resp.Err = errField
			return resp
		}

		model = &database.AuthModel{
			Email_Auth:    req.Email,
			Password_Auth: req.Password,

		}

		err = s.db.Auth().Create(model)
		if err != nil {
			s.log.Error(err)
			resp.Err = rest.ErrFieldResp{
				Meta: rest.ErrFieldRespMeta{
					ErrCode:    500,
					ErrMessage: err.Error(),
				},
			}
			return resp
		}

		resp.Email = req.Email

	return resp
}

func (s *service) SignIn(ctx context.Context, req SignInRequest) (resp *SignInResponse) {
	resp = &SignInResponse{}

	model, err := s.db.Auth().FindByEmail(req.Email)
	if err != nil {
		resp.Err = rest.ErrorInternal("Internal error")
		return resp
	}
	if model == nil {
		resp.Err = rest.ErrorNotFound("User not found")
		return resp
	}

	if req.Password != model.Password_Auth {
		resp.Err = rest.ErrorNotFound("User not found")
		return resp
	}

	sid := rand.RandomAlphaNum(32)

	token, err := s.auth.AccessToken(func(op interface{}) interface{} {
		key, ok := op.(string)
		if !ok || key == "" {
			return ""
		}
		switch key {
		case "raw":
			return map[string]string{
				"id":    string(model.Id_Auth),
				"email": model.Email_Auth,
			}
		case "jti":
			return sid
		case "sub":
			return s.cfg.Sub()
		case "iss":
			return s.cfg.Iss()
		default:
			return ""
		}
	})

	if err != nil {
		resp.Err = rest.ErrorInternal(err.Error())
		return resp
	}

	s.session.Save([]byte(sid), interfaces.SessionActive, 0)

	resp.Id = uint64(model.Id_Auth)
	resp.Token = token
	resp.User = *model

	return resp
}

func (s *service) SignOut(ctx context.Context, req SignOutRequest) (resp *SignOutResponse) {
	resp = &SignOutResponse{}
	s.session.Delete(s.session.SessionId(ctx))
	return resp
}

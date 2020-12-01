package auth

import (
	"context"

	"github.com/nori-io/authentication/internal/domain/entity"

	"github.com/nori-io/authentication/internal/domain/repository"

	rest "github.com/cheebo/gorest"
	"github.com/cheebo/rand"
	serv "github.com/nori-io/authentication/internal/domain/service"
	"github.com/nori-io/authentication/internal/domain/service/database"
	h "github.com/nori-io/interfaces/nori/http"
	s "github.com/nori-io/interfaces/nori/session"
)

type service struct {
	session s.Session
	http    h.Transport
	db      repository.UserRepository
}

func New(sessionInstance s.Session, httpInstance h.Transport, dbInstance repository.UserRepository /*cfg Config,*/) serv.AuthenticationService {

	return &service{
		session: sessionInstance,
		http:    httpInstance,
		db:      dbInstance,
	}

}
func (s *service) SignUp(ctx context.Context, data serv.SignUpData) (*entity.User, error) {
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

func (s *service) SignIn(ctx context.Context, data serv.SignInData) (*entity.Session, error) {
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

	//s.session.Save([]byte(sid), interfaces.SessionActive, 0)

	resp.Id = uint64(model.Id_Auth)
	resp.Token = token
	resp.User = *model

	return resp
}

func (s *service) SignOut(ctx context.Context, data *entity.Session) error {
	resp = &SignOutResponse{}
	//s.session.Delete(s.session.SessionId(ctx))
	return resp
}

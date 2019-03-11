package database

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type authenticationHistory struct {
	db  *sql.DB
	log *log.Logger
}

func (a *authenticationHistory) Create(model *AuthenticationHistoryModel) error {
	_, err := a.db.Exec("INSERT INTO authentification_history (user_id, signin, meta) VALUES(?,?,?)",
		model.UserId, time.Now(), model.Meta)
	return err
}

func (a *authenticationHistory) Update(model *AuthenticationHistoryModel) error {
	if model.Id == 0 {
		return errors.New("Empty model")
	}
	_, err := a.db.Exec("UPDATE authentification_history SET user_id = ?, signin = ?, meta = ?, signout = ?  WHERE id = ? ",
		model.UserId, model.SignIn, model.Meta, model.SignOut, model.Id)
	return err
}
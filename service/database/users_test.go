package database_test

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"

	"github.com/nori-io/auth/service/database"
)

type (
	AnyTime struct{}
)

func TestUsers_Create_userEmail(t *testing.T) {
	mockDatabase, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDatabase.Close()
	defer mock.ExpectClose()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users (status_account, type, created, updated,mfa_type) VALUES(?,?,?,?,?)").
		WithArgs("active", "vendor", AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		RowError(1, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO auth (user_id,  email, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES(?,?,?,?,?,?,?,?)").
		WithArgs(1, "test@mail.ru", "pass", "salt", AnyTime{}, AnyTime{}, false, false).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	d := database.DB(mockDatabase, logrus.New())

	err = d.Users().Create(&database.AuthModel{
		Email:    "test@mail.ru",
		Password: "pass",
		Salt:     "salt",
	}, &database.UsersModel{
		Status_account: "active",
		Type:           "vendor",
	})
	if err != nil {
		t.Error(err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	clear(d)
	d = nil
	debug.SetGCPercent(1)
	runtime.GC()

}

func TestUsers_Create_userPhone(t *testing.T) {
	t.Parallel()

	mockDatabase, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.ExpectClose()
	defer mockDatabase.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users (status_account, type, created, updated,mfa_type) VALUES(?,?,?,?,?)").
		WithArgs("active", "vendor", AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		RowError(1, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO auth (user_id, phone_country_code, phone_number, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES(?,?,?,?,?,?,?,?,?)").
		WithArgs(1, "8", "9191501490", "pass", "salt", AnyTime{}, AnyTime{}, false, false).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	d := database.DB(mockDatabase, logrus.New())
	err = d.Users().Create(&database.AuthModel{
		PhoneCountryCode: "8",
		PhoneNumber:      "9191501490",
		Password:         "pass",
		Salt:             "salt",
	}, &database.UsersModel{
		Status_account: "active",
		Type:           "vendor",
	})
	if err != nil {
		t.Error(err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	clear(d)
	d = nil
	runtime.GC()
}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

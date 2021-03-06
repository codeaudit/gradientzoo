package models

import (
	"database/sql"
	"time"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v3/zero"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

const USER_TABLE = "auth_user"
const BCRYPT_COST = 10

type UserDb struct {
	DB  *runner.DB
	Api *ApiCollection
}

//go:generate counterfeiter $GOFILE UserApi
type UserApi interface {
	ById(id interface{}) (*User, error)
	ByIds(ids []interface{}) ([]*User, error)
	Delete(id interface{}) error
	Save(*User) error
	Hydrate([]*User) error
	Truncate() error

	// TODO: Potentially this should be a separate interface
	ByEmail(email string) (*User, error)
	ByUsername(username string) (*User, error)
}

func NewUserDb(db *runner.DB, api *ApiCollection) *UserDb {
	return &UserDb{
		DB:  db,
		Api: api,
	}
}

type User struct {
	Id               string    `db:"id" json:"id"`
	Email            string    `db:"email" json:"-"`
	Username         string    `db:"username" json:"username"`
	PasswordHash     string    `db:"password_hash" json:"-"`
	StripeCustomerId string    `db:"stripe_customer_id" json:"-"`
	CreatedTime      time.Time `db:"created_time" json:"created_time"`

	// Hydrated fields
	HasStripeCustomerId zero.Bool `json:"has_stripe_customer_id,omitempty"`
}

func NewUser(email, username, password string) *User {
	user := &User{
		Id:          uuid.NewUUID().String(),
		Email:       email,
		Username:    username,
		CreatedTime: time.Now().UTC(),
	}
	user.SetPassword(password)
	return user
}

func (user *User) SetPassword(password string) error {
	hsh, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hsh)
	return nil
}

func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
}

func (db *UserDb) ById(id interface{}) (*User, error) {
	var user User
	err := db.DB.
		Select("*").
		From(USER_TABLE).
		Where("id = $1", id).
		QueryStruct(&user)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}

func (db *UserDb) ByIds(ids []interface{}) ([]*User, error) {
	if len(ids) == 0 {
		return []*User{}, nil
	}
	var users []*User
	err := db.DB.
		Select("*").
		From(USER_TABLE).
		Where("id IN $1", IdStrings(ids)).
		QueryStructs(&users)
	if users == nil {
		users = []*User{}
	}
	return users, err
}

func (db *UserDb) Delete(id interface{}) error {
	_, err := db.DB.
		DeleteFrom(USER_TABLE).
		Where("id = $1", id).
		Exec()
	return err
}

func (db *UserDb) Save(user *User) error {
	cols := []string{
		"id",
		"email",
		"username",
		"password_hash",
		"stripe_customer_id",
		"created_time",
	}
	vals := []interface{}{
		user.Id,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.StripeCustomerId,
		user.CreatedTime,
	}
	_, err := db.DB.
		Upsert(USER_TABLE).
		Columns(cols...).
		Values(vals...).
		Where("id = $1", user.Id).
		Exec()
	return err
}

func (db *UserDb) Hydrate(users []*User) error {
	for _, user := range users {
		user.HasStripeCustomerId = zero.BoolFrom(user.StripeCustomerId != "")
	}
	return nil
}

func (db *UserDb) Truncate() error {
	_, err := db.DB.DeleteFrom(USER_TABLE).Exec()
	return err
}

// -

func (db *UserDb) ByEmail(email string) (*User, error) {
	var user User
	err := db.DB.
		Select("*").
		From(USER_TABLE).
		Where("UPPER(email) = UPPER($1)", email).
		QueryStruct(&user)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}

func (db *UserDb) ByUsername(username string) (*User, error) {
	var user User
	err := db.DB.
		Select("*").
		From(USER_TABLE).
		Where("UPPER(username) = UPPER($1)", username).
		QueryStruct(&user)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}

package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/VikaPaz/pantheon/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	conn *sql.DB
	log  *logrus.Logger
}

func NewUserRepository(conn *sql.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		conn: conn,
		log:  logger,
	}
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	row := r.conn.QueryRow("INSERT INTO users (username) values "+
		"($1) RETURNING *", user.Username)
	if err := row.Err(); err != nil {
		return models.User{}, models.ErrCreateUserResponse
	}
	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return models.User{}, models.ErrCreateUserResponse
	}
	r.log.Debugf("Inserted user: %v", id)
	user.Id = id.String()
	return user, nil
}

func (r *UserRepository) GetById(user models.User) (models.User, error) {
	builder := sq.Select("users").Where(sq.Eq{"id": user.Id})
	builder = builder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return models.User{}, models.ErrUserDeleteResponse
	}

	r.log.Debugf("Executing query: %v", query)
	row := r.conn.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return models.User{}, models.ErrGetUserResponse
	}

	var res models.User
	if err := row.Scan(&user); err != nil {
		return models.User{}, models.ErrCreateUserResponse
	}
	r.log.Debugf("Inserted user: %v", user)
	return res, nil
}

func (r *UserRepository) GetByUsername(user models.User) (models.User, error) {
	builder := sq.Select("users").Where(sq.Eq{"username": user.Username})
	builder = builder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return models.User{}, models.ErrUserDeleteResponse
	}

	r.log.Debugf("Executing query: %v", query)
	row := r.conn.QueryRow(query, args...)
	if err := row.Err(); err != nil {
		return models.User{}, models.ErrGetUserResponse
	}

	var res models.User
	if err := row.Scan(&user); err != nil {
		return models.User{}, models.ErrCreateUserResponse
	}
	r.log.Debugf("Inserted user: %v", user)
	return res, nil
}

func (r *UserRepository) Delete(user models.User) error {
	builder := sq.Delete("users").Where(sq.Eq{"id": user.Id})
	builder = builder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return models.ErrUserDeleteResponse
	}

	r.log.Debugf("Executing query: %v", query)
	_, err = r.conn.Exec(query, args...)
	if err != nil {
		return models.ErrUserDeleteResponse
	}
	return nil
}

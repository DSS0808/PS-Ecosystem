package repository

import (
	"database/sql"
	"fmt"
	"github.com/VikaPaz/pantheon/internal/models"
	_ "github.com/lib/pq"
)

type Config struct{ Host, Port, User, Password, Dbname string }

func Connection(conf Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, models.ErrConnectionDBFailed
	}

	return db, db.Ping()
}

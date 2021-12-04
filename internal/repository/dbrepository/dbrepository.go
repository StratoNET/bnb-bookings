package dbrepository

import (
	"database/sql"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/repository"
)

type mariaDBRepository struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewMariaDBRepository(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepository {
	return &mariaDBRepository{
		App: app,
		DB:  conn,
	}
}

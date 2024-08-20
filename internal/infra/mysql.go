package infra

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"time"
)

func NewMysql(dsn string) (*sqlx.DB, util.CloseFn) {
	db, err := sqlx.Connect("mysql", dsn)
	util.Panic(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	util.Panic(err)

	log.Info().Msg("initialization mysql successfully")
	return db, func(ctx context.Context) (err error) {
		log.Info().Msg("starting close mysql db")

		err = db.Close()
		if err != nil {
			return err
		}

		log.Info().Msg("close mysql db successfully")
		return
	}
}

package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pedropiedade7/mundo-invest-test/internal/app/handler"
	clienthandler "github.com/pedropiedade7/mundo-invest-test/internal/app/handler/client"
	pipefyhandler "github.com/pedropiedade7/mundo-invest-test/internal/app/handler/pipefy"
	"github.com/pedropiedade7/mundo-invest-test/internal/config"
	"github.com/pedropiedade7/mundo-invest-test/internal/infra/database"
	"github.com/pedropiedade7/mundo-invest-test/internal/infra/pipefy"
	clientrepo "github.com/pedropiedade7/mundo-invest-test/internal/infra/repository/client"
	pipefyrepo "github.com/pedropiedade7/mundo-invest-test/internal/infra/repository/pipefy"
	"github.com/pedropiedade7/mundo-invest-test/internal/service"
)

type App struct {
	db     *sql.DB
	router *gin.Engine
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := database.NewPostgresDB(cfg.PostgresConnString())
	if err != nil {
		return nil, err
	}

	clientRepo := clientrepo.New(db)
	pipefyRepo := pipefyrepo.New(db)
	pipefyClient := pipefy.NewClient(cfg.PipefyPipeID)

	createClient := service.NewCreateClientService(clientRepo, pipefyClient)
	processWebhook := service.NewProcessWebhookService(clientRepo, pipefyRepo, pipefyClient)

	clientHandler := clienthandler.New(createClient)
	pipefyHandler := pipefyhandler.New(processWebhook)
	httpHandler := handler.NewHandler(clientHandler, pipefyHandler, db)

	return &App{
		db:     db,
		router: handler.NewRouter(httpHandler),
	}, nil
}

func (a *App) Close() error {
	return a.db.Close()
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}

package cmd

import (
	"fmt"
	"github.com/mirfnsyh/base-service/internal/app/server/http"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"

	"github.com/mirfnsyh/base-service/config"
	"github.com/mirfnsyh/base-service/internal/app/commons"
	"github.com/mirfnsyh/base-service/internal/app/commons/appcontext"
	"github.com/mirfnsyh/base-service/internal/app/domain/repositories"
	"github.com/mirfnsyh/base-service/internal/app/usecases"
)

var rootCmd = &cobra.Command{
	Use:   "base-service",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func start() {
	cfg := config.GetConfig()
	app := appcontext.NewAppContext(cfg)

	var err error

	dbMysql, err := app.GetDBInstance(appcontext.DBDialectMysql)
	if err != nil {
		logrus.Fatalf("Failed to start, error connect to DB MySQL | %v", err)
		return
	}

	opt := commons.Options{
		AppCtx:   app,
		Config:   cfg,
		Database: dbMysql,
	}

	repo := initRepository(opt)
	service := initService(opt, repo)

	httpServer := http.NewServer(opt, service)
	httpServer.StartApp()
}

func initRepository(opt commons.Options) *repositories.Repository {
	repo := repositories.Repository{}

	return &repo
}

func initService(opt commons.Options, repo *repositories.Repository) *usecases.Services {
	healthCheck := usecases.NewHealthCheck(opt)

	svc := usecases.Services{
		HealthCheckService: healthCheck,
	}

	return &svc
}

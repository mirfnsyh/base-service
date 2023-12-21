package cmd

import (
	"database/sql"
	"fmt"
	"github.com/mirfnsyh/base-service/config"
	"github.com/mirfnsyh/base-service/internal/app/commons/appcontext"
	"os"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const MigrationDirectory = "migrations/"

var makeMigrationCommand = &cobra.Command{
	Use:   "migrate/create",
	Short: "Create new migration database file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		migrationDir := MigrationDirectory
		migrationName := args[0]
		err := createMigrationFile(migrationDir, migrationName)

		if err != nil {
			return
		}
	},
}

var migrateUpCommand = &cobra.Command{
	Use:   "migrate/up",
	Short: "Migrate up for database bs-assignment-service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := getFileMigrationSource()

		err := doMigration(source, migrate.Up)

		if err != nil {
			return
		}
	},
}

var migrateDownCommand = &cobra.Command{
	Use:   "migrate/down",
	Short: "Migrate down for database bs-assignment-service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := getFileMigrationSource()

		err := doMigration(source, migrate.Down)
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateUpCommand)
	rootCmd.AddCommand(migrateDownCommand)
	rootCmd.AddCommand(makeMigrationCommand)
}

func createMigrationFile(mDir string, mName string) error {

	var migrationContent = `-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
-- [your SQL script here]

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
-- [your SQL script here]
`
	filename := fmt.Sprintf("%d_%s.sql", time.Now().Unix(), mName)
	filepath := fmt.Sprintf("%s%s", mDir, filename)

	f, err := os.Create(filepath)
	if err != nil {
		logrus.Errorln("Error create migration file", filepath, filename, err)
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logrus.Errorln("Error when closing file")
		}
	}(f)

	_, err = f.WriteString(migrationContent)
	if err != nil {
		return err
	}

	err = f.Sync()

	if err != nil {
		return err
	}

	logrus.Infoln("New migration file has been created", filename)
	return nil
}

func getFileMigrationSource() migrate.FileMigrationSource {
	source := migrate.FileMigrationSource{
		Dir: MigrationDirectory,
	}

	return source
}

func doMigration(fileMigrationSource migrate.FileMigrationSource, direction migrate.MigrationDirection) error {
	cfg := config.GetConfig()
	context := appcontext.NewAppContext(cfg)

	db, err := context.GetDBInstance(appcontext.DBDialectMysql)

	if err != nil {
		logrus.Errorln("Cannot connect to database")
		return err
	}

	sqlDb, _ := db.DB()

	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			return
		}
	}(sqlDb)

	totalMigrated, err := migrate.Exec(sqlDb, appcontext.DBDialectMysql, fileMigrationSource, direction)

	if err != nil {
		logrus.Errorln("Migration failed", err)
		return err
	}

	logrus.Infof("Migrate success, total migrated: %d", totalMigrated)
	return err
}

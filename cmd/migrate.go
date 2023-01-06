package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	migratePath   string
	migrationName string

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate database",
	}

	upArgsValidator = func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("incorrect args length")
		}
		return nil
	}

	downArgsValidator = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("incorrect args length")
		}

		limit, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		if limit < 0 {
			return errors.New("incorrect migrations limit")
		}
		return nil
	}

	createArgsValidator = func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("incorrect args length")
		}
		return nil
	}

	migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "make migration",
		Args:  upArgsValidator,
		RunE:  migrateUpCmdHandler,
	}

	migrateDownCmd = &cobra.Command{
		Use:   "down [limit]",
		Short: "rollback migration",
		Args:  downArgsValidator,
		RunE:  migrateDownCmdHandler,
	}

	migrateCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create migration",
		Args:  createArgsValidator,
		RunE:  migrateCreateCmdHandler,
	}
)

// Command init function.
func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().StringVarP(&migratePath, "migrationsPath", "m", ".bin/app/migrations/pgsql", "Path to migrations")
	migrateCmd.PersistentFlags().StringVarP(&migrationName, "migrationName", "n", "", "New migration name")
}

// Command handler func.
func migrateUpCmdHandler(_ *cobra.Command, args []string) (err error) {
	return makeMigration(migrate.Up, 2)
}

// Command handler func.
func migrateDownCmdHandler(_ *cobra.Command, args []string) (err error) {
	var limit int
	limit, _ = strconv.Atoi(args[0])
	return makeMigration(migrate.Down, limit)
}

func migrateCreateCmdHandler(_ *cobra.Command, args []string) (err error) {
	if migrationName == "" {
		return errors.New("no name provided")
	}
	if _, err = os.Stat(migratePath); os.IsNotExist(err) {
		if err = os.MkdirAll(migratePath, 0755); err != nil {
			return err
		}
	}
	var name = time.Now().Format("20060102_150405_") + migrationName + ".sql"
	var f *os.File
	if f, err = os.Create(migratePath + "/" + name); err != nil {
		return err
	}
	template := `-- +migrate Up

-- +migrate Down
`
	_, err = f.WriteString(template)
	return
}

func makeMigration(direction migrate.MigrationDirection, limit int) (err error) {
	var (
		migrationsList = &migrate.FileMigrationSource{
			Dir: migratePath,
		}
	)
	//формируем конфиги для подключения к базам
	pgSql := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable binary_parameters=yes",
		config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Password, config.DB.Database)
	//создаем конекты к используемым базам
	var connectPG *sql.DB
	if connectPG, err = sql.Open("postgres", pgSql); err != nil {
		return errors.Wrap(err, "error create conn to database")
	}
	log.Info("Create Connect")
	defer connectPG.Close()

	var n int
	if n, err = migrate.ExecMax(connectPG, "postgres", migrationsList, direction, limit); err != nil {
		return err
	}

	log.Printf("Applied migrations count %d", n)
	return nil
}

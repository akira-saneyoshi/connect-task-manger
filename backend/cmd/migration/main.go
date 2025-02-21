package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/a-s/connect-task-manage/internal/infrastructure/config"
	"github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql" //MySQL Driver
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "sql/migrations", "directory with migration files")
)

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Drivers:
    mysql
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	//設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	command := args[0]
	switch command {
	case "create":
		if len(args) < 2 {
			flags.Usage()
			return
		}
		if err := goose.Create(nil, *dir, args[1], "sql"); err != nil { // Create コマンドの修正
			log.Fatalf("goose create: %v", err)
		}
		return

	default:
		db, err := sql.Open("mysql", cfg.DB.DSN)
		if err != nil {
			log.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := goose.SetDialect("mysql"); err != nil {
			log.Fatal(err)
		}

		// コマンドの実行
		if err := goose.Run(command, db, *dir, args[1:]...); err != nil {
			log.Fatalf("goose %v: %v", command, err)
		}
	}
}

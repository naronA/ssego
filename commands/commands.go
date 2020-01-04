package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"ssego"

	"github.com/urfave/cli"
)

var engine *ssego.Engine

func Main() {
	app := cli.NewApp()
	app.Name = "ssego"
	app.Usage = `simple and small search engine for learning`
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		createIndexCommand,
		searchCommand,
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ssego")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	engine = ssego.NewSearchEngine(db)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const (
	exactArgs = iota
	minArgs

	maxArgs
)

func checkArgs(context *cli.Context, expected, checkType int) error {
	var err error
	cmdName := context.Command.Name
	switch checkType {
	case exactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf("%s: %q requires exactly %d argument(s)", os.Args[0], cmdName, expected)
		}
	case minArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q rquires a minimum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	case maxArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q rquires a maximum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	}
	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		cli.ShowCommandHelp(context, cmdName)
		return err
	}
	return nil
}

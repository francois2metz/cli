package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
)

var (
	PgSQLConsoleCommand = cli.Command{
		Name:     "pgsql-console",
		Aliases:  []string{"psql-console", "postgresql-console"},
		Category: "Databases",
		Usage:    "Run an interactive console with your PostgreSQL addon",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "size, s", Value: "", Usage: "Size of the container"},
			cli.StringFlag{Name: "env, e", Value: "", Usage: "Environment variable name to use for the connection to the database"},
		},
		Description: ` Run an interactive console with your PostgreSQL addon.

   Examples
    scalingo --app my-app pgsql-console
    scalingo --app my-app pgsql-console --size L
    scalingo --app my-app pgsql-console --env MY_PSQL_URL

   The --size flag makes it easy to specify the size of the container executing
   the PostgreSQL console. Each container size has different price and performance.
   You can read more about container sizes here:
   http://doc.scalingo.com/internals/container-sizes.html

    # See also 'mongo-console' and 'mysql-console'
`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "pgsql-console")
				return
			}

			err := db.PgSQLConsole(db.PgSQLConsoleOpts{
				App:          detect.CurrentApp(c),
				Size:         c.String("s"),
				VariableName: c.String("e"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "pgsql-console")
		},
	}
)

package main

import (
	"fmt"
	utils "github.com/analogj/go-util/utils"
	"github.com/analogj/tangle/webapp/backend/pkg/config"
	"github.com/analogj/tangle/webapp/backend/pkg/errors"
	"github.com/analogj/tangle/webapp/backend/pkg/version"
	"github.com/analogj/tangle/webapp/backend/pkg/web"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var goos string
var goarch string

func main() {

	config, err := config.Create()
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	//we're going to load the config file manually, since we need to validate it.
	err = config.ReadConfig("tangle.yaml")                // Find and read the config file
	if _, ok := err.(errors.ConfigFileMissingError); ok { // Handle errors reading the config file
		//ignore "could not find config file"
	} else if err != nil {
		os.Exit(1)
	}

	app := &cli.App{
		Name:     "tangle",
		Usage:    "Software Bill-of-Materials storage/query service. Compatible with SPDX/CycloneDX SBOM formats ",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			tangle := "github.com/AnalogJ/tangle"

			var versionInfo string
			if len(goos) > 0 && len(goarch) > 0 {
				versionInfo = fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)
			} else {
				versionInfo = fmt.Sprintf("dev-%s", version.VERSION)
			}

			subtitle := tangle + utils.LeftPad2Len(versionInfo, " ", 65-len(tangle))

			color.New(color.FgGreen).Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
				`
			 ___   ___  ____  __  __  ____  ____  _  _  _  _
			/ __) / __)(  _ \(  )(  )(_  _)(_  _)( \( )( \/ )
			\__ \( (__  )   / )(__)(   )(   _)(_  )  (  \  /
			(___/ \___)(_)\_)(______) (__) (____)(_)\_) (__)
			%s

			`), subtitle))

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start the tangle server",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)
					if c.IsSet("config") {
						err = config.ReadConfig(c.String("config")) // Find and read the config file
						if err != nil {                             // Handle errors reading the config file
							//ignore "could not find config file"
							fmt.Printf("Could not find config file at specified path: %s", c.String("config"))
							return err
						}
					}

					if c.Bool("debug") {
						config.Set("log.level", "DEBUG")
					}

					if c.IsSet("log-file") {
						config.Set("log.file", c.String("log-file"))
					}

					webServer := web.AppEngine{Config: config}

					return webServer.Start()
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "config",
						Usage: "Specify the path to the config file",
					},
					&cli.StringFlag{
						Name:    "log-file",
						Usage:   "Path to file for logging. Leave empty to use STDOUT",
						Value:   "",
						EnvVars: []string{"TANGLE_LOG_FILE"},
					},

					&cli.BoolFlag{
						Name:    "debug",
						Usage:   "Enable debug logging",
						EnvVars: []string{"TANGLE_DEBUG", "DEBUG"},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}

}

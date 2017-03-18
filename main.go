package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	_debug = log.New(ioutil.Discard, "", 0)
	_info  = log.New(ioutil.Discard, "", 0)
	_warn  = log.New(ioutil.Discard, "", 0)
	_error = log.New(ioutil.Discard, "", 0)
)

type users struct {
	Users []struct {
		Name        string
		Gecos       string
		Public_keys []string
		HomeDir     string
		Shell       string
	}
}

func initLogger(level string) {
	switch level {
	case "debug":
		_debug = log.New(os.Stdout, "debug: ", log.Lshortfile)
		_info = log.New(os.Stdout, "info: ", log.Lshortfile)
		_warn = log.New(os.Stdout, "warn: ", log.Lshortfile)
		_error = log.New(os.Stderr, "error: ", log.Lshortfile)
	case "info":
		_info = log.New(os.Stdout, "info: ", log.Lshortfile)
		_warn = log.New(os.Stdout, "warn: ", log.Lshortfile)
		_error = log.New(os.Stderr, "error: ", log.Lshortfile)
	case "warn":
		_warn = log.New(os.Stdout, "warn: ", log.Lshortfile)
		_error = log.New(os.Stderr, "error: ", log.Lshortfile)
	case "error":
		_error = log.New(os.Stderr, "error: ", log.Lshortfile)
	}

}

func main() {
	fmt.Println("START")
	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Name = "uman"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jon Hadfield",
			Email: "jon@lessknown.co.uk",
		},
	}
	app.HelpName = "-"
	app.Usage = "AWS - SSH User Manager"
	app.Description = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Usage: "user file path",
		},
		cli.BoolFlag{
			Name:  "check, c",
			Usage: "check the user file",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "Set log level (debug, info, warn, error)",
			Value: "warn",
		},
	}
	app.Action = func(c *cli.Context) error {
		initLogger(c.String("log-level"))
		fmt.Println(c.String("path"))

		if _, err := os.Stat(c.String("path")); err == nil {
			fmt.Println("Got file")
			file, err := os.Open(c.String("path"))
			if err != nil {
				fmt.Println("OOPS")
			}

			defer file.Close()
			userFile, err := ioutil.ReadFile(c.String("path"))
			fmt.Println(string(userFile))
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}
			var users users
			err = yaml.Unmarshal(userFile, &users)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			} else {
				for _, m := range users.Users {
					fmt.Printf("%+v\n\n", m)
				}
			}
		} else {
			fmt.Printf("Could not find file: %s\n", c.String("path"))
		}

		return nil
	}
	app.Run(os.Args)

}

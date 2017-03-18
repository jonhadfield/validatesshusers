package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
	"strings"
	"encoding/base64"
)

var (
	_debug = log.New(ioutil.Discard, "", 0)
	_info = log.New(ioutil.Discard, "", 0)
	_warn = log.New(ioutil.Discard, "", 0)
	_error = log.New(ioutil.Discard, "", 0)
)

type users struct {
	Users []struct {
		Name        string
		Gecos       string
		Public_keys []string
		Home_dir    string
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
	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Name = "validatesshusers"
	app.Version = "0.0.8"
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
			Name:  "log-level",
			Usage: "Set log level (debug, info, warn, error)",
			Value: "warn",
		},
	}
	app.Action = func(c *cli.Context) error {
		initLogger(c.String("log-level"))
		var usersFile = c.Args().Get(0)
		if _, err := os.Stat(usersFile); err == nil {
			file, err := os.Open(usersFile)
			if err != nil {
				return err
			}
			defer file.Close()
			userFile, err := ioutil.ReadFile(usersFile)
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}
			var users users
			err = yaml.Unmarshal(userFile, &users)
			if err != nil {
				fmt.Println("\nSpecified file is not a valid users file.")
				log.Fatalf("Unmarshal: %v", err)
			} else {
				var allOk bool = true
				for _, user := range users.Users {
					_debug.Printf("%+v\n\n", user)
					// Check no upper case characters in username
					if user.Name != strings.ToLower(user.Name) {
						fmt.Printf("Warning: user: \"%s\". Usernames should only contain lower case characters\n", user.Name)
					}
					// If home directory is specified, warn if not default
					if user.Home_dir != "" && user.Home_dir != "/home/" + user.Name {
						fmt.Printf("Warning: Specified home directory: %s differs from expected default: /home/%s\n", user.Home_dir, user.Name)
						allOk = false
					}
					// Check gecos is specified and contains email address
					if user.Gecos == "" {
						fmt.Printf("Error: gecos (comment) not specified for user: %s\n", user.Name)
						allOk = false
						// Check gecos contains email address
					} else if ! strings.Contains(user.Gecos, "@") {
						fmt.Printf("Error: gecos (comment) for user: %s does not have an email address\n", user.Name)
						allOk = false
					}
					// Check user has public keys defined
					if len(user.Public_keys) == 0 {
						fmt.Printf("Error: user: %s has no public keys specified\n", user.Name)
					} else {
						// Validate public keys
						// TODO: Improve validation
						for _, b64PublicKey := range user.Public_keys {
							publicKey, _ := base64.StdEncoding.DecodeString(b64PublicKey)
							var invalidKeys int = 0
							if ! strings.HasPrefix(string(publicKey), "ssh-rsa ") {
								invalidKeys++
							}
							if invalidKeys > 0 {
								fmt.Printf("Error: user: %s has %d invalid public keys\n", user.Name, invalidKeys)
							}
						}
					}

				}
				if allOk {
					fmt.Println("OK")
				}
			}
		} else {
			fmt.Printf("Could not find file: %s\n", usersFile)
		}

		return nil
	}
	app.Run(os.Args)

}

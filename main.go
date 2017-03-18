package main

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").Bool()
	file  = kingpin.Arg("file", "Users YAML file.").String()
)

func main() {
	kingpin.Version("0.1.2")
	kingpin.Parse()
	if *debug {
		initLogger("debug")
	}
	if *file == "" {
		kingpin.Usage()
		os.Exit(1)
	}
	var usersFile = *file

	if _, err := os.Stat(usersFile); err == nil {
		file, err := os.Open(usersFile)
		if err != nil {
			return
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
			var allOk = true
			for _, user := range users.Users {
				_debug.Printf("%+v\n\n", user)
				// Check no upper case characters in username
				if user.Name != strings.ToLower(user.Name) {
					fmt.Printf("Warning: user: \"%s\". Usernames should only contain lower case characters\n", user.Name)
				}
				// If home directory is specified, warn if not default
				if user.Home_dir != "" && user.Home_dir != "/home/"+user.Name {
					fmt.Printf("Warning: Specified home directory: %s differs from expected default: /home/%s\n", user.Home_dir, user.Name)
					allOk = false
				}
				// Check gecos is specified and contains email address
				if user.Gecos == "" {
					fmt.Printf("Error: gecos (comment) not specified for user: %s\n", user.Name)
					allOk = false
					// Check gecos contains email address
				} else if !strings.Contains(user.Gecos, "@") {
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
						var invalidKeys int
						if !strings.HasPrefix(string(publicKey), "ssh-rsa ") {
							invalidKeys++
						}
						if invalidKeys > 0 {
							fmt.Printf("Error: user: %s has %d invalid public keys\n", user.Name, invalidKeys)
						}
					}
				}

			}
			if allOk {
				fmt.Println("SYNTAX OK")
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
	} else {
		fmt.Printf("Error: Could not open path: %s\n", usersFile)
		os.Exit(1)
	}

	return

}

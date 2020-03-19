package env

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dashjay/logging"
	"github.com/joho/godotenv"
)

var (
	HTTPPort string
	GRPCPort string
	BoltDB   string
)

func parse() {
	HTTPPort = getDefault("HTTPPort", "8080")
	GRPCPort = getDefault("GRPCPort", "8081")
	BoltDB = getDefault("BoltDB", "cookies.db")
}

func init() {
	envFileName := ".env"
	FlagSet := flag.CommandLine
	FlagSet.StringVar(&envFileName, "env", envFileName, "the env file which web app will use to extract its environment variables")
	_ = flag.CommandLine.Parse(os.Args[1:])
	Load(envFileName)
}

// Load loads environment variables that are being used across the whole app.
// Loading from file(s), i.e .env or dev.env
//
// Example of a 'dev.env':
// PORT=8080
// DSN=mongodb://localhost:27017
//
// After `Load` the callers can get an environment variable via `os.Getenv`.
func Load(envFileName string) {
	if args := os.Args; len(args) > 1 && args[1] == "help" {
		_, _ = fmt.Fprintln(os.Stderr, "https://github.com/kataras/iris/blob/master/_examples/tutorials/mongodb/README.md")
		os.Exit(-1)
	}

	logging.Infof("Loading environment variables from file: %s\n", envFileName)
	// If more than one filename passed with comma separated then load from all
	// of these, a env file can be a partial too.
	envFiles := strings.Split(envFileName, ",")
	for i := range envFiles {
		if filepath.Ext(envFiles[i]) == "" {
			envFiles[i] += ".env"
		}
	}

	if err := godotenv.Load(envFiles...); err != nil {
		panic(fmt.Sprintf("error loading environment variables from [%s]: %v", envFileName, err))
	}

	envMap, _ := godotenv.Read(envFiles...)

	for k, v := range envMap {
		logging.Infof("%s=%s", k, v)
	}

	parse()
}

func getDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		_ = os.Setenv(key, def)
		value = def
	}

	return value
}

func getIntDefault(key string, def int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

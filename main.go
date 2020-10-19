package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	routers "github.com/1024casts/fastim/router"
	"github.com/1024casts/snake/pkg/snake"
	v "github.com/1024casts/snake/pkg/version"

	"github.com/1024casts/snake/pkg/conf"
)

var (
	cfg     = pflag.StringP("config", "c", "", "fastim config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	if err := conf.Init(*cfg); err != nil {
		panic(err)
	}

	// init app
	snake.App = snake.New(conf.Conf)

	// Set gin mode.
	gin.SetMode(snake.ModeRelease)
	if viper.GetString("app.run_mode") == snake.ModeDebug {
		gin.SetMode(snake.ModeDebug)
		snake.App.DB.Debug()
	}

	// Create the Gin engine.
	router := snake.App.Router

	// API Routes.
	routers.Load(router)

	// start server
	snake.App.Run()
}

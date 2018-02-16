package main

import (
	"fmt"

	"bytes"

	"github.com/kelseyhightower/envconfig"
	"github.com/mangeshhendre/grpcutils"
	handler "github.com/mangeshhendre/mathsvc/pkg/mathhandler"
	logxi "github.com/mgutz/logxi/v1"
)

// This is my config yo.
type Config struct {
	DSN string `required:"true" desc:"Oracle connection string"`
}

const config_prefix string = ""

func main() {

	logger := logxi.New("mathsvc")

	c := &Config{}

	err := envconfig.Process(config_prefix, c)
	if err != nil {
		buf := bytes.NewBuffer(nil)
		envconfig.Usagef(config_prefix, c, buf, envconfig.DefaultTableFormat)
		logger.Fatal("Unable to load config:", "Config", buf.String())
	}

	//Create the server instance.
	server, err := handler.New(c.DSN)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to create server instance: %v", err))
	}

	defer server.Close()

	grpcManager, err := grpcutils.New("mathsvc.grpc")
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Register your service.
	grpcManager.RegisterHandlers(server)

	// Start the show.
	grpcManager.Run()
}

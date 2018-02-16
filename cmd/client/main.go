package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mangeshhendre/grpcutils"
	pb "github.com/mangeshhendre/models/services_math_v1"
	logxi "github.com/mgutz/logxi/v1"
)

func main() {
	var (
		number1    = flag.Float64("number1", 1, "First Number")
		number2    = flag.Float64("number2", 1, "Second Number")
		insecure   = flag.Bool("insecure", false, "Should I ignore certificate warnings?")
		print      = flag.Bool("print", false, "Should I print the results out?")
		iterations = flag.Int("n", 1, "The number of times to call the service")
	)
	flag.Parse()

	logger := logxi.New("mathsvc.client")

	// Make a suitable grpc client.
	conn, err := grpcutils.MakeGRPCClientConn(logger,
		grpcutils.EnvOrDefault("GRPC_AUTH_URL", "https://authentication."+grpcutils.EnvOrDefault("DOMAIN", "sgtec.io")),
		grpcutils.EnvOrDefault("GRPC_AUTH_USERNAME", "AUTH_URL_UNSET"),
		grpcutils.EnvOrDefault("GRPC_AUTH_PASSWORD", "AUTH_URL_UNSET"),
		grpcutils.EnvOrDefault("GRPC_HOST", "mathsvc.grpc."+grpcutils.EnvOrDefault("DOMAIN", "safeguardproperties.com")),
		grpcutils.EnvOrDefault("GRPC_PORT", "32363"),
		*insecure,
	)

	if err != nil {
		panic(logger.Error("I was unable to create a grpc client."))
	}

	// Create the math client.
	c := pb.NewMathClient(conn)

	// defer un(trace("GRPC Calls"))
	for i := 0; i < *iterations; i++ {
		//getAndMap(*print, logger, c, *propertyLoanID, *orderNumber)
		doItAgain(*print, logger, c, *number1, *number2)
	}

}

func doItAgain(print bool, logger logxi.Logger, client pb.MathClient, number1, number2 float64) {
	request := &pb.MathRequest{
		Number1: number1,
		Number2: number2,
	}

	if number1 == 0 || number2 == 0 {
		return
	}

	logger.Info("Calling AddNumber")
	result, err := client.AddNumber(context.Background(), request)
	if err != nil {
		logger.Info("AddNumber: Call Failed", "Request", request, "Error", err.Error())
		return
	}
	if print {
		dumpJSON("AddNumberResult", result)
	}
}

func getAndMap(print bool, logger logxi.Logger, client pb.MathClient, number1, number2 float64) {
	defer un(trace(fmt.Sprintf("GetAndMap Prop: %f, Order: %f", number1, number2)))

	request := &pb.MathRequest{
		Number1: number1,
		Number2: number2,
	}

	// We should Set these up to be persistent.
	result := &pb.MathResponse{}
	var err error

	logger.Info("Calling DevideNumber")
	result, err = client.DevideNumber(context.Background(), request)
	if err != nil {
		logger.Info("DevideNumber: Call Failed", "Request", request, "Error", err.Error())
		return
	}

	dumpJSON("DevideNumber", result.Result)

	logger.Info("Calling MultiplyNumber")
	result, err = client.MultiplyNumber(context.Background(), request)
	if err != nil {
		logger.Info("MultiplyNumber: Call Failed", "Request", request, "Error", err.Error())
		return
	}
	dumpJSON("MultiplyNumber", result.Result)
}

func dumpJSON(name string, result interface{}) {
	foo, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal to json")
	}
	fmt.Printf("%s\n%s\n", name, foo)
}

func trace(s string) (string, time.Time) {
	// log.Println("START:", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}

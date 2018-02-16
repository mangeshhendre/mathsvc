package mathdb

import (
	"time"

	"github.com/jmoiron/sqlx"
	pb "github.com/mangeshhendre/models/services_math_v1"
	"github.com/mangeshhendre/tracer"
	logxi "github.com/mgutz/logxi/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client is the actual database client.
type Client struct {
	DB     *sqlx.DB
	logger logxi.Logger
	tracer *tracer.Tracer
}

// New is the new database/sql version of the processing.
func New(DB *sqlx.DB) (*Client, error) {

	return &Client{
		DB:     DB,
		logger: logxi.New("sql.go"),
		tracer: tracer.New("graphite:8125", "grpc.mathsvc.adb", 1),
	}, nil
}

// GetWorkOrderDate is a function that given the work order number, it will return the time.Time that represents the audited result for that order.
func (c *Client) getSomeInfoFromDb(in *pb.MathRequest) (mathResp *pb.MathResponse, err error) {
	c.logger.Info("getSomeInfoFromDb")
	defer c.tracer.Statsd("getSomeInfoFromDb", time.Now())

	mathResp = &pb.MathResponse{}

	row := c.DB.QueryRowx(someDBQuery, in.Number1, in.Number2)
	if row.Err() != nil {
		return nil, status.Errorf(codes.Internal, "getSomeInfoFromDb: query error, number1: %f, number2 %f, error: %s", in.Number1, in.Number2, err.Error())
	}

	err = row.StructScan(mathResp)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "getSomeInfoFromDb: failed to query database, Number1: %f, Number2 %f, error: %s", in.Number1, in.Number2, err.Error())
	}

	// This means we did not find a work order.
	if mathResp.Result == 0 {
		return nil, status.Errorf(codes.NotFound, "getSomeInfoFromDb: result not found.query database, Number1: %f, Number2 order number %f", in.Number1, in.Number2)
	}

	return mathResp, err
}

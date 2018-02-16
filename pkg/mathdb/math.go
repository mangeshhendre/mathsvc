package mathdb

import (
	"strconv"
	"time"

	pb "github.com/mangeshhendre/models/services_math_v1"
	context "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddNumber will retrieve database details given the request.
func (c *Client) AddNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	c.logger.Info("AddNumber")
	defer c.tracer.Statsd("AddNumber", time.Now())

	if in.Number1 == 0 || in.Number2 == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Zero is invalid for Number1:%f or Number2:%f", in.Number1, in.Number2)
	}

	//this is sample how to call Db results.
	dbResults, err := c.getSomeInfoFromDb(in)
	if err != nil {
		return nil, err
	}
	c.logger.Info(strconv.FormatFloat(dbResults.Result, 'f', 2, 64))

	response := &pb.MathResponse{}

	response.Result = in.Number1 + in.Number2

	return response, err
}

// MultiplyNumber will retrieve database details given the request.
func (c *Client) MultiplyNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	c.logger.Info("AddNumber")
	defer c.tracer.Statsd("AddNumber", time.Now())

	if in.Number1 == 0 || in.Number2 == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Zero is invalid for Number1:%f or Number2:%f", in.Number1, in.Number2)
	}

	//this is sample how to call Db results.
	dbResults, err := c.getSomeInfoFromDb(in)
	if err != nil {
		return nil, err
	}
	c.logger.Info(strconv.FormatFloat(dbResults.Result, 'f', 2, 64))

	response := &pb.MathResponse{}

	response.Result = in.Number1 * in.Number2

	return response, err
}

// DevideNumber will retrieve database details given the request.
func (c *Client) DevideNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	c.logger.Info("AddNumber")
	defer c.tracer.Statsd("AddNumber", time.Now())

	if in.Number1 == 0 || in.Number2 == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Zero is invalid for Number1:%f or Number2:%f", in.Number1, in.Number2)
	}

	//this is sample how to call Db results.
	dbResults, err := c.getSomeInfoFromDb(in)
	if err != nil {
		return nil, err
	}
	c.logger.Info(strconv.FormatFloat(dbResults.Result, 'f', 2, 64))

	response := &pb.MathResponse{}

	response.Result = in.Number1 / in.Number2

	return response, err
}

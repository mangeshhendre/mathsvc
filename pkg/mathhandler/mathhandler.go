package mathhandler

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mangeshhendre/mathsvc/pkg/mathcache"
	"github.com/mangeshhendre/mathsvc/pkg/mathdb"
	pb "github.com/mangeshhendre/models/services_math_v1"
	"github.com/mangeshhendre/tracer"
	_ "github.com/mattn/go-oci8"
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is the local server handler.
type Server struct {
	Debug         bool
	LibraryDebug  bool
	DB            *sqlx.DB
	cacheInstance pb.MathServer
	dbInstance    pb.MathServer
	tracer        *tracer.Tracer
	logger        log.Logger
}

// New creates a new server handler instance.
func New(dsn string) (*Server, error) {
	// Need a logger.
	logger := log.New("mathsvc.Handler")
	//Create database things here.
	DB, err := sqlx.Connect("oci8", dsn)
	if err != nil {
		return nil, logger.Error("Unable to establish mattn database connection: ", "Error", err)
	}
	DB = DB.Unsafe()

	DB.Mapper = reflectx.NewMapperTagFunc("json", func(str string) string {
		// strToReturn := strings.ToUpper(str)
		// fmt.Println("String", str, "Returning", strToReturn)
		return strings.ToUpper(str)
	},
		func(value string) string {
			//fmt.Println("Value", value)
			var valueToReturn string
			if strings.Contains(value, ",") {
				valueToReturn = strings.Split(value, ",")[0]
			}
			valueToReturn = strings.Replace(valueToReturn, "_", "", -1)
			valueToReturn = strings.ToUpper(valueToReturn)
			//fmt.Println("Value", value, "Returning", valueToReturn)
			return valueToReturn
		})

	dbInstance, err := mathdb.New(DB)
	if err != nil {
		return nil, err
	}

	cacheInstance, err := mathcache.New(dbInstance)
	if err != nil {
		return nil, err
	}

	return &Server{
		DB:            DB,
		cacheInstance: cacheInstance,
		dbInstance:    dbInstance,
		tracer:        tracer.New("graphite:8125", "grpc.mathsvc", 1),
		logger:        logger,
	}, nil
}

// Close will shut it all down.
func (s *Server) Close() {
	s.DB.Close()
}

//AddNumber retrieves the lite math for the specified propertyID
func (s *Server) AddNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	defer s.tracer.Statsd("AddNumber", time.Now())
	return s.cacheInstance.AddNumber(ctx, in)
}

//MultiplyNumber retrieves the lite math for the specified propertyID
func (s *Server) MultiplyNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	defer s.tracer.Statsd("MultiplyNumber", time.Now())
	return s.cacheInstance.MultiplyNumber(ctx, in)
}

//DevideNumber retrieves math for the specified propertyID
func (s *Server) DevideNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	defer s.tracer.Statsd("DevideNumber", time.Now())
	if in.Number2 == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Request: Number2: %f", in.Number2)
	}
	return s.cacheInstance.DevideNumber(ctx, in)
}

// RegisterServices wraps setup of services in the handler library.
func (s *Server) RegisterServices(shim *grpc.Server) {
	defer s.tracer.Statsd("RegisterServices", time.Now())
	pb.RegisterMathServer(shim, s)

}

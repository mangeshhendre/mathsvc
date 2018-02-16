package mathcache

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/mangeshhendre/models/services_math_v1"
	"github.com/mangeshhendre/protocache"
	"github.com/mgutz/logxi/v1"
)

type MathCache struct {
	cache  *protocache.PC
	server pb.MathServer
	logger log.Logger
}

func New(imp pb.MathServer) (pb.MathServer, error) {
	client := &MathCache{
		server: imp,
		cache:  protocache.New(2*time.Second, "Sample", listMemcacheServers()...),
		logger: log.New("MathCache"),
	}
	return client, nil
}

func listMemcacheServers() []string {
	serverString := envOrDefault("MEMCACHE_SERVERS", "mem01:11211;mem01:11212;mem01:11213;mem02:11211;mem02:11212;mem03:11213;mem03:11211;mem03:11212;mem02:11213")
	return strings.Split(serverString, ";")
}

func envOrDefault(input, override string) string {
	foo := os.Getenv(input)
	if len(foo) != 0 {
		return foo
	}
	return override
}

func (s *MathCache) AddNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	return s.getAdditionFromCache(ctx, in)
}

func (s *MathCache) MultiplyNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	return s.getAdditionFromCache(ctx, in) //this needs to implement to get Multiplication
}

func (s *MathCache) DevideNumber(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	return s.getAdditionFromCache(ctx, in) //this needs to implement to get Division
}

func (s *MathCache) getAdditionFromCache(ctx context.Context, in *pb.MathRequest) (*pb.MathResponse, error) {
	response := &pb.MathResponse{}
	primaryContext := fmt.Sprintf("Number1:%d", in.Number1)
	secondaryContext := fmt.Sprintf("Number2:%d", in.Number2)
	cacheKey := "MathAddNumber"

	// Check the cache first.
	err := s.cache.Get(primaryContext, secondaryContext, cacheKey, response)
	if err == nil {
		// Successful result from cache.
		return response, nil
	}
	s.logger.Debug("Unable to get from memcache", "Error", err)

	if in.Number1 != 0 && in.Number2 != 0 {
		response, err = s.server.AddNumber(ctx, in)
	} else {
		err = s.logger.Error("Zero numbers cannot be added")
	}

	if err != nil {
		return nil, err
	}

	memcacheErr := s.cache.Set(primaryContext, secondaryContext, cacheKey, response, 10*time.Second)
	if memcacheErr != nil {
		// We give no sh*ts.
		s.logger.Debug("getAdditionFromCache: Unable to set record in memcache", "Error", memcacheErr)
	}

	return response, nil
}

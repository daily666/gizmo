package api

import (
	"net/http"

	kitserver "github.com/NYTimes/gizmo/server/kit"
	"github.com/NYTimes/gziphandler"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"

	"github.com/NYTimes/gizmo/examples/nyt"
)

type (
	// service will implement kitserver.Service.
	service struct {
		client nyt.Client
	}
	// Config is a struct to contain all the needed
	// configuration for our Service
	Config struct {
		MostPopularToken string `envconfig:"MOST_POPULAR_TOKEN"`
		SemanticToken    string `envconfig:"SEMANTIC_TOKEN"`
	}
)

var _ ApiServiceServer = service{}

// NewService will instantiate a Service
// with the given configuration.
func New(cfg Config) kitserver.Service {
	return service{
		nyt.NewClient(cfg.MostPopularToken, cfg.SemanticToken),
	}
}

func (s service) RouterOptions() []kitserver.RouterOption {
	return nil
}

func (s service) Options() []httptransport.ServerOption {
	return nil
}

// HTTPMiddleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s service) HTTPMiddleware(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// Middleware provides an kit/endpoint.Middleware hook wrapped around all requests.
func (s service) Middleware(e endpoint.Endpoint) endpoint.Endpoint {
	return e
}

func (s service) RPCOptions() []grpc.ServerOption {
	return nil
}

func (s service) ServiceDesc() *grpc.ServiceDesc {
	// snagged from the pb file
	return &_ApiService_serviceDesc
}

// JSONEndpoints is a listing of all endpoints available in the Service.
func (s service) HTTPEndpoints() map[string]map[string]kitserver.Endpoint {
	return map[string]map[string]kitserver.Endpoint{
		"/most-popular/{resourceType:[a-z]+}/{section:[a-z]+}/{timeframe:[0-9]+}": {
			"GET": {
				Endpoint:    s.getMostPopular,
				HTTPDecoder: decodeMostPopularRequest,
			},
		},
		"/cats": {
			"GET": {
				Endpoint: s.getCats,
			},
		},
	}
}
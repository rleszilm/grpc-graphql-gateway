package options

import "google.golang.org/grpc"

// ServerOptions the options for a graphql_resolver.
type ServerOptions struct {
	Host         string
	WithInsecure bool
	Conn         *grpc.ClientConn
}

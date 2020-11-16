package main

var goTemplate = `
// Code generated by proroc-gen-graphql, DO NOT EDIT.
package {{ .RootPackage.Name }}

import (
{{- if .Services }}
	"context"

	"github.com/rleszilm/grpc-graphql-gateway/options"
	"github.com/rleszilm/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
	"github.com/pkg/errors"
{{- end }}
	"github.com/graphql-go/graphql"

{{- range .Packages }}
	{{ if .Path }}{{ .Name }} "{{ .Path }}"{{ end }}
{{ end }}
)

var (
	{{- range .Enums }}
	gql__enum_{{ .Name }} *graphql.Enum // enum {{ .Name }} in {{ .Filename }}
	{{- end }}
	{{- range .Interfaces }}
	gql__interface_{{ .Name }} *graphql.Interface // message {{ .Name }} in {{ .Filename }}
	{{- end }}
	{{- range .Types }}
	gql__type_{{ .TypeName }} *graphql.Object // message {{ .Name }} in {{ .Filename }}
	{{- end }}
	{{- range .Inputs }}
	gql__input_{{ .TypeName }} *graphql.InputObject // message {{ .Name }} in {{ .Filename }}
	{{- end }}
)

{{ range $enum := .Enums -}}
func Gql__enum_{{ .Name }}() *graphql.Enum {
	if gql__enum_{{ .Name }} == nil {
		gql__enum_{{ .Name }} =  graphql.NewEnum(graphql.EnumConfig{
			Name: "{{ $.RootPackage.CamelName }}_Enum_{{ .Name }}",
			Values: graphql.EnumValueConfigMap{
{{- range .Values }}
				"{{ .Name }}": &graphql.EnumValueConfig{
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					Value: {{ $enum.Name }}({{ .Number }}),
				},
{{- end }}
			},
		})
	}
	return gql__enum_{{ .Name }}
}

{{ end }}

{{ range .Interfaces -}}
func Gql__interface_{{ .TypeName }}() *graphql.Interface {
	if gql__interface_{{ .TypeName }} == nil {
		gql__interface_{{ .TypeName }} =  graphql.NewInterface(graphql.InterfaceConfig{
			Name: "{{ $.RootPackage.CamelName }}_Interface_{{ .TypeName }}",
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{- end }}
			Fields: graphql.Fields{
{{- range .Fields }}
			{{- if not .IsCyclic }}
				"{{ .FieldName }}": &graphql.Field{
					Type: {{ .FieldType $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
				},
			{{- end }}
{{- end }}
			},
			ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
				return Gql__type_{{ .TypeName }}()
			},
		})
	}
	return gql__interface_{{ .TypeName }}
}

{{ end }}

{{ range .Types -}}
func Gql__type_{{ .TypeName }}() *graphql.Object {
	if gql__type_{{ .TypeName }} == nil {
		gql__type_{{ .TypeName }} =  graphql.NewObject(graphql.ObjectConfig{
			Name: "{{ $.RootPackage.CamelName }}_Type_{{ .TypeName }}",
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{- end }}
			Fields: graphql.Fields {
{{- range .Fields }}
				{{- if .IsResolve }}
				{{ $query := .ResolveSubField $.Services }}
				"{{ .FieldName }}": &graphql.Field{
						Type: {{ $query.QueryType }},
						{{- if $query.Comment }}
						Description: ` + "`" + `{{ $query.Comment }}` + "`" + `,
						{{- end }}
						Args: graphql.FieldConfigArgument{
						{{- range $query.Args }}
							"{{ .FieldName }}": &graphql.ArgumentConfig{
								Type: {{ .FieldType $.RootPackage.Path }},
								{{- if .Comment }}
								Description: ` + "`" + `{{ .Comment }}` + "`" + `,
								{{- end }}
								{{- if .DefaultValue }}
								DefaultValue: {{ .DefaultValue }},
								{{- end }}
							},
						{{- end }}
						},
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							var req {{ $query.InputType }}
							if err := runtime.MarshalRequest(p.Source, &req, {{ if $query.IsCamel }}true{{ else }}false{{ end }}); err != nil {
								return nil, errors.Wrap(err, "Failed to marshal resolver source for {{ $query.QueryName }}")
							} else if err = runtime.MarshalRequest(p.Args, &req, {{ if $query.IsCamel }}true{{ else }}false{{ end }}); err != nil {
								return nil, errors.Wrap(err, "Failed to marshal resolver request for {{ $query.QueryName }}")
							}
							{{ $s := index $.Services 0 }}
							x := new_graphql_resolver_{{ $s.Name }}(nil)
							conn, closer, err := x.CreateConnection(p.Context)
							if err != nil {
								return nil, errors.Wrap(err, "Failed to create gRPC connection for nested resolver")
							}
							defer closer()
							client := {{ $query.Package }}New{{ $query.Method.Service.Name }}Client(conn)
							resp, err := client.{{ $query.Method.Name }}(p.Context, &req)
							if err != nil {
								return nil, errors.Wrap(err, "Failed to call RPC {{ $query.Method.Name }}")
							}
							{{- if $query.IsPluckResponse }}
								{{- if $query.IsCamel }}
								return runtime.MarshalResponse(resp.Get{{ $query.PluckResponseFieldName }}()), nil
								{{- else }}
								return resp.Get{{ $query.PluckResponseFieldName }}(), nil
								{{- end }}
							{{- else }}
								{{- if $query.IsCamel }}
								return runtime.MarshalResponse(resp), nil
								{{- else }}
								return resp, nil
								{{- end }}
							{{- end }}
						},
				},
				{{- else }}
				"{{ .FieldName }}": &graphql.Field{
					Type: {{ .FieldType $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
				},
				{{- end }}
{{- end }}
			},
			{{- if .Interfaces }}
			Interfaces: []*graphql.Interface{
{{- range .Interfaces }}
			Gql__interface_{{ .TypeName }}(),
{{- end }}
			},
			{{- end }}
		})
	}
	return gql__type_{{ .TypeName }}
}

{{ end }}

{{ range .Inputs -}}
func Gql__input_{{ .TypeName }}() *graphql.InputObject {
	if gql__input_{{ .TypeName }} == nil {
		gql__input_{{ .TypeName }} =  graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "{{ $.RootPackage.CamelName }}_Input_{{ .TypeName }}",
			Fields: graphql.InputObjectConfigFieldMap{
{{- range .Fields }}
				"{{ .FieldName }}": &graphql.InputObjectFieldConfig{
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					Type: {{ .FieldTypeInput $.RootPackage.Path }},
				},
{{- end }}
			},
		})
	}
	return gql__input_{{ .TypeName }}
}

{{ end }}

{{ range $_, $service := .Services -}}

// graphql__resolver_{{ $service.Name }} is a struct for making query, mutation and resolve fields.
// This struct must be implemented runtime.SchemaBuilder interface.
type graphql__resolver_{{ $service.Name }} struct {

	// Automatic connection host
	host string

	// grpc dial options
	dialOptions []grpc.DialOption

	// grpc client connection.
	// this connection may be provided by user
	conn *grpc.ClientConn
}

// new_graphql_resolver_{{ $service.Name }} creates pointer of service struct
func new_graphql_resolver_{{ $service.Name }}(opts *options.ServerOptions) *graphql__resolver_{{ $service.Name }} {
	var conn *grpc.ClientConn
	host := "{{ if .Host }}{{ .Host }}{{ else }}localhost:50051{{ end }}"
	dialOptions := []grpc.DialOption{}

	if opts != nil {
		conn = opts.Conn
		
		if opts.Host != "" {
			host = opts.Host
		}

		if opts.WithInsecure || {{- if .Insecure }} true {{- else }} false {{- end }} {
			dialOptions = append(dialOptions, grpc.WithInsecure())
		}
	}
	
	res := &graphql__resolver_{{ .Name }}{
		conn: conn,
		host: host,
		dialOptions: dialOptions,
	}

	return res
}

// CreateConnection() returns grpc connection which user specified or newly connected and closing function
func (x *graphql__resolver_{{ $service.Name }}) CreateConnection(ctx context.Context) (*grpc.ClientConn, func(), error) {
	// If x.conn is not nil, user injected their own connection
	if x.conn != nil {
		return x.conn, func() {}, nil
	}

	// Otherwise, this handler opens connection with specified host
	conn, err := grpc.DialContext(ctx, x.host, x.dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return conn, func() { conn.Close() }, nil
}

// GetQueries returns acceptable graphql.Fields for Query.
func (x *graphql__resolver_{{ $service.Name }}) GetQueries(conn *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
{{- range .Queries }}
	{{- if not .IsResolver }}
	   "{{ .QueryName }}": &graphql.Field{
			Type: {{ .QueryType }},
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{- end }}
			Args: graphql.FieldConfigArgument{
			{{- range .Args }}
				"{{ .FieldName }}": &graphql.ArgumentConfig{
					Type: {{ .FieldType $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					{{- if .DefaultValue }}
					DefaultValue: {{ .DefaultValue }},
					{{- end }}
				},
			{{- end }}
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var req {{ .InputType }}
				if err := runtime.MarshalRequest(p.Args, &req, {{ if .IsCamel }}true{{ else }}false{{ end }}); err != nil {
					return nil, errors.Wrap(err, "Failed to marshal request for {{ .QueryName }}")
				}
				client := {{ .Package }}New{{ .Method.Service.Name }}Client(conn)
				resp, err := client.{{ .Method.Name }}(p.Context, &req)
				if err != nil {
					return nil, errors.Wrap(err, "Failed to call RPC {{ .Method.Name }}")
				}
				{{- if .IsPluckResponse }}
					{{- if .IsCamel }}
					return runtime.MarshalResponse(resp.Get{{ .PluckResponseFieldName }}()), nil
					{{- else }}
					return resp.Get{{ .PluckResponseFieldName }}(), nil
					{{- end }}
				{{- else }}
					{{- if .IsCamel }}
					return runtime.MarshalResponse(resp), nil
					{{- else }}
					return resp, nil
					{{- end }}
				{{- end }}
			},
	   },
	{{- end }}
{{- end }}
	}
}

// GetMutations returns acceptable graphql.Fields for Mutation.
func (x *graphql__resolver_{{ $service.Name }}) GetMutations(conn *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
{{- range .Mutations }}
		"{{ .MutationName }}": &graphql.Field{
			Type: {{ .MutationType }},
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{ end }}
			Args: graphql.FieldConfigArgument{
			{{- if .InputName }}
				"{{ .InputName }}": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(Gql__input_{{ .Input.TypeName }}()),
				},
			{{- else }}
			{{- range .Args }}
				"{{ .FieldName }}": &graphql.ArgumentConfig{
					Type: {{ .FieldTypeInput $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					{{- if .DefaultValue }}
					DefaultValue: {{ .DefaultValue }},
					{{- end }}
				},
			{{- end }}
			{{- end }}
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var req {{ .InputType }}
				{{- if .InputName }}
				if err := runtime.MarshalRequest(p.Args["{{ .InputName }}"], &req, {{ if .IsCamel }}true{{ else }}false{{ end }}); err != nil {
				{{- else }}
				if err := runtime.MarshalRequest(p.Args, &req, {{ if .IsCamel }}true{{ else }}false{{ end }}); err != nil {
				{{- end }}
					return nil, errors.Wrap(err, "Failed to marshal request for {{ .MutationName }}")
				}
				client := {{ .Package }}New{{ $service.Name }}Client(conn)
				resp, err := client.{{ .Method.Name }}(p.Context, &req)
				if err != nil {
					return nil, errors.Wrap(err, "Failed to call RPC {{ .Method.Name }}")
				}
				{{- if .IsPluckResponse }}
					{{- if .IsCamel }}
					return runtime.MarshalResponse(resp.Get{{ .PluckResponseFieldName }}()), nil
					{{- else }}
					return resp.Get{{ .PluckResponseFieldName }}(), nil
					{{- end }}
				{{- else }}
					{{- if .IsCamel }}
					return runtime.MarshalResponse(resp), nil
					{{- else }}
					return resp, nil
					{{- end }}
				{{- end }}
			},
		},
{{ end }}
	}
}

// WithGRPCAddr sets the address of the grpc server to use.
func (x *graphql__resolver_{{ $service.Name }}) WithGRPCAddr(addr string) {
	x.host = addr
}

// Register package divided graphql handler "without" *grpc.ClientConn,
// therefore gRPC connection will be opened and closed automatically.
// Occasionally you may worry about open/close performance for each handling graphql request,
// then you can call Register{{ .Name }}GraphqlHandler with *grpc.ClientConn manually.
func Register{{ .Name }}Graphql(mux *runtime.ServeMux) error {
	return Register{{ .Name }}GraphqlHandler(mux, nil)
}

// Register package divided graphql handler "with" *grpc.ClientConn.
// this function accepts your defined grpc connection, so that we reuse that and never close connection inside.
// You need to close it maunally when application will terminate.
// Otherwise, you can specify automatic opening connection with ServiceOption directive:
//
// service {{ .Name }} {
//    option (graphql.service) = {
//        host: "host:port"
//        insecure: true or false
//    };
//
//    ...with RPC definitions
// }
func Register{{ .Name }}GraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	opts := &options.ServerOptions{Conn: conn}
	return mux.AddHandler(new_graphql_resolver_{{ .Name }}(opts))
}

// Register{{ .Name }}GraphqlWithOptions registers the service with the given options.
func Register{{ .Name }}GraphqlWithOptions(mux *runtime.ServeMux, opts *options.ServerOptions) error {
	return mux.AddHandler(new_graphql_resolver_{{ .Name }}(opts))
}

{{ end }}
`

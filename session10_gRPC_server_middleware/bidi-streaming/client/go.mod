module github.com/Hocheung1997/bidi-streaming/client

require (
	github.com/Hocheung1997/bidi-streaming/service v0.0.0
	google.golang.org/grpc v1.67.1
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/Hocheung1997/bidi-streaming/service => ../service

go 1.23.0

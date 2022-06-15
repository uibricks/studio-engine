//go:generate protoc -I=. -I=$GOPATH/src -I=.. --gogofaster_out=plugins=grpc,paths=source_relative:. project_service.proto
//go:generate protoc -I=. -I=$GOPATH/src -I=.. --gogofaster_out=plugins=grpc,paths=source_relative:. mapping_service.proto
//go:generate protoc -I=. -I=$GOPATH/src -I=.. --gogofaster_out=plugins=grpc,paths=source_relative:. expression_service.proto
package service

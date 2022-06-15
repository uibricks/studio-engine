//go:generate protoc -I=. -I=$GOPATH/src --gogofaster_out=paths=source_relative:. expression.proto
package expression

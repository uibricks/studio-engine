FROM golang:1.16.0-alpine3.13 as build
RUN GOCACHE=OFF
RUN go env -w GOPRIVATE=github.com/uibricks
WORKDIR /go/src/add
COPY . /go/src/add
RUN apk add git
RUN git config --global url."https://golang:e9cb1d875b843153ea0aadf0eaee8773419a725f@github.com".insteadOf "https://github.com"
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /project

FROM scratch
COPY --from=build /project /project
ENTRYPOINT ["/project","-configFile", "default.yaml"]
EXPOSE 50054
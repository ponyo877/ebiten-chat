FROM golang:1.21-alpine AS build
WORKDIR /src
RUN apk update && apk add git
RUN git clone --depth 1 https://github.com/ponyo877/folks-ui.git /src
RUN go mod download
RUN go build -o /folks-ui main.go

FROM alpine:latest
WORKDIR /
COPY --from=build /folks-ui /folks-ui
ENTRYPOINT ["/folks-ui"]
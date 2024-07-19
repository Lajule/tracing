FROM golang:1.22-alpine AS base
COPY .. /src/tracing
WORKDIR /src/tracing/example

FROM base
RUN go install main.go
CMD ["/go/bin/main"]

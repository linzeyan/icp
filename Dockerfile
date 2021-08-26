ARG dist="/tmp/checkicp"
ARG projectDir="/icp"

FROM golang:1.16-alpine3.14 AS builder
ARG dist
ARG projectDir
WORKDIR ${projectDir}
COPY . .
RUN go mod tidy &&\
  go build -o ${dist} icp.go

FROM alpine
ARG dist
COPY --from=builder ${dist} /usr/local/bin/checkicp
ARG dist="/tmp/checkicp"
ARG projectDir="/icp"

FROM golang:1.16-alpine3.14 AS builder
RUN apk add upx
ARG dist
ARG projectDir
WORKDIR ${projectDir}
COPY . .
RUN go mod tidy &&\
  go build -o icp main.go &&\
  upx -9 -o ${dist} icp


FROM scratch
ARG dist
COPY --from=builder ${dist} /usr/local/bin/checkicp
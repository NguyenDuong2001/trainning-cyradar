FROM golang:1.16-alpine AS build

WORKDIR /

RUN apk add --no-cache bash

ARG CGO_ENABLED=0

ARG GO111MODULE=on

COPY go.mod ./

COPY go.sum ./

RUN go mod download -x
COPY . ./


RUN go build -o /bin/team

FROM alpine:3.15.0

WORKDIR /bin/

COPY --from=build /bin/team / ./

EXPOSE 3334

CMD ["/bin/team"]
######################
# dependencies layer #
######################
FROM golang:1.26.2-alpine AS deps

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


###############
# build layer #
###############
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .

# choose binary: api or grpc
ARG APP=api

RUN CGO_ENABLED=0 GOOS=linux go build -o $APP ./cmd/$APP


#################
# runtime layer #
#################
FROM alpine:3.22

WORKDIR /app

ARG APP=api  # re-declare it here

COPY --from=builder /app/$APP ./server

CMD ["./server"]

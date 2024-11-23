FROM golang:1.22.0-alpine AS build
WORKDIR /project/friendzy/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./cmd/Auth_ServiceExec ./cmd/main.go

FROM alpine:latest
WORKDIR /project/friendzy/

COPY --from=build /project/friendzy/cmd/Auth_ServiceExec ./cmd/Auth_ServiceExec
COPY --from=build /project/friendzy/auth.log ./
COPY --from=build /project/friendzy/dev.env ./dev.env
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apk update && apk add --no-cache ca-certificates tzdata

EXPOSE 50001
ENTRYPOINT ["/project/friendzy/cmd/Auth_ServiceExec"]
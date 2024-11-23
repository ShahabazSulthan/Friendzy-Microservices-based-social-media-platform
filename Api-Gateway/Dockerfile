FROM golang:1.22.0-alpine AS build

WORKDIR /project/friendzy/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./cmd/Api_GatewayExec ./cmd/main.go

FROM alpine:latest

WORKDIR /project/friendzy/

COPY --from=build /project/friendzy/cmd/Api_GatewayExec ./cmd/

COPY --from=build /project/friendzy/dev.env ./dev.env

COPY --from=build /project/friendzy/template ./template

RUN apk update && apk add --no-cache ca-certificates tzdata

EXPOSE 8000

ENTRYPOINT ["/project/friendzy/cmd/Api_GatewayExec"]
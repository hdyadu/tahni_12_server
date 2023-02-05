FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build main.go client.go hub.go room.go -o /tahni-server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /tahni-server /tahni-server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/tahni-server" ]
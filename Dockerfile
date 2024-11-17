FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /tahni-server

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build /tahni-server /tahni-server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/tahni-server" ]
FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build main.go client.go hub.go room.go
CMD [ "/app/main" ]
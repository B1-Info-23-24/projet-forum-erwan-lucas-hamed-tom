FROM golang:1.22.3-alpine
RUN apk update && apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go-api
EXPOSE 8080
EXPOSE 8081
CMD ["/go-api"]

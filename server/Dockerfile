FROM golang:latest
RUN apt-get update && apt-get install -y build-essential
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
EXPOSE 8080
CMD ["./server"]

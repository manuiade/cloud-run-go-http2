FROM golang:1.19 as go-builder

WORKDIR /app

COPY go.mod go.sum app.json index.html ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-http2

FROM alpine:latest
COPY --from=go-builder /app/go-http2 /app/go-http2
EXPOSE 8080
CMD ["/app/go-http2"]
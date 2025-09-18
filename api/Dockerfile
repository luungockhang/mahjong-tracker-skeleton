FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server ./cmd/server


FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]


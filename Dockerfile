FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /backdrop ./cmd/backdrop

FROM scratch
COPY --from=builder /backdrop /backdrop
ENTRYPOINT ["/backdrop"]

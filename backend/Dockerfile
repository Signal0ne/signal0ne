FROM golang:1.22.5 AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o s0backend ./main.go


FROM alpine:3.14

COPY --from=go-builder /app/s0backend .
COPY --from=go-builder /app/.default.env .default.env

# Start both applications
CMD ["./s0backend"]
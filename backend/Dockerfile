FROM golang:1.22.5 AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o s0backend ./main.go


FROM python:3.9-slim

WORKDIR /app
COPY scripts/python_interface/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY scripts/python_interface/ ./python_interface/

COPY --from=go-builder /app/s0backend .
COPY --from=go-builder /app/.default.env .default.env
COPY --from=go-builder /app/scripts/start.sh start.sh
RUN mkdir -p /net && touch /net/socket && chmod 777 /net/socket && chmod 777 start.sh

# Start both applications
CMD ["./start.sh"]
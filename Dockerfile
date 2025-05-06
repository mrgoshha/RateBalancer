FROM golang:alpine AS builder

WORKDIR usr/src/rateBalancer

#dependencies
COPY go.mod go.sum ./
RUN go mod download

#build
COPY . .
RUN go build -o /usr/local/bin/rateBalancer cmd/main.go

FROM alpine AS runner
COPY --from=builder /usr/local/bin/rateBalancer /
COPY config.yaml /config.yaml

CMD ["/rateBalancer"]
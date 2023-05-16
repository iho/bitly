FROM golang:alpine as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o executable ./cmd/shortener

WORKDIR /dist

RUN cp /build/executable .

FROM scratch

COPY --from=builder /dist/executable /app/

WORKDIR /app

CMD ["./executable"]
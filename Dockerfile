FROM golang:latest AS build

WORKDIR /build

COPY . .

RUN go build

FROM golang:latest

WORKDIR /app

COPY --from=build /build/httpServer .

CMD ["./httpServer"]




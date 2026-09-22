FROM golang:1.27 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /book-management-system ./cmd/main

FROM alpine:3.22

COPY --from=build /book-management-system /book-management-system
EXPOSE 8080
ENTRYPOINT ["/book-management-system"]

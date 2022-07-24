FROM golang:1.18-alpine3.15 as build
WORKDIR /app-build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o poker-app

FROM scratch
WORKDIR /app
COPY --from=build /app-build/poker-app ./poker-app
CMD ["./poker-app"]

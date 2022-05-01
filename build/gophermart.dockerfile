## Build
FROM golang:1.18-buster AS build
ENV CGO_ENABLED=0
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -o /gophermart cmd/gophermart/main.go
## Deploy
FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /gophermart /gophermart
COPY --from=build /app/migration /migration
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gophermart"]
FROM golang:1.18.1 AS builder
RUN go env -w GOPROXY=direct
WORKDIR /app/
COPY . .
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o api .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder --chown=nonroot:nonroot /app/api ./api
EXPOSE 8080
ENTRYPOINT [ "./api"]

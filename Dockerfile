FROM golang:1.19
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o kube-shield .

FROM alpine:latest
WORKDIR /
COPY --from=0 /app/kube-shield /
COPY hack/certs/server.crt /var/ssl/server.crt
COPY hack/certs/server.key /var/ssl/server.key
RUN chmod +x /kube-shield
CMD ["/kube-shield"]

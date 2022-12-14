FROM golang:1.19
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o kube-shield .

FROM alpine:latest
WORKDIR /
COPY --from=0 /app/kube-shield /
COPY certs /tmp/certs
RUN chmod +x /kube-shield
CMD ["/kube-shield"]
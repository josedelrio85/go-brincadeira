#Builder stage
FROM golang:alpine AS builder

#Install dependencies
RUN apk update && apk add --no-cache \
  git \
  ca-certificates \
  && update-ca-certificates 

# Add source files and set the proper work dir
COPY . $GOPATH/src/github.com/bysidecar/loaddatareport/
WORKDIR $GOPATH/src/github.com/bysidecar/loaddatareport


# Enable Go modules
ENV GO111MODULE=on
# Build the binary
RUN go build -mod=vendor -o /go/bin/loaddatareport

# Final image
FROM ellerbrock/alpine-mysql-client

# Copy our static executable
COPY --from=builder /go/bin/loaddatareport /go/bin/loaddatareport

# Copy the ca-certificates to be able to perform https requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV DB_PORT=3306
# Run the binary
ENTRYPOINT ["/go/bin/loaddatareport"]
# ENTRYPOINT ["mysql"]
# ENTRYPOINT ["/bin/bash"]
# CMD ["mysql"]

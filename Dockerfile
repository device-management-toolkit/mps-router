#*********************************************************************
# Copyright (c) Intel Corporation 2021
# SPDX-License-Identifier: Apache-2.0
#*********************************************************************/
#build stage
FROM golang:alpine@sha256:43c094ad24b6ac0546c62193baeb3e6e49ce14d3250845d166c77c25f64b0386 AS builder
RUN apk add --no-cache git ca-certificates && update-ca-certificates
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "1000" "scratchuser"
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app -ldflags="-s -w" -v ./cmd/
#final stage
FROM scratch
COPY --from=builder /go/bin/app /app
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER scratchuser
ENTRYPOINT ["/app"]
LABEL Name=mpsrouter Version=1.0.0
LABEL license='SPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2021: Intel'
EXPOSE 8003


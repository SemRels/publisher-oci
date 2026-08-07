# SPDX-License-Identifier: Apache-2.0
# SPDX-FileCopyrightText: 2026 The publisher-oci Authors

FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/plugin ./cmd/plugin && \
    CGO_ENABLED=0 GOBIN=/out go install oras.land/oras/cmd/oras@v1.3.1

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/plugin /usr/local/bin/plugin
COPY --from=build /out/oras /usr/local/bin/oras
ENTRYPOINT ["/usr/local/bin/plugin"]

FROM --platform=$BUILDPLATFORM docker.io/library/golang:alpine AS build
ARG TARGETOS TARGETARCH
ENV GOOS="$TARGETOS" GOARCH="$TARGETARCH" GOFLAGS="-buildvcs=false -trimpath"
WORKDIR /go/src
COPY . .
RUN --mount=type=cache,target=/go/pkg go mod tidy
RUN --mount=type=cache,target=/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build go build -o /igot -ldflags '-w -s -buildid=' ./cmd/cli

FROM docker.io/library/alpine:latest
RUN apk add ca-certificates tzdata --no-cache
COPY --from=build /igot /igot
USER 1000
ENTRYPOINT [ "/igot" ]
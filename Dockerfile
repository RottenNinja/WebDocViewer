FROM golang:1.22rc1  AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /webdocs

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

#COPY www/docs /docs
COPY www /www
COPY --from=build-stage /webdocs /webdocs

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/webdocs"]

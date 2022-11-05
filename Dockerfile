ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS build
# set working directory
WORKDIR /app

# copy files to container
COPY . .

# download dependencies
RUN go mod download

# compile binary
RUN CGO_ENABLED=0 go build \ 
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo -o x cmd/*

FROM gcr.io/distroless/static AS prod
# set container to run as non root
USER nonroot:nonroot

# copy binary from build stage
COPY --from=build --chown=nonroot:nonroot /app/x .

# run 
ENTRYPOINT ["./x"]
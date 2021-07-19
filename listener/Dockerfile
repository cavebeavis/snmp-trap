# syntax=docker/dockerfile:1
#
# To build:
#           $ docker build -t trap-listener .
#
# https://docs.docker.com/develop/develop-images/multistage-build/
# https://docs.bitnami.com/tutorials/optimize-docker-images-multistage-builds/
#########################################################################################
##############################  Stage 1 -- Golang Builder  ##############################
#########################################################################################
FROM golang:1.16 AS builder

WORKDIR /builder/

COPY . .

RUN go mod init listener
RUN go mod tidy

# https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
# https://pkg.go.dev/cmd/link
# https://lukeeckley.com/post/useful-go-build-flags/
# https://medium.com/a-journey-with-go/go-how-to-take-advantage-of-the-symbols-table-360dd52269e5
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app .

#########################################################################################
##                          Stage 2a -- Alpine Final Container                         ##
#########################################################################################
# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /root/

# COPY --from=builder /builder/app .

# CMD ["./app"]

#
# https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# https://github.com/GoogleContainerTools/distroless
#########################################################################################
##                         Stage 2b -- Scratch Final Container                         ##
#########################################################################################
FROM scratch

WORKDIR /var/www/listener/

COPY --from=builder /builder/app .

# Run the app binary.
ENTRYPOINT ["/var/www/listener/app"]
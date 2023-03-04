FROM alpine:latest as alpine
    RUN apk --no-cache add tzdata zip ca-certificates

FROM golang:1.19-alpine as builder
WORKDIR /go/src/app

    COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
    COPY . .

    RUN go mod init "github.com/DevopsGuyXD/Bizapp"
    RUN go get -u
    RUN go mod tidy
    RUN GOOS="linux" go build .

FROM golang:1.19-alpine
WORKDIR /go/src/app

    ENV ACCESS_KEY_ID=""
    ENV ACCESS_KEY_SECRET=""
    ENV AWS_REGION=""
    ENV GIT_REPO_PRIVATE_ENDPOINT=""
    ENV GIT_USERNAME=""
    ENV GIT_EMAIL_ID=""

    COPY --from=builder /go/src/app/Bizapp .
    COPY --from=builder /go/src/app/Assets ./Assets
    COPY --from=builder /go/src/app/go.mod .
    COPY --from=builder /go/src/app/.env .

    RUN apk update
    RUN apk add --no-cache aws-cli
    RUN apk add github-cli

    EXPOSE 8000

    ENTRYPOINT ["/go/src/app/Bizapp"]
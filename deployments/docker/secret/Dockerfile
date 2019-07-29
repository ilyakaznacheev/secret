FROM golang:alpine

RUN mkdir -p /opt/code/
WORKDIR /opt/code/
ADD ./ /opt/code/

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go get

RUN go build  -o bin/secret main.go

FROM alpine

WORKDIR /app

COPY --from=0 /opt/code/bin/secret /app/

ENTRYPOINT ["./secret"]
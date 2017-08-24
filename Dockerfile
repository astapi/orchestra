FROM golang:1.9-alpine as builder

ADD . /go/src/github.com/astapi/orchestra
WORKDIR /go/src/github.com/astapi/orchestra
RUN apk --update add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

ARG bot_id
ARG bot_token
ARG verification_token
ARG channel_id
ARG deploy_all_services
ARG run_tasks

ENV BOT_ID $bot_id
ENV BOT_TOKEN $bot_token
ENV VERIFICATION_TOKEN $verification_token
ENV CHANNEL_ID $channel_id
ENV DEPLOY_ALL_SERVICES $deploy_all_services
ENV RUN_TASKS $run_tasks

ENV GOPATH /go
RUN go build .

FROM alpine:3.6

COPY --from=builder /go/src/github.com/astapi/orchestra/orchestra /orchestra

RUN apk --update add ca-certificates

WORKDIR /
RUN chown -R nobody:nogroup /orchestra
USER nobody

CMD "./orchestra"

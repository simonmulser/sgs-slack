FROM golang:1.16

RUN go get -u github.com/nlopes/slack
RUN go get github.com/robfig/cron
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/sheets/v4
RUN go get github.com/golang/glog
RUN go get github.com/stretchr/testify/mock

RUN apt-get update && apt-get install -y \
  cron

RUN mkdir /root/.credentials
COPY token.json /root/.credentials

RUN mkdir sgs-slackbot
COPY config sgs-slackbot/config
COPY client_secret.json sgs-slackbot/.
COPY sgs-slackbot sgs-slackbot/.

WORKDIR /go/sgs-slackbot
ENTRYPOINT ["./sgs-slackbot"]

FROM golang:1.6

RUN go get github.com/nlopes/slack
RUN go get github.com/robfig/cron
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/sheets/v4
RUN go get github.com/golang/glog
RUN go get github.com/stretchr/testify/mock

RUN apt-get update && apt-get install -y \
  cron

# Copy the application files (needed for production)
ADD . /go/src/github.com/simonmulser/sgs-slackbot
RUN go install github.com/simonmulser/sgs-slackbot

RUN mkdir /root/.credentials
COPY token.json /root/.credentials

ADD crontab /etc/cron.d/cron-task
RUN chmod 0644 /etc/cron.d/cron-task
RUN touch /var/log/cron.log

CMD cron && tail -f /var/log/cron.log

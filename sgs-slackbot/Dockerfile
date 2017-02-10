FROM golang:1.6

# Create the directory where the application will reside
RUN mkdir /app

# Copy the application files (needed for production)
ADD . /go/src/github.com/simonmulser/sgs-slackbot

# Build the sgs-slackbot command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/nlopes/slack
RUN go get github.com/robfig/cron
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/sheets/v4
RUN go get github.com/golang/glog
RUN go install github.com/simonmulser/sgs-slackbot

RUN apt-get update && apt-get install -y \
  cron

RUN mkdir /root/.credentials
COPY token.json /root/.credentials

# cron stuff
ADD crontab /etc/cron.d/cron-task
RUN chmod 0644 /etc/cron.d/cron-task
RUN touch /var/log/cron.log

CMD cron && tail -f /var/log/cron.log

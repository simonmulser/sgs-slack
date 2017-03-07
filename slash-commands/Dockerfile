FROM golang:1.6

RUN go get github.com/nlopes/slack
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/sheets/v4
RUN go get github.com/golang/glog

RUN mkdir /root/.credentials
COPY token.json /root/.credentials

# Copy the application files (needed for production)
COPY . /go/src/github.com/simonmulser/
RUN go install github.com/simonmulser/slash-commands

EXPOSE 80

ENV VIRTUAL_HOST="sgs-slack.simonmulser.com"
ENV LETSENCRYPT_HOST="sgs-slack.simonmulser.com"
ENV LETSENCRYPT_EMAIL="simonmulser@gmail.com"

WORKDIR /go/src/github.com/simonmulser/slash-commands
CMD ["slash-commands", "production"]

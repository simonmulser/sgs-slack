#!/bin/bash
echo "start sgs-slackbot"
export PATH=$PATH:/go/bin:/usr/local/go/bin
cd /go/src/github.com/simonmulser/sgs/sgs-slackbot/
sgs-slackbot production
echo "sgs-slackbot terminated"

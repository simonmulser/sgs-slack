#!/bin/bash
echo "start post command"
export PATH=$PATH:/go/bin:/usr/local/go/bin
cd /go/src/github.com/simonmulser/sgs/post/
post
echo "finish post command"

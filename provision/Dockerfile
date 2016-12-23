# Copyright 2015 bradley childs, All rights reserved.
#

FROM centos:7
MAINTAINER bradley childs, bchilds@gmail.com
# AWS CLI build deps
RUN yum update -y
RUN yum install -y libcurl-devel unzip
RUN yum install -y less


# Provisioner build deps
RUN yum install -y automake gcc-c++ git make openssl-devel golang

# Build the S3FS provisioner
RUN mkdir -p /root/
WORKDIR  /root/

ENV GOPATH /opt/
ENV PATH $PATH:$GOPATH/bin

RUN mkdir -p /opt/src/github.com/childsb

WORKDIR  /opt/src/github.com/childsb/
RUN git clone https://github.com/childsb/s3fs-container.git
RUN go get github.com/kubernetes-incubator/nfs-provisioner
RUN go get github.com/tools/godep

WORKDIR  /opt/src/github.com/childsb/s3fs-container
RUN make


# Build the AWS CLI and create a bootstrap script
RUN mkdir -p /root/aws_cli
WORKDIR  /root/aws_cli

RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
RUN unzip awscli-bundle.zip
RUN ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

# RUN echo $'#!/bin/sh\n \
#          echo "using aws_access_key_id:"\n \
#          echo $aws_access_key_id\n \
#          echo "aws_secret_access_key:"\n \
#          echo $aws_secret_access_key\n \
#          exec /usr/local/bin/aws "$@"' > /root/aws.sh
# RUN chmod +x /root/aws.sh

# ENTRYPOINT ["/root/aws.sh"]
ENTRYPOINT ["/opt/src/github.com/childsb/s3fs-container/provision"]


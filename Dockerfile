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

# RUN mkdir -p /opt/src/github.com/childsb
# WORKDIR  /opt/src/github.com/childsb/
# ENV GOPATH /opt/
# ENV PATH $PATH:$GOPATH/bin

# install the shell flex script that the provisioner uses.
RUN mkdir -p  /opt/go/src/github.com/childsb/s3fs-container/flex/s3fs-container
COPY flex/s3fs-container/s3fs-container /opt/go/src/github.com/childsb/s3fs-container/flex/s3fs-container/

# RUN git clone https://github.com/childsb/s3fs-container.git
# RUN go get github.com/kubernetes-incubator/nfs-provisioner
# RUN go get github.com/tools/godep

# WORKDIR  /opt/src/github.com/childsb/s3fs-container
# RUN make

# install the go kube piece of provisioner
RUN mkdir -p  /opt/go/src/github.com/childsb/s3fs-container/
COPY s3fs-container /opt/go/src/github.com/childsb/s3fs-container/


# ENTRYPOINT ["/root/aws.sh"]
ENTRYPOINT ["/opt/go/src/github.com/childsb/s3fs-container/s3fs-container"]


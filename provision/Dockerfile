# Copyright 2015 bradley childs, All rights reserved.
#

FROM centos:7
MAINTAINER bradley childs, bchilds@gmail.com
RUN yum update -y 
RUN yum install -y libcurl-devel unzip
RUN yum install -y less 

RUN mkdir -p /root
WORKDIR /root

RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
RUN unzip awscli-bundle.zip
RUN ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

RUN echo $'#!/bin/sh\n \
          echo "using aws_access_key_id:"\n \
          echo $aws_access_key_id\n \
          echo "aws_secret_access_key:"\n \
          echo $aws_secret_access_key\n \
          exec /usr/local/bin/aws "$@"' > /root/aws.sh 
RUN chmod +x /root/aws.sh

ENTRYPOINT ["/root/aws.sh"]

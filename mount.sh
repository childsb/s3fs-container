#!/bin/sh

# export S3User=longgeneratedid
# export S3Secret=longnastygernatedsecret

if [[ $# -lt 2 ]]
then
  echo
  echo "set aws_access_key_id and aws_secret_access_key env then.."
  echo 
  echo "Usage: $0 bucket mountpoint"
  echo 
  echo "Example: $0 snuffy /mnt/snuffy"
  echo
  echo "Docker 1.10 or later required"
  exit
fi

CONTAINER_ID=docker run --privileged -d -e S3User=$aws_access_key
                                        -e S3Secret=$aws_secret_access_key
                                        -v $2:/mnt/mountpoint:shared
                                        --cap-add SYS_ADMIN
                        s3fs $1 /mnt/mountpoint -o passwd_file=/etc/passwd-s3fs -d -d -f -o f2 -o curldbg

echo "Container running as ID ${CONTAINER_ID}"
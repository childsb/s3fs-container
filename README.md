# s3fs-mount-container
This docker container uses the s3fs packge to FUSE mount s3 buckets.  The container is setup to export the s3 mount to host.  Just run the container, and mount your s3 bucket with no extra packages!

# Usage
This container uses a new feature in Docker 1.10 which allows a contiainer to share the hosts mount namespace.  Once docker is up and running you can build the container with b.sh.

You'll need to set your s3 username and secret.  Either modify mount.sh or add the following to your ~/.bashrc:
```bash
export S3User=long_generated_id
export S3Secret=long_gernated_secret
```
To mount a s3 bucket run:
```bash
./mount.sh bucket mountpoint
```

Example:
```bash
./mount.sh snuffy /mnt/snuffy
```
The docker container launches and remains running.  To stop you can ctrl+c.  While the container is running the bucket remains mounted.  You can now access the s3 bucket as a local directory!


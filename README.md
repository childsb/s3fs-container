# s3fs-container
This is an S3 volume driver for Kubernetes which uses the FLEX framework.  It also performs provisioning.

The provisioner runs in a POD.

To build the project:

`make`


The s3 FLEX shell script is in flex/s3fs-container/

To create a storage class:
kubectl create -f provision/sc.yaml

To create the provisioner:
kubectl create -f provision/pod.yaml

To create a claim (which will get provisioend into a volume):
kubectl create -f provision/pvc.yaml


# s3fs-fuse
The container uses s3fs-fuse found here: https://github.com/s3fs-fuse/s3fs-fuse


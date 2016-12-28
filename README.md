# s3fs-container
This is an S3 volume driver for Kubernetes which uses the FLEX framework.  It also performs provisioning.

The provisioner runs in a POD.

To build the project:

`make`


The s3 FLEX shell script is in flex/s3fs-container/

To create a storage class:
```bash
kubectl create -f provision/sc.yaml
```
To create the provisioner:
```bash
kubectl create -f provision/pod.yaml
```
To create a claim (which will get provisioend into a volume):
```bash
kubectl create -f provision/pvc.yaml
```

To create an application that uses the claim:

```bash
kubectl create -f provision/pod-application.yaml
```

#Flex driver location
The FLEX driver location is mostly hard coded.  I opened this PR to allow it specified in hack/local-up-cluster.sh

https://github.com/kubernetes/kubernetes/pull/39054

Either copy those changes in yourself, or copy the script into the location which FLEX is expecting the driver.  If using the PR, run hack/local-up-cluster with something like:
```bash
FLEX_VOLUME_PLUGIN_DIR=/opt/go/src/github.com/childsb/s3fs-container/flex hack/local-up-cluster.sh
```


# s3fs-fuse
The container uses s3fs-fuse found here: https://github.com/s3fs-fuse/s3fs-fuse


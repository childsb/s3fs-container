# s3fs-container
This is an S3 volume driver for Kubernetes which uses the FLEX framework.  It also performs dynamic provisioning. The provisioner runs in a POD.

**build the project:**

```bash
mkdir -p $GOPATH/src/github.com/kubernetes-incubator/
cd $GOPATH/src/github.com/kubernetes-incubator/
git clone git@github.com:childsb/s3fs-container.git
cd s3fs-container
make
```
The container images need to be accessible from every node.  Also the FLEX script needs to be installed to proper location 

#Flex driver location
The FLEX driver location is mostly hard coded.  I opened this PR to allow it specified in hack/local-up-cluster.sh

https://github.com/kubernetes/kubernetes/pull/39054

If using hack/local-up-cluster.sh

Either copy those changes in yourself, or copy the script into the location which FLEX is expecting the driver.  If using the PR, run hack/local-up-cluster with something like:
```bash
FLEX_VOLUME_PLUGIN_DIR=/opt/go/src/github.com/childsb/s3fs-container/flex hack/local-up-cluster.sh
```


If you're running from something else, be sure to set the volume plugin path, or use the default path.  

**Copy** 
flex/s3fs-container/s3fs-container to <volume_plugin_path>/s3fs-container/s3fs-container

So that kube can mount the S3 volumes with the FLEX script.


# Run the S3fs provisioner

To create the provisioner:
```bash
kubect
```

To create a storage class:
```bash
kubectl create -f provision/sc.yaml
```
l create -f provision/pod.yaml
```
To create a claim (which will get provisioend into a volume):
```bash
kubectl create -f provision/pvc.yaml
```

To create an application that uses the claim:

```bash
kubectl create -f provision/pod-application.yaml
```





# s3fs-fuse
The container uses s3fs-fuse found here: https://github.com/s3fs-fuse/s3fs-fuse


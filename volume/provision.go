/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package volume

import (

	"github.com/golang/glog"
	"os/exec"
	"github.com/kubernetes-incubator/nfs-provisioner/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/types"
)

const (
	// Name of the file where an s3fsProvisioner will store its identity
	identityFile = "s3fs-provisioner.identity"

	// are we allowed to set this? else make up our own
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "s3fs-dynamic-provisioner"

	annAwsAccessKeyId = "AWS_ACCESS_KEY_ID"
	annAwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	annAwss3bucket = "bucket"

	// A PV annotation for the identity of the s3fsProvisioner that provisioned it
	annProvisionerId = "Provisioner_Id"


)

func News3FSProvisioner(client kubernetes.Interface, execCommand string) controller.Provisioner {
	return newS3fsProvisionerInternal(client, execCommand)
}

func newS3fsProvisionerInternal(client kubernetes.Interface, execCommand string) *s3fsProvisioner {
	var identity types.UID


	provisioner := &s3fsProvisioner{
		client:       client,
		execCommand:	execCommand,
		identity:     identity,
	}

	return provisioner
}

type s3fsProvisioner struct {
	client kubernetes.Interface
	execCommand string
	identity types.UID
}

var _ controller.Provisioner = &s3fsProvisioner{}

// Provision creates a volume i.e. the storage asset and returns a PV object for
// the volume.
func (p *s3fsProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	claim := options.PVC

	s3bucket, err := p.createVolume(options,claim)

	if err != nil {
		return nil, err
	}

	annotations := make(map[string]string)
	annotations[annCreatedBy] = createdBy

	annotations[annProvisionerId] = string(p.identity)

	pv := &v1.PersistentVolume{
		ObjectMeta: v1.ObjectMeta{
			Name:        options.PVName,
			Labels:      map[string]string{},
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                  options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{

				FlexVolume: &v1.FlexVolumeSource{
					Driver: "s3fs-container",
					Options: map[string]string{
						annAwsAccessKeyId:claim.Annotations[annAwsAccessKeyId],
						annAwsSecretAccessKey:claim.Annotations[annAwsSecretAccessKey],
						annAwss3bucket:s3bucket,
					},

					ReadOnly: false,
				},
			},
		},
	}
	glog.Infof("Created PV %s", options.PVName)
	return pv, nil
}

func (p *s3fsProvisioner) createVolume(volumeOptions controller.VolumeOptions, claim *v1.PersistentVolumeClaim) (string, error) {
	s3bucket := claim.Annotations[annAwss3bucket]

	if len(s3bucket)==0{
		s3bucket = volumeOptions.PVName
	}

	cmd := exec.Command(p.execCommand, "provision", s3bucket, claim.Annotations[annAwsAccessKeyId], claim.Annotations[annAwsSecretAccessKey] )
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("Failed to create volume %s, output: %s, error: %s",  s3bucket, output, err.Error())
		//_, err := handleCmdResponse(mountCmd, output)
		return "", err
	}
	glog.Infof("Created s3 bucket %s", s3bucket)
	return s3bucket, nil

}

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
	"encoding/json"
	"github.com/golang/glog"
	"fmt"
	//"io/ioutil"
	//"os"
	"os/exec"
	//"path"
	"strconv"
	"strings"
	//"syscall"

	//"github.com/golang/glog"
	"github.com/childsb/s3fs-container/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/types"
	//"k8s.io/client-go/pkg/util/uuid"

)

const (
	// Name of the file where an s3fsProvisioner will store its identity
	identityFile = "s3fs-provisioner.identity"

	// are we allowed to set this? else make up our own
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "s3fs-dynamic-provisioner"

	// A PV annotation for the entire ganesha EXPORT block or /etc/exports
	// block, needed for deletion.
	// annExportBlock = "EXPORT_block"
	// A PV annotation for the exportId of this PV's backing ganesha/kernel export
	// , needed for ganesha deletion and used for deleting the entry in exportIds
	// map so the id can be reassigned.
	//annExportId = "Export_Id"

	// A PV annotation for the project quota info block, needed for quota
	// deletion.
	//annProjectBlock = "Project_block"
	// A PV annotation for the project quota id, needed for quota deletion
	//annProjectId = "Project_Id"

	// VolumeGidAnnotationKey is the key of the annotation on the PersistentVolume
	// object that specifies a supplemental GID.
	VolumeGidAnnotationKey = "pv.beta.kubernetes.io/gid"

	// A PV annotation for the identity of the s3fsProvisioner that provisioned it
	annProvisionerId = "Provisioner_Id"

	podIPEnv     = "POD_IP"
	serviceEnv   = "SERVICE_NAME"
	namespaceEnv = "POD_NAMESPACE"
	nodeEnv      = "NODE_NAME"
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
		podIPEnv:     podIPEnv,
		serviceEnv:   serviceEnv,
		namespaceEnv: namespaceEnv,
		nodeEnv:      nodeEnv,
	}

	return provisioner
}

type s3fsProvisioner struct {

	// Client, needed for getting a service cluster IP to put as the S3FS server of
	// provisioned PVs
	client kubernetes.Interface
	execCommand string
// Identity of this s3fsProvisioner, generated & persisted to exportDir or
	// recovered from there. Used to mark provisioned PVs
	identity types.UID

	// Environment variables the provisioner pod needs valid values for in order to
	// put a service cluster IP as the server of provisioned S3FS PVs, passed in
	// via downward API. If serviceEnv is set, namespaceEnv must be too.
	podIPEnv     string
	serviceEnv   string
	namespaceEnv string
	nodeEnv      string
}

var _ controller.Provisioner = &s3fsProvisioner{}

// Provision creates a volume i.e. the storage asset and returns a PV object for
// the volume.
func (p *s3fsProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	err := p.createVolume(options)
	if err != nil {
		return nil, err
	}

	annotations := make(map[string]string)
	annotations[annCreatedBy] = createdBy
	//if supGroup != 0 {
	//	annotations[VolumeGidAnnotationKey] = strconv.FormatUint(supGroup, 10)
//	}
	annotations[annProvisionerId] = string(p.identity)

	pv := &v1.PersistentVolume{
		ObjectMeta: v1.ObjectMeta{
			Name:        options.PVName,
			Labels:      map[string]string{},
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.Capacity,
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{

				FlexVolume: &v1.FlexVolumeSource{
					Driver: "s3fs",
					Options: map[string]string{
						"AWS_ACCESS_KEY_ID":"",
						"AWS_SECRET_ACCESS_KEY":"",
						"bucket":"",
					},

					ReadOnly: false,
				},
			},
		},
	}

	return pv, nil
}


func (p *s3fsProvisioner) validateOptions(options controller.VolumeOptions) (string, error) {
	gid := "none"
	for k, v := range options.Parameters {
		switch strings.ToLower(k) {
		case "gid":
			if strings.ToLower(v) == "none" {
				gid = "none"
			} else if i, err := strconv.ParseUint(v, 10, 64); err == nil && i != 0 {
				gid = v
			} else {
				return "", fmt.Errorf("invalid value for parameter gid: %v. valid values are: 'none' or a non-zero integer", v)
			}
		default:
			return "", fmt.Errorf("invalid parameter: %q", k)
		}
	}

	// TODO implement options.ProvisionerSelector parsing
	// pv.Labels MUST be set to match claim.spec.selector
	// gid selector? with or without pv annotation?
	if options.Selector != nil {
		return "", fmt.Errorf("claim.Spec.Selector is not supported")
	}

	return gid, nil
}

func (p *s3fsProvisioner) createVolume(volumeOptions controller.VolumeOptions) ( error) {
	gid, err := p.validateOptions(volumeOptions)
	if err != nil {
		return fmt.Errorf("error validating options for volume: %v", err)
	}

	glog.Infof("createVolume called..%v %v", volumeOptions, gid)

	var options string
	wildBill, err := json.Marshal(volumeOptions)
	if err != nil {
		glog.Errorf("Failed to marshal plugin options, error: %s", err.Error())
		return err
	}
	if len(wildBill) != 0 {
		options = string(wildBill)
	} else {
		options = ""
	}

	// Executable provider command.
	cmd := exec.Command(p.execCommand, options)
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("Failed to mount volume %s, output: %s, error: %s",  output, err.Error())
		//_, err := handleCmdResponse(mountCmd, output)
		return err
	}


	return nil

	/*
		server, err := p.getServer()

		if err != nil {
			return "", "", 0, "", 0, "", 0, fmt.Errorf("error getting S3FS server IP for volume: %v", err)
		}

		path := path.Join(p.exportDir, options.PVName)

		err = p.createDirectory(options.PVName, gid)
		if err != nil {
			return "", "", 0, "", 0, "", 0, fmt.Errorf("error creating directory for volume: %v", err)
		}

		exportBlock, exportId, err := p.createExport(options.PVName)
		if err != nil {
			os.RemoveAll(path)
			return "", "", 0, "", 0, "", 0, fmt.Errorf("error creating export for volume: %v", err)
		}

		projectBlock, projectId, err := p.createQuota(options.PVName, options.Capacity)
		if err != nil {
			os.RemoveAll(path)
			return "", "", 0, "", 0, "", 0, fmt.Errorf("error creating quota for volume: %v", err)
		}

		return server, path, 0, exportBlock, exportId, projectBlock, projectId, nil
		*/
}

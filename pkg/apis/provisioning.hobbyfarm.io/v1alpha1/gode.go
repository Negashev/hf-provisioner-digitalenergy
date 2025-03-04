package v1alpha1

import (
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/kvmx86"
)

func (d *InstanceCreateRequest) ToGode() kvmx86.CreateRequest {

	return kvmx86.CreateRequest{
		RGID:     d.RGID,
		Name:     d.Name,
		CPU:      d.CPU,
		RAM:      d.RAM,
		ImageID:  d.ImageID,
		BootDisk: d.BootDiskSize,
		Start:    true,
	}
}

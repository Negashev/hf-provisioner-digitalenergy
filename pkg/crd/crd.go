package crd

import (
	"github.com/ebauman/crder"
	provisioning_hobbyfarm_io "github.com/negashev/hf-provisioner-digitalenergy/pkg/apis/provisioning.hobbyfarm.io"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
)

func Setup() []crder.CRD {
	instance := crder.NewCRD(v1alpha1.Instance{}, provisioning_hobbyfarm_io.Group, func(c *crder.CRD) {
		c.WithShortNames("basis")
		c.IsNamespaced(true)
		c.AddVersion(v1alpha1.Version, v1alpha1.Instance{}, func(cv *crder.Version) {
			cv.IsStored(true).IsServed(true)
			cv.WithStatus()
			cv.WithPreserveUnknown()
		})
	})

	key := crder.NewCRD(v1alpha1.Key{}, provisioning_hobbyfarm_io.Group, func(c *crder.CRD) {
		c.IsNamespaced(true)
		c.AddVersion(v1alpha1.Version, v1alpha1.Key{}, func(cv *crder.Version) {
			cv.IsStored(true).IsServed(true)
			cv.WithStatus()
			cv.WithPreserveUnknown()
		})
	})

	return []crder.CRD{
		*instance,
		*key,
	}
}

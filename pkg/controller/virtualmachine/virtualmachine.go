package virtualmachine

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/errors"
	labels2 "github.com/negashev/hf-provisioner-digitalenergy/pkg/labels"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

func VMLabelSelector(vmName string) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		labels2.VirtualMachineLabel: vmName,
	})
}

func ProvisionerFinalizer(req router.Request, resp router.Response) error {
	// before deleting the VM, make sure the droplet and key are gone
	instance, err := GetInstance(req)
	// if the droplet is not found, move on. anything else, report!
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("error fetching instance: %s", err.Error())
	}

	if instance != nil {
		// instance exists, delete it
		err = req.Client.Delete(req.Ctx, instance)
		if err != nil {
			return fmt.Errorf("error deleting instance: %s", err.Error())
		}
		resp.RetryAfter(5 * time.Second)
	}

	return nil
}

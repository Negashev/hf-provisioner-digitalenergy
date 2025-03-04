package virtualmachine

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/config"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/errors"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/gode"
	labels2 "github.com/negashev/hf-provisioner-digitalenergy/pkg/labels"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/kvmx86"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func InstanceHandler(req router.Request, resp router.Response) error {
	obj := req.Object.(*v1.VirtualMachine)

	var k8sInstance *v1alpha1.Instance
	k8sInstance, err := GetInstance(req)

	//_, err = client.CloudAPI().Compute().List(ctx, compute.ListRequest{Name: name})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//if errors.IsNotFound(err) {
	name := fmt.Sprintf("%s", obj.Name)

	instance := v1alpha1.Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: obj.Namespace,
			Labels: map[string]string{
				labels2.VirtualMachineLabel: obj.Name,
			},
		},
		Spec: v1alpha1.InstanceSpec{
			Machine: obj.Name,
		},
	}
	dcr := buildInstanceCreateRequest(name, obj, req)

	client, err := gode.GetGodeClient(obj.Name, req)
	if err != nil {
		log.Fatal(err)
	}

	getInstance, err := gode.GetOrCreateInstance(client, name, req, dcr)
	dcrJson, err := json.Marshal(getInstance)
	if err != nil {
		return fmt.Errorf("error marshalling json: %s", err.Error())
	}
	instance.Spec.Instance.Raw = dcrJson
	//if createRequest {
	k8sInstance = &instance
	//}
	//}

	resp.Objects(k8sInstance)

	return nil
}

func GetInstance(req router.Request) (*v1alpha1.Instance, error) {
	instanceList := &v1alpha1.InstanceList{}
	err := req.List(instanceList, &kclient.ListOptions{
		Namespace:     req.Object.GetNamespace(),
		LabelSelector: VMLabelSelector(req.Object.GetName()),
	})

	if err != nil {
		return nil, err
	}

	if len(instanceList.Items) > 0 {
		return &instanceList.Items[0], nil
	}

	return nil, errors.NewNotFoundError("could not find any instance for virtualmachine %s", req.Object.GetName())
}

func buildInstanceCreateRequest(name string, vm *v1.VirtualMachine, req router.Request) kvmx86.CreateRequest {
	// set network
	Interfaces := []kvmx86.Interface{}
	EXTNET := config.ResolveConfigInt(vm, req, "EXTNET")
	if EXTNET != 0 {
		Interfaces = append(Interfaces, kvmx86.Interface{NetType: "EXTNET", NetID: EXTNET})
	}
	VINS := config.ResolveConfigInt(vm, req, "VINS")
	if VINS != 0 {
		Interfaces = append(Interfaces, kvmx86.Interface{NetType: "VINS", NetID: VINS})
	}
	return kvmx86.CreateRequest{
		RGID:       config.ResolveConfigInt(vm, req, "RGID"),
		Name:       name,
		CPU:        config.ResolveConfigInt(vm, req, "CPU"),
		RAM:        config.ResolveConfigInt(vm, req, "RAM"),
		ImageID:    config.ResolveConfigInt(vm, req, "ImageID"),
		BootDisk:   config.ResolveConfigInt(vm, req, "BootDisk"),
		Start:      true,
		Interfaces: Interfaces,
	}
}

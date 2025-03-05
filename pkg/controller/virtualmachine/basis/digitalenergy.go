package basis

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/gode"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/compute"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const EnsureDeletedFinalizer = "provisioner.hobbyfarm.io/digitalenergy-deleted"

func PeriodicUpdate(req router.Request, resp router.Response) error {
	k8sInstance := req.Object.(*v1alpha1.Instance)

	// get the old gode instance from k8s
	var godeInstancr = compute.RecordCompute{}
	if k8sInstance.Status.Instance.Raw == nil {
		return nil
	}
	err := json.Unmarshal(k8sInstance.Status.Instance.Raw, &godeInstancr)
	if err != nil {
		logrus.Errorf("Unable to unmarshal Digital Energy instance from k8s instance object: %s", err.Error())
		return nil
	}
	// get an updated instance object from DE
	dClient, err := gode.GetGodeClient(k8sInstance.Spec.Machine, req)
	GetRequest := compute.GetRequest{ComputeID: godeInstancr.ID}
	newGodeDInstance, err := dClient.CloudAPI().Compute().Get(req.Ctx, GetRequest)
	if err != nil {
		v1alpha1.ConditionInstanceUpdated.False(k8sInstance)
		v1alpha1.ConditionInstanceUpdated.SetError(k8sInstance, "Digital Energy error", err)
		return req.Client.Status().Update(req.Ctx, k8sInstance)
	}

	jsonInstance, err := json.Marshal(newGodeDInstance)
	if err != nil {
		logrus.Errorf("Unable to marshal Digital Energy instance to []byte for storage in k8s: %s", err.Error())
		return nil
	}

	k8sInstance.Status.Instance.Raw = jsonInstance
	switch newGodeDInstance.Status {
	case "ANOTNER????":
		resp.RetryAfter(10 * time.Second)
	case "LOCKED":
		resp.RetryAfter(30 * time.Second)
	}

	return req.Client.Status().Update(req.Ctx, k8sInstance)
}

func UpdateVMStatus(req router.Request, resp router.Response) error {
	k8sInstance := req.Object.(*v1alpha1.Instance)

	instance := compute.RecordCompute{}
	if k8sInstance.Status.Instance.Raw == nil {
		return nil
	}
	err := json.Unmarshal(k8sInstance.Status.Instance.Raw, &instance)
	if err != nil {
		logrus.Errorf("UpdateVMStatus Unable to unmarshal Digital Energy instance from k8s instance object: %s", err.Error())
		return nil
	}

	// get the corresponding VM for the instance
	vm := &v1.VirtualMachine{}
	err = req.Client.Get(req.Ctx, client.ObjectKey{Namespace: k8sInstance.Namespace, Name: k8sInstance.Spec.Machine}, vm)
	if err != nil {
		// could not get VM
		logrus.Errorf("Could not get VirtualMachine with name %s: %s", k8sInstance.Spec.Machine, err.Error())
		return nil
	}

	switch instance.TechStatus {
	case "SCHEDULED":
		vm.Status.Status = v1.VmStatusProvisioned
	case "STARTED":
		vm.Status.Status = v1.VmStatusRunning
	case "STOPPED":
		vm.Status.Status = v1.VmStatusTerminating
	}

	for _, n := range instance.Interfaces {
		switch n.NetType {
		case "EXTNET":
			vm.Status.PublicIP = n.IPAddress
		case "VINS":
			vm.Status.PrivateIP = n.IPAddress
		}
	}
	varValue := os.Getenv("USE_PRIVATE_IP")
	if varValue != "" && vm.Status.PublicIP == "" {
		vm.Status.PublicIP = vm.Status.PrivateIP
	}

	vm.Status.Hostname = instance.Name // usually?

	return req.Client.Status().Update(req.Ctx, vm)
}

func EnsureDeleted(req router.Request, resp router.Response) error {
	instance := req.Object.(*v1alpha1.Instance)

	var godeInstance = compute.RecordCompute{}
	if instance.Status.Instance.Raw == nil {
		return nil
	}
	if err := json.Unmarshal(instance.Status.Instance.Raw, &godeInstance); err != nil || godeInstance.ID == 0 {
		logrus.Errorf("Could not obtain DE instance from object JSON. This is unrecoverable, may result in "+
			"an orphan instance in Digital Energy. Instance (k8s) name was %s", instance.Name)
		return nil
	}

	// instance exists, delete it
	dClient, err := gode.GetGodeClient(instance.Spec.Machine, req)
	if err != nil {
		logrus.Errorf("building digital energy client: %s", err.Error())
		return req.Client.Status().Update(req.Ctx, instance)
	}

	DeleteRequest := compute.DeleteRequest{ComputeID: godeInstance.ID, Permanently: true}
	_, err = dClient.CloudAPI().Compute().Delete(req.Ctx, DeleteRequest)
	if err != nil {
		logrus.Errorf("deleting instance in digital energy: %s", err.Error())
		return req.Client.Status().Update(req.Ctx, instance)
	}

	return nil
}

func EnsureStatus(req router.Request, resp router.Response) error {
	instance := req.Object.(*v1alpha1.Instance)
	if len(instance.Status.Conditions) == 0 {
		v1alpha1.ConditionInstanceExists.SetStatus(instance, "unknown")
		v1alpha1.ConditionConnectionReady.SetStatus(instance, "unknown")

		return req.Client.Status().Update(req.Ctx, instance)
	}

	return nil
}

func InstanceNotCreated(next router.Handler) router.Handler {
	return router.HandlerFunc(func(req router.Request, resp router.Response) error {
		instance := req.Object.(*v1alpha1.Instance)

		if v1alpha1.ConditionInstanceExists.GetStatus(instance) == "unknown" {
			return next.Handle(req, resp)
		}

		return nil
	})
}

func CreateInstance(req router.Request, _ router.Response) error {
	instance := req.Object.(*v1alpha1.Instance)

	var dcr v1alpha1.InstanceCreateRequest
	if instance.Spec.Instance.Raw == nil {
		return nil
	}
	if err := json.Unmarshal(instance.Spec.Instance.Raw, &dcr); err != nil {
		return fmt.Errorf("error unmarshalling instance: %s", err.Error())
	}
	dClient, err := gode.GetGodeClient(instance.Spec.Machine, req)
	if err != nil {
		v1alpha1.ConditionInstanceExists.SetStatus(instance, "false")
		v1alpha1.ConditionInstanceExists.SetError(instance, "error creating digital energy client", err)
		return req.Client.Status().Update(req.Ctx, instance)
	}
	// Могут быть проблемы с сетью при создании новой ВМ
	getInstance, err := gode.GetOrCreateInstance(dClient, instance.Name, req, dcr.ToGode(), "")
	if err != nil {
		v1alpha1.ConditionInstanceExists.SetStatus(instance, "false")
		v1alpha1.ConditionInstanceExists.SetError(instance, "digital energy error", err)
		return req.Client.Status().Update(req.Ctx, instance)
	}

	instanceJson, err := json.Marshal(getInstance)
	if err != nil {
		return fmt.Errorf("error marshalling instance: %s", err.Error())
	}
	instance.Status.Instance.Raw = instanceJson
	v1alpha1.ConditionInstanceExists.SetStatus(instance, "true")
	v1alpha1.ConditionInstanceExists.Reason(instance, "instance created")

	return req.Client.Status().Update(req.Ctx, instance)
}

package controller

import (
	"github.com/acorn-io/baaah/pkg/router"
	hfv1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/controller/virtualmachine"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/controller/virtualmachine/basis"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/controller/virtualmachine/key"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/controller/virtualmachine/secret"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/labels"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/namespace"
)

func routes(r *router.Router) {
	ns := namespace.Resolve()

	//vmRouter := r.Type(&hfv1.VirtualMachine{}).Namespace(ns).Selector(klabels.SelectorFromSet(map[string]string{
	//	labels.ProvisionerLabel: providerregistration.ProviderName(),
	//}))

	vmRouter := r.Type(&hfv1.VirtualMachine{}).Namespace(ns)

	vmRouter.FinalizeFunc(labels.Finalizer, virtualmachine.ProvisionerFinalizer)
	vmRouter.HandlerFunc(secret.SecretHandler)
	vmRouter.Middleware(secret.RequireSecret).HandlerFunc(virtualmachine.KeyHandler)
	vmRouter.Middleware(virtualmachine.RequireKey).HandlerFunc(virtualmachine.InstanceHandler)

	keyRouter := r.Type(&v1alpha1.Key{}).Namespace(namespace.Resolve())

	keyRouter.HandlerFunc(key.EnsureStatus)
	keyRouter.Middleware(key.NotYetCreated).HandlerFunc(key.CreateKey)
	keyRouter.Middleware(key.Created).HandlerFunc(key.WriteVM)
	keyRouter.FinalizeFunc(key.EnsureDeletedFinalizer, key.EnsureDeleted)

	dropletRouter := r.Type(&v1alpha1.Instance{}).Namespace(ns)

	dropletRouter.HandlerFunc(basis.EnsureStatus)
	dropletRouter.Middleware(basis.InstanceNotCreated).HandlerFunc(basis.CreateInstance)
	dropletRouter.HandlerFunc(basis.PeriodicUpdate)
	dropletRouter.HandlerFunc(basis.UpdateVMStatus)
	dropletRouter.FinalizeFunc(basis.EnsureDeletedFinalizer, basis.EnsureDeleted)

}

package v1alpha1

import "github.com/negashev/hf-provisioner-digitalenergy/pkg/retries"

func (k *Key) GetRetries() []retries.GenericRetry {
	return k.Status.Retries
}

func (k *Key) SetRetries(retries []retries.GenericRetry) {
	k.Status.Retries = retries
}

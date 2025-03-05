package gode

import (
	"encoding/json"
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	config "github.com/negashev/hf-provisioner-digitalenergy/pkg/config"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/controller/virtualmachine/secret"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/namespace"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"log"
	decort "repository.basistech.ru/BASIS/decort-golang-sdk"
	decortConfig "repository.basistech.ru/BASIS/decort-golang-sdk/config"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/compute"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/kvmx86"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetGodeClient(vmName string, req router.Request) (*decort.DecortClient, error) {
	// lookup the de token
	app_id := config.ResolveConfigItemName(vmName, req, "app_id")
	app_secret := config.ResolveConfigItemName(vmName, req, "app_secret")
	sso_url := config.ResolveConfigItemName(vmName, req, "sso_url")
	decort_url := config.ResolveConfigItemName(vmName, req, "decort_url")
	ssl_skip_verify := config.ResolveConfigItemName(vmName, req, "ssl_skip_verify")
	if app_id == "" || app_secret == "" {
		// check for token secret?
		secretFile := config.ResolveConfigItemName(vmName, req, "secret-file")
		if secretFile == "" {
			return nil, fmt.Errorf("unable to resolve app_id/app_secret for digital energy api")
		}

		secret := &v1.Secret{}
		err := req.Client.Get(req.Ctx, client.ObjectKey{
			Namespace: namespace.Resolve(),
			Name:      secretFile,
		}, secret)
		if err != nil {
			return nil, fmt.Errorf("error retrieving token secret: %s", err.Error())
		}

		app_id = string(secret.Data["app_id"])
		app_secret = string(secret.Data["app_secret"])
		sso_url = string(secret.Data["sso_url"])
		decort_url = string(secret.Data["decort_url"])
		ssl_skip_verify = string(secret.Data["ssl_skip_verify"])
	}
	var SSLSkipVerify = false
	if ssl_skip_verify == "true" {
		SSLSkipVerify = true
	}
	// Настройка конфигурации
	cfg := decortConfig.Config{
		AppID:         app_id,
		AppSecret:     app_secret,
		SSOURL:        sso_url,
		DecortURL:     decort_url,
		Retries:       5,
		SSLSkipVerify: SSLSkipVerify,
	}

	// Создание клиента
	return decort.New(cfg), nil
}

func GetOrCreateInstance(dClient *decort.DecortClient, instanceName string, req router.Request, dcr kvmx86.CreateRequest, cloud_config string) (*compute.RecordCompute, error) {

	getInstance := &compute.RecordCompute{}

	// провеить что сервер существует
	FindServer, err := dClient.CloudAPI().Compute().List(req.Ctx, compute.ListRequest{Name: instanceName})

	var res = uint64(0)
	if FindServer.EntryCount > 0 {
		res = FindServer.FilterByName(instanceName).Data[0].ID
	} else {
		// add cloud_config
		var yamlMap map[string]interface{}
		if err := yaml.Unmarshal([]byte(cloud_config), &yamlMap); err != nil {
			log.Fatal(err)
		}

		// set cloud init with ssh auth
		secret, err := secret.GetSecret(req)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		public_key_from_secret := string(secret.Data["public_key"])
		public_key := public_key_from_secret[:len(public_key_from_secret)-1]

		// Загрузка второго JSON
		anotherJSON := `{
				"users": [
					{
						"lock-passwd": false,
						"name": "user",
						"ssh-authorized-keys": "` + public_key + `",
						"shell": "/bin/bash",
						"sudo": "ALL=(ALL) NOPASSWD:ALL"
					}
				]
			}`
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(anotherJSON), &jsonMap); err != nil {
			panic(err)
		}
		// Объединение JSON объектов (YAML имеет приоритет)
		merged := mergeMaps(jsonMap, yamlMap)
		// Результат
		result, _ := json.MarshalIndent(merged, "", "  ")
		dcr.Userdata = string(result)
		res, err = dClient.CloudAPI().KVMX86().Create(req.Ctx, dcr)
		if err != nil {
			log.Fatal(err)
		}
	}
	getInstance, err = dClient.CloudAPI().Compute().Get(req.Ctx, compute.GetRequest{ComputeID: res})
	if err != nil {
		return nil, err
	}
	return getInstance, nil
}

// Рекурсивное объединение карт с приоритетом второй карты
func mergeMaps(base, override map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	// Копируем базовые значения
	for k, v := range base {
		res[k] = v
	}

	// Мерджим с переопределениями
	for k, overrideVal := range override {
		if baseVal, exists := res[k]; exists {
			// Проверяем вложенные карты
			baseMap, baseIsMap := baseVal.(map[string]interface{})
			overrideMap, overrideIsMap := overrideVal.(map[string]interface{})
			if baseIsMap && overrideIsMap {
				res[k] = mergeMaps(baseMap, overrideMap)
				continue
			}
		}
		// Перезаписываем значение
		res[k] = overrideVal
	}
	return res
}

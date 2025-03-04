package gode

import (
	"bytes"
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	config "github.com/negashev/hf-provisioner-digitalenergy/pkg/config"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/controller/virtualmachine/secret"
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/namespace"
	"golang.org/x/crypto/ssh"
	v1 "k8s.io/api/core/v1"
	"log"
	decort "repository.basistech.ru/BASIS/decort-golang-sdk"
	decortConfig "repository.basistech.ru/BASIS/decort-golang-sdk/config"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/compute"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/kvmx86"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
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

func GetOrCreateInstance(dClient *decort.DecortClient, instanceName string, req router.Request, dcr kvmx86.CreateRequest) (*compute.RecordCompute, error) {

	getInstance := &compute.RecordCompute{}

	// провеить что сервер существует
	FindServer, err := dClient.CloudAPI().Compute().List(req.Ctx, compute.ListRequest{Name: instanceName})

	var res = uint64(0)
	if FindServer.EntryCount > 0 {
		res = FindServer.FilterByName(instanceName).Data[0].ID
	} else {
		// set cloud init with ssh auth
		secret, err := secret.GetSecret(req)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		public_key_from_secret := string(secret.Data["public_key"])
		public_key := public_key_from_secret[:len(public_key_from_secret)-1]
		dcr.Userdata = `{
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
		res, err = dClient.CloudAPI().KVMX86().Create(req.Ctx, dcr)
		if err != nil {
			log.Fatal(err)
		}
	}
	getInstance, err = dClient.CloudAPI().Compute().Get(req.Ctx, compute.GetRequest{ComputeID: res})
	if err != nil {
		return nil, err
	}
	// START write ssh key
	//secret, err := secret.GetSecret(req)
	//if err != nil {
	//	log.Fatal(err)
	//	return nil, err
	//}
	//
	//public_key_from_secret := string(secret.Data["public_key"])
	//public_key := public_key_from_secret[:len(public_key_from_secret)-1]
	//
	//config := &SSHConfig{
	//	Host:     getInstance.Interfaces[0].IPAddress,
	//	Port:     22,
	//	User:     getInstance.OSUsers[0].Login,
	//	Password: getInstance.OSUsers[0].Password,
	//	Cmd:      "echo '" + public_key + "' > /home/" + getInstance.OSUsers[0].Login + "/.ssh/authorized_keys",
	//}
	//
	//if err := sshConnect(config); err != nil {
	//	panic(err)
	//}
	// STOP
	return getInstance, nil
}

type SSHConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Cmd      string
}

func sshConnect(config *SSHConfig) error {
	// Создаем SSH-конфигурацию
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Не использовать в продакшене!
	}

	var client *ssh.Client
	var err error

	// Попытки подключения с ожиданием
	for attempts := 0; attempts < 100; attempts++ {
		client, err = ssh.Dial("tcp",
			fmt.Sprintf("%s:%d", config.Host, config.Port),
			sshConfig,
		)
		if err == nil {
			break
		}
		fmt.Printf("Подключение не установлено. Повторная попытка через 2 сек... %v\n", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("не удалось установить соединение: %v", err)
	}
	defer client.Close()

	// Создаем сессию
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("не удалось создать сессию: %v", err)
	}
	defer session.Close()

	// Выполняем команду
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(config.Cmd)
	if err != nil {
		return fmt.Errorf("ошибка выполнения команды: %v\nSTDERR: %s", err, stderr.String())
	}

	fmt.Printf("Результат выполнения команды:\nSTDOUT:\n%s\nSTDERR:\n%s", stdout.String(), stderr.String())
	return nil
}

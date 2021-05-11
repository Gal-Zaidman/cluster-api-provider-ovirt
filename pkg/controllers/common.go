package ovirt

import (
	"fmt"
	"github.com/openshift/cluster-api-provider-ovirt/pkg/clients/ovirt"

	"github.com/go-logr/logr"
	ovirtsdk "github.com/ovirt/go-ovirt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NAMESPACE          = "openshift-machine-api"
	CREDENTIALS_SECRET = "ovirt-credentials"
)

type BaseController struct {
	Log             logr.Logger
	Client          client.Client
	ovirtConnection *ovirtsdk.Connection
}

func (b *BaseController) GetConnection(namespace, secretName string) (*ovirtsdk.Connection, error) {
	var err error
	if b.ovirtConnection == nil || b.ovirtConnection.Test() != nil {
		// session expired or some other error, re-login.
		b.ovirtConnection, err = createApiConnection(b.Client, namespace, secretName)
	}
	return b.ovirtConnection, err
}

//createApiConnection returns a a client to oVirt's API endpoint
func createApiConnection(client client.Client, namespace string, secretName string) (*ovirtsdk.Connection, error) {
	creds, err := ovirt.GetCredentialsSecret(client, namespace, secretName)

	if err != nil {
		return nil, fmt.Errorf("failed getting credentials for namespace %s, %s", namespace, err)
	}

	connection, err := ovirtsdk.NewConnectionBuilder().
		URL(creds.URL).
		Username(creds.Username).
		Password(creds.Password).
		CAFile(creds.CAFile).
		Insecure(creds.Insecure).
		Build()
	if err != nil {
		return nil, err
	}

	return connection, nil
}

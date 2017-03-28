package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubectld/helm"
)

func installCatalog(event *events.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	stack := decodeHelmStack(event, cli, false)
	notes, err := helm.InstallHelmStack(stack)
	if err != nil {
		log.Errorf("Error installing helm stack: %s error: %v", stack.Name, err)
	}
	return map[string]interface{}{
		"outputs": map[string]string{
			"notes": notes,
		},
	}, err
}

func upgradeCatalog(event *events.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	stack := decodeHelmStack(event, cli, true)
	notes, err := helm.UpgradeHelmStack(stack)
	if err != nil {
		log.Errorf("Error upgrading helm stack: %s error: %v", stack.Name, err)
	}
	return map[string]interface{}{
		"outputs": map[string]string{
			"notes": notes,
		},
	}, err
}

func removeCatalog(event *events.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	stack := decodeHelmStack(event, cli, false)
	err := helm.DeleteHelmStack(stack)
	if err != nil {
		log.Errorf("Error removing helm stack: %s error:  %v", stack.Name, err)
	}
	return nil, err
}

func rollbackCatalog(event *events.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	stack := decodeHelmStack(event, cli, false)
	err := helm.RollbackHelmStack(stack)
	if err != nil {
		log.Errorf("Error rolling helm stack: %s error: %v", stack.Name, err)
	}
	return nil, err
}

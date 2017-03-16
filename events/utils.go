package events

import (
	"github.com/Sirupsen/logrus"
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubectld/helm"
)

type EventHandlerWithData func(*events.Event, *client.RancherClient) (map[string]interface{}, error)

func wrap(event *events.Event, cli *client.RancherClient, fn EventHandlerWithData) error {
	logrus.Infof("Received event: Name: %s, Event Id: %s, Resource Id: %s", event.Name, event.ID, event.ResourceID)
	logrus.Debugf("event.Data: %+v", event.Data)
	data, err := fn(event, cli)
	resp := newReply(event)
	if err == nil {
		resp.Data = data
	} else {
		resp.TransitioningMessage = err.Error()
		resp.Transitioning = "error"
	}
	return publishReply(resp, cli)
}

func decodeHelmStack(event *events.Event, cli *client.RancherClient, isUpgrade bool) *helm.Stack {
	var templates map[string]string
	if isUpgrade {
		templates = getStringMap(event.Data, "processData", "templates")
	} else {
		templates = getStringMap(event.Data, "environment", "data", "fields", "templates")
	}
	namespace := getString(event.Data, "environment", "data", "fields", "namespace")
	name := getString(event.Data, "environment", "name")

	return &helm.Stack{
		Name:      name,
		Namespace: namespace,
		Files:     templates,
	}
}

func newReply(event *events.Event) *client.Publish {
	return &client.Publish{
		Name:         event.ReplyTo,
		PreviousIds:  []string{event.ID},
		ResourceType: event.ResourceType,
		ResourceId:   event.ResourceID,
	}
}

func publishReply(reply *client.Publish, apiClient *client.RancherClient) error {
	_, err := apiClient.Publish.Create(reply)
	return err
}

func createAndPublishReply(event *events.Event, cli *client.RancherClient) error {
	reply := newReply(event)
	if reply.Name == "" {
		return nil
	}
	err := publishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

func getMap(data map[string]interface{}, keys ...string) map[string]interface{} {
	for _, key := range keys {
		val, ok := data[key]
		if !ok {
			return nil
		}
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}
		data = mapVal
	}

	result := map[string]interface{}{}
	for k, v := range data {
		result[k] = v
	}

	return result
}

func getStringMap(data map[string]interface{}, keys ...string) map[string]string {
	for _, key := range keys {
		val, ok := data[key]
		if !ok {
			return nil
		}
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}
		data = mapVal
	}

	result := map[string]string{}
	for k, v := range data {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}

	return result
}

func getString(data map[string]interface{}, keys ...string) string {
	for i, key := range keys {
		if i > len(keys)-2 {
			break
		}

		val, ok := data[key]
		if !ok {
			return ""
		}
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return ""
		}
		data = mapVal
	}

	val := data[keys[len(keys)-1]]
	if s, ok := val.(string); ok {
		return s
	}

	return ""
}

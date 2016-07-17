package util

import (
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/client"
)

func NewReply(event *events.Event) *client.Publish {
	return &client.Publish{
		Name:         event.ReplyTo,
		PreviousIds:  []string{event.ID},
		ResourceType: event.ResourceType,
		ResourceId:   event.ResourceID,
	}
}

func PublishReply(reply *client.Publish, apiClient *client.RancherClient) error {
	_, err := apiClient.Publish.Create(reply)
	return err
}

func CreateAndPublishReply(event *events.Event, cli *client.RancherClient) error {
	reply := NewReply(event)
	if reply.Name == "" {
		return nil
	}
	err := PublishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

func GetMap(data map[string]interface{}, keys ...string) map[string]interface{} {
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

func GetStringMap(data map[string]interface{}, keys ...string) map[string]string {
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

func GetString(data map[string]interface{}, keys ...string) string {
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

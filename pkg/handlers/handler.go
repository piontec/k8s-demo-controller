/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"github.com/Sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"default": &Default{},
}

// Default handler implements Handler interface,
// print each event with JSON format
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init() error {
	return nil
}

func (d *Default) ObjectCreated(obj interface{}) {
	pod := obj.(*api_v1.Pod)
	logrus.WithFields(logrus.Fields{
		"pod":       pod.Name,
		"namespace": pod.Namespace,
	}).Infof("created")
}

func (d *Default) ObjectDeleted(obj interface{}) {
	pod := obj.(*api_v1.Pod)
	logrus.WithFields(logrus.Fields{
		"pod":       pod.Name,
		"namespace": pod.Namespace,
	}).Infof("deleted")
}

func (d *Default) ObjectUpdated(oldObj, newObj interface{}) {
	pod := oldObj.(*api_v1.Pod)
	logrus.WithFields(logrus.Fields{
		"pod":       pod.Name,
		"namespace": pod.Namespace,
	}).Infof("updated")
}

/*

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

package resources

import (
	"context"

	"github.com/go-logr/logr"
	authorinov1beta1 "github.com/kuadrant/authorino/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AuthConfigHandler interface {
	AddHost(ctx context.Context, log logr.Logger, key types.NamespacedName, host string) error
	RemoveHost(ctx context.Context, log logr.Logger, key types.NamespacedName, host string) error
}

type authConfigHandler struct {
	client client.Client
}

func NewAuthConfigHandler(client client.Client) AuthConfigHandler {
	return &authConfigHandler{
		client: client,
	}
}

func (r *authConfigHandler) AddHost(ctx context.Context, log logr.Logger, key types.NamespacedName, host string) error {
	authConfig := &authorinov1beta1.AuthConfig{}
	err := r.client.Get(ctx, key, authConfig)
	if err != nil && errors.IsNotFound(err) {
		log.V(1).Info("AuthConfig not found", "key", key)
		return nil
	} else if err != nil {
		return err
	}
	log.V(1).Info("Successfully fetch deployed AuthConfig", "key", key)

	for _, h := range r.extractHosts(host) {
		found := false
		for _, ah := range authConfig.Spec.Hosts {
			if ah == h {
				found = true
			}
		}
		if !found {
			authConfig.Spec.Hosts = append(authConfig.Spec.Hosts, h)
		}
	}
	return r.client.Update(ctx, authConfig)
}

func (r *authConfigHandler) RemoveHost(ctx context.Context, log logr.Logger, key types.NamespacedName, host string) error {
	authConfig := &authorinov1beta1.AuthConfig{}
	err := r.client.Get(ctx, key, authConfig)
	if err != nil && errors.IsNotFound(err) {
		log.V(1).Info("AuthConfig not found", "key", key)
		return nil
	} else if err != nil {
		return err
	}
	log.V(1).Info("Successfully fetch deployed AuthConfig", "key", key)

	for _, h := range r.extractHosts(host) {
		foundIndex := -1
		for i, ah := range authConfig.Spec.Hosts {
			if ah == h {
				foundIndex = i
			}
		}
		if foundIndex > -1 {
			authConfig.Spec.Hosts = append(authConfig.Spec.Hosts[:foundIndex], authConfig.Spec.Hosts[foundIndex+1:]...)
		}
	}
	return r.client.Update(ctx, authConfig)
}

func (r *authConfigHandler) extractHosts(host string) []string {
	return []string{host}
}

// inclueLocalHostsIfPresent will add other internal host addresses if an internal address is present
// *.x.svc.cluster.local
// *.x.svc
// *.x
/*
func inclueLocalHostsIfPresent(hosts []string) []string {
	allHosts := []string{}

	for _, val := range hosts {
		if strings.Contains(val, "svc.cluster.local") {
			allHosts = append(allHosts, strings.ReplaceAll(val, ".cluster.local", ""))
			allHosts = append(allHosts, strings.ReplaceAll(val, ".svc.cluster.local", ""))
		}
		allHosts = append(allHosts, val)
	}

	return allHosts
}
*/

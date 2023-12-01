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

package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	kservev1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/opendatahub-io/odh-model-controller/controllers/resources"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KserveAuthConfigReconciler struct {
	client            client.Client
	scheme            *runtime.Scheme
	authConfigHandler resources.AuthConfigHandler
}

func NewKserveAuthConfigReconciler(client client.Client, scheme *runtime.Scheme) *KserveAuthConfigReconciler {
	return &KserveAuthConfigReconciler{
		client:            client,
		scheme:            scheme,
		authConfigHandler: resources.NewAuthConfigHandler(client),
	}
}

func (r *KserveAuthConfigReconciler) Reconcile(ctx context.Context, log logr.Logger, isvc *kservev1beta1.InferenceService) error {

	if isvc.Status.URL == nil {
		log.V(1).Info("Inference Service not ready yet, waiting for URL")
		return nil
	}

	return r.authConfigHandler.AddHost(ctx, log, types.NamespacedName{
		Name:      isvc.GetNamespace() + "-protection",
		Namespace: isvc.GetNamespace(),
	}, isvc.Status.URL.Host)
}

func (r *KserveAuthConfigReconciler) Remove(ctx context.Context, log logr.Logger, isvc *kservev1beta1.InferenceService) error {
	return r.authConfigHandler.RemoveHost(ctx, log, types.NamespacedName{
		Name:      isvc.GetNamespace() + "-protection",
		Namespace: isvc.GetNamespace(),
	}, isvc.Status.URL.Host)
}

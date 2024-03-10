/*
Copyright 2024.

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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	selfservicev1 "github.com/AlonaKaplan/selfserviceoverlay/api/v1"
)

// OverlayNetworkReconciler reconciles a OverlayNetwork object
type OverlayNetworkReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=self.service.ovn.org,resources=overlaynetworks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=self.service.ovn.org,resources=overlaynetworks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=self.service.ovn.org,resources=overlaynetworks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OverlayNetwork object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *OverlayNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("overlaynetwork", req.NamespacedName)

	logger.Info("Reconciling OverlayNetwork")

	overlayNetwork := &selfservicev1.OverlayNetwork{}
	if err := r.Client.Get(ctx, req.NamespacedName, overlayNetwork); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get OverlayNetwrok %q: %v", req.NamespacedName, err)
	}

	netAttachDef, err := renderNetAttachDef(overlayNetwork)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to render NetworkAttachmentDefinition for OverlayNetwrok %q: %v", req.NamespacedName, err)
	}
	if err := r.Create(ctx, netAttachDef); err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("NetworkAttachmentDefinition [%q] already exist", req.NamespacedName)
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OverlayNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&selfservicev1.OverlayNetwork{}).
		Complete(r)
}

func renderNetAttachDef(overlayNet *selfservicev1.OverlayNetwork) (*netv1.NetworkAttachmentDefinition, error) {
	const netAttachDefKind = "NetworkAttachmentDefinition"
	const netAttachDefAPIVer = "v1"

	cniNetConf := map[string]interface{}{
		"type":             "ovn-k8s-cni-overlay",
		"name":             overlayNet.GetName(),
		"netAttachDefName": overlayNet.Namespace + "/" + overlayNet.Name,
		"topology":         "layer2",
	}
	cniNetConfRaw, err := json.Marshal(cniNetConf)
	if err != nil {
		return nil, err
	}

	return &netv1.NetworkAttachmentDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: netAttachDefAPIVer,
			Kind:       netAttachDefKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      overlayNet.Name,
			Namespace: overlayNet.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         overlayNet.APIVersion,
					Kind:               overlayNet.Kind,
					Name:               overlayNet.Name,
					UID:                overlayNet.UID,
					BlockOwnerDeletion: nil,
				},
			},
		},
		Spec: netv1.NetworkAttachmentDefinitionSpec{
			Config: string(cniNetConfRaw),
		},
	}, nil
}

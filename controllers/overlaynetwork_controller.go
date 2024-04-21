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
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	"github.com/AlonaKaplan/selfserviceoverlay/api/v1alpha1"

	"github.com/AlonaKaplan/selfserviceoverlay/pkg/render"
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
	overlayNetwork := &v1alpha1.OverlayNetwork{}
	if err := r.Client.Get(ctx, req.NamespacedName, overlayNetwork); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("failed to get OverlayNetwork %q: %v", req.NamespacedName, err)
		}
		return ctrl.Result{}, nil
	}

	if overlayNetwork.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	desiredNetAttachDef, err := render.NetAttachDef(overlayNetwork)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to render NetworkAttachmentDefinition for OverlayNetwork %q: %v", req.NamespacedName, err)
	}

	actualNetAttachDef := &netv1.NetworkAttachmentDefinition{}
	actualKey := client.ObjectKey{Namespace: overlayNetwork.Namespace, Name: overlayNetwork.Name}
	if gerr := r.Client.Get(ctx, actualKey, actualNetAttachDef); gerr != nil {
		if !errors.IsNotFound(gerr) {
			return ctrl.Result{}, fmt.Errorf("failed to get NetworkAttachmetDefinition: %v", gerr)
		}
		if cerr := r.Client.Create(ctx, desiredNetAttachDef); cerr != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create NetworkAttachmetDefinition: %v", cerr)
		}
		return ctrl.Result{}, nil
	}

	if !hasOwnerReferenceWithUID(overlayNetwork.UID, actualNetAttachDef.OwnerReferences) {
		return ctrl.Result{}, fmt.Errorf("foreign NetworkAttachmetDefinition with the desired name already exist")
	}

	if actualNetAttachDef.Spec.Config != desiredNetAttachDef.Spec.Config {
		return ctrl.Result{}, fmt.Errorf("mutating NetworkAttachmetDefinition is not possible")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OverlayNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OverlayNetwork{}).
		Complete(r)
}

func hasOwnerReferenceWithUID(uid types.UID, ownerRefs []metav1.OwnerReference) bool {
	for _, ownerRef := range ownerRefs {
		if ownerRef.UID == uid {
			return true
		}
	}
	return false
}

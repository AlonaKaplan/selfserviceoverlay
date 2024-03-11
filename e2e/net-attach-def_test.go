package e2e

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	selfservicev1 "github.com/AlonaKaplan/selfserviceoverlay/api/v1"
)

const overlayNetworkKind = "OverlayNetwork"

var _ = Describe("OverlayNetwork controller, net-attach-def lifecycle", func() {
	It("should create and delete net-attach-def according to OverlayNetwork state", func() {
		By("Create test OverlayNetwork instance")
		ovrlyNet := newTestOverlayNetwork(TestsNamespace, "test")
		Expect(RuntimeClient.Create(context.Background(), ovrlyNet)).To(Succeed())

		By("Assert a corresponding net-attach-def has been created")
		nad := &netv1.NetworkAttachmentDefinition{}
		nadKey := types.NamespacedName{Namespace: ovrlyNet.Namespace, Name: ovrlyNet.Name}
		Eventually(func() error {
			return RuntimeClient.Get(context.Background(), nadKey, nad)
		}).Should(Succeed(), "1m", "3s")

		By("Assert overlay-network corresponding net-attach-def object")
		assertOverlayNetworkNetAttachDef(ovrlyNet, nad)

		By("Delete test OverlayNetwork instance")
		Expect(RuntimeClient.Delete(context.Background(), ovrlyNet)).To(Succeed())

		By("Assert the corresponding net-attach-def is deleted by and not exist")
		Eventually(func() bool {
			return errors.IsNotFound(RuntimeClient.Get(context.Background(), nadKey, nad))
		}).Should(BeTrue(), "1m", "3s",
			"the overlay-network corresponding net-attach-def should be disposed")
	})
})

func assertOverlayNetworkNetAttachDef(overlyNet *selfservicev1.OverlayNetwork, nad *netv1.NetworkAttachmentDefinition) {
	expectedNetAttachName := overlyNet.Namespace + "/" + overlyNet.Name
	expectedCniNetConf := fmt.Sprintf(`{	
		"name":"test",
		"type":"ovn-k8s-cni-overlay",				
		"netAttachDefName": "%s",
		"topology":"layer2"
	}`, expectedNetAttachName)

	Expect(nad.Spec.Config).To(MatchJSON(expectedCniNetConf))

	expectedOwnerReference := metav1.OwnerReference{
		APIVersion: selfservicev1.GroupVersion.String(), // "self.service.ovn.org/netv1",
		Kind:       overlayNetworkKind,
		Name:       overlyNet.Name,
		UID:        overlyNet.UID,
	}

	Expect(nad.ObjectMeta.OwnerReferences).To(Equal([]metav1.OwnerReference{expectedOwnerReference}))
}

func newTestOverlayNetwork(namespace, name string) *selfservicev1.OverlayNetwork {
	return &selfservicev1.OverlayNetwork{
		TypeMeta: metav1.TypeMeta{
			Kind:       overlayNetworkKind,
			APIVersion: selfservicev1.GroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: selfservicev1.OverlayNetworkSpec{},
	}
}

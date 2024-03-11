package e2e

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	selfservicev1 "github.com/AlonaKaplan/selfserviceoverlay/api/v1"
)

const TestsNamespace = "overlay-network-tests"

var (
	KubeClient    *kubernetes.Clientset
	RuntimeClient client.Client
)

func TestSelfserviceoverlay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Selfserviceoverlay Suite")
}

var _ = BeforeSuite(func() {
	cfg, err := config.GetConfig()
	Expect(err).ToNot(HaveOccurred(), "failed to get configuration for existing cluster")
	Expect(cfg).ToNot(BeNil())

	err = selfservicev1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = netv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	KubeClient, err = kubernetes.NewForConfig(cfg)
	Expect(KubeClient).ToNot(BeNil())

	RuntimeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(RuntimeClient).ToNot(BeNil())

	By("create tests namespace")
	testNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: TestsNamespace}}
	_, err = KubeClient.CoreV1().Namespaces().Create(context.Background(), testNamespace, metav1.CreateOptions{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("delete tests namespace")
	err := KubeClient.CoreV1().Namespaces().Delete(context.Background(), TestsNamespace, metav1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())

	By("wait for tests namespace to dispose")
	Eventually(func() bool {
		_, err := KubeClient.CoreV1().Namespaces().Get(context.Background(), TestsNamespace, metav1.GetOptions{})
		return errors.IsNotFound(err)
	}, "30s", "1s").Should(BeTrue(), "tests namespace should be disposed")
})

package seldonwrapper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"flag"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	v1alpha2client "github.com/SeldonIO/seldon-operator/pkg/client/clientset/versioned/typed/machinelearning/v1alpha2"

	"github.com/glindsell/seldon/internal/seldonwrapper"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Suite")
}

var _ = Describe("SeldonClient", func() {
	var kubeconfig *string
	//var seldonClient seldonclient.SeldonClient
	var seldonClientset *v1alpha2client.MachinelearningV1alpha2Client

	BeforeEach(func() {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}

		seldonClientset, err = v1alpha2client.NewForConfig(config)
		if err != nil {
			panic(err)
		}
	})

	It("Returns a client with a Seldon Deployment Interface", func() {
		sc, err := seldonwrapper.NewSeldonWrapper(*seldonClientset)
		Expect(err).Should(BeNil())
		_, ok := sc.Deployment.(v1alpha2client.SeldonDeploymentInterface)
		Expect(ok).To(Equal(true))
	})
})

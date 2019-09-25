package seldonwrapper

import (
	"errors"
	"log"
	"time"

	apiv1 "k8s.io/api/core/v1"
	errorsv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha2client "github.com/SeldonIO/seldon-operator/pkg/client/clientset/versioned/typed/machinelearning/v1alpha2"
	v1alpha2api "github.com/seldonio/seldon-operator/pkg/apis/machinelearning/v1alpha2"
)

type SeldonWrapper struct {
	Deployment v1alpha2client.SeldonDeploymentInterface
}

func NewSeldonWrapper(seldonClientset v1alpha2client.MachinelearningV1alpha2Client) (SeldonWrapper, error) {
	seldonDeploymentsClient := seldonClientset.SeldonDeployments(apiv1.NamespaceDefault)
	return SeldonWrapper{seldonDeploymentsClient}, nil
}

func (sw *SeldonWrapper) CreateDeployment(dep *v1alpha2api.SeldonDeployment) error {
	if result, err := sw.Deployment.Create(dep); err != nil {
		return err
	} else {
		log.Printf("%+v\n", result)
	}
	return nil
}

func (sw *SeldonWrapper) DeleteDeployment(name string) error {
	log.Println("Deleting resource...")
	deletePolicy := metav1.DeletePropagationBackground
	if err := sw.Deployment.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

func (sw *SeldonWrapper) WaitForDeploymentState(name string, state string) error {
	serviceStreamWatcher, err := sw.Deployment.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for {
		select {
		case event := <-serviceStreamWatcher.ResultChan():
			if dep, ok := event.Object.(*v1alpha2api.SeldonDeployment); ok {
				if dep.Name == name && dep.Status.State == state {
					log.Printf("*** Success! *** Event for deployment: \"%v\" received. Required state: \"%v\" achieved.", dep.Name, dep.Status.State)
					return nil
				}
				log.Printf("Event for deployment: \"%v\" received. Current state: \"%v\", waiting for state: \"%v\", timeout in 120s...\n", dep.Name, dep.Status.State, state)
			}
		case <-time.After(120 * time.Second):
			return errors.New("Timed out waiting for deployment state to change")
		}
	}
}

func (sw *SeldonWrapper) WaitForDeploymentNotFound(name string) error {
	for {
		dep, err := sw.Deployment.Get(name, metav1.GetOptions{})
		if err != nil {
			if serr, ok := err.(*errorsv1.StatusError); ok && serr.ErrStatus.Reason == metav1.StatusReasonNotFound {
				log.Printf("*** Success! *** Deployment: \"%v\" %v.", name, serr.ErrStatus.Reason)
				return nil
			} else {
				return err
			}
		}
		log.Printf("Deployment: \"%v\" still exists in state: \"%v\". Please wait for deletion to finish. Timeout in 120s...\n", dep.Name, dep.Status.State)
		time.Sleep(5 * time.Second)
	}
}

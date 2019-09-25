# Seldon Deployment Manager
This app...  
**1. takes in a Seldon Core Custom Resource and creates it over the Kubernetes API**  
*it then...*  
**2. watches the created resource to wait for it to become available**  
*and finally...*  
**3. when it is available deletes the Custom Resource.**  

## Prerequisites
Minikube [https://kubernetes.io/docs/tasks/tools/install-minikube/](https://kubernetes.io/docs/tasks/tools/install-minikube/)  
Helm [https://helm.sh/docs/using_helm/](https://helm.sh/docs/using_helm/)  
Go [https://golang.org/doc/install](https://golang.org/doc/install)  

## Steps to run

### Start Minikube cluster and install seldon core helmcharts:

```minikube start```

```helm init --history-max 200```

```helm install seldon-core-operator --name seldon-core --repo https://storage.googleapis.com/seldon-charts --set usageMetrics.enabled=true --namespace seldon-system```

### Build/Test/Run/Clean with make:

```make build```

```make test```

```make run```

```make clean```

### Use as standalone binary:

```make```

```./seldondm model.json```

### When finished:  

```make clean```

```minikube stop && minikube delete```

## Notes & Further Work
**1. A workaround has been implemented to avoid a "null string error" in JSON unmarshalling (see TODO in main.go)**  
Example error without workaround:
```
panic: SeldonDeployment.machinelearning.seldon.io "seldon-model" is invalid: ... : validation failure list:
spec.predictors.componentSpecs.metadata.creationTimestamp in body must be of type string: "null"
```
Relevant issues:
[https://github.com/kubernetes/kubernetes/issues/58311](https://github.com/kubernetes/kubernetes/issues/58311)

[https://github.com/kubernetes/kubernetes/issues/66899](https://github.com/kubernetes/kubernetes/issues/66899)

[https://github.com/coreos/prometheus-operator/issues/2399](https://github.com/coreos/prometheus-operator/issues/2399)

[https://github.com/coreos/prometheus-operator/blob/f1392d5430147b033d01ebb6a67ded18d2f2e6fc/test/e2e/alertmanager_test.go](https://github.com/coreos/prometheus-operator/blob/f1392d5430147b033d01ebb6a67ded18d2f2e6fc/test/e2e/alertmanager_test.go)

**2. More tests need to be added**  
Current test is only a basic example to show use of behavioural testing framework.  
[https://onsi.github.io/ginkgo/](https://onsi.github.io/ginkgo/)

**3. Make timeout configurable**  
Timeout is currently hard coded to 120s while waiting for resource to be available and again for it to be deleted. This may not be suitable for all environments.

**4. Add help flag to binary**  
Add help flag which displays command help and command line arguments.

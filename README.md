## Kubernetes app auto-scaling with Prometheus and KEDA

![](keda-prometheus.jpg)

Scalability is a key requirement for cloud native applications. With [Kubernetes](https://kubernetes.io), scaling your application is as simple as increasing the number of replicas for the corresponding `Deployment` or `ReplicaSet` - but, this is a manual process. Kubernetes makes it possible to automatically scale your applications (i.e. `Pod`s in a `Deployment` or `ReplicaSet`) in a declarative manner using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.16/#horizontalpodautoscaler-v1-autoscaling) specification.

[This blog post](https://dev.to/azure/how-to-auto-scale-your-kubernetes-apps-with-prometheus-and-keda-39km) demonstrates how you can external metrics to auto-scale a Kubernetes application. For demonstration purposes, we will use HTTP access request metrics that are exposed using [Prometheus](https://prometheus.io). Instead of using the `Horizontal Pod Autoscaler` directly, we will leverage [**Kubernetes Event Driven Autoscaling** aka **KEDA**](https://github.com/kedacore/keda) - an open source Kubernetes operator which integrates natively with the `Horizontal Pod Autoscaler` to provide fine grained autoscaling (including to/from zero) for event-driven workloads.

## Setup (Updated for Helm3)
Ensure you have a local k8s cluster running, `kubectl` configured and helm 3 installed
1. Install KEDA via helm
    ```bash
    helm repo add kedacore https://kedacore.github.io/charts
    helm repo update
    helm install keda kedacore/keda --namespace keda
    ```
2. Ensure KEDA is installed properly and Running
    ```bash
    kubectl get pods -n keda -w
    ```
3. Install Redis
    ```bash
    helm install redis --set architecture=standalone --set auth.enabled=false bitnami/redis
    ```
4. Start application
    ```bash
    kubectl apply -f go-app.yaml
    ```
5. Start prometheus
    ```bash
    kubectl apply -f prometheus.yaml
    ```
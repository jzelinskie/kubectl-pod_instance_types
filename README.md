# kubectl pod-instance-types

A kubectl plugin to list pods alongside the cloud instance type of the node each pod is running on.

## Installation

Use [krew](https://krew.sigs.k8s.io) plugin manager to install:

```text
kubectl krew install pod-instance-types
kubectl pod-instance-types --help
```

## Demo

```text
$ kubectl pod-instance-types -n kube-system
NAME                                       READY   STATUS    RESTARTS   AGE   INSTANCE-TYPE
coredns-5dd5756b68-4xnkf                   1/1     Running   0          2d    m5.large
coredns-5dd5756b68-r7ftn                   1/1     Running   0          2d    m5.large
etcd-control-plane                         1/1     Running   0          2d    m5.xlarge
kube-apiserver-control-plane               1/1     Running   0          2d    m5.xlarge
kube-controller-manager-control-plane      1/1     Running   0          2d    m5.xlarge
kube-proxy-6fgjz                           1/1     Running   0          2d    m5.large
kube-proxy-92xjp                           1/1     Running   0          2d    c5.2xlarge
kube-scheduler-control-plane               1/1     Running   0          2d    m5.xlarge
metrics-server-7db4fb59f9-2q8cx            1/1     Running   0          2d    c5.2xlarge
```

The `INSTANCE-TYPE` column is resolved from the node's `node.kubernetes.io/instance-type` label (falling back to the legacy `beta.kubernetes.io/instance-type` label for older clusters).
Pods not yet scheduled to a node show `<none>`.

All standard output formats are supported by injecting the instance-type as an annotation:

```text
$ kubectl pit -n kube-system -o json | jq '.[].metadata.annotations["node.kubernetes.io/instance-type"]'
"m5.large"
"m5.large"
...
```

## How it works

For human-readable output, the plugin requests pods from the API server as a server-side table (the same mechanism `kubectl get` uses) and appends an `INSTANCE-TYPE` column populated by looking up each pod's node.
For structured output formats, the instance type is injected as the `node.kubernetes.io/instance-type` annotation on each pod object before serialization.

## License

Apache 2.0. See [LICENSE](./LICENSE).

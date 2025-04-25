# sidecarGracefulShutdown
Understanding Graceful shutdown of containers in Pods

# Building and Running in Minikube
Build the Docker image:
```bash
    docker build -t graceful-shutdown-demo:latest .
```

If using Minikube, load the image:
```bash
    minikube image load graceful-shutdown-demo:latest
```

Apply the Kubernetes manifests:
```bash
    kubectl apply -f pod.yaml
    kubectl apply -f service.yaml
```

Observe the running pod:
```bash
    kubectl get pods
```

Test the application:
```bash
    kubectl port-forward pod/graceful-shutdown-demo 8080:8080
```
Then visit http://localhost:8080 in your browser or use curl.

# Demonstrating Graceful Shutdown
Trigger pod deletion:
```bash
    kubectl delete pod graceful-shutdown-demo
```

Observe the shutdown process in logs:
```bash
    kubectl logs -f graceful-shutdown-demo -c main-app
    kubectl logs -f graceful-shutdown-demo -c sidecar
```

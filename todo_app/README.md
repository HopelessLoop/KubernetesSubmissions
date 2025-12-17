# Todo App - K3d Deployment

This project is a Todo List application written in Go, containerized and configured for deployment on a local Kubernetes cluster using **K3d** (K3s in Docker).

## Prerequisites

Before running the project, ensure you have the following tools installed:

* [Docker](https://www.docker.com/)
* [K3d](https://k3d.io/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Quick Start

### 1. Create K3d Cluster with Port Mapping

Since K3d runs Kubernetes nodes inside Docker containers, we must explicitly map the host machine's port to the K3s node's `NodePort`.

* **Host Port**: `30880` (The port you will use in your browser).
* **NodePort**: `30880` (Must match the `nodePort` defined in `service.yaml`).

Run the following command to create the cluster and map the ports:

```bash
# Creates a cluster named "todo-cluster"
# Maps localhost:8081 -> K3d Container:30880
k3d cluster create -p "30880:30880@server:0 --agents 2"

```

### 2. Build and Import Image

Since we are using a local Docker image (instead of pulling from a public registry like Docker Hub), you must import it into the cluster so K3s can access it.

```bash
# 1. Build the image
docker build -t todo-app:1.0 .

# 2. Import the image into the K3d cluster
k3d image import todo-app:1.0 -c k3s-default
```

### 3. Deploy Resources

Apply the Kubernetes configurations. Ensure your `deployment.yaml` and `service.yaml` are in the current directory.

```bash
# Apply the Deployment
kubectl apply -f deployment.yaml

# Apply the Service
kubectl apply -f service.yaml
```

### 4. Verify and Access

Check if the pods are running:

```bash
kubectl get pods
kubectl get svc
```

Once the pod status is `Running`, access the application in your browser:

```
http://localhost:30880
```
---

## Port Configuration & Architecture

Understanding the network flow is crucial for troubleshooting. The following table explains how traffic flows from your machine to the application container.

| Component | Port | Description |
| --- | --- | --- |
| **Host Machine** | **30880** | Entry point for your browser/client. Mapped via `k3d create -p`. |
| **Node (K3d)** | **30880** | The `nodePort` defined in `service.yaml`. The K3s node listens here. |
| **Service** | **20880** | The `port` defined in `service.yaml`. Internal virtual port (ClusterIP). |
| **Pod** | **30880** | The `targetPort` defined in `service.yaml`. The actual port the Go app listens on. |

**Configuration Note:**
If you change the port in your Go application code, you must update the `targetPort` in `service.yaml`. If you change the `nodePort` in `service.yaml`, you must recreate the K3d cluster with the updated `-p` mapping.

## Clean Up

To remove the cluster and release system resources:

```bash
k3d cluster delete
```
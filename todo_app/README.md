# Todo App Project

This project contains a simple Todo application. This guide outlines the steps to build the Docker image, run it locally, and deploy it to a local k3d Kubernetes cluster.

## Prerequisites

Ensure you have the following tools installed:

* **Docker**
* **k3d** (Lightweight Kubernetes that runs in Docker)
* **kubectl** (Kubernetes command-line tool)

## Getting Started

### 1. Build the Docker Image

First, build the container image from the Dockerfile in the current directory.

```bash
docker build . -t todo_app:1.0
```

### 2. Run Locally (Docker Standalone)

To test the application without Kubernetes, run the container directly using Docker. This maps host port `30880` to container port `30880`.

```bash
docker run --name todo_app -p 30880:30880 -d todo_app:1.0
```

*You can now access the application at `http://localhost:30880`.*

---

## Kubernetes Deployment (k3d)

### 3. Import Image to k3d Cluster

Since the image is built locally, you must import it into the k3d cluster so Kubernetes can access it.
*(Note: This assumes your cluster is named `k3s-default`).*

```bash
k3d image import todo_app:1.0 -c k3s-default
```

### 4. Deploy to Kubernetes

Create a deployment in the cluster using the imported image.

```bash
kubectl apply -f ./manifests/deployment.yaml
```

### 5. Verification & Logs

Check the status of your pods to ensure the application is running.

```bash
kubectl get pods
```

Once you have the pod name from the previous command, you can stream the application logs using the following command (replace `<your-pod-name>` with the actual ID):

```bash
kubectl logs -f <your-pod-name>
```

### 6. Access the Application (Port Forwarding)

To access the application running inside the Kubernetes cluster from your local machine, use `port-forward`.

Use the pod name retrieved in Step 5 (e.g., `todo-app-dep-xxxxx-xxxxx`) to replace `<your-pod-name>` below:

```bash
kubectl port-forward <your-pod-name> 30880:30880
```

Once the forwarding is active, open your browser and visit: `http://localhost:30880`

**Expected Output:**
You should see the following JSON response:

```json
{"message":"Hello World"}
```
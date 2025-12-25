# Hourly Image Application

This project is a containerized application designed to demonstrate a deployment workflow using **k3d**, **Docker**, and **Kubernetes**. It serves a web application that generates a timestamped UUID.

## Prerequisites

Before running this project, ensure you have the following tools installed on your machine:

* [Docker](https://docs.docker.com/get-docker/)
* [k3d](https://k3d.io/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Getting Started

Follow the steps below to set up the cluster, build the image, and deploy the application.

### 1. Create a k3d Cluster

Create a local Kubernetes cluster with 2 agent nodes and map the cluster's LoadBalancer port 80 to your local port 8081.

```bash
k3d cluster create -p 8081:80@loadbalancer --agents 2
```

### 2. Build the Docker Image

Build the application image using the local Dockerfile. We will tag this version as `1.0`.

```bash
docker build -t todo-app-backend:1.0 .
```

### 3. Import Image to k3d

Since k3d runs in containers, it cannot access your local Docker daemon's images by default. Import the built image into the cluster named `k3s-default`.

```bash
k3d image import todo-app-backend:1.0
```

### 4. Deploy Manifests

Apply the Kubernetes manifests (Deployment, Service, Ingress, etc.) located in the `manifests` folder.

```bash
kubectl apply -f ./manifests/
```

> **Note:** Please allow a few moments for the pods to initialize and reach the `Running` state. You can check the status using `kubectl get pods`.

## Verification

Once the deployment is successful, open your browser and visit:

```
http://localhost:8081/image
```

### Expected Output

You should see a image.

## Cleanup

To stop and remove the local cluster created for this project, run:

```bash
k3d cluster delete
```
# Deployment Process  

## Build Docker Image  
Execute the following command in the root directory of the application to build an image tagged 1.0:  
```  
docker build . -t log_output:1.0  
```  
## Import Image to K3s Cluster  
Import the locally built image into the K3s cluster (cluster name: k3s-default):
```
k3d image import log_output:1.0 -c k3s-default
```

## Create Deployment Resource
Create a deployment via the kubectl command and specify the image built above:
```
kubectl create deployment log-output --image=log_output:1.0
```

## Check Deployment Status and Application Logs

### Check Pod Running Status

Execute the following command to confirm that the Pod related to log-output has started normally:

```
kubectl get pods
```

Expected Output Example (Pod name is dynamically generated, subject to the actual environment):

```
NAME                                 READY   STATUS    RESTARTS   AGE
log-output-585b8c89bb-rqpxb          1/1     Running   0          26s
```

### View Application Logs

Replace <your-pod-name> in the following command with the actual Pod name starting with "log-output" to view the application output logs:

```
kubectl logs -f <your-pod-name>
```

Example Command (based on the above output example):

```
kubectl logs -f log-output-585b8c89bb-rqpxb
```

# Notes

- Pod names are dynamically generated. You need to obtain the actual name through the kubectl get pods command before replacing it.

- If the Pod status does not show Running, you can use kubectl describe pod <your-pod-name> to troubleshoot deployment issues.

- The image tag (1.0) should be adjusted according to actual version requirements. Ensure the tag is consistent when importing into the cluster and creating the deployment.

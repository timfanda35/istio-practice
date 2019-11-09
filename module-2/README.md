# An Web Application

## Requirement

* Docker

Move work directory to `cd module-2/`

```
cd module-2/
```

## Run Sample Web Applicaton with Docker Compose

```
docker-compose up
```

Access the `http://localhost:8080` with web browserss.

Looks like:

![](./images/developing.png)

## Build the container image and push to registry

```
# color-backend
docker build ./color-backend/ \
  -f ./color-backend/dockerfile \
  -t "gcr.io/$(gcloud config get-value core/project)/color-backend:v1"

# color-frontend
docker build ./color-frontend/ \
  -f ./color-frontend/dockerfile \
  -t "gcr.io/$(gcloud config get-value core/project)/color-frontend:v1"
```

Authorizate to Container Registry

```
gcloud auth configure-docker
```

Push images to Container Registry

```
docker push "gcr.io/$(gcloud config get-value core/project)/color-backend:v1"
docker push "gcr.io/$(gcloud config get-value core/project)/color-frontend:v1"
```

Check images on Container Registry

```
gcloud container images list
```

## Run Sample Web Applicaton on Kubernetes

Deploy application to kubernetes cluster.

```
# color-backend
cat ./kubernetes/v1-color-backend-deployment.yaml \
 | sed "s|<REPLACE_WITH_YOUR_PROJECT_ID>|$(gcloud config get-value core/project)|g" \
 | kubectl apply -f -

kubectl apply -f ./kubernetes/v1-color-backend-service.yaml

# color-frontend
cat ./kubernetes/v1-color-frontend-deployment.yaml \
 | sed "s|<REPLACE_WITH_YOUR_PROJECT_ID>|$(gcloud config get-value core/project)|g" \
 | kubectl apply -f -

kubectl apply -f ./kubernetes/v1-color-frontend-service.yaml
```

Check the containers in Pod

```
kubectl get pod
```

Looks like:

```
NAME                                 READY   STATUS    RESTARTS   AGE
color-backend-v1-566b649bfb-dn6fd    1/1     Running   0          8s
color-frontend-v1-849bcbc8dd-tvddz   1/1     Running   0          2s
```

The container in Pod in only one. It means there's no Istio sidecar injection.

## Enable the auto sidecar injection

```
kubectl label namespace default istio-injection=enabled
```

Redeploy application

```
kubectl delete pod -l workload=color
```

Check the containers in Pod

```
kubectl get pod
```

Looks like:

```
NAME                                 READY   STATUS    RESTARTS   AGE
color-backend-v1-566b649bfb-qkphd    2/2     Running   0          14s
color-frontend-v1-849bcbc8dd-sgpjw   2/2     Running   0          8s
```

The container in Pod is two now.

## Expose the web application with Istio

```
kubectl apply -f ./istio
```

Get ingress external IP

```
echo "$(kubectl -n istio-system get services -l app=istio-ingressgateway \
   -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')"
```

Access the `http://<The Ingress external IP>` with web browserss.

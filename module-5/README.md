# Timout/Retry/Fault Injection

Before start this module, you have to finish Module-4.

Move work directory to `cd module-5/`

```
cd module-5/
```

## Update Destination Rule for new versions

```
kubectl apply -f ./istio/color-frontend-destination-rule.yaml
```

```
Expose the ingress ip

export INGRESS_IP=<The Ingress external IP>
```

## Handle Timeout

### Add delay to color-frontend:v3

Modify `color-frontend-with-delay/src/main.go`, We had prepared the modified files.

```
...

import (
	  "fmt"
		"io/ioutil"
		"html/template"
		"log"
		"net/http"
		"os"
		"time"
)

...

func index(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second) // Delay 5 second
...

```

Build and push new container image

```
docker build ./color-frontend-with-delay/ \
  -f ./color-frontend-with-delay/dockerfile \
  -t "gcr.io/$(gcloud config get-value core/project)/color-frontend:v3"

docker push "gcr.io/$(gcloud config get-value core/project)/color-frontend:v3"
```

Deploy new versions to kubernetes cluster

```
cat ./kubernetes/v3-color-frontend-deployment.yaml \
 | sed "s|<REPLACE_WITH_YOUR_PROJECT_ID>|$(gcloud config get-value core/project)|g" \
 | kubectl apply -f -
```

Update Istio config

```
kubectl apply -f ./istio/delay-color-frontend-virtual-service.yaml
```

Try with curl

```
time curl -s http://${INGRESS_IP}/
```

You will get

```
real	0m10.145s
```

### Add timeout setting

Update Istio config

```
kubectl apply -f ./istio/timeout-color-frontend-virtual-service.yaml
```

Try with curl

```
time curl -s http://${INGRESS_IP}/
```

You will get

```
upstream request timeout
real	0m3.076s
```

## Handle Retry

### Set error response to color-frontend:v4

Modify `color-frontend-with-error/src/main.go`, We had prepared the modified files.

```
import (
	"fmt"
	"io/ioutil"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

...

func index(w http.ResponseWriter, r *http.Request) {
...
}
```

to

```
...

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

...

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Something bad happened!"))
}
```

Build and push new container image

```
docker build ./color-frontend-with-error/ \
  -f ./color-frontend-with-error/dockerfile \
  -t "gcr.io/$(gcloud config get-value core/project)/color-frontend:v4"

docker push "gcr.io/$(gcloud config get-value core/project)/color-frontend:v4"
```

Deploy new versions to kubernetes cluster

```
cat ./kubernetes/v4-color-frontend-deployment.yaml \
 | sed "s|<REPLACE_WITH_YOUR_PROJECT_ID>|$(gcloud config get-value core/project)|g" \
 | kubectl apply -f -
```

Update Istio config

```
kubectl apply -f ./istio/error-color-frontend-virtual-service.yaml
```

Try with curl

```
curl -s http://${INGRESS_IP}/
```

You will get

```
500 - Something bad happened!
```

### Add retry setting

Update Istio config

```
kubectl apply -f ./istio/retry-color-frontend-virtual-service.yaml
```

Enable the proxy log

```
# In Istio directory
helm template install/kubernetes/helm/istio \
  --namespace=istio-system -x templates/configmap.yaml \
  --set global.proxy.accessLogFile="/dev/stdout" \
  | kubectl replace -f -
```

Try with curl

```
curl -s http://${INGRESS_IP}/
```

You will still get

```
500 - Something bad happened!
```

But let we see the proxy log of color-frontend-v4

```
kubectl logs \
  $(kubectl get pod -l version=v4 -o jsonpath='{.items[0].metadata.name}') \
  -c istio-proxy
```

You can see the access log has appended 4 request records:

```
[2019-09-26T06:32:03.850Z] "GET / HTTP/1.1" 500 - "-" "-" 0 29 0 0 "10.20.1.1" "curl/7.60.0" "02dc03d0-3978-4b3e-b550-a0aad476f552" "<INGRESS_IP>" "127.0.0.1:8080" inbound|80|http|color-frontend.default.svc.cluster.local - 10.20.2.7:8080 10.20.1.1:0 - default
[2019-09-26T06:32:03.858Z] "GET / HTTP/1.1" 500 - "-" "-" 0 29 0 0 "10.20.1.1" "curl/7.60.0" "02dc03d0-3978-4b3e-b550-a0aad476f552" "<INGRESS_IP>" "127.0.0.1:8080" inbound|80|http|color-frontend.default.svc.cluster.local - 10.20.2.7:8080 10.20.1.1:0 - default
[2019-09-26T06:32:03.912Z] "GET / HTTP/1.1" 500 - "-" "-" 0 29 0 0 "10.20.1.1" "curl/7.60.0" "02dc03d0-3978-4b3e-b550-a0aad476f552" "<INGRESS_IP>" "127.0.0.1:8080" inbound|80|http|color-frontend.default.svc.cluster.local - 10.20.2.7:8080 10.20.1.1:0 - default
[2019-09-26T06:32:03.951Z] "GET / HTTP/1.1" 500 - "-" "-" 0 29 0 0 "10.20.1.1" "curl/7.60.0" "02dc03d0-3978-4b3e-b550-a0aad476f552" "<INGRESS_IP>" "127.0.0.1:8080" inbound|80|http|color-frontend.default.svc.cluster.local - 10.20.2.7:8080 10.20.1.1:0 - default
```

The retries default rely on:
```
connect-failure,refused-stream,unavailable,cancelled,resource-exhausted,retriable-status-codes
```

Reference: https://istio.io/docs/reference/config/networking/v1alpha3/virtual-service/#HTTPRetry

## Fault Injection

### Reverse Istio config to health color-health v1 and v2

```
kubectl apply -f ./istio/color-frontend-virtual-service.yaml
```

### Inject Fualt delay to color-backend:v1

```
kubectl apply -f ./istio/delay-color-backend-virtual-service.yaml
```

Try with curl

```
time curl -s http://${INGRESS_IP}/
```

You will get

```
real	0m5.319s
```

### Inject Fualt abort to color-backend:v1

```
kubectl apply -f ./istio/abort-color-backend-virtual-service.yaml
```

Try with curl

```
curl -s http://${INGRESS_IP}/ | grep fault | wc -l
```

You will get `1` means get fail.

# kube-shield

A dynamic multi-purpose admission webhook.


Local development with minikube:


**Create minikube cluster:**

`minikube start`

**Create tunnel for apiserver <-> webhook connections:**

`minikube tunnel`

**Apply the webhook configuration:**

`kubectl apply -f hack/manifests/webhook.yaml`

**Start webhook locally:**

`go run main.go -k ~/.kube/config -p hack/certs/server.crt -i hack/certs/server.key --ipv4 -d`

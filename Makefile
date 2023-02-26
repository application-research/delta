SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.19

.PHONY: all
all: build

.PHONY: build
build:
	git submodule update --init --recursive
	make -C extern/filecoin-ffi
	go build -tags netgo -ldflags '-s -w' -o delta

.PHONE: clean
clean:
	rm -f delta
	git submodule deinit --all -f

install:
	install -C -m 0755 delta /usr/local/bin

# delta kubernetes makefile vars
NAME=delta
VERSION=v0.0.1
VERSION_KUBE_DASHBOARD=v2.4.0
K8S_TOOLS=kubectl kind metrics-server
KIND_CONFIG=k8s/config/kind/kind.config

# Generate TLS keys, deploy to kubernetes, wait for the pods/services/deployments to start and map the frontend to ports 80 and 443
k8s.all: generate.keys k8s.deploy k8s.wait 

# Generate TLS keys, deploy to kubernetes
k8s.install: generate.keys k8s.deploy

# setup
k8s.setup:
	curl -ksLS https://get.arkade.dev | sudo sh;
	grep -qxF PATH=\$$PATH:\$$HOME/.arkade/bin/ ~/.profile || echo PATH=\$$PATH:\$$HOME/.arkade/bin/ >> ~/.profile
	source ~/.profile;
	arkade get $(K8S_TOOLS);
	docker build -t 0utercore/delta:v0.0.1 .

# Delete the deployment from kubernetes
k8s.uninstall: k8s.delete

# delete the generated keys and remove the deployment from kubernetes
k8s.clean: clean-local k8s.delete

# Delete the generated keys
k8s.clean-local:
	@-rm localhost.key localhost.crt;

# run in a clean development container, mounting the current directory to /workspace
k8s.dev:
	docker run -it -v $$(pwd):/workspace golang /bin/bash;

# open the default portainer dashboard view
k8s.dash:
	kubectl port-forward --address 0.0.0.0 svc/portainer 9000:9000;

# start the local container isntance
cluster.up: cluster.start

# stop the cluster
cluster.down: cluster.delete

# start the local container instance
cluster.start:
	kind create cluster --config $(KIND_CONFIG);
	arkade install portainer;
	arkade install kubernetes-dashboard;
	arkade install metrics-server;
#	kubectl apply -f k8s/config/kubernetes-dashboard/admin-user.yml;
#	kubectl apply -f k8s/config/kind/nginx-ingress-controller.yml;
	kubectl wait --for=condition=available --timeout=10m --all deployments;
#	docker build -t 0utercore/delta:v0.0.1 .
#   docker push 0utercore/delta:v0.0.1;
	
docker.deploy:
	docker build -t 0utercore/delta:v0.0.1 .
    docker push 0utercore/delta:v0.0.1;

# clean up the entire container cluster from the system
cluster.delete:
	kind delete cluster;

# generate tls keys and deploy them as a k8s tls secret for use in the nginx pod
generate.keys:
	MSYS_NO_PATHCONV=1 openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout localhost.key -out localhost.crt -subj "/CN=foo.bar.com";

# apply the kubernetes artifacts to the cluster
k8s.deploy:
	@-kubectl apply -f k8s/delta/;

# remove all components from kubernetes
k8s.delete:
	@-kubectl delete -f k8s/delta/;

# map the local port 8080 to cluster service port
k8s.start:
	@-sudo setcap CAP_NET_BIND_SERVICE=+eip $(which kubectl);
	kubectl port-forward --address 0.0.0.0 --namespace=ingress-nginx service/ingress-nginx-controller 80:80 443:443;

# start a local port-forward service mapping the k8s service to external port 8080
k8s.startd:
	@-setcap CAP_NET_BIND_SERVICE=+eip $(which kubectl);
	kubectl port-forward --address 0.0.0.0 --namespace=ingress-nginx service/ingress-nginx-controller 80:80 443:443 & echo $$! > k8s.PID

# stop the local service
k8s.stopd:
	@-kill `cat k8s.PID` || true;
	@-rm k8s.PID || true;

# wait time for all services to start up
k8s.wait:
	kubectl wait --for=condition=available --timeout=10m --all deployments

# synonym for all target
k8s.up: k8s.deploy

# synonym for uninstall target
k8s.down: k8s.delete
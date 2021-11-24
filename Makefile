up:
	minikube start --kubernetes-version 1.22.4
	skaffold dev

clean:
	skaffold delete
	minikube delete

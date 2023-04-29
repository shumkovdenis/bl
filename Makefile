namespace := denis

install-dapr:
	dapr init -k -n $(namespace)

add-bitnami:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update

install-redis:
	helm install redis bitnami/redis --set architecture=standalone -n $(namespace)

kube-apply:
	kubectl apply -f ./deploy/ -n $(namespace)

get-pods:
	kubectl get pods -n $(namespace)

namespace := denis

kube-apply:
	kubectl apply -f ./deploy -n $(namespace)

kube-apply-services:
	kubectl apply -f ./deploy/services -n $(namespace)

kube-apply-services:
	kubectl apply -f ./deploy/services -n $(namespace)

get-pods:
	kubectl get pods -n $(namespace)

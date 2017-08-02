# Guestbook CI/CD Demo with Jenkins
This app demonstrates a Jenkins pipeline with a Golang based web app and a SQL Server on Linux container. Containers are deployed to ACS Kubernetes via Helm charts and images are stored in Azure Container Registry

## Demo Setup

### Fork or Clone this repo

### Setup ACS Kubernetes

Use standard Azure Container Service instructions [here.](https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-kubernetes-walkthrough)  

### Azure Container Registry

Use standard Azure Container Registry instructions [here.](https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-tutorial-kubernetes-prepare-acr) 

* Add Kubernetes secret with ACR creds base64 encoded. Update secret-update.yaml with your values
```
kubectl create -f secret.yaml
```

### Install Jenkins

* Update jenkins-values.yaml
* Install Jenkins helm chart
```
helm --namespace jenkins --name jenkins -f ./jenkins-values.yaml install stable/jenkins

watch kubectl get svc --namespace jenkins # wait for external ip
export JENKINS_IP=$(kubectl get svc jenkins-jenkins --namespace jenkins --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
export JENKINS_URL=http://${JENKINS_IP}:8080

kubectl get pods --namespace jenkins # wait for running
open ${JENKINS_URL}/login
```
* Add ACR creds in Jenkins Global Credentials

### Database setup

* Helm chart install
```
helm install --name=guestbook-db ./charts/guestbook-db
```

* Wait for IP: ```watch kubectl get svc guestbook-db-guestbook-db```
* Get IP
```
export SQLDB_IP=$(kubectl get svc guestbook-db-guestbook-db --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
```

* Seed data (use sqlcmd CLI tool)
```
sqlcmd -S $SQLDB_IP,10433 -U sa -P 'Your@Password'
CREATE DATABASE sql_guestbook;
USE sql_guestbook;
CREATE TABLE guestlog (entrydate DATETIME, name NVARCHAR(30), phone NVARCHAR(30), message TEXT, sentiment_score NVARCHAR(30));
INSERT INTO guestlog VALUES ('2017-5-2 23:59:59', 'anonymous', '12158379120', 'Get busy living, or get busy dying', '0.9950121');
```

### Golang web app

Simple web page that connects to SQL Server and builds a table. Uses the following environment variables:
* SQLSERVER
* SQLPORT
* SQLID
* SQLPWD
* SQLDB
* GIT_SHA

Installation/upgrade will occur with Jenkins pipeline.

### Setup Jenkins Pipeline

* Open Jenkins Blue Ocean
* Add Github organization

## Running the Demo

* Make code changes in dev branch
* Git commit/add
* Submit PR
* Merge with Master
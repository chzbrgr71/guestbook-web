# Guestbook CI/CD Demo with Jenkins
This app demonstrates a Jenkins pipeline with a Golang based web app and a SQL Server on Linux container. Containers are deployed to ACS Kubernetes via Helm charts and images are stored in Azure Container Registry

## Demo Setup

### Setup ACS Kubernetes

### Azure Container Registry

### Install Jenkins

### Database setup

* Helm chart install: ```helm install --name=guestbook-db ./charts/guestbook-db```
* Wait for IP: ```watch kubectl get svc guestbook-db-guestbook-db```
* Get IP: ```export SQLDB_IP=$(kubectl get svc guestbook-db-guestbook-db --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")```
* Seed data (use sqlcmd CLI tool)
```
sqlcmd -S $SQLDB_IP,10433 -U sa -P 'Your@Password'
CREATE DATABASE sql_guestbook;
USE sql_guestbook;
CREATE TABLE guestlog (entrydate DATETIME, name NVARCHAR(30), phone NVARCHAR(30), message TEXT, sentiment_score NVARCHAR(30));
INSERT INTO guestlog VALUES ('2017-5-2 23:59:59', 'anonymous', '12158379120', 'Get busy living, or get busy dying', '0.9950121');
```

## Golang web app

Helm chart install: ```helm install --name=guestbook-web ./charts/guestbook-web```

Simple web page that connects to SQL Server and builds a table. Uses the following environment variables:
* SQLSERVER
* SQLPORT
* SQLID
* SQLPWD
* SQLDB
* GIT_SHA
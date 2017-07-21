# Guestbook CI/CD Demo with Jenkins
This app demonstrates a Jenkins pipeline with a Golang based web app and a SQL Server on Linux container. Containers are deployed to ACS Kubernetes via Helm charts and images are stored in Azure Container Registry

## Demo Setup


## Golang web app

Simple web page that connects to SQL Server and builds a table. Uses the following environment variables:

* SQLSERVER
* SQLPORT
* SQLID
* SQLPWD
* SQLDB
* GIT_SHA

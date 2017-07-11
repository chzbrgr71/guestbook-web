# Guestbook Web App - Go
This simple go app works with some other demos. Connects to a SQL Server DB running in Kubernetes and is deployed via Helm chart

## Container Setup

Image: chzbrgr71/go-guestbook:<tag>

Environment variables:

* SQLSERVER
* SQLPORT
* SQLID
* SQLPWD
* SQLDB

## Docker

```
docker build --build-arg VCS_REF=brian999 -t chzbrgr71/go-guestbook .

docker run -d -e "SQLSERVER=23.99.10.5" -e "SQLPORT=10433" -e "SQLID=sa" -e "SQLPWD=Pass@word" -e "SQLDB=sql_guestbook" --name web -p 80:8080 chzbrgr71/go-guestbook
```

## Helm Chart
test

# Guestbook CI/CD Demo with Jenkins
This app demonstrates a Jenkins pipeline with a Golang based web app and a SQL Server on Linux container. Containers are deployed to ACS Kubernetes via Helm charts and images are stored in Azure Container Registry.

Thank you to Lachie Evenson for helping with this. Much of the demo is reverse engineered from the infamous [Croc Hunter demo.](https://github.com/lachie83/croc-hunter)

## Demo Setup

1. Fork or clone this repo

2. Setup ACS Kubernetes/Helm

    * Use standard Azure Container Service instructions [here.](https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-kubernetes-walkthrough)  
    * Ensure kubectl is installed on your local machine and you have the proper kube config file (~/.kube/config)
    * Install Helm 
        ```
        # MAC OS
        brew install kubernetes-helm 
        helm init
        helm repo update
        ```

3. Azure Container Registry

    * Use standard Azure Container Registry instructions [here.](https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-tutorial-kubernetes-prepare-acr) 
    * Add Kubernetes secret with ACR creds base64 encoded. Update secret-update.yaml with your values
        ```
        kubectl create -f secret.yaml
        ```
4. Add Infrastructure Stuff
    * Install Kube Lego chart
        ```
        helm install stable/kube-lego --set config.LEGO_EMAIL=<valid-email>,config.LEGO_URL=https://acme-v01.api.letsencrypt.org/directory
        ```
    * Install Nginx ingress chart
        ```
        helm install stable/nginx-ingress

        Follow the notes from helm status to determine the external IP of the nginx-ingress service
        ```
    * Add a DNS entry with your provider and point it do the external IP
        ```
        blah.test.com in A <nginx ingress svc external-IP>

        or *.test.com in A <nginx ingress svc external-IP>
        ```

5. Install Jenkins

    * Update jenkins-values.yaml. Replace brianredmond.co with the domain name setup above.
    * Install Jenkins helm chart
        ```
        helm --namespace jenkins --name jenkins -f ./jenkins-values.yaml install stable/jenkins
        ```
    * Once pod is up and running, browse to http://jenkins.brianredmond.co/login [replace with your domain name]
    * Get Jenkins password:
        ```
        kubectl get secret --namespace jenkins jenkins-jenkins -o jsonpath="{.data.jenkins-admin-password}" | base64 --decode
        ```
    * Add ACR creds in Jenkins Global Credentials. Credentials > Jenkins > Global credentials > Add Credentials

6. Database setup

    * Helm chart install
        ```
        helm install --name=guestbook-db ./charts/guestbook-db
        ```

    * Get IP for DB endpoint (wait patiently)
        ```
        watch kubectl get svc guestbook-db-guestbook-db
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

7. Setup Jenkins Pipeline

    * Open Jenkins Blue Ocean
    * Add Github organization. A Github access token will need to be created for this step

8. Setup Github webhook in your repo

    * The webhook will direct to http://jenkins-url/github-webhook/ 


## Running the Demo

The process for the demo is to make a change in the dev branch to kickoff an initial CD build and test. From there a PR is submitted and the CI process runs.

    * Checkout dev branch
    ```
    git checkout dev
    git merge master #if out of date
    ```
    * Make code changes and commit
    ```
    # After a simple update to the sqlguestbook.go code
    git add .
    git commit -m "Updated web UI"
    git push
    ```
    * After pipeline is run for dev branch, submit a PR in Github repo
    * PR pipeline will run
    * Merge PR in Github and the master branch will build and deploy
    * Validate updated web app
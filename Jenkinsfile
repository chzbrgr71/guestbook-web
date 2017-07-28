#!/usr/bin/groovy

podTemplate(label: 'jenkins-pipeline', containers: [
    containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:2.62', args: '${computer.jnlpmac} ${computer.name}', workingDir: '/home/jenkins', resourceRequestCpu: '200m', resourceLimitCpu: '200m', resourceRequestMemory: '256Mi', resourceLimitMemory: '256Mi'),
    containerTemplate(name: 'docker', image: 'docker:1.12.6',       command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'golang', image: 'golang:1.7.5', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'helm', image: 'lachlanevenson/k8s-helm:v2.4.1', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'kubectl', image: 'lachlanevenson/k8s-kubectl:v1.4.8', command: 'cat', ttyEnabled: true)
],
volumes:[
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
])
    {
        node ('jenkins-pipeline') {
            println "DEBUG: Pipeline starting"
            def pwd = pwd()
            def chart_dir = "${pwd}/charts/guestbook-web"

            checkout scm
            
            // read in required jenkins workflow config values
            def inputFile = readFile('Jenkinsfile.json')
            def config = new groovy.json.JsonSlurperClassic().parseText(inputFile)
            println "DEBUG: pipeline config ==> ${config}"
            
            // set additional git envvars for image tagging
            gitEnvVars()
            
            def acct = getContainerRepoAcct(config)

            // tag image with version, and branch-commit_id
            def image_tags_map = getContainerTags(config)

            // compile tag list
            def image_tags_list = getMapValues(image_tags_map)

            stage ('BUILD: code compile and test') {

                container('golang') {
                    sh "go get github.com/denisenkom/go-mssqldb"
                    sh "go build"
                    sh "go test -v"
                }
            }

            stage ('TEST: k8s deployment') {

                container('helm') {
                    // run helm chart linter
                    helmLint(chart_dir)
                }
            }

            stage ('BUILD: containerize and publish') {

                container('docker') {
                    // Login to ACR
                    withCredentials([[$class          : 'UsernamePasswordMultiBinding', credentialsId: config.container_repo.jenkins_creds_id,
                                    usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                        sh "docker login briarregistry-microsoft.azurecr.io -u briarregistry -p 5XuGNf=zidx=/K46X7ig/BKKvO9nGIrE"
                    }

                    // build and publish container
                    containerBuildPub(
                        dockerfile: config.container_repo.dockerfile,
                        host      : config.container_repo.host,
                        acct      : acct,
                        repo      : config.container_repo.repo,
                        tags      : image_tags_list,
                        auth_id   : config.container_repo.jenkins_creds_id
                    )
                }
            }

            stage ('SECURE: scan container images for vulnerabilities') {
                println "DEBUG: Run vulnerability scan of container images in repo"
    
            }
            
            stage ('DEPLOY: helm release to k8s') {
      
                container('helm') {
                    // Deploy using Helm chart
                    helmDeploy(
                        dry_run       : false,
                        name          : config.app.name,
                        namespace     : config.app.name,
                        version_tag   : image_tags_list.get(0),
                        chart_dir     : chart_dir,
                        replicas      : config.app.replicas,
                        cpu           : config.app.cpu,
                        memory        : config.app.memory,
                        hostname      : config.app.hostname
                    )
                }
            }
            
            println "DEBUG: FINISHED"
        }
    
        
        
             
    }

def kubectlTest() {
    // Test that kubectl can correctly communication with the Kubernetes API
    println "checking kubectl connnectivity to the API"
    sh "kubectl get nodes"

}

def helmLint(String chart_dir) {
    // lint helm chart
    println "running helm lint ${chart_dir}"
    sh "helm lint ${chart_dir}"

}

def helmConfig() {
    //setup helm connectivity to Kubernetes API and Tiller
    println "initiliazing helm client"
    sh "helm init"
    println "checking client/server version"
    sh "helm version"
}


def helmDeploy(Map args) {
    //configure helm client and confirm tiller process is installed
    //helmConfig()

    if (args.dry_run) {
        println "Running dry-run deployment"
        helmConfig()
        //sh "helm upgrade --dry-run --install ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory} --namespace=${args.namespace}"
        sh "helm upgrade --dry-run --install ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory}"
    } else {
        println "Running deployment"
        println "CMD: helm upgrade --install --wait ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory}"
        //sh "helm upgrade --install --wait ${args.name} ${args.chart_dir} --set imageTag=${args.version_tag},replicas=${args.replicas},cpu=${args.cpu},memory=${args.memory} --namespace=${args.namespace}"
        sh "helm upgrade guestbook ${args.chart_dir} --set imageTag=${args.version_tag}"

        echo "Application ${args.name} successfully deployed. Use helm status ${args.name} to check"
    }
}

def helmDelete(Map args) {
        println "Running helm delete ${args.name}"

        sh "helm delete ${args.name}"
}

def helmTest(Map args) {
    println "Running Helm test"

    sh "helm test ${args.name} --cleanup"
}

def gitEnvVars() {
    // create git envvars
    println "Setting envvars to tag container"

    sh 'git rev-parse HEAD > git_commit_id.txt'
    try {
        env.GIT_COMMIT_ID = readFile('git_commit_id.txt').trim()
        env.GIT_SHA = env.GIT_COMMIT_ID.substring(0, 7)
    } catch (e) {
        error "${e}"
    }
    println "env.GIT_COMMIT_ID ==> ${env.GIT_COMMIT_ID}"

    sh 'git config --get remote.origin.url> git_remote_origin_url.txt'
    try {
        env.GIT_REMOTE_URL = readFile('git_remote_origin_url.txt').trim()
    } catch (e) {
        error "${e}"
    }
    println "env.GIT_REMOTE_URL ==> ${env.GIT_REMOTE_URL}"
}


def containerBuildPub(Map args) {

    println "Running Docker build/publish: ${args.host}/${args.acct}/${args.repo}:${args.tags}"
    //Running Docker build/publish: briarregistry-microsoft.azurecr.io/chzbrgr71/go-guestbook:[master-ff6bc56, latest]
      
    docker.withRegistry("https://${args.host}", "${args.auth_id}") {

        // def img = docker.build("${args.acct}/${args.repo}", args.dockerfile)
        def img = docker.image("${args.acct}/${args.repo}")
        //println "${args.acct}/${args.repo}"
        //sh "docker build --build-arg VCS_REF=${env.GIT_SHA} --build-arg BUILD_DATE=`date -u +'%Y-%m-%dT%H:%M:%SZ'` -t ${args.acct}/${args.repo} ${args.dockerfile}"
        sh "docker build --build-arg VCS_REF=${env.GIT_SHA} -t ${args.acct}/${args.repo} ${args.dockerfile}"        
        //for (int i = 0; i < args.tags.size(); i++) {
            //img.push(args.tags.get(i))
        //}
        img.push(args.tags.get(0))
        return img.id
    }
}

def getContainerTags(config, Map tags = [:]) {

    println "getting list of tags for container"
    def String commit_tag
    def String version_tag

    try {
        // if PR branch tag with only branch name
        if (env.BRANCH_NAME.contains('PR')) {
            commit_tag = env.BRANCH_NAME
            tags << ['commit': commit_tag]
            return tags
        }
    } catch (Exception e) {
        println "WARNING: commit unavailable from env. ${e}"
    }

    // commit tag
    try {
        // if branch available, use as prefix, otherwise only commit hash
        if (env.BRANCH_NAME) {
            commit_tag = env.BRANCH_NAME + '-' + env.GIT_COMMIT_ID.substring(0, 7)
        } else {
            commit_tag = env.GIT_COMMIT_ID.substring(0, 7)
        }
        tags << ['commit': commit_tag]
    } catch (Exception e) {
        println "WARNING: commit unavailable from env. ${e}"
    }

    // master tag
    try {
        if (env.BRANCH_NAME == 'master') {
            tags << ['master': 'latest']
        }
    } catch (Exception e) {
        println "WARNING: branch unavailable from env. ${e}"
    }

    // build tag only if none of the above are available
    if (!tags) {
        try {
            tags << ['build': env.BUILD_TAG]
        } catch (Exception e) {
            println "WARNING: build tag unavailable from config.project. ${e}"
        }
    }

    return tags
}

def getContainerRepoAcct(config) {

    println "setting container registry creds according to Jenkinsfile.json"
    def String acct

    if (env.BRANCH_NAME == 'master') {
        acct = config.container_repo.master_acct
    } else {
        acct = config.container_repo.alt_acct
    }

    return acct
}

@NonCPS
def getMapValues(Map map=[:]) {
    // jenkins and workflow restriction force this function instead of map.values(): https://issues.jenkins-ci.org/browse/JENKINS-27421
    def entries = []
    def map_values = []

    entries.addAll(map.entrySet())

    for (int i=0; i < entries.size(); i++){
        String value =  entries.get(i).value
        map_values.add(value)
    }

    return map_values
}  
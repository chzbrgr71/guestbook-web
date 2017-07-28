#!/usr/bin/groovy

// load pipeline functions
@Library('./PipelineFx.groovy')
def pipeline = new co.brianredmond.jenkins.PipelineFx()

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
            println "Pipeline starting"
            def pwd = pwd()
            def chart_dir = "${pwd}/charts/guestbook-web"

            checkout scm
            
            // read in required jenkins workflow config values
            def inputFile = readFile('Jenkinsfile.json')
            def config = new groovy.json.JsonSlurperClassic().parseText(inputFile)
            println "pipeline config ==> ${config}"
            
            // set additional git envvars for image tagging
            pipeline.gitEnvVars()
            
            println "FINISHED"
            
        }
    
        
        
        
        
        
        
        
    }
    
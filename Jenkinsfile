// Declarative syntax
pipeline {
    agent any
    stages {
        stage ('Info') {
            parallel {
                stage ('System') {
                    steps {
                        sh "hostname"
                        sh "uptime"
                        sh "free -mh"
                        sh "df -h"
                        sh "pwd"
                        sh "git --version"
                        sh "make --version"
                    }
                }
                stage ('Docker') {
                    steps {
                        sh "docker --version"
                        sh "docker-compose --version"
                        sh "docker system df"
                        sh "docker images"
                        sh "docker ps -a"
                    }
                }
            }
        }
        stage ('Release') {
            when {
                anyOf {
                    branch 'master'
                    buildingTag()
                }
            }
            steps {
                script {
                    docker.withRegistry('', 'kisiodigital-user-dockerhub') {
                        sh "make release"
                    }
                }
            }
        }
        stage ('CD (dev)') {
            when {
                branch 'master'
            }
            steps {
                script {
                    def handle = triggerRemoteJob(
                        remoteJenkinsName: 'jenkins-deployment',
                        job: 'pad_deploy',
                        parameters: 'ENVIRONMENT=dev\nVERSION=latest',
                        blockBuildUntilComplete: true,
                        preventRemoteBuildQueue: true,
                    )
                    echo 'Remote Status: ' + handle.getBuildStatus().toString()
                }
            }
        }
    }
    post {
        always {
            deleteDir()
        }
    }
}

pipeline {
    agent any

    post {
        failure {
            updateGitlabCommitStatus name: 'build', state: 'failed'

        }
        success {
            updateGitlabCommitStatus name: 'build', state: 'success'
            mail bcc: '', body: "<b>Build Succeeded</b><br>Project: ${env.JOB_NAME} <br>Build Number: ${env.BUILD_NUMBER} <br> URL de build: ${env.BUILD_URL}", cc: '', charset: 'UTF-8', from: 'jenkins@stickybits.red', mimeType: 'text/html', replyTo: '', subject: "ERROR CI: Project name -> ${env.JOB_NAME}", to: "jenkins";  
        }
    }
    options {
        gitLabConnection('Jenkins')
    }

    stages {
        stage('Prepare') {
            steps {
                updateGitlabCommitStatus name: 'build', state: 'running'
                sh 'rm -rf releases'
                sh 'mkdir -p releases'
                echo "[*] Created releases directory"
            }
        }
        stage('Build for linux/amd64') {
            steps {
                sh 'GOOS=linux GOARCH=amd64 go build -o releases/agent main.go '
                echo "[*] Agent has been built"
            }
        }
    }
}

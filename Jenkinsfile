pipeline {
    agent any

    post {
        failure {
            updateGitlabCommitStatus name: 'build', state: 'failed'
        }
        success {
            updateGitlabCommitStatus name: 'build', state: 'success'
        }
    }
    options {
        gitLabConnection('GitLab')
    }

    stages {
        stage('Prepare') {
            steps {
                cleanWs()
                updateGitlabCommitStatus name: 'build', state: 'running'
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
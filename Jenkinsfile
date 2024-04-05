pipeline {
    agent any

    post {
        failure {
            updateGitlabCommitStatus name: 'build', state: 'failed'

        }
        success {
            updateGitlabCommitStatus name: 'build', state: 'success'
            matrixSendMessage hostname: 'matrix.org', accessTokenCredentialsId: 'syt_czR0bw_GUrAjypEJJhOeDWypmqB_1Hrqre', roomId: 'Jenkins', body: '[Jenkins] Agent build successful!', formattedBody: '<b>[Jenkins]</b> Agent build successful!'
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
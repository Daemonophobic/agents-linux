pipeline {
    agent any

    stages {
        stage('Prepare') {
            steps {
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
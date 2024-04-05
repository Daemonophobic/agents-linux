pipeline {
    agent any

    stages {
        stage('Prepare') {
            steps {
                sh 'mkdir releases'
                echo "Created releases direction"
            }
        }
        stage('Build for linux/amd64') {
            steps {
                sh 'GOOS=linux GOARCH=amd64 go build main.go -o releases/agent'
                echo "Agent has been built"
            }
        }
    }
}
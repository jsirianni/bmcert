#!/usr/bin/env groovy
import groovy.json.JsonOutput
def disableSlack = false
def slackNotificationChannel = 'devops-jenkins'     // ex: = "builds"
def notifySlack(text, channel, attachments, slackURL) {

    def jenkinsIcon = 'https://wiki.jenkins-ci.org/download/attachments/2916393/logo.png'

    def payload = JsonOutput.toJson([text: text,
                                     channel: channel,
                                     username: "Jenkins",
                                     icon_url: jenkinsIcon,
                                     attachments: attachments
    ])

    sh "curl -X POST --data-urlencode \'payload=${payload}\' ${slackURL}"
}

pipeline {
    agent {
        label "docker"
    }
    environment {
        SLACK_HOOK_URL = credentials('SLACK_HOOK_URL')
        VAULT_ADDR = credentials('VAULT_ADDR')
        VAULT_SKIP_VERIFY = 1
        VAULT_CERT_URL = credentials('VAULT_CERT_URL')
        VAULT_GITHUB_TOKEN = credentials('GITHUB_TOKEN')
    }
    stages {
        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }
        stage('Run build script') {
            steps {
                    sh '''
            rm -f *.zip
            ./build.sh
            '''
            }
        }
        stage('Test create certificate') {
            steps {
                    sh '''
            export VAULT_TOKEN=`vault login -no-store -token-only -method=github token=$VAULT_GITHUB_TOKEN`
            unzip bmcert-*-linux-amd64.zip
            ./bmcert create --hostname test.bluemedora.localnet --tls-skip-verify
            openssl x509 -in zk-ref-c1-2.bluemedora.localnet.pem -text -noout >> /dev/null
            '''
            }
        }
    }
    post {
        always {
            script {
                //populateGlobalVariables()
                def buildColor = currentBuild.result == null ? "good" : "warning"
                def buildStatus = currentBuild.result == null ? "Success" : currentBuild.result
                def jobName = "${env.JOB_NAME}"
                def slack_hook_url = "${env.SLACK_HOOK_URL}"

                // Strip the branch name out of the job name (ex: "Job Name/branch1" -> "Job Name")
                jobName = jobName.getAt(0..(jobName.indexOf('/') - 0))

                echo jobName
                if (disableSlack == false){
                if (buildStatus == "Failed") {
                    buildStatus = "Failed"
                    buildColor = "danger"
                    notifySlack("", slackNotificationChannel, [[
                                                                       title      : "${jobName}, build #${env.BUILD_NUMBER}",
                                                                       title_link : "${env.BUILD_URL}",
                                                                       color      : "${buildColor}",
                                                                       text       : "${buildStatus}\n",
                                                                       "mrkdwn_in": ["fields"],
                                                                       fields     : [
                                                                               [
                                                                                       title: "Node Name",
                                                                                       value: "${env.NODE_NAME}",
                                                                                       short: true
                                                                               ]

                                                                       ]
                                                               ],
                                                               [
                                                                       title      : "Failed Tests",
                                                                       color      : "${buildColor}",
                                                                       text       : "Build Failed",
                                                                       "mrkdwn_in": ["text"],
                                                               ]], "${slack_hook_url}")
                } else {
                    notifySlack("", slackNotificationChannel, [[
                                                                       title     : "${jobName}, build #${env.BUILD_NUMBER}",
                                                                       title_link: "${env.BUILD_URL}",
                                                                       color     : "${buildColor}",
                                                                       text      : "${buildStatus}\n",
                                                                       fields    : [
                                                                               [
                                                                                       title: "Node Name",
                                                                                       value: "${env.NODE_NAME}",
                                                                                       short: true
                                                                               ]
                                                                       ]
                                                               ]], "${slack_hook_url}")
                }
                }else{
                 echo "slack disabled"
                echo "Build status ${buildStatus}"
                }
            }

        }
    }
}

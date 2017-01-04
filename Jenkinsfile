#!/usr/bin/env groovy

// Based on https://stackoverflow.com/questions/36837683/how-to-perform-actions-for-failed-builds-in-jenkinsfile
def task_wrapper(String label, Closure body) {
     node(label) {
        wrap([$class: 'TimestamperBuildWrapper']) {
            // This borks:
            // "Scripts not permitted to use staticMethod org.codehaus.groovy.runtime.DefaultGroovyMethods minus java.lang.String java.lang.Object"
            //
            // def branch_url = "${env.BUILD_URL}" - "${env.BUILD_NUMBER}/"
            def branch_url = env.BUILD_URL.substring(0, env.BUILD_URL.length() - (env.BUILD_NUMBER.length() +1))
            def msg_body = "Build url: ${env.BUILD_URL}\nBranch url: ${branch_url}"
            def msg = "Build number `${env.BUILD_NUMBER}` succeded."
            def color = "good"
            def prefix = "[SUCCESS]"

            try {
                body.call()
            } catch (Exception e) {
                msg = "Build number `${env.BUILD_NUMBER}` failed."
                color = "danger"
                prefix = "[FAILURE]"

                throw e; // rethrow so the build is considered failed
            } finally {
                withCredentials([[$class: 'StringBinding',
                                  credentialsId: '8b793652-f26a-422f-a9ba-0d1e47eb9d89',
                                  variable: 'SLACK_TOKEN']
                                 ]) {
                    slackSend (channel: '#dcos-security-ci',
                        message: "`${env.JOB_NAME}` ${msg}\n\n${msg_body}",
                        teamDomain: 'mesosphere',
                        token: "${env.SLACK_TOKEN}",
                        color: color,
                        )
                }

                mail (subject: "[${env.JOB_NAME}]${prefix} ${msg}",
                    body: "${msg_body}",
                    to: 'dcos-security-ci@mesosphere.io')
            }
        }
    }
}

task_wrapper ('mesos'){
    // Unfortunatelly we cannot use flyweight executor here:
    // https://issues.jenkins-ci.org/browse/JENKINS-27386
    // and esp. https://issues.jenkins-ci.org/browse/JENKINS-32314
    stage "Verify author"
        def authed_users = ["ktf", "dberzano"]

        echo "Changeset from `" + env.CHANGE_AUTHOR + "`"

        timeout(time: 24, unit: 'HOURS') {
            if (authed_users.contains(env.CHANGE_AUTHOR)) {
                // Let's not spam our slack channels with too much info.
                echo "PR comes from authorized user, testing it now!"
            } else {
                withCredentials([[$class: 'StringBinding',
                                    credentialsId: '8b793652-f26a-422f-a9ba-0d1e47eb9d89',
                                    variable: 'SLACK_TOKEN']
                                    ]) {
                    def branch_url = env.BUILD_URL.substring(0, env.BUILD_URL.length() - (env.BUILD_NUMBER.length() +1))
                    slackSend (channel: '#dcos-security-ci',
                        message: "`${env.JOB_NAME}` has a new build from user `${env.CHANGE_AUTHOR}` waiting for ACK\n Build URL: `${branch_url}`",
                        teamDomain: 'mesosphere',
                        token: "${env.SLACK_TOKEN}",
                        color: "danger",
                        )
                    userInput = input(message: 'Build the PR ?')
                }
            }
        }

    stage 'Cleanup workspace'
        deleteDir()

    stage 'Checkout'
        checkout scm

    load 'Jenkinsfile-insecure.groovy'
}

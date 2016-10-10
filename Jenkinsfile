#!/usr/bin/env groovy

node ('mesos'){
    stage 'Checkout'
        checkout scm

        // http://stackoverflow.com/questions/35554983/git-variables-in-jenkins-workflow-plugin
        // https://issues.jenkins-ci.org/browse/JENKINS-35230
        def gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

        def nameByBranch = "mesosphere/vault-devkit:${env.BRANCH_NAME}"
        def nameByCommit = "mesosphere/vault-devkit:${gitCommit}"

        // Some debugging:
        sh 'env | sort '
        echo "Name of the container to be published, by branch: ${nameByBranch}"
        echo "Name of the container to be published, by commit: ${nameByCommit}"

    stage 'Prepare devkit'
        sh 'make update-devkit'

    try {
        stage 'Prepare aux-containers'
            sh 'make aux'

        stage 'make testplain'
            sh 'make testplain'

        stage 'make testrace'
            sh 'make testplain'

        stage 'make build'
            sh 'make build'

        stage 'Build mesosphere/vault container'

            sh "docker build --rm --force-rm -t ${nameByBranch} -f ./docker/Dockerfile.publish ./"
            sh "docker tag ${nameByBranch} ${nameByCommit}"

        stage 'Push to docker registry'
            sh 'echo noop'

    } finally {
        stage 'Cleanup docker containers'
            sh 'make clean-containers'
            sh "docker rmi -f ${nameByBranch} || true"
            sh "docker rmi -f ${nameByCommit} || true"
    }
}

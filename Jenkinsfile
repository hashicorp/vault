#!/usr/bin/env groovy

@Library('sec_ci_libs@gh_syncing')
import task_wrapper
import user_is_authorized


task_wrapper('mesos'){
    stage("Verify author") {
        user_is_authorized()
    }

    stage('Cleanup workspace') {
        deleteDir()
    }

    stage('Checkout') {
        checkout scm
    }

    load 'Jenkinsfile-insecure.groovy'
}

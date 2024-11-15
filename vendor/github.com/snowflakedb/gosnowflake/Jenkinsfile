import groovy.json.JsonOutput


timestamps {
  node('regular-memory-node') {
    stage('checkout') {
      scmInfo = checkout scm
      println("${scmInfo}")
      env.GIT_BRANCH = scmInfo.GIT_BRANCH
      env.GIT_COMMIT = scmInfo.GIT_COMMIT
    }
    params = [
      string(name: 'svn_revision', value: 'main'),
      string(name: 'branch', value: 'main'),
      string(name: 'client_git_commit', value: scmInfo.GIT_COMMIT),
      string(name: 'client_git_branch', value: scmInfo.GIT_BRANCH),
      string(name: 'TARGET_DOCKER_TEST_IMAGE', value: 'go-centos7-go1.21'),
      string(name: 'parent_job', value: env.JOB_NAME),
      string(name: 'parent_build_number', value: env.BUILD_NUMBER)
    ]
    stage('Test') {
      build job: 'RT-LanguageGo-PC',parameters: params
    }
  }
}


pipeline {
  agent { label 'regular-memory-node' }
  options { timestamps() }
  environment {
    COMMIT_SHA_LONG = sh(returnStdout: true, script: "echo \$(git rev-parse " + "HEAD)").trim()

    // environment variables for semgrep_agent (for findings / analytics page)
    // remove .git at the end
    // remove SCM URL + .git at the end

    BASELINE_BRANCH = "${env.CHANGE_TARGET}"
  }
  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }
  }
}

def wgetUpdateGithub(String state, String folder, String targetUrl, String seconds) {
    def ghURL = "https://api.github.com/repos/snowflakedb/gosnowflake/statuses/$COMMIT_SHA_LONG"
    def data = JsonOutput.toJson([state: "${state}", context: "jenkins/${folder}",target_url: "${targetUrl}"])
    sh "wget ${ghURL} --spider -q --header='Authorization: token $GIT_PASSWORD' --post-data='${data}'"
}

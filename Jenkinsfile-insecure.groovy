stage 'Calculate git-related vars'
    // http://stackoverflow.com/questions/35554983/git-variables-in-jenkins-workflow-plugin
    // https://issues.jenkins-ci.org/browse/JENKINS-35230
    def gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

    def nameByBranch = "mesosphereci/vault:${env.BRANCH_NAME}"
    def nameByCommit = "mesosphereci/vault:${gitCommit}"

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
        sh 'make testrace'

    stage 'make build'
        sh 'make build'

    stage 'Build mesosphereci/vault container'

        sh "docker build --rm --force-rm -t ${nameByBranch} -f ./Dockerfile.publish ./"
        sh "docker tag ${nameByBranch} ${nameByCommit}"

    stage 'Push to docker registry'
        withCredentials(
        [[$class: 'StringBinding',
            credentialsId: '7bdd2775-2911-41ba-918f-59c8ae52326d',
            variable: 'DOCKER_HUB_USERNAME'],
            [$class: 'StringBinding',
            credentialsId: '42f2e3fb-3f4f-47b2-a128-10ac6d0f6825',
            variable: 'DOCKER_HUB_PASSWORD'],
            [$class: 'StringBinding',
            credentialsId: '4551c307-10ae-40f9-a0ac-f1bb44206b5b',
            variable: 'DOCKER_HUB_EMAIL']
        ]) {
            sh "docker login -u '${env.DOCKER_HUB_USERNAME}' -p '${env.DOCKER_HUB_PASSWORD}'"
        }
        sh "docker push ${nameByBranch}"
        sh "docker push ${nameByCommit}"

} finally {
    stage 'Cleanup docker containers'
        sh 'make clean'
        sh "docker rmi -f ${nameByBranch} || true"
        sh "docker rmi -f ${nameByCommit} || true"
}

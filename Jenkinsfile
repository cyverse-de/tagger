#!groovy

jenkinsCredentials = [[$class:           'UsernamePasswordMultiBinding',
                       credentialsId:    'jenkins-docker-credentials',
                       passwordVariable: 'DOCKER_PASSWORD',
                       usernameVariable: 'DOCKER_USERNAME']]

// Build the service.
def build(dockerRepo, git_commit, descriptive_version) {
    sh """docker build --pull --no-cache --rm \\
                 --build-arg git_commit=${git_commit} \\
                 --build-arg descriptive_version=${descriptive_version} \\
                 -t ${dockerRepo} ."""

    image_sha = sh(returnStdout: true, script: "docker inspect -f '{{ .Config.Image }}' ${dockerRepo}").trim()
    echo image_sha

    writeFile(file: "${dockerRepo}.docker-image-sha", text: "${image_sha}")
    fingerprint "${dockerRepo}.docker-image-sha"
}

// Run the service test.
def test(dockerRepo, dockerTestRunner, service) {
    sh """docker run --rm --name ${dockerTestRunner} \\
                 --entrypoint 'go' \\
                 ${dockerRepo} test github.com/cyverse-de/${service.repo}/..."""
}

// Push the docker image.
def dockerPush(dockerRepo, dockerPusher, dockerPushRepo, service) {
    lock("docker-push-${dockerPushRepo}") {
        milestone 101

        sh "docker tag ${dockerRepo} ${dockerPushRepo}"

        withCredentials(jenkinsCredentials) {
            sh """docker run -e DOCKER_USERNAME -e DOCKER_PASSWORD \\
                                -v /var/run/docker.sock:/var/run/docker.sock \\
                                --rm --name ${dockerPusher} \\
                                docker:\$(docker version --format '{{ .Server.Version }}') \\
                                sh -e -c \\
                         'docker login -u \"\$DOCKER_USERNAME\" -p \"\$DOCKER_PASSWORD\" && \\
                          docker push ${dockerPushRepo} && \\
                          docker rmi ${dockerPushRepo} && \\
                          docker logout'"""
        }
    }
}

// Build, test and push the docker image.
def buildJob(slackJobDescription) {
    checkout scm

    dockerRepo = "test-${env.BUILD_TAG}"
    service = readProperties file: 'service.properties'

    git_commit = sh(returnStdout: true, script: "git rev-parse HEAD").trim()
    echo git_commit

    descriptive_version = sh(returnStdout: true, script: 'git describe --long --tags --dirty --always').trim()
    echo descriptive_version

    stage("Build") {
        build(dockerRepo, git_commit, descriptive_version)
    }

    dockerTestRunner = "test-${env.BUILD_TAG}"
    dockerPusher = "push-${env.BUILD_TAG}"
    try {
        stage("Test") {
            test(dockerRepo, dockerTestRunner, service)
        }

        milestone 100
        stage("Docker Push") {
            dockerPushRepo = "${service.dockerUser}/${service.repo}:${env.BRANCH_NAME}"
            dockerPush(dockerRepo, dockerPusher, dockerPushRepo, service)
        }
    } finally {
        sh returnStatus: true, script: "docker kill ${dockerTestRunner}"
        sh returnStatus: true, script: "docker rm ${dockerTestRunner}"

        sh returnStatus: true, script: "docker kill ${dockerPusher}"
        sh returnStatus: true, script: "docker rm ${dockerPusher}"

        sh returnStatus: true, script: "docker rmi ${dockerRepo}"

        sh returnStatus: true, script: "docker rmi \$(docker images -qf 'dangling=true')"

        step([$class:        'hudson.plugins.jira.JiraIssueUpdater',
              issueSelector: [$class: 'hudson.plugins.jira.selector.DefaultIssueSelector'],
              scm:           scm,
              labels:        [ "${service.repo}-${descriptive_version}" ]])
    }
}

node('docker') {
    slackJobDescription = "job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})"
    try {
        buildJob(slackJobDescription)
    } catch (InterruptedException e) {
        currentBuild.result = 'ABORTED'
        slackSend color: 'warning', message: "ABORTED: ${slackJobDescription}"
        throw e
    } catch (e) {
        currentBuild.result = 'FAILED'
        sh "echo ${e}"
        slackSend color: 'danger', message: "FAILED: ${slackJobDescription}"
        throw e
    }
}

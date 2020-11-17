node {
     checkout scm
     env.VERSION_STRING = "0.01-dev-${env.BUILD_ID}"
     GIT_TAG = sh(returnStdout: true, script: "git tag --sort version:refname | tail -1").trim()
     if(GIT_TAG) {
         env.VERSION_STRING = GIT_TAG
     }
     docker.image('golang:1.15').inside {
         withEnv(["HOME=${env.WORKSPACE}"]) {
             stage 'build'
             sh 'go build -o echoapp'
             stash includes: 'echoapp', name: 'binary'
         }
     }
     docker.image('cdrx/fpm-ubuntu:18.04').inside {
         stage 'package'
         unstash 'binary'
         sh 'fpm -s dir -t deb -n echoapp -v \$VERSION_STRING dist/linux-amd64=/opt/jwt-tokenissuer'
         stash includes: '*.deb', name: 'package'
         archive includes: '*.deb'
     }
}

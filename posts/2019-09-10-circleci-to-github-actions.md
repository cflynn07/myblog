I decided to switch the CI/CD pipeline of this blog from CircleCI to Github
Actions just to check it out. I know it's basically overkill to have a CI/CD
setup for a blog, but I did it for curiosity. One of the most obvious
advantages of Github Actions is its tight integration with Github. If nothing
else, it's nice just to have your repository and CI information together on one
website.

My CI/CD flow is:
<ol>
<li>Run tests & upload coverage reports to CodeCov</li>
<li>Build & push a docker image</li>
<li>Deploy to kubernetes cluster, either staging or production, using helm</li>
</ol>

Previously, with CircleCI my configuration file looked like this.
#### .circleci/config.yml
<pre class="prettyprint linenums">
version: 2
jobs:
  test:
    docker:
      - image: circleci/golang:1.12
    working_directory: /go/src/myblog
    steps:
      - checkout #checks out the source code to working directory
      - run:
          name: Run E2E tests
          command: go test ./... -coverprofile=coverage.txt -covermode=atomic -coverpkg=myblog/app
      - run:
          name: Upload coverage report to codecov
          command: bash <(curl -s https://codecov.io/bash)
  build_and_deploy:
    docker:
      - image: cflynnus/google-cloud-sdk-helm:v1
    environment:
      - PROJECT_NAME: "blog"
      - GOOGLE_PROJECT_ID: "blog-44444"
      - GOOGLE_COMPUTE_ZONE: "us-central1-b"
      - GOOGLE_CLUSTER_NAME: "blog-cluster"
    steps:
      - checkout
      - run:
          name: Setup Google Cloud SDK
          command: |
            apt-get install -qq -y gettext
            echo $GCLOUD_SERVICE_KEY > ${HOME}/gcloud-service-key.json
            gcloud auth activate-service-account --key-file=${HOME}/gcloud-service-key.json
            gcloud --quiet config set project ${GOOGLE_PROJECT_ID}
            gcloud --quiet config set compute/zone ${GOOGLE_COMPUTE_ZONE}
            gcloud --quiet container clusters get-credentials ${GOOGLE_CLUSTER_NAME}
      - setup_remote_docker
      - run:
          name: Docker build and push
          command: |
            echo $DOCKER_HUB_PASSWORD | docker login --username cflynnus --password-stdin
            docker build -t cflynnus/blog:${CIRCLE_BRANCH}-${CIRCLE_SHA1} .
            docker push cflynnus/blog:${CIRCLE_BRANCH}-${CIRCLE_SHA1}
            docker tag cflynnus/blog:${CIRCLE_BRANCH}-${CIRCLE_SHA1} cflynnus/blog:latest
            docker push cflynnus/blog:latest
      - run:
          name: Deploy to Kubernetes
          command: |
            if [[ "$CIRCLE_BRANCH" == "develop" ]]; then
              helm upgrade --install blog --reuse-values --debug \
              --set develop_image=cflynnus/blog:${CIRCLE_BRANCH}-${CIRCLE_SHA1} \
              ./helm
            elif [[ "$CIRCLE_BRANCH" == "master" ]]; then
              helm upgrade --install blog --reuse-values --debug \
              --set master_image=cflynnus/blog:${CIRCLE_BRANCH}-${CIRCLE_SHA1} \
              ./helm
            fi
            kubectl rollout status deployment/${PROJECT_NAME}-${CIRCLE_BRANCH}-app
workflows:
  version: 2
  build_test_deploy:
    jobs:
      - test
      - build_and_deploy:
          requires:
            - test
          filters:
            branches:
              only:
                - develop
                - master
</pre>

Pretty simple. It has two jobs, `test` and `build_and_deploy`.
`build_and_deploy` runs conditionally on `test` completing successfully.
`$GCLOUD_SERVICE_KEY` and `$DOCKER_HUB_PASSWORD` are "secrets" that are stored
and encrypted with CircleCI.

And after some hacking around, here's the solution I came up with for Github
Actions.
#### .github/workflows/test_build_deploy.yml
<pre class="prettyprint linenums">
name: Test, Build and Deploy
on: [push]
jobs:
  test:
    name: Go Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - name: Go Test
        uses: cedrickring/golang-action/go1.12@1.3.0
        with:
          args: go test ./... -coverprofile=coverage.txt -covermode=atomic -coverpkg=myblog/app
        env:
          GO111MODULE: on
          GOFLAGS: -mod=vendor
      - name: Upload Coverage
        run: curl -s https://codecov.io/bash | bash -s -- -t ${{ "{{" }}secrets.CODECOV_TOKEN}} -f ./coverage.txt
  build_and_push:
    name: Build and Push
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@master
      - name: Docker Login
        uses: actions/docker/login@master
        env:
          DOCKER_USERNAME: ${{ "{{" }} secrets.DOCKER_USERNAME {{ "}}" }}
          DOCKER_PASSWORD: ${{ "{{" }} ssecrets.DOCKER_PASSWORD {{ "}}" }}
      - name: Docker Build
        uses: actions/docker/cli@master
        with:
          args: build . --file Dockerfile -t cflynnus/blog:`echo ${GITHUB_REF} | cut -d'/' -f3`-${GITHUB_SHA}
      - name: Docker Tag Latest
        uses: actions/docker/cli@master
        with:
          args: tag cflynnus/blog:`echo ${GITHUB_REF} | cut -d'/' -f3`-${GITHUB_SHA} cflynnus/blog:latest
      - name: Docker Push Hash Tag
        uses: actions/docker/cli@master
        with:
          args: push cflynnus/blog:`echo ${GITHUB_REF} | cut -d'/' -f3`-${GITHUB_SHA} 
      - name: Docker Push Latest
        uses: actions/docker/cli@master
        with:
          args: push cflynnus/blog:latest
  deploy_staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    needs: build_and_push
    steps:
      - uses: actions/checkout@master
        if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop'
      - name: Google Cloud Authenticate
        if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop'
        uses: actions/gcloud/auth@master
        env:
          GCLOUD_AUTH: ${{ "{{" }}secrets.GCLOUD_AUTH}}
      - name: GC Set Project ID
        if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop'
        uses: actions/gcloud/cli@master
        env:
          GOOGLE_PROJECT_ID: blog-229516
        with:
          args: config set project ${GOOGLE_PROJECT_ID}
      - name: GC set compute/zone
        if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop'
        uses: actions/gcloud/cli@master
        env:
          GOOGLE_COMPUTE_ZONE: us-central1-b
        with:
          args: config set compute/zone ${GOOGLE_COMPUTE_ZONE}
      - name: GC get-credentials
        if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop'
        uses: actions/gcloud/cli@master
        env:
          GOOGLE_CLUSTER_NAME: blog-cluster
        with:
          args: container clusters get-credentials ${GOOGLE_CLUSTER_NAME}
      - name: GCP List Containers
        if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop'
        uses: actions/gcloud/cli@master
        with:
          args: container clusters list
      - name: Helm Deploy Develop
        uses: stefanprodan/gh-actions/helm@master
        if: github.ref == 'refs/heads/develop'
        with:
          args: |
            upgrade --install blog --reuse-values --debug \
            --set develop_image=cflynnus/blog:`echo ${GITHUB_REF} | cut -d'/' -f3`-${GITHUB_SHA} \
            ./helm
      - name: Helm Deploy Master
        uses: stefanprodan/gh-actions/helm@master
        if: github.ref == 'refs/heads/master'
        with:
          args: |
            upgrade --install blog --reuse-values --debug \
            --set master_image=cflynnus/blog:`echo ${GITHUB_REF} | cut -d'/' -f3`-${GITHUB_SHA} \
            ./helm
</pre>

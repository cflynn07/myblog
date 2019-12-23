I decided to switch the CI/CD pipeline of this blog from CircleCI to Github
Actions just to check it out. I know it's basically overkill to have a CI/CD
setup for a blog (and to use Kubernetes), but I did it for curiosity. One of
the most obvious advantages of Github Actions is tight integration with
Github. If nothing else, it's nice just to have your repository and CI
information together on one website.

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
        run: curl -s https://codecov.io/bash | bash -s -- -t ${{secrets.CODECOV_TOKEN}} -f ./coverage.txt
  build_and_push:
    name: Build and Push
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@master
      - name: Docker Login
        run: echo "${{secrets.DOCKER_PASSWORD}}" | docker login -u ${{secrets.DOCKER_USERNAME}} --password-stdin
      - name: Docker Build
        run: |
          echo "GITHUB_REF: $GITHUB_REF"; \
          BRANCH_NAME="`echo $GITHUB_REF | awk -F'/' '{print $3}'`"; \
          echo "BRANCH_NAME: $BRANCH_NAME"; \
          IMAGE_TAG_NAME="cflynnus/blog:${BRANCH_NAME}-${GITHUB_SHA}"; \
          echo "IMAGE_TAG_NAME: $IMAGE_TAG_NAME"; \
          echo "::set-env name=IMAGE_TAG_NAME::$IMAGE_TAG_NAME"; \
          docker build . --file Dockerfile -t $IMAGE_TAG_NAME;
      - name: Docker Tag Latest
        run: docker tag "$IMAGE_TAG_NAME" cflynnus/blog:latest
      - name: Docker Push Tags
        run: |
          echo "IMAGE_TAG_NAME: $IMAGE_TAG_NAME"; \
          docker push "$IMAGE_TAG_NAME"; \
          docker push cflynnus/blog:latest;
  deploy:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    needs: build_and_push
    if: github.ref == 'refs/heads/develop' || github.ref == 'refs/heads/master'
    steps:
      - uses: actions/checkout@master
      - name: Set Envs
        run: |
          echo "GITHUB_REF: $GITHUB_REF"; \
          BRANCH_NAME="`echo $GITHUB_REF | awk -F'/' '{print $3}'`"; \
          echo "BRANCH_NAME: $BRANCH_NAME"; \
          IMAGE_TAG_NAME="cflynnus/blog:${BRANCH_NAME}-${GITHUB_SHA}"; \
          echo "IMAGE_TAG_NAME: $IMAGE_TAG_NAME"; \
          echo "::set-env name=IMAGE_TAG_NAME::$IMAGE_TAG_NAME"; \
      - name: Google Cloud Authenticate
        uses: actions/gcloud/auth@master
        env:
          GCLOUD_AUTH: ${{secrets.GCLOUD_AUTH}}
      - name: GC Set Project ID
        uses: actions/gcloud/cli@master
        env:
          GCLOUD_PROJECT_ID: ${{secrets.GCLOUD_PROJECT_ID}}
        with:
          args: config set project ${GCLOUD_PROJECT_ID}
      - name: GC set compute/zone
        uses: actions/gcloud/cli@master
        env:
          GCLOUD_COMPUTE_ZONE: ${{secrets.GCLOUD_COMPUTE_ZONE}}
        with:
          args: config set compute/zone ${GCLOUD_COMPUTE_ZONE}
      - name: GC get-credentials
        uses: actions/gcloud/cli@master
        env:
          GCLOUD_CLUSTER_NAME: ${{secrets.GCLOUD_CLUSTER_NAME}}
        with:
          args: container clusters get-credentials ${GCLOUD_CLUSTER_NAME}
      - name: GCP List Containers
        uses: actions/gcloud/cli@master
        with:
          args: container clusters list
      - name: Helm Deploy Develop
        uses: stefanprodan/gh-actions/helm@master
        if: github.ref == 'refs/heads/develop'
        with:
          args: |
            upgrade --install blog --reuse-values --debug \
            --set develop_image="$IMAGE_TAG_NAME" \
            ./helm
      - name: Helm Deploy Master
        uses: stefanprodan/gh-actions/helm@master
        if: github.ref == 'refs/heads/master'
        with:
          args: |
            upgrade --install blog --reuse-values --debug \
            --set master_image="$IMAGE_TAG_NAME" \
            ./helm
</pre>

Pretty similar. My Github Action has 3 jobs and runs on a push event (to any
branch). Actions can run on events, on a cron schedule, or manually. For my
purposes, running on push events works. My first two jobs run on pushes to all
branches. The third job, which deploys my code to my kubernetes cluster running
on google cloud, only runs on the develop and master branches.


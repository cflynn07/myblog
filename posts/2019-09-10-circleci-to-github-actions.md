I decided to switch the CI/CD pipeline of this blog from CircleCI to Github
Actions just for fun. It's overkill to have a CI/CD setup for a blog (and to
use Kubernetes). I tried it just to play around a bit with GH Actions and
Kubernetes. One of the most obvious advantages of Github Actions is tight
integration with Github. If nothing else, it's nice just to have your
repository and CI/CD together on one platform.

#### My Github Actions Workflow:
##### For pushes to all branches
<ol>
<li>Run tests & upload coverage reports to CodeCov</li>
<li>Build & push a docker image (tagged with reference to branch name for use in local development)</li>
</ol>

##### For pushes to the `master` and `develop` branches:
<ol>
<li>Deploy to kubernetes cluster, either staging or production, using helm</li>
</ol>

This is the CircleCI configuration file that defines the CI/CD pipeline I
ported to Github Actions
###### .circleci/config.yml
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
`build_and_deploy` runs conditionally on the `test` job completing
successfully. `$GCLOUD_SERVICE_KEY` and `$DOCKER_HUB_PASSWORD` are "secrets"
that are stored and encrypted with CircleCI.

And after some hacking around, here's the solution I came up with for Github
Actions.
###### .github/workflows/test_build_deploy.yml
<pre class="prettyprint linenums">
name: Test, Build and Deploy
on: [push]
jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@master

      - name: Test
        uses: cedrickring/golang-action/go1.12@1.3.0
        with:
          args: go test ./... -coverprofile=coverage.txt -covermode=atomic -coverpkg=myblog/app
        env:
          GO111MODULE: on
          GOFLAGS: -mod=vendor

      - name: Upload Coverage to CodeCov
        run: curl -s https://codecov.io/bash | bash -s -- -t ${{secrets.CODECOV_TOKEN}} -f ./coverage.txt
  build_and_push_image:
    name: Build and Push Image
    runs-on: ubuntu-latest
    needs: test
    steps:
      - name: Checkout
        uses: actions/checkout@master

      - name: Set Envs
        # | ::set-env explanation
        # Bit of non-dry code here, also copied to the "Deploy" job to calculate ENVs
        # https://help.github.com/en/actions/automating-your-workflow-with-github-actions/development-tools-for-github-actions#set-an-environment-variable-set-env
        run: |
          BRANCH_NAME="`echo $GITHUB_REF | awk -F'/' '{print $3}'`";
          IMAGE_TAG_NAME="cflynnus/blog:$BRANCH_NAME-$GITHUB_SHA";
          echo "::set-env name=IMAGE_TAG_NAME::$IMAGE_TAG_NAME";

      - name: Docker Login, Build, Tag and Push
        run: |
          echo "${{secrets.DOCKER_PASSWORD}}" | docker login -u ${{secrets.DOCKER_USERNAME}} --password-stdin;
          docker build . --file Dockerfile -t "$IMAGE_TAG_NAME";
          docker tag "$IMAGE_TAG_NAME" cflynnus/blog:latest;
          docker push "$IMAGE_TAG_NAME";
          docker push cflynnus/blog:latest;
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    needs: build_and_push_image
    if: github.ref == 'refs/heads/develop' || github.ref == 'refs/heads/master'
    steps:
      - name: Checkout
        uses: actions/checkout@master

      - name: Set Envs & Install helm3 Client
        # | ::set-env explanation
        # Bit of non-dry code here, also copied to the "Deploy to Staging" job to calculate ENVs
        # https://help.github.com/en/actions/automating-your-workflow-with-github-actions/development-tools-for-github-actions#set-an-environment-variable-set-env
        run: |
          BRANCH_NAME="`echo $GITHUB_REF | awk -F'/' '{print $3}'`";
          IMAGE_TAG_NAME="cflynnus/blog:$BRANCH_NAME-$GITHUB_SHA";

          HELM_3_FILE="helm-v3.0.0-linux-amd64.tar.gz";
          HELM_URL="https://get.helm.sh/$HELM_3_FILE";
          curl -Ls "$HELM_URL" | tar xvz;
          mkdir -p "$GITHUB_WORKSPACE/bin";
          mv linux-amd64/helm "$GITHUB_WORKSPACE/bin/helm3";
          echo "::add-path::$GITHUB_WORKSPACE/bin";

          echo "::set-env name=BRANCH_NAME::$BRANCH_NAME";
          echo "::set-env name=IMAGE_TAG_NAME::$IMAGE_TAG_NAME";

      - name: Install gcloud cli
        uses: GoogleCloudPlatform/github-actions/setup-gcloud@master
        with:
          version: '270.0.0'
          service_account_email: ${{ secrets.GCLOUD_SA_EMAIL }}
          service_account_key: ${{ secrets.GCLOUD_SA_KEY }}

      - name: gcloud configure
        run: |
          gcloud config set project ${{secrets.GCLOUD_PROJECT_ID}};
          gcloud config set compute/zone ${{secrets.GCLOUD_COMPUTE_ZONE}};
          gcloud container clusters get-credentials ${{secrets.GCLOUD_CLUSTER_NAME}};

      - name: Deploy
        run: |
          overrides=(
            "develop_deployment_sha=$GITHUB_SHA"
            "develop_deployment_time=`TZ=Asia/Taipei date`"
          )
          if [[ $BRANCH_NAME == "develop" ]]; then
            overrides+=("develop_image=$IMAGE_TAG_NAME")
          elif [[ $BRANCH_NAME == "master" ]]; then
            overrides+=("master_image=$IMAGE_TAG_NAME")
          fi
          overrides=$(for i in "${overrides[@]}"; do echo -n "$i,"; done)
          overrides=${overrides:0:${#overrides}-1}
          helm3 upgrade blog ./helm \
            --install \
            --debug \
            --reuse-values \
            --set-string "$overrides"

      - name: Rollout Status
        run: kubectl rollout status "deployment/blog-$BRANCH_NAME-app"
</pre>

My Github Action has 1 workflow that contains 3 jobs.

The first job `test` runs on a push event to any branch. It runs the tests,
generates coverage reports, and uploads those reports to the CodeCov service.

The second job `build_and_push_image` runs on a push event to any branch but
only if the `test` job completes successfully. The job builds a docker image
and tags it with a few custom tags based on the branch name and commit SHA. These
tags are pushed to the remote docker registry.

The last job `deploy` runs a push event to either the `master` or `develop`
branches and only if the `build_and_push_image` completes successfully
(implicitly requiring the `test` job to complete successfully). This job is the
most complex. Github's hosted runners have many tools and packages
pre-installed. I use helm3 to deploy this blog and as of Jan 2020 the hosted
runners only have the helm2 client pre-installed. I looked around for a custom
action I could use to easily get helm3 installed in the PATH of my runner but I
couldn't find one. I took a stab at creating my own action
[cflynn07/gha-helm3](https://github.com/cflynn07/gha-helm3) but decided it more
simple to just install helm3 as a step. I also wanted to compute some values
and use those in later steps so I used the `::set-env` development tool to set
environment variables for subsequent jobs.

<pre class="prettyprint linenums">
- name: Set Envs & Install Helm3 Client
  # | ::set-env explanation
  # Bit of non-dry code here, also copied to the "Deploy to Staging" job to calculate ENVs
  # https://help.github.com/en/actions/automating-your-workflow-with-github-actions/development-tools-for-github-actions#set-an-environment-variable-set-env
  run: |
    BRANCH_NAME="`echo $GITHUB_REF | awk -F'/' '{print $3}'`";
    IMAGE_TAG_NAME="cflynnus/blog:$BRANCH_NAME-$GITHUB_SHA";

    HELM_3_FILE="helm-v3.0.0-linux-amd64.tar.gz";
    HELM_URL="https://get.helm.sh/$HELM_3_FILE";
    curl -Ls "$HELM_URL" | tar xvz;
    mkdir -p "$GITHUB_WORKSPACE/bin";
    mv linux-amd64/helm "$GITHUB_WORKSPACE/bin/helm3";
    echo "::add-path::$GITHUB_WORKSPACE/bin";

    echo "::set-env name=BRANCH_NAME::$BRANCH_NAME";
    echo "::set-env name=IMAGE_TAG_NAME::$IMAGE_TAG_NAME";
</pre>

[Software Installed on Github Hosted Runners](https://help.github.com/en/actions/automating-your-workflow-with-github-actions/software-installed-on-github-hosted-runners)

[Issue to add helm3 to Github Hosted Runners](https://github.com/actions/virtual-environments/issues/108)

[Actions developer tools](https://help.github.com/en/actions/automating-your-workflow-with-github-actions/development-tools-for-github-actions)

I also needed the [Google Cloud CLI](https://cloud.google.com/sdk/gcloud/),
which does not come pre-installed on the runner. Google provides an official
action
[GoogleCloudPlatform/github-actions](https://github.com/GoogleCloudPlatform/github-actions/blob/master/setup-gcloud/README.md)
that can be used to install and provision the gcloud cli in the runner.

<pre class="prettyprint linenums">
- name: Install gcloud cli
  uses: GoogleCloudPlatform/github-actions/setup-gcloud@master
  with:
    version: '270.0.0'
    service_account_email: ${{ secrets.GCLOUD_SA_EMAIL }}
    service_account_key: ${{ secrets.GCLOUD_SA_KEY }}
</pre>

The last three steps in the job configure kubectl (pre-installed on runner),
update the deployment object's container image in my k8s cluster with the new
image to deploy, and check the status of the rollout. 

<pre class="prettyprint linenums">
- name: gcloud configure
  run: |
    gcloud config set project ${{secrets.GCLOUD_PROJECT_ID}};
    gcloud config set compute/zone ${{secrets.GCLOUD_COMPUTE_ZONE}};
    gcloud container clusters get-credentials ${{secrets.GCLOUD_CLUSTER_NAME}};

- name: Deploy
  run: |
    overrides=(
      "develop_deployment_sha=$GITHUB_SHA"
      "develop_deployment_time=`TZ=Asia/Taipei date`"
    )
    if [[ $BRANCH_NAME == "develop" ]]; then
      overrides+=("develop_image=$IMAGE_TAG_NAME")
    elif [[ $BRANCH_NAME == "master" ]]; then
      overrides+=("master_image=$IMAGE_TAG_NAME")
    fi
    overrides=$(for i in "${overrides[@]}"; do echo -n "$i,"; done)
    overrides=${overrides:0:${#overrides}-1}
    helm3 upgrade blog ./helm \
      --install \
      --debug \
      --reuse-values \
      --set-string "$overrides"

- name: Rollout Status
  run: kubectl rollout status "deployment/blog-$BRANCH_NAME-app"
</pre>

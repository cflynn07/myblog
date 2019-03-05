I use helm and circleci to update the production and staging deployments
of this blog (yes I actually have a staging deployment for a blog) in my
kubernetes cluster running in google cloud. In my helm chart I create
two deployments, two services, and one ingress resource to route to
either service based on the http host of the request.
<pre class="prettyprint">
$ tree ./helm
./helm
├── Chart.yaml
├── templates
│   ├── NOTES.txt
│   ├── _helpers.tpl
│   ├── deployment-develop.yaml
│   ├── deployment-master.yaml
│   ├── ingress.yaml
│   ├── service-develop.yaml
│   └── service-master.yaml
└── values.yaml
</pre>

<pre class="prettyprint">
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
</pre>

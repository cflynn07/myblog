I use helm and circleci to update the production and staging deployments of
this blog (yes I actually have a staging deployment for a blog) in my
kubernetes cluster running on google cloud. The helm chart for this project has
two deployments, two services, and one ingress resource that routes requests to
either production or staging service based on the http host header of the
request.

#### Two deployments, two services, one ingress
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

#### CI script that conditionally updates helm deployment
<pre class="prettyprint linenums">
// .circleci/config.yml
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

I ran into a problem when I added a new variable to my helm chart values.yaml
file.
<pre class="prettyprint linenums">
// values.yaml
develop_image: "cflynnus/blog:latest"
master_image: "cflynnus/blog:latest"
google_analytics: "UA-000000000-1" # New variable
</pre>

My CI uses `helm upgrade` with the `--reuse-values` flag to update either my
production or my staging environment depending on whether I've just pushed to
my master or develop branches. This flag instructs tiller to generate the yaml
kubernetes resources using the computed values from my last deployment and
accepting values passed to `helm upgrade` via the `--set` flag. This works most
of the time but becomes an issue when you want to add a value to your
values.yaml file in your helm chart. When running `helm upgrade --reuse-values
--set ...` tiller **will not** use new values in values.yaml. This leads to
template errors, where the template expects certain values to be present but
are not.

The way I get around this, is to do idempotent deploy from my local machine
with the new variable, then push code to my remote git/github repo and have my
CI run the deployment job. Since the value is already present in tiller,
using `--reuse-values` will lead to the correct yaml being generated.

<pre class="prettyprint">
// This wont cause any changes to kuberenetes resources since no resources use 'google_analytics' yet
$ helm upgrade --install YOUR_DEPLOYMENT --set google_analytics="xxxxxxxxxxxxxxx" .
</pre>

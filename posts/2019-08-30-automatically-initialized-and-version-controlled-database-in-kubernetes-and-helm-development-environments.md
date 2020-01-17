Today I was hacking together a kitchen sink project to play with development
using kubernetes, skaffold and helm. I spent a bit of time thinking about how
to set up MySQL. At my previous company we used docker-compose and created a
`MySQL` container along with a container called `migrations` that initialized
our database with schema and sample data. The migrations container would run
once when the development environment started and exit when complete.

For this latest project, I decided to go the route of creating a custom docker
image based on the `mysql:5.6` image (https://hub.docker.com/_/mysql). This
custom image simply copies a schema dump file to `/docker-entrypoint-initdb.d/`.
My goal is to create a single source of truth for the latest database schema in
a project as well as use version control to document the history of changes.

#### Dockerfile
https://github.com/cflynn07/rgbm-mysql
<pre class="prettyprint">
FROM mysql:5.6
COPY rgbm.sql /docker-entrypoint-initdb.d/rgbm.sql
</pre>

I use Github Actions to automatically build & push to the docker hub image
registry on any push to the master branch of the repository.
https://github.com/cflynn07/rgbm-mysql/blob/master/.github/workflows/dockerimage.yml
#### Github Action
<pre class="prettyprint linenums">
name: Docker Image CI
on:
  push:
    paths:
      - 'rgbm.sql'
      - 'Dockerfile-mysql'
    branches:
      - master
jobs:
  build_and_push:
    name: build and push
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - uses: 'actions/docker/login@master'
        name: 'docker login'
        env:
          DOCKER_USERNAME: ${{ "{{" }} secrets.DOCKER_USERNAME {{ "}}" }}
          DOCKER_PASSWORD: ${{ "{{" }} ssecrets.DOCKER_PASSWORD {{ "}}" }}
      - uses: 'actions/docker/cli@master'
        name: 'docker build'
        with:
          args: 'build . --file Dockerfile-mysql --tag cflynnus/rgbm-mysql:${{ "{{" }} github.sha {{ "}}" }}'
      - uses: 'actions/docker/cli@master'
        name: 'docker push hash'
        with:
          args: 'push cflynnus/rgbm-mysql:${{ "{{" }} github.sha {{ "}}" }}'
      - uses: 'actions/docker/cli@master'
        name: 'docker tag latest'
        with:
          args: 'tag cflynnus/rgbm-mysql:${{ "{{" }} github.sha {{ "}}" }} cflynnus/rgbm-mysql:latest'
      - uses: 'actions/docker/cli@master'
        name: 'docker push latest'
        with:
          args: 'push cflynnus/rgbm-mysql:latest'
</pre>

Then in my helm templates I add a deployment object that references this image.
Note I set the `imagePullPolicy` to `Always` so that each time I start my
development environment docker will pull the latest image. The default value of
imagePullPolicy is `IfNotPresent`. 
https://github.com/cflynn07/rgbm/blob/master/helm/templates/deployment.yaml
#### Kubernetes (helm managed) Deployment Object
<pre class="prettyprint linenums">
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - image: cflynnus/mysql-rgbm
          name: mysql
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: password
            - name: MYSQL_DATABASE
              value: rgbm
            - name: MYSQL_USER
              value: user
            - name: MYSQL_PASSWORD
              value: password
          ports:
            - containerPort: 3306
              name: mysql
          volumeMounts:
            - name: mysql-persistent-storage
              mountPath: /var/lib/mysql
      volumes:
        - name: mysql-persistent-storage
          persistentVolumeClaim:
            claimName: mysql-pv-claim
</pre>

When I modify the database I perform a MySQL dump and overwrite the .sql file
in my cflynn07/rgbm-mysql repository. Then I commit and push my changes to the
remote repository.
<pre class="prettyprint">
$ docker ps | grep mysql | awk '{print $1}'
bb4f90b37b92
$ docker exec -i bb4f90b37b92 sh -c 'exec mysqldump --all-databases -uuser -ppassword' > ./rgbm.sql
</pre>

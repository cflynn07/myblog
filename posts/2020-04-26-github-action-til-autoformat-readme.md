[TIL Auto-Format REAME Action on Github Marketplace][4]  
[cflynn07/til][5]  

While doing my daily browsing of Hacker News, I came across [this][1] post
([Hacker News][2]) by Simon Wilson on the merits of writing small, actionable
TILs "Today I Learned's." and how he leverages GitHub Actions to automatically
generate a README file based on the contents of his [TIL
repository (simonw/til)][3].

I like the idea of a TIL repo and using GitHub Actions to automate indexing.
Other people also have had the idea to use GitHub actions to index their TIL
repository READMEs as well, however all the examples I could find used GitHub
Actions to run a script that was included in their repository. This works but
it seemed like a reusable GitHub Action that could be quickly dropped into a
TIL repo would be useful for many people.

I've been using GitHub actions for a few months and I'm enjoying the product.
Free and easy to use CI/CD platforms that integrate with GitHub event hooks to
run arbitrary code on push events have existed for years. The nice part about
GitHub Actions is the tight integration with GitHub and the focus on
encouraging users to create an ecosystem of small, reusable discrete "Actions"
that can be dropped into others' workflows. The fact that these actions can
essentially be docker containers that recieve a repository mounted as a volume
inside the container makes it very easy to build reliably consistent actions.

This is how the action can be added to a TIL repo.

###### .github/workflows/action.yml
<pre class="prettyprint">
name: Build README
on:
  push:
    branches:
    - master
    paths-ignore:
    - README.md
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Check out repo
      uses: actions/checkout@v2
      with:
        # necessary for github-action-til-autoformat-readme
        fetch-depth: 0
    - name: Autoformat README
      uses: cflynn07/github-action-til-autoformat-readme@1.1.0
      with:
        description: |
          A collection of concrete writeups of small things I learn daily while working
          and researching. My goal is to work in public. I was inspired to start this
          repository after reading Simon Wilson's [Hacker News post][1], and he was
          apparently inspired by Josh Branchaud's [TIL collection][2].
        footer: |
          [1]: https://simonwillison.net/2020/Apr/20/self-rewriting-readme/
          [2]: https://github.com/jbranchaud/til
        list_most_recent: 2
</pre>

The action writes the README.md file to the repo, commits and pushes it back.

I built the action with Go and used a [multi-stage build][6] to keep the
resulting image as small as possible by using an [alpine][8] image with the
built binary.

<pre class="prettyprint">
FROM golang:1.14 as builder
WORKDIR /go/src/app
COPY . .
RUN go build -mod=vendor -o /go/bin/main .

FROM alpine:latest
WORKDIR /root
RUN apk update && apk add git
COPY --from=builder /go/bin/main ./main 
COPY --from=builder /go/src/app/README.md.tmpl ./README.md.tmpl
COPY --from=builder /go/src/app/entrypoint.sh ./entrypoint.sh
ENTRYPOINT [ "/root/entrypoint.sh" ]
</pre>

When the GitHub Actions running runs a reusable action, it essentially builds
the docker image from scratch each time based on the action's Dockerfile.
Building the above Dockerfile for each GitHub Action run is a little slow and
also isn't totally necessary. Instead, in the repository that holds my GitHub
action I have a [GitHub Action workflow][7] that builds, tags and pushes a base
image to Docker Hub whenever a new release/tag is created. That base image is
referenced in the `Dockerfile` in the action repository root:
<pre class="prettyprint">
FROM cflynnus/github-action-til-autoformat-readme:1.1.0 #tagged base image
ENV TEMPLATE_PATH "/root/README.md.tmpl"
ENV REPO_PATH "/github/workspace"
</pre>

This way, when this action is run in a workflow the runner only has to build
this workflow which is much faster than the multi-stage build. The resulting
[image][9] is ~14Mb.

<img src="/static/images/2020-04-26/Screen_Shot_2020-04-26_at_2.01.45_PM.png" />

[1]: https://simonwillison.net/2020/Apr/20/self-rewriting-readme/
[2]: https://news.ycombinator.com/item?id=22920437
[3]: https://github.com/simonw/til
[4]: https://github.com/marketplace/actions/til-auto-format-readme
[5]: https://github.com/cflynn07/til
[6]: https://docs.docker.com/develop/develop-images/multistage-build/
[7]: https://github.com/cflynn07/github-action-til-autoformat-readme/blob/master/.github/workflows/tag_test_push.yml
[8]: https://hub.docker.com/_/alpine
[9]: https://hub.docker.com/layers/cflynnus/github-action-til-autoformat-readme/1.1.0/images/sha256-42570a0bcdf96ab66ff555c267bd9129d660b186762976cad6c27be76fbf7323?context=repo

MyBlog
======

[![CircleCI](https://circleci.com/gh/cflynn07/myblog/tree/master.svg?style=svg)](https://circleci.com/gh/cflynn07/myblog/tree/master)
[![codecov](https://codecov.io/gh/cflynn07/myblog/branch/master/graph/badge.svg)](https://codecov.io/gh/cflynn07/myblog)
![](https://img.shields.io/github/last-commit/cflynn07/myblog.svg)
[![Maintainability](https://api.codeclimate.com/v1/badges/ddb5503e282c7693f9f5/maintainability)](https://codeclimate.com/github/cflynn07/myblog/maintainability)

My blog website, written in golang. Deployed to google cloud and managed with
kubernetes (overkill for a blog).

Attempts to follow golang project standard layout
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

Development Instructions
------------------------
```bash
# Build and run
$ docker build . -t myblog
$ docker run -it --name myblog -p 3001:3001 --rm -e "PORT=3001" myblog

# Develop using skaffold, image will be rebuild and cluster updated every time code is changed
$ skaffold dev

# This project uses gulp to run build tasks (sass)
$ npm install -g gulp && npm install # install gulp globally and local dev dependencies
$ gulp
```

CI/CD Pipeline


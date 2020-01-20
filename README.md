MyBlog
======

[![CircleCI](https://circleci.com/gh/cflynn07/myblog/tree/master.svg?style=svg)](https://circleci.com/gh/cflynn07/myblog/tree/master)
[![codecov](https://codecov.io/gh/cflynn07/myblog/branch/master/graph/badge.svg)](https://codecov.io/gh/cflynn07/myblog)
![](https://img.shields.io/github/last-commit/cflynn07/myblog.svg)
[![Maintainability](https://api.codeclimate.com/v1/badges/ddb5503e282c7693f9f5/maintainability)](https://codeclimate.com/github/cflynn07/myblog/maintainability)

My blog website, written in golang. Deployed to google cloud and managed with
kubernetes and helm (overkill for a blog).

Attempts to follow golang project standard layout
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

Development Instructions
------------------------
This project uses skaffold and kubernetes for local development.

```bash
# Set up kubernetes cluster and configure kubectl to use desired context
$ skaffold dev

# This project uses gulp to run build tasks (sass)
$ npm install -g gulp && npm install # install gulp globally and local dev dependencies
$ gulp
```

CI/CD Pipeline
- Commits/merges to `develop` are automatically deployed to staging
  - http://cflynn-blog.com (requires editing /etc/hosts)
- Commits/merges to `master` are automatically deployed to production
  - https://cflynn.us

Authors
-------

* [Casey Flynn](http://github.com/cflynn07) - <https://cflynn.us/>

License
-------

myblog free and unencumbered public domain software. For more information, see
<http://unlicense.org/> or the accompanying UNLICENSE file.

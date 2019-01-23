MyBlog
======

[![CircleCI](https://circleci.com/gh/cflynn07/myblog/tree/master.svg?style=svg)](https://circleci.com/gh/cflynn07/myblog/tree/master)

My blog website, written in golang. Deployed to google cloud and managed with
kubernetes (overkill for a blog).

Attempts to follow golang project standard layout
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

Instructions
------------
```bash
$ docker build . -t myblog
$ docker run -it --name myblog -p 3001:3001 --rm -e "PORT=3001" myblog
```

# [HOWTO] Maintain tech-meetup.com

## Auto Release Is Ready
Since commit https://github.com/lamuguo/blog/commit/5ccae9d354f5b1d2c98c3bed6ba69da1af9bc3fa, automatically push is ready. For additional event, we just need to commit the modified/added files, the result will be automatically pushed to tech-meetup.com. Details can be found in .travis.yml file.

For debugging, you can use command below:
```
$ go build -o /tmp/server . && /tmp/server
```

## Instructions
```sh
$ export MEETUP_IMAGE=lamuguo/meetup:20170214

$ docker build -t $MEETUP_IMAGE .
$ docker push $MEETUP_IMAGE
    # if not login, please “docker login -u lamuguo”

$ docker rm -f meetup && \
  docker run -d --name meetup \
  -p 80:80 \
  --entrypoint /go/bin/blog \
  $MEETUP_IMAGE \
  -http="0.0.0.0:80" \
  -vhost_map="blog.tech-meetup.com=testing/|lamuguo.tech-meetup.com=lamuguo/"
```

Notes
- For using the command, remember to update MEETUP_IMAGE.
- vhost_map is used for virtual hosting, the string format is: "\<vhost1>=\<dir1>|\<vhost2>=\<dir2>|..."

## Useful Links
- Make sure sites below works correctly:
  * http://tech-meetup.com/
  * http://blog.tech-meetup.com/
  * http://lamuguo.tech-meetup.com/
- Code repository: https://github.com/lamuguo/blog
- Sample changelist: https://github.com/lamuguo/blog/commit/cdeadb9857dfb347a72e7cb13a2dd01a39d6da5e

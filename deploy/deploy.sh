MEETUP_IMAGE=lamuguo/meetup:live

rm -rf /root/blog
git clone https://github.com/lamuguo/blog /root/blog

cd /root/blog \
    && docker build -t $MEETUP_IMAGE . \
    && docker rm -f meetup \
    && docker run -d --name meetup \
	      -p 80:80 \
	      --entrypoint /go/bin/blog \
	      $MEETUP_IMAGE \
	      -http="0.0.0.0:80" \
	      -vhost_map="blog.tech-meetup.com=testing/|lamuguo.tech-meetup.com=lamuguo/"

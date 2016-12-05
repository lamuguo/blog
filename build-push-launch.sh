export VERSION=20161224

docker build -t lamuguo/meetup:$VERSION .
docker push lamuguo/meetup:$VERSION
docker rm -f meetup && docker run -d --name meetup  -p 80:8080 lamuguo/meetup:$VERSION

FROM ubuntu:latest

# Essential for using tls
RUN apt-get update
RUN apt-get install ca-certificates -y
RUN update-ca-certificates

# web port
EXPOSE 8081
# grpc port
EXPOSE 50051

ADD build/notificationGo /app/notificationGo
ADD templates /app/templates
RUN ls -l

WORKDIR /app
CMD ./notificationGo

FROM alpine:3.12

ADD docker-sync /usr/bin
ADD skopeo /usr/bin

ADD env.sh /
RUN chmod +x /env.sh && chmod +x /usr/bin/docker-sync && chmod +x /usr/bin/skopeo

WORKDIR /github/workspace
ENTRYPOINT ["/env.sh"]
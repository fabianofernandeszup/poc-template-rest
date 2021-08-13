FROM ritclizup/rit-node-runner

USER root

RUN apk add git jq

RUN mkdir /rit
COPY . /rit
RUN sed -i 's/\r//g' /rit/set_umask.sh
RUN sed -i 's/\r//g' /rit/run.sh
RUN chmod +x /rit/set_umask.sh

WORKDIR /rit
ENTRYPOINT ["/rit/set_umask.sh"]
CMD ["/rit/run.sh"]

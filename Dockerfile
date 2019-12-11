# build stage
FROM golang:1.13 AS build
ADD . /src
RUN cd /src \
  && go build ./cmd/subscribe \
  && go build ./cmd/publish


# final stage
FROM ubuntu:bionic

RUN apt-get update && apt-get install -y curl \
  && apt-get install -y gnupg2 \
  && apt-get install -y apt-transport-https \
  && curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
  && echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list \
  && apt-get update \
  && apt-get install -y kubectl \
  && apt-get install -y jq \
  && apt-get remove -y --auto-remove apt-transport-https \
  && apt-get remove -y --auto-remove gnupg2 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

ADD scripts/* /riff/dev-utils/

COPY --from=build /src/subscribe /riff/dev-utils
COPY --from=build /src/publish /riff/dev-utils

WORKDIR /riff/dev-utils/

ENV PATH="/riff/dev-utils/:${PATH}"

CMD ["entrypoint.sh"]

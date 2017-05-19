FROM discoenv/golang-base:master

ENV CONF_TEMPLATE=/go/src/github.com/cyverse-de/tagger/tagger.yaml.tmpl
ENV CONF_FILENAME=tagger.yaml
ENV PROGRAM=tagger

COPY . /go/src/github.com/cyverse-de/tagger/

RUN git clone https://github.com/swagger-api/swagger-ui.git /tmp/swagger-ui \
    && cd /tmp/swagger-ui \
    && git checkout v2.2.10 \
    && mkdir -p /docs \
    && cp -pr dist/* /docs/ \
    && cd / \
    && rm -rf /tmp/swagger-ui \
    && cp /go/src/github.com/cyverse-de/tagger/index.html /docs/index.html

RUN go install github.com/cyverse-de/tagger/... \
    && cp /go/bin/tagger-server /bin/tagger

WORKDIR /
EXPOSE 60000

ARG git_commit=unknown
ARG version="2.9.0"
ARG descriptive_version=unknown

LABEL org.cyverse.git-ref="$git_commit"
LABEL org.cyverse.version="$version"
LABEL org.cyverse.descriptive-version="$descriptive_version"
LABEL org.label-schema.vcs-ref="$git_commit"
LABEL org.label-schema.vcs-url="https://github.com/cyverse-de/tagger"
LABEL org.label-schema.version="$descriptive_version"

FROM --platform=linux/amd64 zenika/alpine-chrome:with-chromedriver
###############################################################################
#                                LiveDog-build
###############################################################################

LABEL maintainer="shichen437 <shichen437@126.com>"
LABEL version="0.0.8"
LABEL description="LiveDog 监控平台"
LABEL license="MIT"
LABEL source="https://github.com/shichen437/live-dog"

ENV WORKDIR=/LiveDog
ENV TZ=Asia/Shanghai
ENV PROJECT_UPLOAD=$WORKDIR/upload \
    PROJECT_OUTPUT=$WORKDIR/video

USER root

RUN mkdir -p $PROJECT_OUTPUT && \
    mkdir -p $PROJECT_UPLOAD && \
    apk --no-cache add nginx ffmpeg libc6-compat tzdata coreutils && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

ADD web/nginx.conf /etc/nginx/nginx.conf

COPY i18n $WORKDIR/i18n
COPY manifest/config/config_docker.yaml $WORKDIR/config.yaml
COPY manifest/docker/entrypoint.sh /entrypoint.sh
COPY manifest/migrate $WORKDIR/manifest/migrate
COPY temp/linux_amd64/main $WORKDIR/main
COPY web/dist /usr/share/nginx/html

RUN chmod +x /entrypoint.sh

VOLUME $PROJECT_OUTPUT
VOLUME $PROJECT_UPLOAD

ENV DATABASE_DEFAULT_LINK="" \
    PROJECT_SM4KEY="abcdefghijklmnopqrstuvwxyz123456"

EXPOSE 9876

WORKDIR $WORKDIR
ENTRYPOINT ["sh"]
CMD ["/entrypoint.sh"]
FROM --platform=linux/amd64 ubuntu:20.04

LABEL maintainer="shichen437 <shichen437@126.com>"
LABEL version="centos"
LABEL description="LiveDog 监控平台"
LABEL license="MIT"
LABEL source="https://github.com/shichen437/live-dog"

ENV WORKDIR=/LiveDog
ENV TZ=Asia/Shanghai
ENV PROJECT_UPLOAD=$WORKDIR/upload \
    PROJECT_OUTPUT=$WORKDIR/video

USER root

# Set non-interactive mode for apt-get
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y wget gnupg2 nginx && \
    wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list' && \
    apt-get update && \
    apt-get install -y google-chrome-stable ffmpeg tzdata && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    useradd -r -d /var/cache/nginx -s /sbin/nologin nginx

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
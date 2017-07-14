FROM alpine:3.6
RUN apk --update add nodejs nodejs-npm

# ADD ./fwatchdog /usr/bin
ADD https://github.com/alexellis/faas/releases/download/0.5.6-alpha/fwatchdog /usr/bin

RUN chmod +x /usr/bin/fwatchdog

COPY package.json .
COPY main.js .
RUN npm i
ENV fprocess="node main.js"

HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["fwatchdog"]

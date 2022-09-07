#!/bin/sh

docker run \
--name nginx \
--label "sui.app.icon=web" \
--label "sui.app.name=nginx" \
--label "sui.app.url=nginx.mydomain.xyz" \
nginx
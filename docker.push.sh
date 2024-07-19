#!/bin/bash
STEEDOS_VERSION="v1.638.9";

aws ecr get-login-password --region cn-northwest-1 | docker login --username AWS --password-stdin 252208178451.dkr.ecr.cn-northwest-1.amazonaws.com.cn

docker tag steedos/casdoor:$STEEDOS_VERSION 252208178451.dkr.ecr.cn-northwest-1.amazonaws.com.cn/dockerhub/steedos/casdoor:$STEEDOS_VERSION

docker push 252208178451.dkr.ecr.cn-northwest-1.amazonaws.com.cn/dockerhub/steedos/casdoor:$STEEDOS_VERSION
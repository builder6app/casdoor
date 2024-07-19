# Docker Build
```
docker build --target STANDARD -t steedos/casdoor:v1.638.2 .
```

buyProcuct接口支持传入PaymentEnv='InWechatMiniProgram:小程序Id'


# 推国内镜像准备
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
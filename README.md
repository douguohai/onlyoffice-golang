# onlyoffice-golang
基于golang go语言（beego框架）下的ONLYOFFICE Document Server二次开发。
主要功能为文档的上传、预览、覆盖、回调等功能。

```bash
docker run -i -t -d --restart=always --name onlyoffice-documentServer-server -p 30080:80 -e REDIS_SERVER_HOST=192.168.10.239 -e REDIS_SERVER_PORT=6379 -e REDIS_SERVER_PASS=redis2020! -e DB_TYPE=postgres -e DB_HOST=192.168.10.240 -e DB_PORT=5432 -e DB_NAME=document -e DB_USER=postgres -e DB_PWD=Xtm@123456 douguohai/onlyoffice-documentserver:7.1.1.76

docker run --name postgres -e POSTGRES_PASSWORD=woshidoudou -p 5432:5432 -d postgres:9.6

docker build -t douguohai/onlyoffice-golang:v11 . 

docker run -d -p 30081:8080 -e serverUrl=https://q.sss-xtm.com:30081 -e wsServer=wss://q.sss-xtm.com:30081/ws -e documentServer=https://q.sss-xtm.com:30080 -e dbHost=192.168.10.240 -e dbPassword=Xtm@123456 douguohai/onlyoffice-golang:v12
```

#### 上传文件，获取访问连接
![img.png](document/img/img_upload.png)

#### 获取下载文档连接
![img.png](document/img/img_download.png)

#### 强行覆盖文件内容
![img.png](document/img/img_overwrite.png)




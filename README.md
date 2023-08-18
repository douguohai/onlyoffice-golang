# ONLYOFFICE-golang
基于golang go语言（beego框架）下的ONLYOFFICE Document Server二次开发。主要功能为文档的管理（上传、列表、删除等）和根据ONLYOFFICE Document Server中的文档进行更新和回调、历史版本数据管理。这个项目来自3xxx/EngineerCMS，完整的代码也可参考EngineerCMS。

```bash
docker run -i -t -d --restart=always --name onlyoffice-documentServer-server -p 30080:80 -e REDIS_SERVER_HOST=192.168.10.239 -e REDIS_SERVER_PORT=6379 -e REDIS_SERVER_PASS=redis2020! -e DB_TYPE=postgres -e DB_HOST=192.168.10.240 -e DB_PORT=5432 -e DB_NAME=document -e DB_USER=postgres -e DB_PWD=Xtm@123456 douguohai/onlyoffice-documentserver:7.1.1.76

docker build -t douguohai/onlyoffice-golang:v8 . 


docker run -d -p 18080:8080 -e serverUrl=http://192.168.1.183:18080 -e wsServer=192.168.1.183:18080 -e documentServer=http://192.168.1.183:30080  douguohai/onlyoffice-golang:v9

docker run --name postgres -e POSTGRES_PASSWORD=woshidoudou -p 5432:5432 -d postgres:9.6

```



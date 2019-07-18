# demo3
## 服务注册

需要自己实现注册和健康检查的整个逻辑~~

默认consul地址：127.0.0.1:8500

## step
1. run registry service
```
cd registry
go build -o regi
./regi
```
2. run discovery service
```
cd discovery
go build -o dis
./dis
```
3. test
```
curl -d '{"str":"world"}' localhost:9001/calculate
```
# grpc-gateway
grpc-gateway是一个基于grpc和protobuf开发的个人gateway项目，封装比较优雅，API友好，源码注释比较明确，具有快速灵活，容错方便等特点，让你快速了解golang项目中grpc和protobuf的使用

### 环境依赖
node v0.10.28+  
redIs ~  

### 部署步骤
1. 创建protobuf中间文件protos/goods/goods.proto    
   ```go
    syntax = "proto3";
         package services;
         import "goods_model.proto";
    
         message GoodsRequest {
             // @inject_tag: json:"size",form:"size"
             int32 size = 1;
             // @inject_tag: uri:"goods_id"
             int32 goods_id=2;
         }
    
         message GoodsListResponse {
             repeated GoodsModel data=1;
         }
    
         message GoodsDetailResponse {
             GoodsModel data=1;
         }
    
         service GoodsService {
             rpc GetGoodsList(GoodsRequest) returns(GoodsListResponse);
             rpc GetGoodsDetail(GoodsRequest) returns(GoodsDetailResponse);
         }
    ```   

2. proto文件编译  //protoc编译
    protoc --go_out=./ goods.proto

3. 编译项目   //编译  
    go build

4. 启动项目  
    go run main.go  


###目录结构描述
├── Readme.md                   // help  
├── services                    // 应用服务  
├── config                      // 配置  
│   ├── default.json  
│   ├── dev.json                // 开发环境  
│   ├── experiment.json         // 实验  
│   ├── index.js                // 配置控制  
│   ├── local.json              // 本地  
│   ├── production.json         // 生产环境  
│   └── test.json               // 测试环境  
├── data  
├── protos                      // protos文档  
├── environment  
├── gulpfile.js  
├── locales  
├── logger-service.js           // 启动日志配置  
├── node_modules  
├── package.json  
├── app-service.js              // 启动应用配置  
├── static                      // web静态资源加载  
│   └── initjson  
│       └── config.js         // 提供给前端的配置  
├── test  
├── test-service.js  
└── tools  

###V1.0.0 版本内容更新
1. 新功能     获取商品信息
2. 新功能     获取优惠券信息
3. 新功能     获取会员信息
4. 新功能     服务熔断降级 
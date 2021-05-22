> 第一次使用golang的记录
>
### 使用GO Modules管理依赖

初始化Mod文件 ` go mod init GGFramework `

检测代码中依赖，并写入go.mod中 `go mod tidy`

下载依赖 `go mod download`

### gin 框架使用

1. 添加router 
   - GET/POST
   - REST（路径作为参数）
2. 解析参数
   - query/postform
   - params
   - bind & validator
3. 处理

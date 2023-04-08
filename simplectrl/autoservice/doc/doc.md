### 如何开启 auto service
为 deployment 指定 `"enable-auto-service": "true"` label，则创建后会自动生成对应的 service

### 如果将 enable-auto-service 改为 false 或者删除，如何处理已经创建出来的 service
auto service 的 name 和 namespace 与其所属的 deployment 相同，可以先查一下是否
存在这个 service，存在的话可以将其删除

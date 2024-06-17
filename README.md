# 项目描述
* 这是一个高性能短码服务，可在此基础上扩展成短链接服务

# 概要设计
## 容量
* 6位短码通过Base52编码，可创建200亿个短码，每个短码的有效期2年
* 200亿个短码存储约是500G，按照3副本存储设计约占磁盘1500G
* 每个短码绑定的数据约是256Bit，约5500G，按照3副本存储谁约占磁盘16500G
* Redis会缓存最近30天的短码绑定数据，如果满负荷，那么预计需要460G内存，请按照实际情况合理配置过期时间
* 本地会缓存500万个短码绑定数据，有效期7天，占用应用内存约1.25G,建议单个应用内存不低于2G

## 设计简述
* 短码、短码绑定数据使用ClickHouse集群存储
* 数据绑定短码：本地会预加载20000个未使用的短码，从本地缓存获取一个未使用的短码，绑定后数据暂存本地内存，绑定关系刷新到Redis，暂存数据待一定条件后持久化到ClickHouse
  * 当本地缓存的未使用的短码低于10000的时候会异步去数据库再加载20000条数据追加到本地缓存
  * 当暂存的数据超过200条会将数据刷写进ClickHouse
  * 定时任务会每20秒将暂存数据持久化到ClickHouse
* 通过短码获取绑定的数据：先从本地缓存读取，未命中去Redis读取，Redis未命中再去ClickHouse查询
  * 未命中的情况数据会反写到缓存中
  * 查询ClickHouse会进行限流配置，避免因大量未命中拖垮数据库

## 可用性
* 数据绑定短码：
  * Redis不可用：绑定数据会写入ClickHouse
  * ClickHouse不可用：绑定数据会写入Redis，暂存数据写入一份到Redis，等待ClickHouse恢复暂存数据持久化到ClickHouse
  * Redis和ClickHouse都不可用：用户无法提价绑定，服务不可用
  * Redis和ClickHouse都不可用并且服务宕机：如果Redis宕机数据丢失，则暂存的数据丢失不可恢复。如果Redis数据还在，则应用会读取Redis中的暂存数据持久化到ClickHouse

* 通过短码获取绑定的数据：
  * Redis不可用：本地缓存无法命中则读取ClickHouse
  * ClickHouse不可用：两个缓存都无法命中则获取数据失败
  * Redis和ClickHouse都不可用：本地缓存没有命中则获取数据失败
  * Redis和ClickHouse都不可用并且服务宕机：服务不可用

## 构建项目
* 搭建数据库集群
* [初始化表结构](doc/install.sql)
* [执行短码库脚本](build_code.go)：生成CSV短码数据，将数据导入到ClickHouse。200亿数据预计需要1.5T磁盘，执行脚本的时候请先确保磁盘够用，建议使用最少16G内存的机器
* 修改配置文件 [dev.yaml](initialize%2Fdev.yaml)
* 启动服务：[执行](Web_QueryAndBind.go) go run Web_QueryAndBind.go dev
* [Web_QueryAndBind.go](Web_QueryAndBind.go)提供了查询和绑定接口；[Web_Bind.go](Web_Bind.go)只提供了绑定接口；[web_Query.go](web_Query.go)只提供了查询接口


# 压测情况
* ClickHouse：8C64G * 4
* Redis：4C8G * 6
* 应用服务：8C16G * 1
* 压力肉鸡：8C16G * 1
* 使用jmeter压测
* 数据绑定短码压测 线程数:256 循环次数:200 ![绑定数据压测.jpg](doc%2Fimages%2F%E7%BB%91%E5%AE%9A%E6%95%B0%E6%8D%AE%E5%8E%8B%E6%B5%8B.jpg)
* 通过短码获取绑定的数据 线程数:512 循环次数:250 50%缓存命中 ![查询_50%命中.jpg](doc%2Fimages%2F%E6%9F%A5%E8%AF%A2_50%25%E5%91%BD%E4%B8%AD.jpg)
* 通过短码获取绑定的数据 线程数:512 循环次数:250 60%缓存命中 ![查询_60%命中.jpg](doc%2Fimages%2F%E6%9F%A5%E8%AF%A2_60%25%E5%91%BD%E4%B8%AD.jpg)
* 通过短码获取绑定的数据 线程数:512 循环次数:250 70%缓存命中 ![查询_70%命中.jpg](doc%2Fimages%2F%E6%9F%A5%E8%AF%A2_70%25%E5%91%BD%E4%B8%AD.jpg)
* 通过短码获取绑定的数据 线程数:512 循环次数:250 80%缓存命中 ![查询_80%命中.jpg](doc%2Fimages%2F%E6%9F%A5%E8%AF%A2_80%25%E5%91%BD%E4%B8%AD.jpg)
* 通过短码获取绑定的数据 线程数:512 循环次数:250 90%缓存命中 ![查询_90%命中.jpg](doc%2Fimages%2F%E6%9F%A5%E8%AF%A2_90%25%E5%91%BD%E4%B8%AD.jpg)

## 限制
* 200亿个短码，2年的有效期；平均每天绑定不得超过2730万个，平均每秒的绑定不得超过317个；超过限制建议预警
* 通过短码获取绑定的数据 建议单台QPS控制在3000以内，

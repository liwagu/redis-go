

## Menu

### 基础库
* ./example/common/concurrent_routine.go 示例使用的并发执行器
* ./example/common/concurrent_event_logger.go 示例使用的日志搜集及打印工具
* ./example/redis_client.go redis_client初始化工具

### 案例
#### 签到案例
* ./example/ex01_checkin.go 
> 执行命令
```shell
go run main.go Ex01 1165894833417101
```

#### 基于setnx的简化版分布式锁
* ./example/ex02_setnx.go 
> 执行命令：
```shell
go run main.go Ex02
```

#### 基于incr,decr的简单限流
* ./example/ex03_limiter.go
> 执行命令：
```shell
go run main.go Ex03
```

#### 消息通知
* ./example/ex04_list.go
> 执行命令：
```shell
go run main.go Ex04
```

#### 用户各类计数保存
* ./example/ex05_hash.go
> 执行命令：
```shell
go run main.go Ex05 init
go run main.go Ex05 get 1556564194374926
go run main.go Ex05 incr_collect 1556564194374926
go run main.go Ex05 decr_collect 1556564194374926
```

#### 排行榜
* ./example/ex06_zset.go
> 执行命令：
```shell
go run main.go init # 初始化计数
go run main.go Ex06 rev_order # 按score逆序输出排行榜
go run main.go Ex06 get_score user2 # 获取user2的积分
go run main.go Ex06 add_user_score user2 10 # 为用户user2增加10分
```


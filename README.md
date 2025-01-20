

# 测试环境





## 测试的数据库和表

```sql
#业务测试表
create database if not exists dtm_test
    DEFAULT CHARACTER SET utf8mb4;


#balance是金额
#trading_balance 是被冻结的金额
CREATE TABLE dtm_test.`user_account` (
 `id` int(11) AUTO_INCREMENT PRIMARY KEY,
 `user_id` int(11) not NULL UNIQUE ,
 `balance` decimal(10,2) NOT NULL DEFAULT '0.00',
 `trading_balance` decimal(10,2) NOT NULL DEFAULT '0.00',
 `create_time` datetime DEFAULT now(),
 `update_time` datetime DEFAULT now()
);

#插入测试数据
insert into dtm_test.user_account(user_id, balance)
values (1, 1000),(2, 1000);
```



## 需要修改的环境

在configs/config.dev.yaml里

修改etcd的连接，并且dtm需要注册到etcd中

修改dtmConf的连接，连接到业务数据库

然后就无需修改了

在test目录里可以进行测试












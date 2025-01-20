

# dtm-kratos-xa

## 结合kraots

我们的dtm需要指定一个api地址，在分布式环境下，我们的肯定要从etcd中取api地址，kratos服务发现结合了dtm，可以很方便的使，并且dtm自身也要注册到etcd中，在配置部分可以手动配置

```go
//参考
https://dtm.pub/ref/kratos.html
```





## 注册resolver

在initEtcd里加入下面的代码，不然dtm无法通过服务名称去访问我们的服务

```go

// 导入 kratos 的 dtm 驱动
_ "github.com/dtm-labs/driver-kratos"

//注册全局的resolver，现在我们的业务可以使用discovery:///dtmservice 来访问dtm和服务名称
//  记得引入driver-kratos的驱动
//我们在kratos里调用其它微服务，也是是通过discovery:///来的
//  这一点kratos给我们实现了，然后dtm相当于是接入kratos的这种方式
resolver.Register(
    discovery.NewBuilder(r, discovery.WithInsecure(true)))
```








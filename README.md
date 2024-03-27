# iGen

超级代码生成器，致力于生成标准化模式代码

## decorate

### 思想
借鉴了 mockgen，并基于其代码删繁就简改造而来。在我们的实践开发中，通常会遇到非常广泛的外部调用，这些外部调用有可能是 mysql、redis 这种，也有一些内部第三方调用，譬如资源管理与 admin 调度扩展 腾讯云之间的互相调用。这些调用，由于每个项目的设计规范不同，我们很难直接采用统一的抽象去审计这些接口调用，此外，在每个接口调用函数中单独定义一套，实现起来又很脏，而且在代码迭代中，可能出现不同函数审计日志规范不同的问题。

所以我决定定义这个工具。

它可以基于 golang interface，使用装饰器模式，包装具体的业务实现，在每个业务实现的函数上，定义打印日志和上报 Promtheus（可选）的逻辑。

### 日志
生成代码中可以自带函数的 Req 和 Rsp，以及函数执行耗时。你需要保证引入的 log 包，包含如下函数
```golang
log.Infof(format string, args ...interface{})
```

为在生成的代码中自动 import 对应的 log 包，你需要指定如下参数
```bash
--imports=github.com/chengyumeng/igen/log
```

对于 Req 和 Rsp 参数，为了提高性能和保护隐私，用户需要自行定义不同的 Req 和 Rsp 类型，指定不同的解析函数，igen 给了一个示例，请参见 sample/decorate/explain.json
举几个例子：
- "" 代表原值识别，适用于 go 的基本类型
- "#json" 会使用 golang 标准库的 json.Mashal 忽略掉 err 解析
- ".func" 代表该数据结构存在导出函数可用于解析 ，生成代码为 req.(func)
- "func" 代表可使用特定函数解析它，生成代码为 func(req)

生成的代码示例
```golang
// QueryTKECVMOrder audit base method.
func (m *AuditInterface) QueryTKECVMOrder(ctx context.Context, orderID string) (*QueryOrderRsp, error) {
	begin := time.Now()
	ret0, err := m.client.QueryTKECVMOrder(ctx, orderID)
	ret0_2, _ := json.Marshal(ret0)
	log.Infof(context.Background(), "function: QueryTKECVMOrder, time: %s req:  %v rsp:  %v %v", time.Since(begin), orderID, string(ret0_2), err)
	return ret0, err
}
```

### Promtheus
igen 要求你的 promtheus 上报函数满足如下格式的要求：
```golang
func Prom(packageName string,function string,err error)
```
由于该函数是自定义的，所以你需要传递如下参数指定你的函数模板,如果没有指定 --prom 将不会生成 prom 上报代码
```bash
--prom==prom.Create
--imports=github.com/chengyumeng/igen/pkg/basic/prom
```
生成的代码示例
```golang
// QueryTKECVMOrder audit base method.
func (m *AuditInterface) QueryTKECVMOrder(ctx context.Context, orderID string) (*QueryOrderRsp, error) {
	begin := time.Now()
	ret0, err := m.client.QueryTKECVMOrder(ctx, orderID)
	ret0_2, _ := json.Marshal(ret0)
	prom.Create("cloudapi", "QueryTKECVMOrder", err)
	log.Infof(context.Background(), "function: QueryTKECVMOrder, time: %s req:  %v rsp:  %v %v", time.Since(begin), orderID, string(ret0_2), err)
	return ret0, err
}
```

### 最佳实践
生成的代码应该放置到与原代码同一包内，且文件名后缀建议为 _audit.go 为了保证代码无问题，实现这个目标需要指定如下变量
```
--self-package="github.com/chengyumeng/igen"
```

### 生产环境示例
```bash
igen decorate -s=github.com/chengyumeng/igen/interface.go  \
--explain=./sample/decorate/explain.json \
--self-package=github.com/chengyumeng/igen/pkg/kube \
--imports=github.com/chengyumeng/igen/log,github.com/chengyumeng/igen/pkg/basic/prom \
--prom=prom.CreateRequestCounter \
--destination=./sample/decorate/gen.txt \
```

### 其他

##### interface 下函数如何不走审计逻辑？
使用注释 //igen:skip func1,func2,func3,func4

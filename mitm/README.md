# ！本README由AI编写
# MITM 代理模块

这个模块提供了一个简洁的 MITM (Man-in-the-Middle) 代理服务，用于拦截和处理网络请求。

## 架构设计

```
mitm/
├── model.go           # 数据模型和接口定义
├── manager.go         # 高级管理接口
├── proxy.go          # 代理核心逻辑
├── login_processor.go # 登录处理器实现
├── setup.go          # 系统代理设置
└── README.md         # 文档说明
```

## 核心组件

### 1. MitmManager (管理器)
提供高级的代理管理接口，简化使用流程：

```go
manager := mitm.NewMitmManager()

// 捕获登录信息（带超时）
loginInfo, err := manager.CaptureLogin(30 * time.Second)

// 通用捕获
processors := []mitm.MitmProcessor{processor1, processor2}
err := manager.StartCapture(processors)
```

### 2. MitmProcessor (处理器接口)
定义数据处理逻辑的标准接口：

```go
type MitmProcessor interface {
    Match(url string) bool
    Process(capture *MitmCapture) (interface{}, error)
    GetName() string
}
```

### 3. ProxyManager (代理管理器)
管理底层代理服务的生命周期和状态。

### 4. SystemProxyManager (系统代理管理器)
处理 Windows 系统代理设置。

## 使用示例

### 基本使用
```go
// 使用默认管理器
loginInfo, err := mitm.DefaultManager.CaptureLogin(30 * time.Second)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("登录成功: %s\n", loginInfo.AccessToken)
```

### 自定义处理器
```go
type CustomProcessor struct {
    pattern *regexp.Regexp
}

func (p *CustomProcessor) Match(url string) bool {
    return p.pattern.MatchString(url)
}

func (p *CustomProcessor) Process(capture *mitm.MitmCapture) (interface{}, error) {
    // 自定义处理逻辑
    var data MyDataType
    if err := capture.ParseJSON(&data); err != nil {
        return nil, err
    }
    return &data, nil
}

func (p *CustomProcessor) GetName() string {
    return "CustomProcessor"
}
```

## 特性

- **简洁的 API**: 提供高级和低级两套 API，满足不同需求
- **类型安全**: 强类型的数据结构和接口设计
- **错误处理**: 完善的错误处理和资源清理
- **并发安全**: 线程安全的状态管理
- **可扩展**: 通过处理器接口轻松扩展功能

## 注意事项

1. 需要管理员权限来修改系统代理设置
2. 确保证书文件存在于 `./cert` 目录
3. 代理监听端口为 9080，确保端口未被占用
4. 使用完毕后会自动清理系统代理设置
   - 专注于MITM代理的启动、停止和管理
   - 支持添加多种抓取目标处理器
   - 提供通用的`StartMitmProxy()`接口

3. **login_processor.go** - 登录抓取处理器
   - 独立的登录数据处理逻辑
   - 实现`MitmProcessor`接口

4. **setup.go** - 系统代理设置
   - Windows系统代理的设置和清理

## 主要改进

### 1. 模块化设计
- 将登录抓取独立为单独的处理器
- proxy模块专注于代理管理
- 每个处理器都是独立的，可以单独测试和维护

### 2. 扩展性增强
- 通过`MitmProcessor`接口，可以轻松添加新的抓取目标
- 支持同时运行多个处理器
- 结果通过通用的`interface{}`类型传递，支持不同数据类型

### 3. 代码复用
- 代理启动/停止逻辑复用
- 通用的抓包数据结构
- 统一的错误处理

## 使用方式

### 单目标抓取（登录）
```go
loginInfo, err := StartCaptureForLogin()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("登录成功: %s\n", loginInfo.AccessToken)
```

### 多目标抓取
```go
resultChan := make(chan interface{}, 10)
processors := []MitmProcessor{
    NewLoginProcessor(),
    NewGachaProcessor(),
}

err := StartMitmProxy(processors, resultChan)
if err != nil {
    log.Fatal(err)
}

// 处理结果
for result := range resultChan {
    switch data := result.(type) {
    case *LoginInfo:
        // 处理登录信息
    case *GachaInfo:
        // 处理抽卡记录
    }
}
```

### 添加新的处理器
1. 创建新的处理器文件（如`xxx_processor.go`）
2. 实现`MitmProcessor`接口的三个方法：
   - `GetName()`: 返回处理器名称
   - `Match(url string)`: 判断是否处理该URL
   - `Process(capture *MitmCapture)`: 处理抓包数据
3. 在使用时将处理器添加到processors列表中

## 接口设计

### MitmProcessor接口
```go
type MitmProcessor interface {
    Match(url string) bool                              // 检查是否需要处理该URL
    Process(capture *MitmCapture) (interface{}, error) // 处理抓包数据
    GetName() string                                    // 获取处理器名称
}
```

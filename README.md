# fabric-performance
基于go sdk 的性能测试工具

## fabric 网络性能测试

#### 项目介绍

```
主要测试fabric吞吐量和并发，基于fabric-sdk-go，测试工具使用wrk
```

#### 项目依赖


```
- fabric-sdk-go （它本身有很多依赖）
- git clone https://github.com/wg/wrk.git
```

#### 安装

```
go get -u learnergo/fabric-performance-test

cd $GOPATH/src/github.com/learnergo/fabric-performance-test
```

#### chaincode 

测试链码（官网），主要做了存(put)取(get)操作，存的过程加入了加解密操作增加复杂度，cli 操纵示例：


```
peer chaincode query -C mychannel -n mycc -c '{"Args":["get","a"]}'
```

```
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n mycc -c '{"Args":["put","a",b"]}'
```

#### 实现思路

```
因为主要侧重吞吐和并发测试，对通道和链码安装部分不是很侧重。
本项目fabric网络有mychannel通道和testcc链码（名字可以自己确定，并在程序和配置中对应修改）。
在用命令行创建链码时，先存入一个值对（"a":"b"）,取的测试是取a值；存的测试是存入当前时间戳。

**特别提示**：
为了避免日志打印对性能影响，只打印了error日志。运行正常的标志也就是没有日志打印
```

#### 操作步骤
- 配置fixtures下证书密钥和配置文(只配置一个peer和orderer即可)
- 运行main.go
- 在新窗口用wrk进行测试（调整-t 和-c 值即可，-d 越大越准确）

#### 测试环境

```
多机
Linux VM-0-17-ubuntu 4.4.0-91-generic #114-Ubuntu SMP Tue Aug 8 11:56:56 UTC 2017 x86_64 x86_64 x86_64 GNU/Linux
```


#### 测试结果
读写并发都在1000以上，读的tps在400以上，写在100以上

本人tps最佳参数：

```
./wrk -t4 -c150 -d10 --timeout 10 http://localhost:8026/v1/gettest

./wrk -t4 -c100 -d10 --timeout 30 http://localhost:8080/puttest
```

#### 影响因素


```
- 节点数量
- 服务器配置（cpu 内存 网络等等）
- 日志级别（级别越低性能越低）
- 是否启用tls（不启用tps高）
- solo or kafka （solo 高）
- leveldb couchdb选择（leveldb 高）
- orderer 出块配置（自己研究吧）
```
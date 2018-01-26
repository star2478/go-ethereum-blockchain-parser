# 基于geth对以太坊区块链的解析
* 本地[下载安装geth](https://github.com/ethereum/go-ethereum)
* 本地设置环境变量$GOPATH，可以指定为任意目录
* 下载依赖的开源库，依次执行以下命令
    * `go get github.com/ethereum/go-ethereum`
    * `go get github.com/deckarep/golang-set`
* 执行`geth --rpc`启动geth的http服务器
* 执行golang脚本`nohup go run xx.go &`
    * getAccount.go 负责并发解析多个block，获取这些block里所有交易account，并写入accounts目录的目标文件
    * getBalance.go 负责获取一个源文件里每个account的余额，将结果写入目标文件
    * getBalanceMulti.go 负责并发获取多个源文件里每个account的余额，将结果写入多个目标文件
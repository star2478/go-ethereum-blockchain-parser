# 基于geth对以太坊区块链的解析
* 本地[下载安装geth](https://github.com/ethereum/go-ethereum)，执行`geth --rpc`启动geth的http服务器
* 下载开源库go-ethereum到本地当前目录：进入本文件同级目录，执行git clone https://github.com/ethereum/go-ethereum
* 执行golang脚本`nohup go run xx.go &`
    * getAccount.go 负责并发解析多个block，获取这些block里所有交易account，并写入accounts目录的目标文件
    * getBalance.go 负责获取一个源文件里每个account的余额，将结果写入目标文件
    * getBalanceMulti.go 负责并发获取多个源文件里每个account的余额，将结果写入多个目标文件
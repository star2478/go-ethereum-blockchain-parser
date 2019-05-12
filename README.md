# 基于geth对以太坊区块链的解析
* 本地[下载安装geth](https://github.com/ethereum/go-ethereum)
* 本地设置环境变量$GOPATH，可以指定为任意目录
* 下载依赖的开源库，依次执行以下命令
    * `go get github.com/ethereum/go-ethereum`
    * `go get github.com/deckarep/golang-set`
* 执行`geth --rpc`启动geth的http服务器
* 可选操作：执行`go run synBlockTime.go`，将本地所有block的编号和时间写入blocktime文件。当我们使用后面功能时，脚本会从blocktime中读取出生成于[timeFrom]到[timeTo]间的所有block编号，再根据block编号获取block详细信息。因此，如果不执行该操作，blocktime文件可能未记录最新block，使用后面功能时就无法获取最新block的详细信息
* 功能
    * 获取一段时间内所有交易账户：`go run getAccount.go [timeFrom] [timeTo]`, 比如`go run getAccount.go 2018-01-01-00-00-00 2018-02-01-00-00-00`。结果会存到accounts目录下，文件名为[timeFrom]-[timeTo]
    * 获取一段时间内所有交易账户及其余额：`go run getBalance.go [timeFrom] [timeTo]`, 比如`go run getBalance.go 2018-01-01-00-00-00 2018-02-01-00-00-00`。结果会存到balance目录下，文件名为[timeFrom]-[timeTo]。注意：执行此命令前需要先执行`go run getAccount.go [timeFrom] [timeTo]`
    * 获取一段时间内所有交易明细：`go run getTxByTime.go [timeFrom] [timeTo]`, 比如`go run getTxByTime.go 2018-01-01-00-00-00 2018-02-01-00-00-00`。结果会存到tx目录下，文件名为[timeFrom]-[timeTo]。同时还会产生另外两个文件[timeFrom]-[timeTo]-from-sort和[timeFrom]-[timeTo]-to-sort，分别存放以交易卖出账户和买进账户排序后的结果
    * 获取一段时间内每个交易账户的交易明细：`go run getTxTimelineGroupByAccount.go [timeFrom] [timeTo]`, 比如`go run getTxTimelineGroupByAccount.go 2018-01-01-00-00-00 2018-02-01-00-00-00`。结果会存到tx目录下，文件名为[timeFrom]-[timeTo]-timeline。注意：执行此命令前需要先执行`go run getTxByTime.go [timeFrom] [timeTo]`
    * 获取一段时间内每个交易账户的交易总数（含进出总数和进出交易量）：`go run getTxCountGroupByAccount.go [timeFrom] [timeTo]`, 比如`go run getTxCountGroupByAccount.go 2018-01-01-00-00-00 2018-02-01-00-00-00`。结果会存到tx目录下，文件名为[timeFrom]-[timeTo]-count。注意：执行此命令前需要先执行`go run getTxByTime.go [timeFrom] [timeTo]`


# wiki
https://github.com/star2478/go-ethereum-blockchain-parser/wiki/%E4%BB%A5%E5%A4%AA%E5%9D%8AGo-Ethereum(Geth)%E8%83%BD%E5%81%9A%E4%BB%80%E4%B9%88

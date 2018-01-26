curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x51cd215bB9Aa24484870b91a5E61B8F6aB693A0f","latest"],"id":1}' -H "Content-type: application/json;charset=UTF-8"  localhost:8545
awk 'BEGIN{FS=" "}{keycount[$1]++;valsum[$1]+=$2}END{for(i in keycount){printf("%s %d\n", i,valsum[i]);}}' accountTxCnt 

package main

import (
  "bytes"
//  "fmt"
  "log"
  "strconv"
  "strings"
  "os/exec" 
  "math"
  "time"
  "runtime"
  "encoding/json"
  "os"
//  "./go-ethereum/common"
//  "./go-ethereum/common/hexutil"
  //"./golang-set"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/common/hexutil"
//  "github.com/deckarep/golang-set"
)

var MAX_NUM_SHARE = 10
var MAX_NUM_RETRY = 3
var APPEND_FILE_NUM = 100
var c = make(chan string, 100)

type ResultOfGetBlockByNumber struct {
    Jsonrpc    string
    Id    int
    Result	Block
}
type RPCTransaction struct {
        BlockHash        common.Hash     `json:"blockHash"`
        BlockNumber      *hexutil.Big    `json:"blockNumber"`
        From             *common.Address  `json:"from"`
        Gas              hexutil.Uint64  `json:"gas"`
        GasPrice         *hexutil.Big    `json:"gasPrice"`
        Hash             common.Hash     `json:"hash"`
        Input            hexutil.Bytes   `json:"input"`
        Nonce            hexutil.Uint64  `json:"nonce"`
        To               *common.Address `json:"to"`
        TransactionIndex hexutil.Uint    `json:"transactionIndex"`
        Value            *hexutil.Big    `json:"value"`
        V                *hexutil.Big    `json:"v"`
        R                *hexutil.Big    `json:"r"`
        S                *hexutil.Big    `json:"s"`
}

type Block struct {
        UncleHashes  []common.Hash    `json:"uncles"`
        Hash         common.Hash      `json:"hash"`
        Transactions []RPCTransaction `json:"transactions"`
}

var tenToAny map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}
 
 
func decimalToAny(num, n int) string {
 new_num_str := ""
 var remainder int
 var remainder_string string
 for num != 0 {
  remainder = num % n
  if 76 > remainder && remainder > 9 {
   remainder_string = tenToAny[remainder]
  } else {
   remainder_string = strconv.Itoa(remainder)
  }
  new_num_str = remainder_string + new_num_str
  num = num / n
 }
 return new_num_str
}

func findkey(in string) int {
 result := -1
 for k, v := range tenToAny {
  if in == v {
   result = k
  }
 }
 return result
}
func anyToDecimal(num string, n int) int {
 var new_num float64
 new_num = 0.0
 nNum := len(strings.Split(num, "")) - 1
 for _, value := range strings.Split(num, "") {
  tmp := float64(findkey(value))
  if tmp != -1 {
   new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
   nNum = nNum - 1
  } else {
   break
  }
 }
 return int(new_num)
}

func exec_shell(s string, blockNumber int, retry int, fileContent *string) {
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer

    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        //log.Fatal("cmd="+s+" run fail!", err)
	retry--
        log.Print("Fatal!!!retry=",retry," cmd="+s+" run fail!", err)
	if(retry > 0) {
		exec_shell(s, blockNumber, retry, fileContent)
    	}
    }
    
    stb := &ResultOfGetBlockByNumber{}
    //fmt.Printf("%s\n", out.String())
    err = json.Unmarshal([]byte(out.String()), &stb)
    
    if err != nil {
	retry--;
        log.Print("umarshal fail!!!retry=",retry," block=",blockNumber," cmd=", s, " input=", out.String()," ", err)
	if(retry > 0) {
		exec_shell(s, blockNumber, retry, fileContent)
    	}
    } 
    //fmt.Printf("----result=%s:%d\n", stb.Jsonrpc, stb.Id)
    blockInfo := stb.Result
    if len(blockInfo.Transactions) <= 0 {
	//log.Print("block=", blockNumber, " transaction len <= 0")
	return
    }
    //accountSet := mapset.NewSet() 
    for _, value := range blockInfo.Transactions {
	    //from := ""
	    if value.From != nil {
 	    //	from = value.From.String()
//		accountSet.Add(value.From.String())
	    	*fileContent += value.From.String() + "\n";
	    }
	    //to := ""
	    if value.To != nil {
 	    //	to = value.To.String()
//		accountSet.Add(value.To.String())
	    	*fileContent += value.To.String() + "\n";
	    }
	    //fmt.Printf("----block=%d i=%d fr=%s, to=%s\n", blockNumber, index, from, to)
    }
    //log.Print(accountSet)
}

func pathExists(path string) (bool) {
        _, err := os.Stat(path)
        if err == nil {
                return true
        }
        return false
}

func getAccount(fromBlockNumber int, toBlockNumber int, fileName string) {
	taskId := strconv.Itoa(fromBlockNumber) + "-" + strconv.Itoa(toBlockNumber)
	log.Print(taskId," begin: file=", fileName)
        fd,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
        if err != nil {
		log.Fatal(fileName, " open fail! ", err)
		return
        }
        fileContent := ""
	i := 0
	for blockNumber := fromBlockNumber; blockNumber <= toBlockNumber; blockNumber++ {
  		blockNumberHex := "0x" + decimalToAny(blockNumber, 16);
    		command := "curl -X POST --data '{\"jsonrpc\":\"3.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"" + blockNumberHex + "\", true],\"id\":1}' -H \"Content-type: application/json;charset=UTF-8\"  localhost:8545";

    	 	exec_shell(command, blockNumber, MAX_NUM_RETRY, &fileContent)
		if i%APPEND_FILE_NUM == 0 {
			if len(fileContent) > 0 {
                        	buf:=[]byte(fileContent)
                        	fd.Write(buf)
                        	fileContent = ""
				log.Print("write accounts into file=", fileName, " taskId=", taskId, ", i=",i)
			} else {
                        	fileContent = ""
				log.Print("fileContent len <= 0, not need to write accounts into file=", fileName, " taskId=", taskId, ", i=",i)
			}
                }
		i++
	}
	if len(fileContent) > 0 {
                buf:=[]byte(fileContent)
                fd.Write(buf)
		log.Print("write accounts into file=", fileName, " taskId=", taskId, ", i=",i)
        }
        fd.Close() 
	c <- taskId
}

func getAndCheckDir() string {
	curDir, _ := os.Getwd()  //
        newDir := curDir +"/accounts"
        if pathExists(newDir) == false {
                err := os.Mkdir(newDir, os.ModePerm)  //
                if err != nil {
                        log.Fatal(newDir, " create fail! ", err)
                        //return nil
                }
        }
	return newDir
}

func main() {

  if len(os.Args) < 3 {
        log.Fatal("Param Invalid!!! go run getAccount.go [blockNumberBegin] [blockNumberEnd]")
  }
  log.Print("getAccount begin==================");
  timeBegin := time.Now().Unix()  
  MULTICORE := runtime.NumCPU()
  runtime.GOMAXPROCS(MULTICORE)
  //blockNumber := 4927600;
  blockNumberBegin,err1 := strconv.Atoi(os.Args[1])
  blockNumberEnd,err2 := strconv.Atoi(os.Args[2]) 
  if err1 != nil || err2!=nil || blockNumberBegin > blockNumberEnd {
	log.Fatal("Param Fail!!! blockNumberBegin > blockNumberEnd or err1=", err1, " or err2=", err2);
  }
  totalBlockNum := blockNumberEnd - blockNumberBegin + 1;
  share := totalBlockNum / MAX_NUM_SHARE + 1
  loopCount := 0
  dir := getAndCheckDir()
  for i := blockNumberBegin; i <= blockNumberEnd; i++ {
	from := i
	if ((share+i) <= blockNumberEnd) {
		i += share 
	} else {
		i = blockNumberEnd
	}
	to := i
	loopCount++
        fileName := dir + "/" + strconv.Itoa(loopCount) + ".txt"
	go getAccount(from, to, fileName);
  }
  for i := 0; i < loopCount; i++ {
	taskId := <- c
	log.Print(taskId, " finish");
  }

//    command := "curl -X POST --data '{\"jsonrpc\":\"3.0\",\"method\":\"eth_getBalance\",\"params\":[\"0x51cd215bB9Aa24484870b91a5E61B8F6aB693A0f\",\"latest\"],\"id\":1}' -H \"Content-type: application/json;charset=UTF-8\"  localhost:8545";

  timeEnd := time.Now().Unix()  
  log.Print("getAccount finish, cost=", (timeEnd - timeBegin), "s")
}

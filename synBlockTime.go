package main

import (
  "bytes"
  //"sort"
  "fmt"
  "log"
  "strconv"
  "strings"
  "time"
  "runtime"
  "encoding/json"
  "os"
  "io"
  "math"
  "os/exec"
  "bufio"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/common/hexutil"
)

var MAX_NUM_SHARE = 10
var MAX_NUM_RETRY = 3
var APPEND_FILE_ACCOUNT_MAX_NUM = 100
var APPEND_FILE_LOOP_MAX_NUM = 1000
var APPEND_FILE_ACCOUNT_TX_MAX_NUM = 10000
var txMap = map[string]string{}
var tenToAny map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29:"t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62:"M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}

type ResultOfBlockNumber struct {
    Jsonrpc    string
    Id    int
    Result     string 
}
type ResultOfGetBlockByNumber struct {
    Jsonrpc    string
    Id    int
    Result      Block
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
        Timestamp    string      `json:"timestamp"`
        Transactions []RPCTransaction `json:"transactions"`
}

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
func anyToDecimal(num string, n int) float64 {
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
 return float64(new_num)
}
func getCurrentBlockNumber() int {
	 command := "curl -X POST --data '{\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\": 84}' -H \"Content-type: application/json;charset=UTF-8\" localhost:8545"
	cmd := exec.Command("/bin/bash", "-c", command)
    	var out bytes.Buffer
    	cmd.Stdout = &out
    	err := cmd.Run()
    	if err != nil {
        	log.Print("Fatal!!!cmd="+command+" run fail!", err)
	}
	stb := &ResultOfBlockNumber{}
    	err = json.Unmarshal([]byte(out.String()), &stb)

    	if err != nil {
        	log.Fatal("umarshal fail!!! cmd=", command, " input=", out.String()," ", err)
    	}
    	hexBlockNumber := stb.Result
	hexBlockNumber = hexBlockNumber[2 : len(hexBlockNumber)] // get rid of 0x prefix
	return int(anyToDecimal(hexBlockNumber, 16))
}

func getBlockTime(blockNumber int) (string,int) {
	blockNumberHex := "0x" + decimalToAny(blockNumber, 16);
        command := "curl -X POST --data '{\"jsonrpc\":\"3.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"" + blockNumberHex + "\", false],\"id\":1}' -H \"Content-type: application/json;charset=UTF-8\"  localhost:8545";

	cmd := exec.Command("/bin/bash", "-c", command)
    	var out bytes.Buffer
    	cmd.Stdout = &out
    	err := cmd.Run()
    	if err != nil {
        	log.Print("Fatal!!!cmd="+command+" run fail!", err)
	}
	stb := &ResultOfGetBlockByNumber{}
    err = json.Unmarshal([]byte(out.String()), &stb)

    if err != nil {
        log.Fatal("umarshal fail!!! block=",blockNumber," cmd=", command, " input=", out.String()," ", err)
    }
    //fmt.Printf("----result=%s:%d\n", stb.Jsonrpc, stb.Id)
    blockInfo := stb.Result
	if len(blockInfo.Timestamp) <= 0 {
		return "0", blockNumber
	}
    timeLayout := "2006-01-02-15-04-05"
    timestamp := blockInfo.Timestamp
    timestamp = timestamp[2 : len(timestamp)] // get rid of 0x prefix
    txTime := time.Unix(int64(anyToDecimal(timestamp, 16)), 0).Format(timeLayout)
	return txTime, blockNumber
}

func main() {
  	fmt.Println("synBlockTime begin==================");
  	timeBegin := time.Now().Unix()  
  	MULTICORE := runtime.NumCPU()
  	runtime.GOMAXPROCS(MULTICORE)
  	fileName := "blocktime" 
        fd,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		log.Fatal("file=", fileName," open fail!!", err);
	}
        defer fd.Close() 
	lastBlockInFile := ""
	readBuf := bufio.NewReader(fd)
	for {		
		line, errRead := readBuf.ReadString('\n')
		if errRead == io.EOF {
			break
		}
		lastBlockInFile = line
	}
	lastBlockInFile = strings.TrimSpace(lastBlockInFile)
	firstToFindBlockNumber := -1
	if len(lastBlockInFile) > 0 {
		arrBlockInfo := strings.Split(lastBlockInFile, " ")
		firstToFindBlockNumber,_ = strconv.Atoi(arrBlockInfo[1])
	}

	curMaxBlockNumber := getCurrentBlockNumber()
	log.Print("begin to write result file, curMaxBlockNumber=", curMaxBlockNumber, " firstToFindBlockNumber=", firstToFindBlockNumber+1)
	loopCount := 0
	fileContent := bytes.Buffer{}
	for i := firstToFindBlockNumber+1; i <= curMaxBlockNumber; i++ {
		txTime, blockNumber := getBlockTime(i)
	
		fileContent.WriteString(txTime + " " + strconv.Itoa(blockNumber) + "\n")
		if loopCount % APPEND_FILE_LOOP_MAX_NUM == 0 {
			wBuf:=[]byte(fileContent.String())
                	fd.Write(wBuf)
			log.Print("write block time into file=", fileName, " loopCount=", loopCount)
			//log.Print("write block time into =", txTime, blockNumber)
			fileContent = bytes.Buffer{}
		}
		loopCount++
    	}
	if len(fileContent.String()) > 0 {
		wBuf:=[]byte(fileContent.String())
               	fd.Write(wBuf)
		log.Print("write block time into file=", fileName, " loopCount=", loopCount)
	}	

  	
	timeEnd := time.Now().Unix()  
  	log.Print("synBlockTime finish, cost=", (timeEnd - timeBegin), "s")
}

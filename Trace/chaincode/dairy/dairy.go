package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"time"
)

// 定义结构体, 继承ChainCode接口
type DairyFarm struct {
}

// 定义数据结构体
type FarmInfo struct {
	Id      string
	// 奶牛场名字
	Name    string
	// 生产日期
	Date    string
	// 质量等级
	Quality string
	// 牛奶产量, 单位t
	Yield   int
}

// 方法实现
func (t* DairyFarm)Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.init(stub)
}

func (t* DairyFarm)init(stub shim.ChaincodeStubInterface) peer.Response {
	// 数据初始化
	infos := []FarmInfo{
		FarmInfo{Id:"DF-001", Name:"东郊农场", Date:"2018-12-11", Quality:"优", Yield:5},
		FarmInfo{Id:"DF-002", Name:"东郊农场", Date:"2018-12-12", Quality:"优", Yield:5},
		FarmInfo{Id:"DF-003", Name:"东郊农场", Date:"2018-12-13", Quality:"良", Yield:5},
		FarmInfo{Id:"DF-004", Name:"东郊农场", Date:"2018-12-14", Quality:"优", Yield:5},
		FarmInfo{Id:"DF-005", Name:"东郊农场", Date:"2018-12-15", Quality:"良", Yield:5},
		FarmInfo{Id:"DF-006", Name:"西郊农场", Date:"2018-12-11", Quality:"优", Yield:6},
		FarmInfo{Id:"DF-007", Name:"西郊农场", Date:"2018-12-12", Quality:"良", Yield:6},
		FarmInfo{Id:"DF-008", Name:"西郊农场", Date:"2018-12-13", Quality:"优", Yield:6},
		FarmInfo{Id:"DF-009", Name:"西郊农场", Date:"2018-12-14", Quality:"良", Yield:6},
		FarmInfo{Id:"DF-010", Name:"西郊农场", Date:"2018-12-15", Quality:"优", Yield:6},
		FarmInfo{Id:"DF-011", Name:"南郊农场", Date:"2018-12-11", Quality:"良", Yield:8},
		FarmInfo{Id:"DF-012", Name:"南郊农场", Date:"2018-12-12", Quality:"优", Yield:8},
		FarmInfo{Id:"DF-013", Name:"南郊农场", Date:"2018-12-13", Quality:"优", Yield:8},
		FarmInfo{Id:"DF-014", Name:"南郊农场", Date:"2018-12-14", Quality:"良", Yield:8},
		FarmInfo{Id:"DF-015", Name:"南郊农场", Date:"2018-12-15", Quality:"优", Yield:8},
		FarmInfo{Id:"DF-016", Name:"北郊农场", Date:"2018-12-11", Quality:"良", Yield:3},
		FarmInfo{Id:"DF-017", Name:"北郊农场", Date:"2018-12-12", Quality:"良", Yield:3},
		FarmInfo{Id:"DF-018", Name:"北郊农场", Date:"2018-12-13", Quality:"优", Yield:3},
		FarmInfo{Id:"DF-019", Name:"北郊农场", Date:"2018-12-14", Quality:"良", Yield:3},
		FarmInfo{Id:"DF-020", Name:"北郊农场", Date:"2018-12-15", Quality:"良", Yield:3},
	}
	// 遍历, 写入账本中
	i := 0
	for i < len(infos) {
		jsontext, error := json.Marshal(infos[i])
		if error != nil {
			return shim.Error("init error, json marshal fail...")
		}
		// 数据写入账本中
		stub.PutState(infos[i].Id, jsontext)
		i++
	}
	return shim.Success([]byte("init ledger OK!!!"))
}

func(t* DairyFarm)Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()

	if funcName == "setvalue" {
		return t.setvalue(stub, args)
	}else if funcName == "query" {
		return t.query(stub, args)
	}else if funcName == "gethistory" {
		return t.gethistory(stub, args)
	}
	return shim.Success([]byte("invoke OK"))
}
func(t* DairyFarm)setvalue(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	error := stub.PutState(keyID, []byte(args[1]))
	if error != nil {
		return shim.Error("PutState fail...")
	}
	return shim.Success([]byte("SetValue Sucess !!!"))
}

func(t* DairyFarm)query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	value, error := stub.GetState(keyID)
	if error != nil {
		return shim.Error("GetState fail...")
	}
	return shim.Success(value)
}

// 根据keyID查询历史记录
func(t* DairyFarm)gethistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyiter, error := stub.GetHistoryForKey(args[0])
	if error != nil {
		return shim.Error("GetHistoryForKey fail...")
	}
	defer keyiter.Close()
	// 通过迭代器对象遍历结果
	var myList []string
	for keyiter.HasNext() {
		// 获取当前值
		result, error := keyiter.Next()
		if error != nil {
			return shim.Error("keyiter.Next() fail...")
		}
		// 获取需要的信息
		txID := result.TxId
		txValue := result.Value
		txTime := result.Timestamp
		txStatus := result.IsDelete
		tm := time.Unix(txTime.Seconds, 0)
		datastr := tm.Format("2006-01-02 15:04:05")
		all := fmt.Sprintf("%s, %s, %s, %t", txID, txValue, datastr, txStatus)
		myList = append(myList, all)
	}
	// 数据格式化为json
	jsonText, error := json.Marshal(myList)
	if error != nil {
		return shim.Error("json.Marshal(myList) fail...")
	}
	return shim.Success(jsonText)
}

func main() {
	error := shim.Start(new(DairyFarm))
	if error != nil {
		println("程序启动失败...")
		return
	}
	fmt.Println("程序启动成功...")
}

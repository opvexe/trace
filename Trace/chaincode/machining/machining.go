package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"time"
)

// 定义结构体, 继承ChainCode接口
type Machining struct {
}

// 定义数据结构体
type MachinInfo struct {
	Id       string		// 加工厂ID
	FromId   string		// 奶场来源ID
	Date     string		// 生产日期
	Name     string		// 名字
	Validity int		// 保质期
}

// 方法实现
func (t* Machining)Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.init(stub)
}

func (t* Machining)init(stub shim.ChaincodeStubInterface) peer.Response {
	// 数据初始化
	infos := []MachinInfo{
		MachinInfo{Id:"JGC-001", Name:"三鹿", Date:"2018-12-11", FromId:"DF-001", Validity:45},
		MachinInfo{Id:"JGC-002", Name:"三鹿", Date:"2018-12-12", FromId:"DF-002", Validity:45},
		MachinInfo{Id:"JGC-003", Name:"三鹿", Date:"2018-12-13", FromId:"DF-020", Validity:45},
		MachinInfo{Id:"JGC-004", Name:"三鹿", Date:"2018-12-14", FromId:"DF-003", Validity:45},
		MachinInfo{Id:"JGC-005", Name:"三鹿", Date:"2018-12-15", FromId:"DF-004", Validity:45},
		MachinInfo{Id:"JGC-006", Name:"蒙牛", Date:"2018-12-11", FromId:"DF-005", Validity:45},
		MachinInfo{Id:"JGC-007", Name:"蒙牛", Date:"2018-12-12", FromId:"DF-006", Validity:45},
		MachinInfo{Id:"JGC-008", Name:"蒙牛", Date:"2018-12-13", FromId:"DF-007", Validity:45},
		MachinInfo{Id:"JGC-009", Name:"蒙牛", Date:"2018-12-14", FromId:"DF-008", Validity:45},
		MachinInfo{Id:"JGC-010", Name:"蒙牛", Date:"2018-12-15", FromId:"DF-009", Validity:45},
		MachinInfo{Id:"JGC-011", Name:"伊利", Date:"2018-12-11", FromId:"DF-010", Validity:45},
		MachinInfo{Id:"JGC-012", Name:"伊利", Date:"2018-12-12", FromId:"DF-011", Validity:45},
		MachinInfo{Id:"JGC-013", Name:"伊利", Date:"2018-12-13", FromId:"DF-012", Validity:45},
		MachinInfo{Id:"JGC-014", Name:"伊利", Date:"2018-12-14", FromId:"DF-013", Validity:45},
		MachinInfo{Id:"JGC-015", Name:"伊利", Date:"2018-12-15", FromId:"DF-014", Validity:45},
		MachinInfo{Id:"JGC-016", Name:"三元", Date:"2018-12-11", FromId:"DF-015", Validity:45},
		MachinInfo{Id:"JGC-017", Name:"三元", Date:"2018-12-12", FromId:"DF-016", Validity:45},
		MachinInfo{Id:"JGC-018", Name:"三元", Date:"2018-12-13", FromId:"DF-017", Validity:45},
		MachinInfo{Id:"JGC-019", Name:"三元", Date:"2018-12-14", FromId:"DF-018", Validity:45},
		MachinInfo{Id:"JGC-020", Name:"三元", Date:"2018-12-15", FromId:"DF-019", Validity:45},
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

func(t* Machining)Invoke(stub shim.ChaincodeStubInterface) peer.Response {
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
func(t* Machining)setvalue(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	error := stub.PutState(keyID, []byte(args[1]))
	if error != nil {
		return shim.Error("PutState fail...")
	}
	return shim.Success([]byte("SetValue Sucess !!!"))
}

func(t* Machining)query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	value, error := stub.GetState(keyID)
	if error != nil {
		return shim.Error("GetState fail...")
	}
	return shim.Success(value)
}

// 根据keyID查询历史记录
func(t* Machining)gethistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
	error := shim.Start(new(Machining))
	if error != nil {
		println("程序启动失败...")
		return
	}
	fmt.Println("程序启动成功...")
}

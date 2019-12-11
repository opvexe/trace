package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 定义结构体, 继承ChainCode接口
type Sale struct {
}

// 定义数据结构体
type SaleInfo struct {
	Id       string		// 当前交易ID
	FromId   string		// 和加工厂的交易ID
	Price    string		// 价格
	Name     string		// 当前销售商名字
}

// 定义数据结构体
type MachinInfo struct {
	Id       string		// 加工厂ID
	FromId   string		// 奶场来源ID
	Date     string		// 生产日期
	Name     string		// 名字
	Validity int		// 保质期
}

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
func (t* Sale)Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.init(stub)
}

func (t* Sale)init(stub shim.ChaincodeStubInterface) peer.Response {
	// 数据初始化
	infos := []SaleInfo{
		SaleInfo{Id:"XSS-001", Name:"淘宝", Price:"55", FromId:"JGC-001"},
		SaleInfo{Id:"XSS-002", Name:"淘宝", Price:"55", FromId:"JGC-002"},
		SaleInfo{Id:"XSS-003", Name:"淘宝", Price:"55", FromId:"JGC-003"},
		SaleInfo{Id:"XSS-004", Name:"淘宝", Price:"55", FromId:"JGC-004"},
		SaleInfo{Id:"XSS-005", Name:"天猫", Price:"66", FromId:"JGC-005"},
		SaleInfo{Id:"XSS-006", Name:"天猫", Price:"66", FromId:"JGC-006"},
		SaleInfo{Id:"XSS-007", Name:"天猫", Price:"66", FromId:"JGC-007"},
		SaleInfo{Id:"XSS-008", Name:"天猫", Price:"66", FromId:"JGC-008"},
		SaleInfo{Id:"XSS-009", Name:"京东", Price:"55", FromId:"JGC-009"},
		SaleInfo{Id:"XSS-010", Name:"京东", Price:"55", FromId:"JGC-010"},
		SaleInfo{Id:"XSS-011", Name:"京东", Price:"55", FromId:"JGC-011"},
		SaleInfo{Id:"XSS-012", Name:"京东", Price:"55", FromId:"JGC-012"},
		SaleInfo{Id:"XSS-013", Name:"家乐福", Price:"45", FromId:"JGC-013"},
		SaleInfo{Id:"XSS-014", Name:"家乐福", Price:"78", FromId:"JGC-014"},
		SaleInfo{Id:"XSS-015", Name:"家乐福", Price:"78", FromId:"JGC-015"},
		SaleInfo{Id:"XSS-016", Name:"家乐福", Price:"45", FromId:"JGC-016"},
		SaleInfo{Id:"XSS-017", Name:"沃尔玛", Price:"45", FromId:"JGC-017"},
		SaleInfo{Id:"XSS-018", Name:"沃尔玛", Price:"78", FromId:"JGC-018"},
		SaleInfo{Id:"XSS-019", Name:"沃尔玛", Price:"78", FromId:"JGC-019"},
		SaleInfo{Id:"XSS-020", Name:"沃尔玛", Price:"78", FromId:"JGC-020"},
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

func(t* Sale)Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()

	if funcName == "setvalue" {
		return t.setvalue(stub, args)
	}else if funcName == "query" {
		return t.query(stub, args)
	}else if funcName == "trace" {
		return t.trace(stub, args)
	}
	return shim.Success([]byte("invoke OK"))
}
func(t* Sale)setvalue(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	error := stub.PutState(keyID, []byte(args[1]))
	if error != nil {
		return shim.Error("PutState fail...")
	}
	return shim.Success([]byte("SetValue Sucess !!!"))
}

func(t* Sale)query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyID := args[0]
	value, error := stub.GetState(keyID)
	if error != nil {
		return shim.Error("GetState fail...")
	}
	return shim.Success(value)
}

// 根据keyID查询历史记录
func(t* Sale)trace(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	text := t.query(stub, args)
	if text.Status != shim.OK {
		return shim.Error("query error ......")
	}
	var result string
	var fromID string
	// 获取需要的信息
	var saleItem SaleInfo
	json.Unmarshal(text.Payload, &saleItem)
	fromID = saleItem.FromId
	result += fmt.Sprintf("销售商:%s, ID:%s, FromID:%s <--- ", saleItem.Name, saleItem.Id, saleItem.FromId)

	// 找加工厂
	myArgs := [][]byte{[]byte("query"), []byte(fromID)}
	response := stub.InvokeChaincode("machincc", myArgs, "tracechannel")
	if response.Status != shim.OK {
		return shim.Error("InvokeChaincode error ......"+ string(response.Payload))
	}
	var machItem MachinInfo
	json.Unmarshal(response.Payload, &machItem)
	fromID = machItem.FromId
	result += fmt.Sprintf("加工厂:%s, ID:%s, FromID:%s <--- ", machItem.Name, machItem.Id, machItem.FromId)

	// 搜索奶牛场信息
	myArgs = [][]byte{[]byte("query"), []byte(fromID)}
	response = stub.InvokeChaincode("dairycc", myArgs, "tracechannel")
	if response.Status != shim.OK {
		return shim.Error("InvokeChaincode error ......")
	}
	var farmIitem FarmInfo
	json.Unmarshal(response.Payload, &farmIitem)
	result += fmt.Sprintf("奶牛场:%s, 牛奶质量:%s 。 ", farmIitem.Name, farmIitem.Quality)

	// 数据格式化为json
	//jsonText, error := json.Marshal(myList)
	//if error != nil {
	//	return shim.Error("json.Marshal(myList) fail...")
	//}
	return shim.Success([]byte(result))
}

func main() {
	error := shim.Start(new(Sale))
	if error != nil {
		println("程序启动失败...")
		return
	}
	fmt.Println("程序启动成功...")
}

package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var loggertest = shim.NewLogger("SmartContract labhlf-<xx>-sc UNIT_TEST")

var chaincodeVersionTest = "2019-03-nn"

//This function call the init function of the smart contrat and check the returned value
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
	if string(res.Payload) != chaincodeVersionTest {
		fmt.Println("Init returned wrong payload, \nexpected: ", chaincodeVersionTest, " \nreturned: ", string(res.Payload))
		t.FailNow()
	}
}

// checkState This function checks the value of an attribute of the worldstate
func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		//fmt.Println("State", name, "failed to get value")
		loggertest.Error(stub.Name, " State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		//fmt.Println("State value", name, "was not", value, "as expected")
		loggertest.Error(stub.Name, " State value", name, " \nvalue: ", string(bytes), " \nexpected: ", value)
		t.FailNow()
	}
}

//function to test the init of labhlf
func Test_HhlfscInit(t *testing.T) {
	scc := new(Labhlfsc)
	stub := shim.NewMockStub("Test_labhlfsc_Init", scc)
	loggertest.Infof(" Starting Test_labhlfsc_Init...")

	checkInit(t, stub, [][]byte{[]byte("init1"), []byte("INITTEST")})
	loggertest.Infof("Test_labhlfsc_Init completed")
}

//function to test the transaction OrderShipment
func Test_Labhlfsc_OrderShipment(t *testing.T) {
	scc := new(Labhlfsc)
	stub := shim.NewMockStub("Test_Labhlfsc_OrderShipment", scc)
	loggertest.Infof(" Starting Test_Labhlfsc_OrderShipment...")

	var packageValue string
	packageValue = `{"packageID":"P1","description":"Package for product P001","destination": "Montpellier, FRANCE, 34006"}`
	//var txId string
	txID := "PACK_1"
	res := stub.MockInvoke(txID, [][]byte{[]byte("OrderShipment"), []byte(packageValue)})
	if res.Status != shim.OK {
		//fmt.Println(string(res.Message))
		fmt.Println("Invoke OrderShipment failed ", res.Message)
		t.FailNow()
	}

	// check returned value
	var testvalue string
	testvalue = `{"txid":"PACK_1","err":null}`
	// check returned payload
	if res.Payload == nil {
		fmt.Println("Invoke OrderShipment returned wrong payload, expected txid => nil returned !!!")
		t.FailNow()
	}
	if string(res.Payload) != testvalue {
		fmt.Println("Invoke OrderShipment returned wrong payload, \nexpected: ", testvalue, " \nreturned: ", string(res.Payload))
		t.FailNow()
	}

	// check stored value in worldstate
	testvalue = `{"packageId":"P1","description":"Package for product P001","status":"READY","destination":"Montpellier, FRANCE, 34006"}`
	checkState(t, stub, "P1", testvalue)

	loggertest.Infof(" Test_Labhlfsc_OrderShipment completed")

}

//function to test the transaction OrderShipment
func Test_Labhlfsc_Ship(t *testing.T) {
	scc := new(Labhlfsc)
	stub := shim.NewMockStub("Test_Labhlfsc_Ship", scc)
	loggertest.Infof(" Starting Test_Labhlfsc_Ship...")

	var packageValue string
	packageValue = `{"packageID":"P1","description":"Package for product P001","destination": "Montpellier, FRANCE, 34006"}`
	//var txId string
	txID := "PACK_1"
	res := stub.MockInvoke(txID, [][]byte{[]byte("OrderShipment"), []byte(packageValue)})
	if res.Status != shim.OK {
		//fmt.Println(string(res.Message))
		fmt.Println("Invoke OrderShipment for Ship function has failed ", res.Message)
		t.FailNow()
	}

	var status, packid string
	packid = "P1"
	status = "SHIPMENT"
	//var txId string
	txID = "PACK_2"
	res = stub.MockInvoke(txID, [][]byte{[]byte("Ship"), []byte(packid), []byte(status)})
	if res.Status != shim.OK {
		//fmt.Println(string(res.Message))
		fmt.Println("Invoke Ship failed ", res.Message)
		t.FailNow()
	}

	// check returned value
	var testvalue string
	testvalue = `{"txid":"PACK_2","err":null}`
	// check returned payload
	if res.Payload == nil {
		fmt.Println("Invoke Ship returned wrong payload, expected txid => nil returned !!!")
		t.FailNow()
	}
	if string(res.Payload) != testvalue {
		fmt.Println("Invoke Ship returned wrong payload, \nexpected: ", testvalue, " \nreturned: ", string(res.Payload))
		t.FailNow()
	}

	// check stored value in worldstate
	testvalue = `{"packageId":"P1","description":"Package for product P001","status":"SHIPMENT","destination":"Montpellier, FRANCE, 34006"}`
	checkState(t, stub, "P1", testvalue)

	loggertest.Infof(" Test_Labhlfsc_Ship completed")

}

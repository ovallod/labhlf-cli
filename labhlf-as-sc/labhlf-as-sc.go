package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Labhlfsc Describes the structure of the labhlf
type Labhlfsc struct {
}

// Package Describes the package that will be shipped to the destination
type Package struct {
	PackageID   string `json:"packageId"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Destination string `json:"destination"`
}

var chaincodeVersion = "2019-03-nn"
var logger = shim.NewLogger("SmartContract Labhlfsc")

func (t *Labhlfsc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte(chaincodeVersion))
}

func (t *Labhlfsc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var args = stub.GetArgs()
	var err error

	logger.Info("Function: ", string(args[0]))

	// test the first argument which correspond to the transaction name
	switch string(args[0]) {
	case "OrderShipment":
		// Create a Package asset
		logger.Info(" OrderShipment function, Value:  ", string(args[1]))

		// Check input format
		var pack Package
		err = json.Unmarshal(args[1], &pack)
		if err != nil {
			fmt.Println(err)
			jsonResp := "Failed to Unmarshal package input data"
			return shim.Error(jsonResp)
		}
		// verify mandatory fields
		missingFields := false
		missingFieldsList := "Missing or empty attributes: "
		if pack.PackageID == "" {
			missingFieldsList += "PackageID, "
			missingFields = true
		}
		if pack.Destination == "" {
			missingFieldsList += "Destination, "
			missingFields = true
		}
		// ..... Complete the fields checking
		if missingFields {
			fmt.Println(missingFieldsList)
			return shim.Error(missingFieldsList)
		}
		pack.Status = "READY"

		var packageBytes []byte
		packageBytes, err = json.Marshal(pack)
		if err != nil {
			fmt.Println(err)
			return shim.Error("Failed to marshal Package object")
		}

		err = stub.PutState(pack.PackageID, packageBytes)
		if err != nil {
			fmt.Println(err)
			return shim.Error("Failed to write packageBytes data")
		}

		logger.Info(" OrderShipment function, txid: ", stub.GetTxID())
		return shim.Success([]byte(`{"txid":"` + stub.GetTxID() + `","err":null}`))

	case "Ship":

		// Change the status of the package
		// parameters : packageId, status
		logger.Info(" Ship function, Package id:  ", string(args[1]), " - Status ", string(args[2]))

		// Retrieve the parameters
		var packageID, status string
		var packageBytes []byte
		var pack Package
		packageID = string(args[1])
		status = string(args[2])

		// Test the value of the status : it should be SHIPMENT, SHIPPED, DELIVERED
		// ... TBC .... Implement the test

		// Get the Package in the ledger based on the PackageId
		packagebytes, err := stub.GetState(packageID)
		if err != nil {
			return shim.Error("Failed to get package " + packageID)
		}

		err = json.Unmarshal(packagebytes, &pack)
		if err != nil {
			fmt.Println(err)
			jsonResp := "Failed to Unmarshal package data " + packageID
			return shim.Error(jsonResp)
		}
		// Test the value of the current status:
		// if the new status is SHIPMENT, the current one should be READY
		// if the new status is SHIPPED, the current one should be SHIPMENT
		// if the new status is DELIVERED, the current one should be SHIPPED
		// ... TBC .... Implement the test
		if (status == "SHIPMENT" && pack.Status != "READY") ||
			(status == "SHIPPED" && pack.Status != "SHIPMENT") ||
			(status == "DELIVERED" && pack.Status != "SHIPPED") {
			jsonResp := "Bad status : " + status + " compared to the current one : " + pack.Status
			return shim.Error(jsonResp)

		}

		//Update the status
		pack.Status = status
		packageBytes, err = json.Marshal(pack)
		if err != nil {
			fmt.Println(err)
			return shim.Error("Failed to marshal Package object")
		}
		// Update the package in the ledger
		err = stub.PutState(pack.PackageID, packageBytes)
		if err != nil {
			fmt.Println(err)
			return shim.Error("Failed to write packageBytes data")
		}
		// Submit an event to inform about the status change
		err = stub.SetEvent("Shipp", packageBytes)
		if err != nil {
			fmt.Println(err)
			return shim.Error("Failed to raise enregistrePalette event!!!")
		}
		logger.Info(" shipp function, txid: ", stub.GetTxID())
		return shim.Success([]byte(`{"txid":"` + stub.GetTxID() + `","err":null}`))

	case "GetPackageStatus":

		// Get the status of the package
		// parameters : packageId
		logger.Info(" GetPackageStatus function, Package id:  ", string(args[1]))

		// Retrieve the parameters
		var packageID string
		var pack Package
		packageID = string(args[1])

		// Get the Package in the ledger based on the PackageId
		packagebytes, err := stub.GetState(packageID)
		if err != nil {
			return shim.Error("Failed to get package " + packageID)
		}

		err = json.Unmarshal(packagebytes, &pack)
		if err != nil {
			fmt.Println(err)
			jsonResp := "Failed to Unmarshal package data " + packageID
			return shim.Error(jsonResp)
		}

		return shim.Success([]byte(`{"PackageID":"` + pack.PackageID + `","Status":"` + pack.Status + `"}`))

	case "Acknowledgement":

		return shim.Success(nil)
	default:
		//The transaction name is not known
		return shim.Error("unkwnon function")

	}
}

func main() {
	err := shim.Start(new(Labhlfsc))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

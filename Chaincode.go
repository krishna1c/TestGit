package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type CBTC_Chaincode struct {
}

func (t *CBTC_Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	return CreateDatabase(stub, args)

}

func (t *CBTC_Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "InsertContractDetails" {
		return InsertContractDetails(stub, args)
	} else if function == "UpdateContractStatus" {
		return UpdateContractStatus(stub, args)
	} else if function == "InsertUniqueIdCountryRelation" {
		return InsertUserDetails(stub, args[0], args[1], args[2])
	} else if function == "UpdateOrderShipment" {
		return UpdateOrderShipment(stub, args)
	} else if function == "UploadDocument" {
		return UploadDocument(stub, args)
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *CBTC_Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetContractsDetails" {
		UniqueId := args[0]
		Country := args[1]
		return GetContractsDetails(stub, UniqueId, Country)
	} else if function == "GetDocument" {
		return GetDocument(stub, args)
	}
	return nil,nil
}

func main() {
	shim.Start(new(CBTC_Chaincode))
}

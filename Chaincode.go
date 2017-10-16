package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type KYC_Chaincode struct {
}

func (t *KYC_Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	return CreateDatabase(stub, args)

}

func (t *KYC_Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "InsertUserDetails" {
		return InsertUserDetails(stub, args)
	} else if function == "InsertPartyDetails" {
		return InsertPartyDetails(stub, args)
	} else if function == "InsertUserDocuments" {
		return InsertUserDocuments(stub, args)
	} else if function == "InsertPartyDocuments" {
		return InsertPartyDocuments(stub, args)
	} else if function == "UpdateDocumentStatus" {
		return UpdateDocumentStatus(stub, args)
	} else if function == "UpdateDocument" {
		return UpdateDocument(stub, args)
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KYC_Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "GetUserDocuments" {
		return GetUserDocuments(stub, args)
	} else if function == "GetDocument" {
		return GetDocument(stub, args)
	} else if function == "GetUserDetails" {
		return GetUserDetails(stub, args)
	} else if function == "GetPartyDetails" {
		return GetPartyDetails(stub, args)
	} else if function == "GetListOfBanks" {
		return GetListOfBanks(stub, args)
	} else if function == "GetListOfDocuments" {
		return GetListOfDocuments(stub, args)
	} else if function == "GetListOfUserIds" {
		return GetListOfUserIds(stub, args)
	} else if function == "GetListOfUsers" {
		return GetListOfUsers(stub, args)
	} else if function == "GetUserConsortium" {
		return GetUserConsortium(stub, args)
	} else if function == "GetBankKey" {
		return GetBankKey(stub, args)
	}
	return nil, nil
}

func main() {
	shim.Start(new(KYC_Chaincode))
}

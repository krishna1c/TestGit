package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Region Chaincode implementation
type KycChaincode struct {
}


type KycData struct {
	USER_ID    string `json:"USER_ID"`
	USER_NAME    string `json:"USER_NAME"`
	USER_KYC_PDF string `json:"USER_KYC_PDF"`
	KYC_DATE string `json:"KYC_DATE"`
	EXPIRY_DATE string `json:"EXPIRY_DATE"`
	BANK_NAME string `json:"BANK_NAME"`
}

func (t *KycChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Add user KYC data in Blockchain
func (t *KycChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "write" {
		return t.RegisterKYC(stub, args)
	}
	return nil, nil
}

func (t *KycChaincode) RegisterKYC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var KycDataObj KycData
	var err error
	var userId string

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Need 6 arguments")
	}

	userId = args[0]
	KycDataObj.USER_ID=args[0]
	KycDataObj.USER_NAME = args[1]
	KycDataObj.USER_KYC_PDF = args[2]
	KycDataObj.KYC_DATE= args[3]
	KycDataObj.EXPIRY_DATE= args[4]
	KycDataObj.BANK_NAME= args[5]

	fmt.Printf("Input from user:%s\n", KycDataObj)

	jsonAsBytes, _ := json.Marshal(KycDataObj)

	err = stub.PutState(userId, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var resAsBytes []byte
	var userId string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	userId = args[0]
	resAsBytes, err = t.GetKycDetails(stub, userId)
	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *KycChaincode) GetKycDetails(stub shim.ChaincodeStubInterface, userId string) ([]byte, error) {

	
	KycTxAsBytes, err := stub.GetState(userId)
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	return KycTxAsBytes, nil

}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

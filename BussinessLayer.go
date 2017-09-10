package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const dateFormat string = "2006-01-02"

func CreateDatabase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//Create table "ContractDetails"
	err = stub.CreateTable("ContractDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ContractId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "CreateDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ContractDueDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ContractStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "AdditionalInf", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Invoice", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PackagingList", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "LineOfCredit", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "BillOfShipment", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating ContractDetails table.")
	}

	//Create table "OrderDetails"
	err = stub.CreateTable("OrderDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "OrderId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Articles", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Buyer", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Seller", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Shipment", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Amount", Type: shim.ColumnDefinition_BYTES, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating OrderDetails table.")
	}

	//Create table "Document"
	err = stub.CreateTable("Document", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Data", Type: shim.ColumnDefinition_STRING, Key: false},
	})

	err = stub.CreateTable("UserContracts", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UniqueId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Country", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractId", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating UserContract table.")
	}
	return nil, nil

}

func InsertContractDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var ContractDetails Contract
	var err error
	var ok bool
	
	json.Unmarshal([]byte(args[0]), &ContractDetails)
	
	ok, err = stub.InsertRow("ContractDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractId}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.CreateDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractDueDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractStatus}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.AdditionalInf}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.Invoice}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.PackagingList}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.LineOfCredit}},			
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.BillOfShipment}},						
		},
	})
	if !ok || err != nil {
		return nil, errors.New("Error in adding ContractDetails record.")
	}

	//Order Details
	ArticlesBytes, _ := json.Marshal(ContractDetails.Order.Articles)
	BuyerBytes, _ := json.Marshal(ContractDetails.Order.Buyer)
	SellerBytes, _ := json.Marshal(ContractDetails.Order.Seller)
	ShipmentBytes, _ := json.Marshal(ContractDetails.Order.Shipment)
	AmountBytes, _ := json.Marshal(ContractDetails.Order.Amount)

	ok, err = stub.InsertRow("OrderDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractId}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: ArticlesBytes}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: BuyerBytes}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: SellerBytes}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: ShipmentBytes}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: AmountBytes}},
		},
	})
	if !ok || err != nil {
		return nil, errors.New("Error in adding OrderDetails record.")
	}

	
	ok, err = stub.InsertRow("Document", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.PackagingList}},
				&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
			},
		})

	if err != nil {
			return nil, errors.New("Error in Adding User Documents.")
		}

	//var seller Seller

	var BuyerDetails User
	BuyerDetails = ContractDetails.Order.Buyer
	InsertUserDetails(stub, BuyerDetails.UniqueId, BuyerDetails.Country, ContractDetails.ContractId)

	var SellerDetails User
	SellerDetails = ContractDetails.Order.Seller
	InsertUserDetails(stub, SellerDetails.UniqueId, SellerDetails.Country, ContractDetails.ContractId)

	return nil, nil
}

func InsertUserDetails(stub shim.ChaincodeStubInterface, UniqueId string,Country string, ContractId string) ([]byte, error) {

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: UniqueId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: Country}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	var ContractIds []interface{}
	ContractIds = append(ContractIds, ContractId)
	CIdBytes, _ := json.Marshal(ContractIds)
	ok, err := stub.InsertRow("UserContracts", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UniqueId}},
			&shim.Column{Value: &shim.Column_String_{String_: Country}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: CIdBytes}},
		},
	})

	if err != nil {
		return nil, errors.New("Error in Adding User Contracts.")
	}
	if !ok && err == nil {
		row, er := stub.GetRow("UserContracts", columns)
		if er != nil {
			return nil, errors.New("Error in Adding User Contracts.")
		}
		json.Unmarshal(row.Columns[2].GetBytes(), &ContractIds)
		ContractIds = append(ContractIds, ContractId)
		CIdBytes, _ = json.Marshal(ContractIds)
		ok, err = stub.ReplaceRow("UserContracts", shim.Row{

			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: UniqueId}},
				&shim.Column{Value: &shim.Column_String_{String_: Country}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: CIdBytes}},
			},
		})

		if !ok || err != nil {
			return nil, errors.New("Error in getting User Details.")
		}
	}
	return nil, nil

}
func GetContractDetails(stub shim.ChaincodeStubInterface, cid string) (Contract, error) {
	var ContractDetails Contract
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: cid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ContractDetails", columns)
	if err != nil {
		return ContractDetails, errors.New("Error in getting Contract Details.")
	}

	ContractDetails.ContractId = row.Columns[0].GetString_()
	ContractDetails.CreateDate = row.Columns[1].GetString_()
	ContractDetails.ContractDueDate = row.Columns[2].GetString_()
	ContractDetails.ContractStatus = row.Columns[3].GetString_()
	ContractDetails.AdditionalInf = row.Columns[4].GetString_()
	ContractDetails.Invoice = row.Columns[5].GetString_()
	ContractDetails.PackagingList = row.Columns[6].GetString_()
	ContractDetails.LineOfCredit = row.Columns[7].GetString_()
	ContractDetails.BillOfShipment = row.Columns[8].GetString_()
	return ContractDetails, nil

}

func GetDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var columns []shim.Column
	var Document Document

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Document", columns)
	if err != nil {
		return nil, errors.New("Error in getting Document.")
	}
	
	Document.ID = row.Columns[0].GetString_()
	Document.Data = row.Columns[1].GetString_()
	DocumentBytes, _ := json.Marshal(Document)
	return DocumentBytes, nil

}

func GetOrderDetails(stub shim.ChaincodeStubInterface, cid string, ContractDetails Contract) (Contract, error) {
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: cid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("OrderDetails", columns)
	if err != nil {
		return ContractDetails, errors.New("Error in getting Order Details.")
	}
	var Order Order

	json.Unmarshal(row.Columns[1].GetBytes(), &Order.Articles)
	json.Unmarshal(row.Columns[2].GetBytes(), &Order.Buyer)
	json.Unmarshal(row.Columns[3].GetBytes(), &Order.Seller)
	json.Unmarshal(row.Columns[4].GetBytes(), &Order.Shipment)
	json.Unmarshal(row.Columns[5].GetBytes(), &Order.Amount)
	ContractDetails.Order = Order

	return ContractDetails, nil
}

func GetContractsDetails(stub shim.ChaincodeStubInterface, UniqueId string, Country string) ([]byte, error) {
	var Contracts []Contract
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: UniqueId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: Country}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("UserContracts", columns)

	if err != nil {
		return nil, errors.New("Error in Getting Contract Details.")
	}
	var CIds []interface{}
	
	if( row.Columns == nil || len(row.Columns) != 3){
		return nil, errors.New("Error in Getting Contract Details.")
	}
	json.Unmarshal(row.Columns[2].GetBytes(), &CIds)

	var ContractDetails Contract
	for _, Id := range CIds {
		ContractDetails, _ = GetContractDetails(stub, Id.(string))
		ContractDetails, _ = GetOrderDetails(stub, Id.(string), ContractDetails)
		Contracts = append(Contracts, ContractDetails)
	}
	responseBytes, _ := json.Marshal(Contracts)
	return responseBytes, nil
}

func UpdateContractStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)
	row, err := stub.GetRow("ContractDetails", columns)
	if err != nil {
		return nil, errors.New("Error in Adding User Contracts.")
	}
	row.Columns[3] = &shim.Column{Value: &shim.Column_String_{String_: args[1]}}
	ok, er := stub.ReplaceRow("ContractDetails", row)
	if !ok || er != nil {
		return nil, errors.New("Error in Getting Contract Details.")
	}
	return nil, nil
}

func UpdateOrderShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)
	row, err := stub.GetRow("OrderDetails", columns)
	if err != nil {
		return nil, errors.New("Error in Adding User Contracts.")
	}
	
	
	row.Columns[4] = &shim.Column{Value:  &shim.Column_Bytes{Bytes: []byte(args[1])}}
	ok, er := stub.ReplaceRow("OrderDetails", row)
	if !ok || er != nil {
		return nil, errors.New("Error in Getting Contract Details.")
	}
	return nil, nil
}

func UploadDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	ok, err := stub.InsertRow("Document", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: args[3]}},
				&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			},
		})

	if !ok || err != nil {
			return nil, errors.New("Error in Adding User Documents.")
		}
	
	Number,e := strconv.Atoi(args[2])
	
	if e != nil {
			return nil, errors.New("Error in Adding User Documents.")
		}
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: args[1]}}
	columns = append(columns, col1)
	row, err1 := stub.GetRow("ContractDetails", columns)
	if err1 != nil {
		return nil, errors.New("Error in Adding User Contracts.")
	}
	row.Columns[Number] = &shim.Column{Value: &shim.Column_String_{String_: args[3]}}
	okk, er := stub.ReplaceRow("ContractDetails", row)
	if !okk || er != nil {
		return nil, errors.New("Error in Getting Contract Details.")
	}
	return nil, nil
}
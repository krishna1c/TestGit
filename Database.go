package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CreateDatabase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//Create table "UserDetails"
	err = stub.CreateTable("UserDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UniqueID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Phone", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Email", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating UserDetails table.")
	}

	//Create table "PartyDetails"
	err = stub.CreateTable("PartyDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UniqueID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Phone", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Email", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DocumentTypes", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Country", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating UserDetails table.")
	}

	//Create table "DocumentDetails"
	err = stub.CreateTable("DocumentDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Status", Type: shim.ColumnDefinition_STRING, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating DocumentDetails table.")
	}

	//Create table "Document"
	err = stub.CreateTable("Document", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Data", Type: shim.ColumnDefinition_STRING, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating Document table.")
	}

	//Create table "UserDocuments"
	err = stub.CreateTable("UserDocuments", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UniqueIDParty", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UniqueIDUser", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DocumentID", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "LastDate", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating UserDocuments table.")
	}

	LisOfBanksBytes, _ := json.Marshal("[]")
	stub.PutState("ListOfBanks", LisOfBanksBytes)

	ListOfDocumentTypesBytes, _ := json.Marshal("[]")
	stub.PutState("ListOfDocumentTypes", ListOfDocumentTypesBytes)

	return nil, nil
}

func InsertUserDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var User UserDetails
	var err error
	var ok bool

	json.Unmarshal([]byte(args[0]), &User)

	ok, err = stub.InsertRow("UserDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: User.UniqueID}},
			&shim.Column{Value: &shim.Column_String_{String_: User.Name}},
			&shim.Column{Value: &shim.Column_String_{String_: User.Phone}},
			&shim.Column{Value: &shim.Column_String_{String_: User.Email}},
		},
	})
	if !ok || err != nil {
		return nil, errors.New("Error in adding UserDetails record.")
	}
	arguments := []string{args[1], User.UniqueID, "new"}
	InsertUserDocuments(stub, arguments)
	return nil, nil
}

func InsertPartyDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var Party PartyDetails
	var err error
	var ok bool

	json.Unmarshal([]byte(args[0]), &Party)

	DocumentTypesBytes, _ := json.Marshal(Party.DocumentTypes)
	ok, err = stub.InsertRow("PartyDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: Party.UniqueID}},
			&shim.Column{Value: &shim.Column_String_{String_: Party.Name}},
			&shim.Column{Value: &shim.Column_String_{String_: Party.Phone}},
			&shim.Column{Value: &shim.Column_String_{String_: Party.Email}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: DocumentTypesBytes}},
			&shim.Column{Value: &shim.Column_String_{String_: Party.Country}},
		},
	})
	if !ok || err != nil {
		return nil, errors.New("Error in adding PartyDetails record.")
	}

	var ListOfBanks []PartysDetails
	ListOfBanksBytes, _ := stub.GetState("ListOfBanks")
	json.Unmarshal(ListOfBanksBytes, &ListOfBanks)
	var PartysDetails PartysDetails
	PartysDetails.Name = Party.Name
	PartysDetails.DocumentTypes = Party.DocumentTypes
	ListOfBanks = append(ListOfBanks, PartysDetails)
	ListOfBanksBytes, _ = json.Marshal(ListOfBanks)
	stub.PutState("ListOfBanks", ListOfBanksBytes)

	var ListOfDocumentTypes []string
	ListOfDocumentTypesBytes, _ := stub.GetState("ListOfDocumentTypes")
	json.Unmarshal(ListOfDocumentTypesBytes, &ListOfDocumentTypes)
	for _, Id := range Party.DocumentTypes {
		if !contains(ListOfDocumentTypes, Id) {
			ListOfDocumentTypes = append(ListOfDocumentTypes, Id)
		}
	}
	ListOfDocumentTypesBytes, _ = json.Marshal(ListOfDocumentTypes)
	stub.PutState("ListOfDocumentTypes", ListOfDocumentTypesBytes)

	var ListOfUsers []string
	ListOfUsersBytes, _ := json.Marshal(ListOfUsers)
	stub.PutState(Party.UniqueID, ListOfUsersBytes)
	stub.PutState(Party.Name+args[1], []byte(Party.UniqueID))
	return nil, nil
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func InsertUserDocuments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var DocumentDetail DocumentDetails
	var err error
	var ok bool
	var DocumentIDs []interface{}
	var columns []shim.Column

	if args[2] == "new" {
		DocumentIDBytes, _ := json.Marshal(DocumentIDs)
		ok, err = stub.InsertRow("UserDocuments", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
				&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: DocumentIDBytes}},
			},
		})
	} else {
		json.Unmarshal([]byte(args[2]), &DocumentDetail)

		col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: args[1]}}
		columns = append(columns, col1)
		columns = append(columns, col2)

		DocumentIDs = append(DocumentIDs, DocumentDetail.ID)
		DocumentIDBytes, _ := json.Marshal(DocumentIDs)
		ok, err = stub.InsertRow("UserDocuments", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
				&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: DocumentIDBytes}},
				&shim.Column{Value: &shim.Column_String_{String_: args[3]}},
			},
		})

		if err != nil {
			return nil, errors.New("Error in Adding User Documents.")
		}
		if !ok && err == nil {
			row, er := stub.GetRow("UserDocuments", columns)
			if er != nil {
				return nil, errors.New("Error in Adding User Documents.")
			}
			json.Unmarshal(row.Columns[2].GetBytes(), &DocumentIDs)
			DocumentIDs = append(DocumentIDs, DocumentDetail.ID)
			DocumentIDBytes, _ = json.Marshal(DocumentIDs)
			ok, err = stub.ReplaceRow("UserDocuments", shim.Row{

				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
					&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: DocumentIDBytes}},
					&shim.Column{Value: &shim.Column_String_{String_: args[3]}},
				},
			})

			if !ok || err != nil {
				return nil, errors.New("Error in getting User Documents.")
			}
		}
		ok, err = stub.InsertRow("DocumentDetails", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: DocumentDetail.ID}},
				&shim.Column{Value: &shim.Column_String_{String_: DocumentDetail.Type}},
				&shim.Column{Value: &shim.Column_String_{String_: DocumentDetail.Status}},
			},
		})

		if err != nil {
			return nil, errors.New("Error in Adding User Documents.")
		}
		ok, err = stub.InsertRow("Document", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: DocumentDetail.ID}},
				&shim.Column{Value: &shim.Column_String_{String_: args[3]}},
			},
		})

		if err != nil {
			return nil, errors.New("Error in Adding User Documents.")
		}

	}
	return nil, nil

}
func InsertPartyDocuments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var DocumentIds []interface{}
	json.Unmarshal([]byte(args[2]), &DocumentIds)
	DocumentIDBytes, _ := json.Marshal(DocumentIds)
	ok, err := stub.InsertRow("UserDocuments", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: args[0]}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: DocumentIDBytes}},
			&shim.Column{Value: &shim.Column_String_{String_: args[3]}},
		},
	})

	if !ok || err != nil {
		return nil, errors.New("Error in Adding Party Documents.")
	}

	var ListOfUsers []string
	ListOfUsersBytes, _ := stub.GetState(args[0])
	json.Unmarshal(ListOfUsersBytes, &ListOfUsers)
	ListOfUsers = append(ListOfUsers, args[1])
	ListOfUsersBytes, _ = json.Marshal(ListOfUsers)
	stub.PutState(args[0], ListOfUsersBytes)

	return nil, nil
}
func UpdateDocumentStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var Document DocumentDetails
	json.Unmarshal([]byte(args[0]), &Document)
	ok, err := stub.ReplaceRow("DocumentDetails", shim.Row{

		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: Document.ID}},
			&shim.Column{Value: &shim.Column_String_{String_: Document.Type}},
			&shim.Column{Value: &shim.Column_String_{String_: Document.Status}},
		},
	})

	if !ok || err != nil {
		return nil, errors.New("Error in Updating DocumentDetails.")
	}
	return nil, nil
}
func UpdateDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var Document Document
	var DocumentDetail DocumentDetails
	json.Unmarshal([]byte(args[0]), &Document)
	ok, err := stub.ReplaceRow("Document", shim.Row{

		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: Document.ID}},
			&shim.Column{Value: &shim.Column_String_{String_: Document.Data}},
		},
	})

	if !ok || err != nil {
		return nil, errors.New("Error in Updating Document.")
	}
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: Document.ID}}
	columns = append(columns, col1)

	row, er := stub.GetRow("DocumentDetails", columns)

	if er != nil {
		return nil, errors.New("Error in getting Document Details.")
	}

	DocumentDetail.ID = row.Columns[0].GetString_()
	DocumentDetail.Type = row.Columns[1].GetString_()
	DocumentDetail.Status = row.Columns[2].GetString_()

	ok, err = stub.ReplaceRow("DocumentDetails", shim.Row{

		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: DocumentDetail.ID}},
			&shim.Column{Value: &shim.Column_String_{String_: DocumentDetail.Type}},
			&shim.Column{Value: &shim.Column_String_{String_: args[1]}},
		},
	})
	if !ok || err != nil {
		return nil, errors.New("Error in Updating Document.")
	}

	return nil, nil
}

//Get Functions
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
func GetDocumentDetails(stub shim.ChaincodeStubInterface, DocumentID string) (DocumentDetails, error) {
	var columns []shim.Column
	var Document DocumentDetails

	col1 := shim.Column{Value: &shim.Column_String_{String_: DocumentID}}
	columns = append(columns, col1)

	row, err := stub.GetRow("DocumentDetails", columns)

	if err != nil {
		return Document, errors.New("Error in getting Document Details.")
	}

	Document.ID = row.Columns[0].GetString_()
	Document.Type = row.Columns[1].GetString_()
	Document.Status = row.Columns[2].GetString_()
	return Document, nil
}
func GetUserDocuments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var UserDocument UserDocuments
	var Documents []DocumentDetails
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: args[1]}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("UserDocuments", columns)

	if err != nil {
		return nil, errors.New("Error in Getting Contract Details.")
	}
	var DocumentIds []interface{}
	json.Unmarshal(row.Columns[2].GetBytes(), &DocumentIds)

	var DocumentDetails DocumentDetails
	for _, Id := range DocumentIds {
		DocumentDetails, _ = GetDocumentDetails(stub, Id.(string))
		Documents = append(Documents, DocumentDetails)
	}
	UserDocument.UniqueIDParty = args[0]
	UserDocument.UniqueIDUser = args[1]
	UserDocument.DocumentDetails = Documents
	UserDocument.LastDate = row.Columns[3].GetString_();

	responseBytes, _ := json.Marshal(UserDocument)
	return responseBytes, nil
}

func GetUserDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var columns []shim.Column
	var User UserDetails

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)

	row, err := stub.GetRow("UserDetails", columns)
	if err != nil {
		return nil, errors.New("Error in getting UserDetails.")
	}

	User.UniqueID = row.Columns[0].GetString_()
	User.Name = row.Columns[1].GetString_()
	User.Phone = row.Columns[2].GetString_()
	User.Email = row.Columns[3].GetString_()
	UserBytes, _ := json.Marshal(User)
	return UserBytes, nil
}

func GetPartyDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var columns []shim.Column
	var Party PartyDetails

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)

	row, err := stub.GetRow("PartyDetails", columns)
	if err != nil {
		return nil, errors.New("Error in getting Party Details.")
	}

	Party.UniqueID = row.Columns[0].GetString_()
	Party.Name = row.Columns[1].GetString_()
	Party.Phone = row.Columns[2].GetString_()
	Party.Email = row.Columns[3].GetString_()
	Party.Email = row.Columns[5].GetString_()

	var DocumentTypes []interface{}
	json.Unmarshal([]byte(row.Columns[4].GetString_()), &DocumentTypes)

	var DocumentTypesStrings []string
	for _, Id := range DocumentTypes {
		DocumentTypesStrings = append(DocumentTypesStrings, Id.(string))
	}
	Party.DocumentTypes = DocumentTypesStrings
	PartyBytes, _ := json.Marshal(Party)
	return PartyBytes, nil
}
func GetListOfBanks(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	ListOfBanksBytes, _ := stub.GetState("ListOfBanks")
	return ListOfBanksBytes, nil
}
func GetListOfDocuments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	ListOfDocumentTypesBytes, _ := stub.GetState("ListOfDocumentTypes")
	return ListOfDocumentTypesBytes, nil
}
func GetListOfUserIds(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	ListOfUsers, _ := stub.GetState(args[0])
	return ListOfUsers, nil
}
func GetListOfUsers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var UserPartys []UserParty
	ListOfUsersBytes, _ := stub.GetState(args[0])
	var ListOfUsers []string
	json.Unmarshal(ListOfUsersBytes, &ListOfUsers)

	var UserParty UserParty
	var UserDetailsBytes []byte
	var User UserDetails
	var UserDocuments UserDocuments
	var UserDocumentsBytes []byte
	for _, Id := range ListOfUsers {
		arguments := []string{Id}
		UserDetailsBytes, _ = GetUserDetails(stub, arguments)
		json.Unmarshal(UserDetailsBytes, &User)
		arguments = []string{args[0], Id}
		UserDocumentsBytes, _ = GetUserDocuments(stub, arguments)
		json.Unmarshal(UserDocumentsBytes, &UserDocuments)
		UserParty.UserDetails = User
		UserParty.DocumentDetails = UserDocuments.DocumentDetails
		UserPartys = append(UserPartys, UserParty)
	}
	responseBytes, _ := json.Marshal(UserPartys)
	return responseBytes, nil
}
func GetUserConsortium(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var UserConsortium UserConsortium
	var UserDetailsBytes []byte
	var User UserDetails
	var UserDocuments UserDocuments
	var UserDocumentsBytes []byte
	arguments := []string{args[1]}
	UserDetailsBytes, _ = GetUserDetails(stub, arguments)
	json.Unmarshal(UserDetailsBytes, &User)
	arguments = []string{args[0], args[1]}
	UserDocumentsBytes, _ = GetUserDocuments(stub, arguments)
	json.Unmarshal(UserDocumentsBytes, &UserDocuments)
	UserConsortium.UserDetails = User
	UserConsortium.DocumentDetails = UserDocuments.DocumentDetails
	var Documents []string
	var PartysDetails []PartysDetails
	GetListOfBanksBytes, _ := GetListOfBanks(stub, arguments)
	json.Unmarshal(GetListOfBanksBytes, &PartysDetails)
	GetListOfDocumentsBytes, _ := GetListOfDocuments(stub, arguments)
	json.Unmarshal(GetListOfDocumentsBytes, &Documents)
	UserConsortium.AllDocuments = Documents
	UserConsortium.PartysDetails = PartysDetails
	responseBytes, _ := json.Marshal(UserConsortium)
	return responseBytes, nil
}
func GetBankKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	responseBytes, _ := stub.GetState(args[0])
	return responseBytes, nil
}

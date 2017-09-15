package main

type UserDetails struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	UniqueID string `json:"uniqueId"`
	Password string `json:"password"`
}

type UserDocuments struct {
	UniqueIDParty   string            `json:"uniqueIdParty"`
	UniqueIDUser    string            `json:"uniqueIdUser"`
	DocumentDetails []DocumentDetails `json:"documentDetails"`
	LastDate		string			  `json:"lastDate"`				
}

type DocumentDetails struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type Document struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type PartyDetails struct {
	UniqueID      string   `json:"uniqueId"`
	DocumentTypes []string `json:"documentTypes"`
	Name          string   `json:"name"`
	Password      string   `json:"password"`
	Phone         string   `json:"phone"`
	Email         string   `json:"email"`
	Country       string   `json:"country"`
}

type PartysDetails struct {
	DocumentTypes []string `json:"documentTypes"`
	Name          string   `json:"name"`
	Country       string   `json:"country"`
}

type UserParty struct {
	UserDetails     UserDetails       `json:"userDetails"`
	DocumentDetails []DocumentDetails `json:"documentDetails"`
	LastDate		string			  `json:"lastDate"`
}

type UserConsortium struct {
	UserDetails     UserDetails       `json:"userDetails"`
	DocumentDetails []DocumentDetails `json:"documentDetails"`
	AllDocuments    []string          `json:"allDocuments"`
	PartysDetails   []PartysDetails   `json:"partysDetails"`
}

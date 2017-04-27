/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var myLogger = logging.MustGetLogger("asset_mgm")

//     Eric: 1 6avZQLwcUe9b bank_a
//     Jessica: 1 6avZQLwcUe9b bank_a
//     Simon: 1 6avZQLwcUe9b bank_a
//     Oliver: 1 6avZQLwcUe9b bank_a
//     Mary: 1 6avZQLwcUe9b bank_a
//     Tomas: 1 6avZQLwcUe9b bank_a
//     Cassie: 1 6avZQLwcUe9b bank_a
//   attribute-entry-12: Eric;bank_a;account;Eric;2015-02-02T00:00:00-03:00;;
//   attribute-entry-13: Eric;bank_a;role;Controller;2015-02-02T00:00:00-03:00;;
//   attribute-entry-14: Eric;bank_a;organization;Company1;2015-02-02T00:00:00-03:00;;
//   attribute-entry-15: Jessica;bank_a;account;Jessica;2015-02-02T00:00:00-03:00;;
//   attribute-entry-16: Jessica;bank_a;role;Buyer;2015-02-02T00:00:00-03:00;;
//   attribute-entry-17: Jessica;bank_a;organization;Company1;2015-02-02T00:00:00-03:00;;
//   attribute-entry-18: Simon;bank_a;account;Simon;2015-02-02T00:00:00-03:00;;
//   attribute-entry-19: Simon;bank_a;role;Buyer;2015-02-02T00:00:00-03:00;;
//   attribute-entry-20: Simon;bank_a;organization;Company1;2015-02-02T00:00:00-03:00;;
//   attribute-entry-21: Oliver;bank_a;account;Oliver;2015-02-02T00:00:00-03:00;;
//   attribute-entry-22: Oliver;bank_a;role;Controller;2015-02-02T00:00:00-03:00;;
//   attribute-entry-23: Oliver;bank_a;organization;Company2;2015-02-02T00:00:00-03:00;;
//   attribute-entry-24: Mary;bank_a;account;Mary;2015-02-02T00:00:00-03:00;;
//   attribute-entry-25: Mary;bank_a;role;Buyer;2015-02-02T00:00:00-03:00;;
//   attribute-entry-26: Mary;bank_a;organization;Company2;2015-02-02T00:00:00-03:00;;
//   attribute-entry-27: Tomas;bank_a;account;Tomas;2015-02-02T00:00:00-03:00;;
//   attribute-entry-28: Tomas;bank_a;role;Buyer;2015-02-02T00:00:00-03:00;;
//   attribute-entry-29: Tomas;bank_a;organization;Company2;2015-02-02T00:00:00-03:00;;
//   attribute-entry-30: Cassie;bank_a;account;Cassie;2015-02-02T00:00:00-03:00;;
//   attribute-entry-31: Cassie;bank_a;role;Regulator;2015-02-02T00:00:00-03:00;;
//   attribute-entry-32: Cassie;bank_a;organization;Company0;2015-02-02T00:00:00-03:00;;

// RiskManagementChaincode example simple Asset Management Chaincode implementation
type RiskManagementChaincode struct {
}

// RiskInfo ...
type RiskInfo struct {
	UUID            string `json:"uuid"`            //UUID
	SupplierID      string `json:"supplierId"`      //supplierId
	SupplierName    string `json:"supplierName"`    //supplierName
	AssessCompany   string `json:"assessCompany"`   //assessCompany
	RiskEvent       string `json:"riskEvent"`       //riskEvent
	RiskType        string `json:"riskType"`        //riskType
	RiskLevel       string `json:"riskLevel"`       //riskLevel
	RiskDate        string `json:"riskDate"`        //riskDate
	RiskLocation    string `json:"riskLocation"`    //riskLocation
	CreateBy        string `json:"createBy"`        //createBy
	CreateTimestamp string `json:"createTimestamp"` //createTimestamp
	UpdateBy        string `json:"updateBy"`        //updateBy
	UpdateTimestamp string `json:"updateTimestamp"` //updateTimestamp
	UsefulCount     string `json:"usefulCount"`     //UsefulCount
	UselessCount    string `json:"uselessCount"`    //UselessCount
}

// RiskHistory ...
type RiskHistory struct {
	UUID          string `json:"uuid"`          //UUID
	Timestamp     string `json:"timestamp"`     //Timestamp
	OperationType string `json:"operationType"` //OperationType
	Operation     string `json:"operation"`     //Operation
}

// AuthInfo ...
type AuthInfo struct {
	UUID        string `json:"uuid"`        //UUID
	UserName    string `json:"userName"`    //UserName
	UserCompany string `json:"userCompany"` //UserCompany
}

// Init initialization
func (t *RiskManagementChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.Info("[RiskManagementChaincode] Init")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	//Create risk information table
	err := stub.DeleteTable("RI")
	err = stub.CreateTable("RI", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "uuid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "supplierId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "supplierName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "assessCompany", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "riskEvent", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "riskType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "riskLevel", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "riskDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "riskLocation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "createBy", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "createTimestamp", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "updateBy", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "updateTimestamp", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "usefulCount", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "uselessCount", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed creating AssetsOwnership table, [%v]", err)
	}

	//Create authorization table
	err = stub.DeleteTable("RA")
	err = stub.CreateTable("RA", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "uuid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "userName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "userCompany", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "dummy", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed creating Risk Authorization table, [%v]", err)
	}

	// Create risk history table
	err = stub.DeleteTable("RH")
	err = stub.CreateTable("RH", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "uuid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "timestamp", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "operationType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "operation", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed creating Risk History table, [%v]", err)
	}

	// Create risk useful mark table
	err = stub.DeleteTable("RM")
	err = stub.CreateTable("RM", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "uuid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "userName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "userCompany", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "mark", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed creating risk useful mark table, [%v]", err)
	}
	fmt.Printf("Init succeed! \n")
	return nil, nil
}

func (t *RiskManagementChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Registering risk information...")

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}
	uuid := args[0]
	supplierID := args[1]
	supplierName := args[2]
	riskEvent := args[3]
	riskType := args[4]
	riskLevel := args[5]
	riskDate := args[6]
	riskLocation := args[7]

	creator, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}
	role, err := stub.ReadCertAttribute("role")
	if err != nil {
		fmt.Printf("Error reading attribute 'role' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller role. Error was [%v]", err)
	}
	if string(role) != "Controller" {
		fmt.Printf("Error don not have Controller role")
		return nil, fmt.Errorf("Error don not have Controller role")
	}
	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}
	// if len(creator) == 0 {
	// 	fmt.Printf("Account length is 0.")
	// 	return nil, errors.New("Account length is 0.")
	// }
	fmt.Printf("Register info. supplierId : [%s], supplierName : [%s], riskEvent : [%s]", supplierID, supplierName, riskEvent)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ok, err := stub.InsertRow("RI", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: uuid}},
			&shim.Column{Value: &shim.Column_String_{String_: supplierID}},
			&shim.Column{Value: &shim.Column_String_{String_: supplierName}},
			&shim.Column{Value: &shim.Column_String_{String_: string(organization)}},
			&shim.Column{Value: &shim.Column_String_{String_: riskEvent}},
			&shim.Column{Value: &shim.Column_String_{String_: riskType}},
			&shim.Column{Value: &shim.Column_String_{String_: riskLevel}},
			&shim.Column{Value: &shim.Column_String_{String_: riskDate}},
			&shim.Column{Value: &shim.Column_String_{String_: riskLocation}},
			&shim.Column{Value: &shim.Column_String_{String_: string(creator)}},
			&shim.Column{Value: &shim.Column_String_{String_: timestamp}},
			&shim.Column{Value: &shim.Column_String_{String_: string(creator)}},
			&shim.Column{Value: &shim.Column_String_{String_: timestamp}},
			&shim.Column{Value: &shim.Column_String_{String_: "0"}},
			&shim.Column{Value: &shim.Column_String_{String_: "0"}}},
	})
	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("Risk information was already published.")
	}
	stub.InsertRow("RA", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: uuid}},
			&shim.Column{Value: &shim.Column_String_{String_: string(creator)}},
			&shim.Column{Value: &shim.Column_String_{String_: string(organization)}},
			&shim.Column{Value: &shim.Column_String_{String_: "true"}}},
	})
	t.addHistory(stub, uuid, "Create", "")
	return []byte(supplierID), err
}

func (t *RiskManagementChaincode) update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Updating risk information...")

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	uuid := args[0]
	supplierID := args[1]
	supplierName := args[2]
	riskEvent := args[3]
	riskType := args[4]
	riskLevel := args[5]
	riskDate := args[6]
	riskLocation := args[7]

	creator, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}
	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}
	// At this point, the proof of ownership is valid, then register transfer

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("RI", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". Error " + err.Error() + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if string(creator) != row.Columns[9].GetString_() || string(organization) != row.Columns[3].GetString_() {
		jsonResp := "{\"Error\":\"No permission. \"}"
		return nil, errors.New(jsonResp)
	}
	err = stub.DeleteRow(
		"RI", columns,
	)
	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}
	//fmt.Printf("Update info. supplierId : [%s], supplierName : [%s], riskEvent : [%s]", supplierID, supplierName, riskEvent)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ok, err := stub.InsertRow("RI", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[0].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: supplierID}},
			&shim.Column{Value: &shim.Column_String_{String_: supplierName}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[3].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: riskEvent}},
			&shim.Column{Value: &shim.Column_String_{String_: riskType}},
			&shim.Column{Value: &shim.Column_String_{String_: riskLevel}},
			&shim.Column{Value: &shim.Column_String_{String_: riskDate}},
			&shim.Column{Value: &shim.Column_String_{String_: riskLocation}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[9].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[10].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: string(creator)}},
			&shim.Column{Value: &shim.Column_String_{String_: timestamp}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[13].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[14].GetString_()}}},
	})
	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("Risk information was already published.")
	}
	updateFields := ""
	if supplierID != row.Columns[1].GetString_() {
		updateFields = "[supplierID] "
	}
	if supplierName != row.Columns[2].GetString_() {
		updateFields += "[supplierName] "
	}
	if riskEvent != row.Columns[4].GetString_() {
		updateFields += "[riskEvent] "
	}
	if riskType != row.Columns[5].GetString_() {
		updateFields += "[riskType] "
	}
	if riskLevel != row.Columns[6].GetString_() {
		updateFields += "[riskLevel] "
	}
	if riskDate != row.Columns[7].GetString_() {
		updateFields += "[riskDate] "
	}
	if riskLocation != row.Columns[8].GetString_() {
		updateFields += "[riskLocation] "
	}

	t.addHistory(stub, uuid, "Update", updateFields)
	return []byte(supplierID), err
}

func (t *RiskManagementChaincode) authorize(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	uuid := args[0]
	toUserName := args[1]
	toUserCompany := args[2]

	account, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}
	role, err := stub.ReadCertAttribute("role")
	if err != nil {
		fmt.Printf("Error reading attribute 'role' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller role. Error was [%v]", err)
	}
	if string(role) != "Controller" {
		fmt.Printf("Error don not have Controller role")
		return nil, fmt.Errorf("Error don not have Controller role")
	}
	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("RI", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". Error " + err.Error() + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". \"}"
		return nil, errors.New(jsonResp)
	}
	hasRight := false

	// if string(account) == row.Columns[9].GetString_() && string(organization) == row.Columns[3].GetString_() {
	// 	hasRight = true
	// }
	if hasRight == false {
		var columnsAuth []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: string(account)}}
		col3 := shim.Column{Value: &shim.Column_String_{String_: string(organization)}}
		columnsAuth = append(columnsAuth, col1)
		columnsAuth = append(columnsAuth, col2)
		columnsAuth = append(columnsAuth, col3)
		rowAuth, err := stub.GetRow("RA", columnsAuth)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed retrieving risk authorize. Error " + err.Error() + ". \"}"
			return nil, errors.New(jsonResp)
		}
		if len(rowAuth.Columns) != 0 {
			hasRight = true
		}
	}
	if hasRight == true {
		ok, err := stub.InsertRow("RA", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: uuid}},
				&shim.Column{Value: &shim.Column_String_{String_: toUserName}},
				&shim.Column{Value: &shim.Column_String_{String_: toUserCompany}},
				&shim.Column{Value: &shim.Column_String_{String_: "true"}}},
		})
		if !ok && err == nil {
			fmt.Println("Error inserting row")
			return nil, errors.New("Risk information was already published.")
		}
	}
	t.addHistory(stub, uuid, "Authorize", string(account)+"->"+toUserName)
	return nil, nil
}

func (t *RiskManagementChaincode) thumbUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	uuid := args[0]
	mark := args[1]

	account, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}
	role, err := stub.ReadCertAttribute("role")
	if err != nil {
		fmt.Printf("Error reading attribute 'role' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller role. Error was [%v]", err)
	}
	if string(role) != "Buyer" {
		fmt.Printf("Error don not have Buyer role")
		return nil, fmt.Errorf("Error don not have Buyer role")
	}
	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}

	ok, err := stub.InsertRow("RM", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: uuid}},
			&shim.Column{Value: &shim.Column_String_{String_: string(account)}},
			&shim.Column{Value: &shim.Column_String_{String_: string(organization)}},
			&shim.Column{Value: &shim.Column_String_{String_: mark}}},
	})
	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("Risk useful info was already published.")
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("RI", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". Error " + err.Error() + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if string(account) == row.Columns[9].GetString_() && string(organization) == row.Columns[3].GetString_() {
		jsonResp := "{\"Error\":\"No permission. \"}"
		return nil, errors.New(jsonResp)
	}
	err = stub.DeleteRow(
		"RI", columns,
	)
	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}
	//fmt.Printf("Update info. supplierId : [%s], supplierName : [%s], riskEvent : [%s]", supplierID, supplierName, riskEvent)
	usefulCount := row.Columns[13].GetString_()
	uselessCount := row.Columns[14].GetString_()
	if mark == "true" {
		b, _ := strconv.Atoi(usefulCount)
		b++
		usefulCount = strconv.Itoa(b)
	} else {
		b, _ := strconv.Atoi(uselessCount)
		b++
		uselessCount = strconv.Itoa(b)
	}

	ok, err = stub.InsertRow("RI", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[0].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[1].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[2].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[3].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[4].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[5].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[6].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[7].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[8].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[9].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[10].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[11].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[12].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: usefulCount}},
			&shim.Column{Value: &shim.Column_String_{String_: uselessCount}}},
	})
	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("Risk information was already published.")
	}
	return nil, nil
}

func (t *RiskManagementChaincode) getAllRisks(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0 arguments.")
	}
	var columns []shim.Column
	rows, err := stub.GetRows("RI", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk information. \"}"
		return nil, errors.New(jsonResp)
	}

	var riskInfos []RiskInfo
	for row := range rows {
		if len(row.Columns) == 0 {
			jsonResp := "{\"Error\":\"Failed retrieving risk information. \"}"
			return nil, errors.New(jsonResp)
		}
		riskInfo := RiskInfo{UUID: row.Columns[0].GetString_(), SupplierID: row.Columns[1].GetString_(), SupplierName: row.Columns[2].GetString_(), AssessCompany: row.Columns[3].GetString_(), RiskType: row.Columns[5].GetString_(), UsefulCount: row.Columns[13].GetString_(), UselessCount: row.Columns[14].GetString_()}
		riskInfos = append(riskInfos, riskInfo)
	}
	riskInfosBytes, err := json.Marshal(&riskInfos)
	if err != nil {
		return nil, errors.New("Error getting risk information")
	}
	return riskInfosBytes, nil
}

func (t *RiskManagementChaincode) getAllRisksByOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0 arguments.")
	}

	account, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}

	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}

	var columns []shim.Column

	rows, err := stub.GetRows("RI", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk information. \"}"
		return nil, errors.New(jsonResp)
	}

	var riskInfos []RiskInfo
	for row := range rows {
		if len(row.Columns) == 0 {
			jsonResp := "{\"Error\":\"Failed retrieving risk information. \"}"
			return nil, errors.New(jsonResp)
		}
		hasRight := false
		// if string(account) == row.Columns[9].GetString_() && string(organization) == row.Columns[3].GetString_() {
		// 	hasRight = true
		// }
		if hasRight == false {
			var columnsAuth []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: row.Columns[0].GetString_()}}
			col2 := shim.Column{Value: &shim.Column_String_{String_: string(account)}}
			col3 := shim.Column{Value: &shim.Column_String_{String_: string(organization)}}
			columnsAuth = append(columnsAuth, col1)
			columnsAuth = append(columnsAuth, col2)
			columnsAuth = append(columnsAuth, col3)
			rowAuth, err := stub.GetRow("RA", columnsAuth)
			if err != nil {
				jsonResp := "{\"Error\":\"Failed retrieving risk authorize. Error " + err.Error() + ". \"}"
				return nil, errors.New(jsonResp)
			}
			if len(rowAuth.Columns) != 0 {
				hasRight = true
			}
		}
		if hasRight == true {
			riskInfo := RiskInfo{UUID: row.Columns[0].GetString_(), SupplierID: row.Columns[1].GetString_(), SupplierName: row.Columns[2].GetString_(), AssessCompany: row.Columns[3].GetString_(), RiskEvent: row.Columns[4].GetString_(), RiskType: row.Columns[5].GetString_(), RiskLevel: row.Columns[6].GetString_(), RiskDate: row.Columns[7].GetString_(), RiskLocation: row.Columns[8].GetString_(), CreateBy: row.Columns[9].GetString_(), CreateTimestamp: row.Columns[10].GetString_(), UpdateBy: row.Columns[11].GetString_(), UpdateTimestamp: row.Columns[12].GetString_(), UsefulCount: row.Columns[13].GetString_(), UselessCount: row.Columns[14].GetString_()}
			riskInfos = append(riskInfos, riskInfo)
		}
	}
	riskInfosBytes, err := json.Marshal(&riskInfos)
	if err != nil {
		return nil, errors.New("Error getting risk information")
	}
	return riskInfosBytes, nil
}

func (t *RiskManagementChaincode) getRiskByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID and company of an risk to query")
	}

	uuid := args[0]

	account, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}

	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("RI", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". Error " + err.Error() + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving risk " + uuid + ". \"}"
		return nil, errors.New(jsonResp)
	}
	hasRight := false

	// if string(account) == row.Columns[9].GetString_() && string(organization) == row.Columns[3].GetString_() {
	// 	hasRight = true
	// }
	if hasRight == false {
		var columnsAuth []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: string(account)}}
		col3 := shim.Column{Value: &shim.Column_String_{String_: string(organization)}}
		columnsAuth = append(columnsAuth, col1)
		columnsAuth = append(columnsAuth, col2)
		columnsAuth = append(columnsAuth, col3)
		rowAuth, err := stub.GetRow("RA", columnsAuth)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed retrieving risk authorize. Error " + err.Error() + ". \"}"
			return nil, errors.New(jsonResp)
		}
		if len(rowAuth.Columns) != 0 {
			hasRight = true
		}
	}
	var riskInfoBytes []byte
	if hasRight == true {
		riskInfo := RiskInfo{UUID: row.Columns[0].GetString_(), SupplierID: row.Columns[1].GetString_(), SupplierName: row.Columns[2].GetString_(), AssessCompany: row.Columns[3].GetString_(), RiskEvent: row.Columns[4].GetString_(), RiskType: row.Columns[5].GetString_(), RiskLevel: row.Columns[6].GetString_(), RiskDate: row.Columns[7].GetString_(), RiskLocation: row.Columns[8].GetString_(), CreateBy: row.Columns[9].GetString_(), CreateTimestamp: row.Columns[10].GetString_(), UpdateBy: row.Columns[11].GetString_(), UpdateTimestamp: row.Columns[12].GetString_(), UsefulCount: row.Columns[13].GetString_(), UselessCount: row.Columns[14].GetString_()}
		riskInfoBytes, err = json.Marshal(&riskInfo)
		if err != nil {
			//fmt.Println("error creating account")
			return nil, errors.New("Error creating Risk Info  " + riskInfo.SupplierID)
		}
	}
	return riskInfoBytes, nil
}

func (t *RiskManagementChaincode) getAuthByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	uuid := args[0]

	account, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}

	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}

	var columnsAuth []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: uuid}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: string(account)}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: string(organization)}}
	columnsAuth = append(columnsAuth, col1)
	columnsAuth = append(columnsAuth, col2)
	columnsAuth = append(columnsAuth, col3)
	rowAuth, err := stub.GetRow("RA", columnsAuth)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk authorize. Error " + err.Error() + ". \"}"
		return nil, errors.New(jsonResp)
	}
	if len(rowAuth.Columns) == 0 {
		fmt.Printf("No permission.")
		return nil, fmt.Errorf("No permission.")
	}

	var columns []shim.Column

	rows, err := stub.GetRows("RA", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk auth information. \"}"
		return nil, errors.New(jsonResp)
	}

	var authInfos []AuthInfo
	for row := range rows {
		if len(row.Columns) == 0 {
			jsonResp := "{\"Error\":\"Failed retrieving risk auth information. \"}"
			return nil, errors.New(jsonResp)
		}
		if uuid == row.Columns[0].GetString_() {
			authInfo := AuthInfo{UUID: row.Columns[0].GetString_(), UserName: row.Columns[1].GetString_(), UserCompany: row.Columns[2].GetString_()}
			authInfos = append(authInfos, authInfo)
		}
	}
	authInfosBytes, err := json.Marshal(&authInfos)
	if err != nil {
		return nil, errors.New("Error getting risk information")
	}
	return authInfosBytes, nil
}

func (t *RiskManagementChaincode) getAllHistoryByOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0 arguments.")
	}

	account, err := stub.ReadCertAttribute("account")
	if err != nil {
		fmt.Printf("Error reading attribute 'account' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller account. Error was [%v]", err)
	}

	role, err := stub.ReadCertAttribute("role")
	if err != nil {
		fmt.Printf("Error reading attribute 'role' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller role. Error was [%v]", err)
	}

	organization, err := stub.ReadCertAttribute("organization")
	if err != nil {
		fmt.Printf("Error reading attribute 'organization' [%v] \n", err)
		return nil, fmt.Errorf("Failed fetching caller organization. Error was [%v]", err)
	}

	var columns []shim.Column

	rows, err := stub.GetRows("RH", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving risk information. \"}"
		return nil, errors.New(jsonResp)
	}

	var riskHistories []RiskHistory
	for row := range rows {
		if len(row.Columns) == 0 {
			jsonResp := "{\"Error\":\"Failed retrieving risk history. \"}"
			return nil, errors.New(jsonResp)
		}
		hasRight := false

		if string(role) == "Regulator" {
			hasRight = true
		}
		if hasRight == false {
			var columnsAuth []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: row.Columns[0].GetString_()}}
			col2 := shim.Column{Value: &shim.Column_String_{String_: string(account)}}
			col3 := shim.Column{Value: &shim.Column_String_{String_: string(organization)}}
			columnsAuth = append(columnsAuth, col1)
			columnsAuth = append(columnsAuth, col2)
			columnsAuth = append(columnsAuth, col3)
			rowAuth, err := stub.GetRow("RA", columnsAuth)
			if err != nil {
				jsonResp := "{\"Error\":\"Failed retrieving risk authorize. Error " + err.Error() + ". \"}"
				return nil, errors.New(jsonResp)
			}
			if len(rowAuth.Columns) != 0 {
				hasRight = true
			}
		}
		if hasRight == true {
			fmt.Println(row.Columns[0].GetString_() + row.Columns[1].GetString_())
			riskHistory := RiskHistory{UUID: row.Columns[0].GetString_(), Timestamp: row.Columns[1].GetString_(), OperationType: row.Columns[2].GetString_(), Operation: row.Columns[3].GetString_()}
			riskHistories = append(riskHistories, riskHistory)
		}
	}
	riskHistoriesBytes, err := json.Marshal(&riskHistories)
	if err != nil {
		return nil, errors.New("Error getting risk history")
	}
	return riskHistoriesBytes, nil
}

// Invoke runs callback representing the invocation of a chaincode
func (t *RiskManagementChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//fmt.Println("In Invoke!")
	// Handle different functions
	if function == "register" {
		return t.register(stub, args)
	} else if function == "authorize" {
		return t.authorize(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "thumbUp" {
		return t.thumbUp(stub, args)
	} else if function == "init" {
		return t.Init(stub, "", args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (t *RiskManagementChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "getRiskByID" {
		return t.getRiskByID(stub, args)
	} else if function == "getAllRisks" {
		return t.getAllRisks(stub, args)
	} else if function == "getAllRisksByOwner" {
		return t.getAllRisksByOwner(stub, args)
	} else if function == "getAllHistoryByOwner" {
		return t.getAllHistoryByOwner(stub, args)
	} else if function == "getAuthByID" {
		return t.getAuthByID(stub, args)
	}
	return nil, errors.New("Received unknown function query")
}

func (t *RiskManagementChaincode) addHistory(stub shim.ChaincodeStubInterface, uuid string, opType string, operation string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ok, err := stub.InsertRow("RH", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: uuid}},
			&shim.Column{Value: &shim.Column_String_{String_: timestamp}},
			&shim.Column{Value: &shim.Column_String_{String_: opType}},
			&shim.Column{Value: &shim.Column_String_{String_: operation}}},
	})
	if !ok && err == nil {
		fmt.Println("Error inserting history." + uuid + timestamp + opType + operation)
	}
}

func main() {
	err := shim.Start(new(RiskManagementChaincode))
	if err != nil {
		fmt.Printf("Error starting RiskManagementChaincode: %s", err)
	}
}

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

// This program is an erroneous chaincode program that attempts to put state in query context - query should return error
package main

import (
	//"errors"
	"errors"
	"fmt"
	"strconv"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// OPChaincode example simple Chaincode implementation
type OPChaincode struct {
}

// Init takes a string and int. These are stored as a key/value pair in the state
func (t *OPChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// var err error
	// err = stub.PutState("test", []byte("123"))
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

// Invoke is a no-op
func (t *OPChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "recharge" {
		// Transfer ownership
		return t.recharge(stub, args)
	} else if function == "withdraw" {
		// Transfer ownership
		args[1] = "-" + args[1]
		return t.recharge(stub, args)
	}

	return nil, nil
}
func (t *OPChaincode) recharge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]
	amount := args[1]
	var balance int
	balance = 0
	Avalbytes, err := stub.GetState(account)
	if err == nil || Avalbytes != nil {
		balance, err = strconv.Atoi(string(Avalbytes))
	}
	iAmount, err := strconv.Atoi(amount)
	balance = balance + iAmount
	// Write the state back to the ledger
	err = stub.PutState(account, []byte(strconv.Itoa(balance)))
	if err != nil {
		fmt.Printf("Error starting OPChaincode: %s", err)
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *OPChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "getBalance" {
		// Transfer ownership
		return t.getBalance(stub, args)
	}

	return nil, nil
}
func (t *OPChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("getBalance\n")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	account := args[0]
	Avalbytes, err := stub.GetState(account)
	if err != nil {
		fmt.Printf("Error starting OPChaincode: %s", err)
		return nil, err
	}
	if Avalbytes == nil {
		Avalbytes = []byte("0")
	}
	return Avalbytes, nil
}
func main() {
	err := shim.Start(new(OPChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

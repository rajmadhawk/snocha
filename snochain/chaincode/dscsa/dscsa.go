/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright Expiryship.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the LPN structure, with 4 properties.  Structure tags are used by encoding/json library
type LPN struct {
	GTIN   string `json:"GTIN"`
	Serial  string `json:"Serial"`
	Lot string `json:"Lot"`
	Expiry  string `json:"Expiry"`
}

/*
 * The Init method is called when the Smart Contract "fabLPN" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabLPN"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryLPN" {
		return s.queryLPN(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createLPN" {
		return s.createLPN(APIstub, args)
	} else if function == "queryAllLPNs" {
		return s.queryAllLPNs(APIstub)
	} else if function == "changeLPNExpiry" {
		return s.changeLPNExpiry(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryLPN(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	LPNAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(LPNAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	LPNs := []LPN{
		LPN{GTIN: "0030591255012", Serial: "102010203333", Lot: "APQPM15", Expiry: "20171030"},
		LPN{GTIN: "0030591255029", Serial: "102010203334", Lot: "APQPM15", Expiry: "20171030"},
		LPN{GTIN: "0030591255013", Serial: "102010203335", Lot: "APQPM15", Expiry: "20171030"},
		LPN{GTIN: "0030591255015", Serial: "102010203336", Lot: "APQPM15", Expiry: "20171030"},
	}

	i := 0
	for i < len(LPNs) {
		fmt.Println("i is ", i)
		LPNAsBytes, _ := json.Marshal(LPNs[i])
		APIstub.PutState("LPN"+strconv.Itoa(i), LPNAsBytes)
		fmt.Println("Added", LPNs[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createLPN(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var LPN = LPN{GTIN: args[1], Serial: args[2], Lot: args[3], Expiry: args[4]}

	LPNAsBytes, _ := json.Marshal(LPN)
	APIstub.PutState(args[0], LPNAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllLPNs(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "LPN0"
	endKey := "LPN999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllLPNs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeLPNExpiry(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	LPNAsBytes, _ := APIstub.GetState(args[0])
	LPN := LPN{}

	json.Unmarshal(LPNAsBytes, &LPN)
	LPN.Expiry = args[1]

	LPNAsBytes, _ = json.Marshal(LPN)
	APIstub.PutState(args[0], LPNAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

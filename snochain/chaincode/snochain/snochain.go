 package main
 
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
 
 // Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
 type Lpn struct {
	 Barcode   string `json:"barcode"`
 }
 
 /*
  * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
  * Best practice is to have any Ledger initialization in separate function -- see initLedger()
  */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 /*
  * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
  * The calling application program has also specified the particular smart contract function to be called, with arguments
  */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested Smart Contract function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "queryLpn" {
		 return s.queryLpn(APIstub, args)
	 } else if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "createLpn" {
		 return s.createLpn(APIstub, args)
	 } else if function == "queryAllLpns" {
		 return s.queryAllLpns(APIstub)
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) queryLpn(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 lpnAsBytes, _ := APIstub.GetState(args[0])
	 return shim.Success(lpnAsBytes)
 }
 
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 lpns := []Lpn{
		 Lpn{Barcode: "(01)00356891234567(21)1000000000(10)APN3457(17)201231"},
		 Lpn{Barcode: "(01)00356891234567(21)1000000003(10)APN3457(17)201231"},
		 Lpn{Barcode: "(01)00356891234567(21)1000000004(10)APN3457(17)201231"},
		 Lpn{Barcode: "(01)00356891234567(21)1000000005(10)APN3457(17)201231"},
		 Lpn{Barcode: "(01)00356891234567(21)1000000006(10)APN3457(17)201231"},
		 Lpn{Barcode: "(01)00356891234567(21)1000000007(10)APN3457(17)201231"},

	 }
 
	 i := 0
	 for i < len(lpns) {
		 fmt.Println("i is ", i)
		 lpnAsBytes, _ := json.Marshal(lpns[i])
		 APIstub.PutState("LPN"+strconv.Itoa(i), lpnAsBytes)
		 fmt.Println("Added", lpns[i])
		 i = i + 1
	 }
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) createLpn(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 var lpn = Lpn{Barcode: args[1]}
 
	 lpnAsBytes, _ := json.Marshal(lpn)
	 APIstub.PutState(args[0], lpnAsBytes)
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) queryAllLpns(APIstub shim.ChaincodeStubInterface) sc.Response {
 
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
 
	 fmt.Printf("- queryAllLpns:\n%s\n", buffer.String())
 
	 return shim.Success(buffer.Bytes())
 }
 
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 
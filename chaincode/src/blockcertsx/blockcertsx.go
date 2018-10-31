
package main

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"strconv"
	//"strings"
	//"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Blockcertx struct {
	context			 []string					`json:"@context"`
	badge  			 BadgeObj					`json:"badge"`
	id					 string	 					`json:"id"`
	issuedOn 		 string 					`json:"issuedOn"`
	recipient 	 RecipientObj			`json:"recipient"`
	recP 				 RecipientProfile `json:"recipientProfile"`
	certType 		 string 					`json:"type"`
	verification VerificationObj 	`json:"verification"`
}

type BadgeObj struct{

	criteria 		string `json:"criteria"`
	description string `json:"description"`
	id 					string `json:"id"`
	image				string `json:"image"`
	issuer			string `json:"issuer"` 	// TODO : Issuer struct
	name 				string `json:"name"`
	signLines		string `json:"signatureLines"` //TODO: signatureLines struct/Array
	badgeType 	string `json:"type"`
}

type RecipientObj struct{
	hashed 			bool	`json:"hashed"`
	identity		string `json:"identity"`
	RecType 		string `json:"type"`
}

type RecipientProfile struct {
	name string `json:"name"`
	publicKey	string `json:"publicKey"`
	rpType 		string `json:"type"`  //TODO it's a string array
}

type VerificationObj struct {

	publicKey	string `json:"publicKey"`
	verType		string `json:"type"` // TODO string array
}








type cert struct {
	Name       string 		`json:"certName"`    //the fieldtags are needed to keep case from bouncing around
	certBytes  Blockcertx `json:"certJSON"` //docType is used to distinguish the various types of objects in state database
}





// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting blockcerts chaincode: %s", err)
	}else{
		fmt.Println("Successfully Started BlockcertsX")
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("-Init function-")

	fmt.Println("BlockcertsX Is Starting Up")
	funcName, args := stub.GetFunctionAndParameters()
	var number int
	var err error
	txId := stub.GetTxID()

	fmt.Println("Init() is running")
	fmt.Println("Transaction ID:", txId)
	fmt.Println("  GetFunctionAndParameters() function:", funcName)
	fmt.Println("  GetFunctionAndParameters() args count:", len(args))
	fmt.Println("  GetFunctionAndParameters() args found:", args)

	// expecting 1 arg for instantiate or upgrade
	if len(args) == 1 {
		fmt.Println("  GetFunctionAndParameters() arg[0] length", len(args[0]))

		// expecting arg[0] to be length 0 for upgrade
		if len(args[0]) == 0 {
			fmt.Println("  Uh oh, args[0] is empty...")
		} else {
			fmt.Println("  Great news everyone, args[0] is not empty")

			// convert numeric string to integer
			number, err = strconv.Atoi(args[0])
			if err != nil {
				return shim.Error("Expecting a numeric string argument to Init() for instantiate")
			}

			// this is a very simple test. let's write to the ledger and error out on any errors
			// it's handy to read this right away to verify network is healthy if it wrote the correct value
			err = stub.PutState("selftest", []byte(strconv.Itoa(number)))
			if err != nil {
				return shim.Error(err.Error())                  //self-test fail
			}
		}
	}

	// showing the alternative argument shim function
	alt := stub.GetStringArgs()
	fmt.Println("  GetStringArgs() args count:", len(alt))
	fmt.Println("  GetStringArgs() args found:", alt)

	// store compatible marbles application version
	err = stub.PutState("marbles_ui", []byte("4.0.1"))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Ready for action")                          //self-test pass
	return shim.Success(nil)



}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is run" + function)

	// Handle different functions
	if function == "initCert" { // Save a given certificate
		return t.initCert(stub, args)
	} else if function == "readCert" { //read a certificate
		return t.readCert(stub, args)
	} //else if function == "saveCert" { // Save certificate
		//return t.saveCert(stub, args)
	//}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initMarble - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initCert(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// ==== Input sanitation ====
	fmt.Println("start init blockcertx")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	certName := args[0]
	//certString := args[1]


	// ==== Create cert object and marshal to JSON ====
	//objectType := "cert"
	//cert := &cert{certName, certBytes}

	//var cert blockcertx

	//err = json.Unmarshal(certString, &cert)
	//if err != nil {
		//return shim.Error(err.Error())
	//}
	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//marbleJSONasBytes := []byte(str)

	// === Save cert to state ===
	err = stub.PutState(certName, []byte("TEST_STRING"))
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init certificate")
	return shim.Success(nil)
}

// ===============================================
// readCert - read a certificate from chaincode state
// ===============================================
func (t *SimpleChaincode) readCert(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	fmt.Println("Starting readCert")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the certificate to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("About to return Value")
	fmt.Println(valAsbytes)

	return shim.Success(valAsbytes)
}

package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"chaincode/chaincode"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

/*func TestCreatePatient(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.CreatePatient(transactionContext, "", "", 0, "", "", "", "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.CreatePatient(transactionContext, "Patient1", "", 0, "", "", "", "", "")
	require.EqualError(t, err, "the patient Patient1 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve Patient"))
	err = assetTransfer.CreatePatient(transactionContext, "patient1", "", 0, "", "", "", "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve Patient")
}*/

func TestReadPatientMedicalInfo(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.PatientData{ID: "Patient1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	Patient, err := assetTransfer.ReadPatientMedicalInfo(transactionContext, "", "")
	require.NoError(t, err)
	require.Equal(t, expectedAsset.Record, Patient)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve Patient"))
	_, err = assetTransfer.ReadPatientMedicalInfo(transactionContext, "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve Patient")

	chaincodeStub.GetStateReturns(nil, nil)
	Patient, err = assetTransfer.ReadPatientMedicalInfo(transactionContext, "Doctor1", "Patient1")
	require.EqualError(t, err, "the Patient Patient1 does not exist")
	require.Nil(t, Patient)
}

func TestReadPatientAllInfo(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedPatient := &chaincode.PatientData{ID: "Patient1"}
	bytes, err := json.Marshal(expectedPatient)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	Patient, err := assetTransfer.ReadPatientAllInfo(transactionContext, "", "")
	require.NoError(t, err)
	require.Equal(t, expectedPatient, Patient)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve Patient"))
	_, err = assetTransfer.ReadPatientAllInfo(transactionContext, "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve Patient")

	chaincodeStub.GetStateReturns(nil, nil)
	Patient, err = assetTransfer.ReadPatientAllInfo(transactionContext, "Doctor1", "Patient1")
	require.EqualError(t, err, "the Patient Patient1 does not exist")
	require.Nil(t, Patient)
}

func TestUpdateMedicalpatientrecords(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedPatient := &chaincode.PatientData{ID: "Patient1"}
	bytes, err := json.Marshal(expectedPatient)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.UpdateMedicalpatientrecords(transactionContext, "", "", "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.UpdateMedicalpatientrecords(transactionContext, "Doctor1", "Patient1", "", "")
	require.EqualError(t, err, "the Patient Patient1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve Patient"))
	err = assetTransfer.UpdateMedicalpatientrecords(transactionContext, "Doctor1", "Patient1", "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve Patient")
}

func TestDeletePatient(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	Patient := &chaincode.PatientData{ID: "Patient1"}
	bytes, err := json.Marshal(Patient)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.DeletePatient(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.DeletePatient(transactionContext, "Patient1")
	require.EqualError(t, err, "the Patient Patient1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve Patient"))
	err = assetTransfer.DeletePatient(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve Patient")
}

func TestGrantAccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}

	// Define test patient data
	patientData := &chaincode.PatientData{
		ID:     "Patient1",
		Name:   "John Doe",
		Age:    35,
		Gender: "Male",
		Access: map[string]bool{},
	}

	// Marshal patient data into JSON format
	patientDataJSON, err := json.Marshal(patientData)
	require.NoError(t, err)

	// Set up chaincode stub to return patient data for patient ID "patient1"
	chaincodeStub.GetStateReturns(patientDataJSON, nil)

	// Call GrantAccess with valid patient and doctor IDs
	err = assetTransfer.GrantAccess(transactionContext, "patient1", "doctor1")
	require.NoError(t, err)

	// Verify that access has been granted for doctor1
	updatedPatientDataJSON, err := chaincodeStub.GetState("patient1")
	require.NoError(t, err)
	var updatedPatientData chaincode.PatientData
	err = json.Unmarshal(updatedPatientDataJSON, &updatedPatientData)
	require.NoError(t, err)
	require.True(t, updatedPatientData.Access["doctor1"])

	// Call GrantAccess with invalid doctor ID
	err = assetTransfer.GrantAccess(transactionContext, "patient1", "")
	require.EqualError(t, err, "doctor ID cannot be empty")

	// Call GrantAccess with invalid patient ID
	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.GrantAccess(transactionContext, "", "doctor1")
	require.EqualError(t, err, "patient ID cannot be empty")

	// Call GrantAccess with patient data not found in ledger
	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.GrantAccess(transactionContext, "patient2", "doctor1")
	require.EqualError(t, err, "patient data with ID patient2 does not exist")

	// Call GrantAccess with error reading patient data from ledger
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve patient data"))
	err = assetTransfer.GrantAccess(transactionContext, "patient1", "doctor1")
	require.EqualError(t, err, "failed to read patient data from world state: unable to retrieve patient data")
}

func TestRevokeAccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}

	// Define test patient data
	patientData := &chaincode.PatientData{

		ID:     "Patient1",
		Name:   "John Doe",
		Age:    35,
		Gender: "Male",
		Access: map[string]bool{
			"doctor1": true,
			"doctor2": true,
		},
	}

	// Marshal patient data into JSON format
	patientDataJSON, err := json.Marshal(patientData)
	require.NoError(t, err)

	// Set up chaincode stub to return patient data for patient ID "patient1"
	chaincodeStub.GetStateReturns(patientDataJSON, nil)

	// Call RevokeAccess with valid patient and doctor IDs
	err = assetTransfer.RevokeAccess(transactionContext, "patient1", "doctor1")
	require.NoError(t, err)

	// Verify that access has been revoked for doctor1
	updatedPatientDataJSON, err := chaincodeStub.GetState("patient1")
	require.NoError(t, err)
	var updatedPatientData chaincode.PatientData
	err = json.Unmarshal(updatedPatientDataJSON, &updatedPatientData)
	require.NoError(t, err)
	require.False(t, updatedPatientData.Access["doctor1"])

	// Call RevokeAccess with invalid doctor ID
	err = assetTransfer.RevokeAccess(transactionContext, "patient1", "doctor3")
	require.EqualError(t, err, "Doctor with ID doctor3 does not have access to patient data with ID patient1")

	// Call RevokeAccess with invalid patient ID
	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.RevokeAccess(transactionContext, "patient2", "doctor1")
	require.EqualError(t, err, "patient data with ID patient2 does not exist")

	// Call RevokeAccess with error reading patient data from ledger
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve patient data"))
	err = assetTransfer.RevokeAccess(transactionContext, "patient1", "doctor1")
	require.EqualError(t, err, "failed to read patient data from world state: unable to retrieve patient data")

}

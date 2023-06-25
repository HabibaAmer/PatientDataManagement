package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Patient
type SmartContract struct {
	contractapi.Contract
}

// PatientData describes basic details of what makes up a simple Patient

type PatientData struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Age       int             `json:"age"`
	Gender    string          `json:"gender"`
	BloodType string          `json:"bloodType"`
	Allergies string          `json:"allergies"`
	Access    map[string]bool `json:"Doctor Access"`
	Record    MedicalRecord   `json:"record"`
}
type MedicalRecord struct {
	Diagnose           string   `json:"diagnose"`
	Medications        string   `json:"medications"`
	DiagnosesHistory   []string `json:"diagnoseshistory"`
	MedicationsHistory []string `json:"medicationhistory"`
}

/*
// Define the access control struct
type AccessControl struct {
	providerID string
	authorized bool
	Access map[string]bool

} */

// InitLedger adds a base set of Patients to the ledger --> The init function is called when the smart contract is first deployed to the network
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	Patients := []PatientData{
		{ID: "Patient1", Name: "test1", Age: 5, Gender: "male", BloodType: "B+", Allergies: "xx", Record: MedicalRecord{"Diagnose1", "Medications1", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor1": true}},
		{ID: "Patient2", Name: "test2", Age: 5, Gender: "male", BloodType: "A+", Allergies: "yy", Record: MedicalRecord{"Diagnose2", "Medications2", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor2": true}},
		{ID: "Patient3", Name: "test3", Age: 10, Gender: "female", BloodType: "AB", Allergies: "cc", Record: MedicalRecord{"Diagnose3", "Medications3", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor1": true}},
		{ID: "Patient4", Name: "test4", Age: 10, Gender: "female", BloodType: "O+", Allergies: "nn", Record: MedicalRecord{"Diagnose4", "Medications4", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor2": true}},
		{ID: "Patient5", Name: "test5", Age: 15, Gender: "female", BloodType: "O-", Allergies: "mm", Record: MedicalRecord{"Diagnose5", "Medications5", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor1": true}},
		{ID: "Patient6", Name: "test6", Age: 15, Gender: "female", BloodType: "B-", Allergies: "jj", Record: MedicalRecord{"Diagnose6", "Medications6", []string{"diagnose11", "diagnose12"}, []string{"medication11", "medication12"}}, Access: map[string]bool{"Doctor2": true}},
	}

	fmt.Errorf("ledger is initialed successfuly")

	for _, Patient := range Patients {
		PatientJSON, err := json.Marshal(Patient) // take each patient and convert it to json format then store this format in PatientJSON file and check for error
		if err != nil {                           //if error field is not empty
			return err //then return the error and return to main and exit the init function
		}

		err = ctx.GetStub().PutState(Patient.ID, PatientJSON) //ctx:context-interface that used to access the BC ledger , PutState: generate the key-value pair==(Patient.ID-PatientJSON) , GetStub: provide APIs access to the world state
		// = is used for assigment while := is used for decleration
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

/*
func (s *SmartContract) CreatePatient(ctx contractapi.TransactionContextInterface, patientID string, name string, age int, gender string, bloodType string, allergies string, diagnose string, medication string) (interface{}, error) {
	exists, err := s.PatientExists(ctx, patientID)

	fmt.Println("we will search if patient exists or no!!")

	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("the Patient %s already exists", patientID)
	}

	NewPatient := PatientData{
		ID:        patientID,
		Name:      name,
		Age:       age,
		Gender:    gender,
		BloodType: bloodType,
		Allergies: allergies,
		Record:    MedicalRecord{Diagnose: diagnose, Medications: medication},
	}
	PatientJSON, err := json.Marshal(NewPatient)
	if err != nil {
		return nil, err
	}
	fmt.Printf("The Patient is added to the ledger successfully")
	return ctx.GetStub().PutState(patientID, PatientJSON), nil
}
*/

func (s *SmartContract) PatientExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	PatientJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return PatientJSON != nil, nil
}

func (s *SmartContract) UpdateMedicalpatientrecords(ctx contractapi.TransactionContextInterface, providerID string, patientID string, diagnose string, medication string) error {

	access, err := s.hasPermission(ctx, providerID, patientID)

	if access == false && err != nil {
		return fmt.Errorf("Provider does not have permission to share data for patient %s", patientID)
	}

	patientRecord, err := ctx.GetStub().GetState(patientID) // get state of ID from the ledger and store it in PatientRecords , if an error occur store it in err
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if patientRecord == nil {
		return fmt.Errorf("the patient with ID %s does not exist", patientID)
	}

	var patientData PatientData //declare updated data as a variable
	err = json.Unmarshal(patientRecord, &patientData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal patient record: %v", err)
	}

	//overwrite the diagnose and medication fields
	patientData.Record.Diagnose = diagnose
	patientData.Record.Medications = medication

	//add the new diagnose and medication to the history
	patientData.Record.DiagnosesHistory = append(patientData.Record.DiagnosesHistory, diagnose)
	patientData.Record.MedicationsHistory = append(patientData.Record.MedicationsHistory, medication)

	updatedpatientData, err := json.Marshal(patientData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient record: %v", err)
	}

	err = ctx.GetStub().PutState(patientID, updatedpatientData)
	if err != nil {
		return fmt.Errorf("failed to update patient record: %v", err)
	}

	return nil

}

// ReadPatient returns the Medical info only stored in the world state with given id.
func (s *SmartContract) ReadPatientMedicalInfo(ctx contractapi.TransactionContextInterface, providerID string, patientID string) (*MedicalRecord, error) {

	access, err := s.hasPermission(ctx, providerID, patientID)

	if access == false && err != nil {
		return nil, fmt.Errorf("Provider does not have permission to share data for patient %s", patientID)
	}
	PatientRecordJSON, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if PatientRecordJSON == nil {
		return nil, fmt.Errorf("the patient %s does not exist", patientID)
	}

	var Patient PatientData
	err = json.Unmarshal(PatientRecordJSON, &Patient)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Diagnosis: %s  , Medication: %s \n", Patient.Record.Diagnose, Patient.Record.Medications)
	fmt.Println("Diagnoses history:")
	for _, diagnosis := range Patient.Record.DiagnosesHistory {
		fmt.Printf("- %s\n", diagnosis)
	}
	fmt.Println("Medication history:")
	for _, history := range Patient.Record.MedicationsHistory {
		fmt.Printf("- %s\n", history)
	}
	return &Patient.Record, nil
}

// ReadPatient returns the all patient info stored in the world state with given id.
func (s *SmartContract) ReadPatientAllInfo(ctx contractapi.TransactionContextInterface, providerID string, patientID string) (*PatientData, error) {

	access, err := s.hasPermission(ctx, providerID, patientID)

	if access == false && err != nil {
		return nil, fmt.Errorf("Provider does not have permission to share data for patient %s", patientID)
	}

	PatientJSON, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if PatientJSON == nil {
		return nil, fmt.Errorf("the patient %s does not exist", patientID)
	}

	var Patient PatientData
	err = json.Unmarshal(PatientJSON, &Patient)
	if err != nil {
		return nil, err
	}

	return &Patient, nil
}

// DeletePatient deletes an given patient from the world state.
func (s *SmartContract) DeletePatient(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.PatientExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Patient %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// Define the grantAccess function

func (p *SmartContract) GrantAccess(ctx contractapi.TransactionContextInterface, patientID string, providerID string) error {
	// Check if the caller is authorized to grant access
	creator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}
	if creator != patientID {
		return fmt.Errorf("caller is not authorized to grant access for this patient")
	}

	// Retrieve patient data from the ledger
	patientDataJSON, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return fmt.Errorf("failed to read patient data from world state: %v", err)

	}
	if patientDataJSON == nil {
		return fmt.Errorf("patient data with ID %s does not exist", patientID)
	}

	// Unmarshal patient data JSON into struct
	var patientData PatientData
	err = json.Unmarshal(patientDataJSON, &patientData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal patient data JSON: %v", err)
	}

	access, ok := patientData.Access[providerID]
	if access == false || ok == false {
		patientData.Access = map[string]bool{providerID: true}
	}

	updatedpatientData, err := json.Marshal(patientData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient record: %v", err)
	}

	err = ctx.GetStub().PutState(patientID, updatedpatientData)
	if err != nil {
		return fmt.Errorf("failed to update patient record: %v", err)
	}

	/*
		// Check if the provider exists
		providerExists, err := p.ProviderExists(ctx, providerID)
		if err != nil {
			return fmt.Errorf("failed to check if provider exists: %v", err)
		}
		if !providerExists {
			return fmt.Errorf("provider does not exist")
		}*/

	/*
		// Check if the access control record exists
		accessControlKey, err := ctx.GetStub().CreateCompositeKey("accessControl", []string{patientID, providerID})
		if err != nil {
			return fmt.Errorf("failed to create composite key: %v", err)
		}
		accessControlBytes, err := ctx.GetStub().GetState(accessControlKey) //get the acl from ledger and stor it into accesscontrolbytes
		if err != nil {
			return fmt.Errorf("failed to read access control record: %v", err)
		}

		// If the access control record does not exist, create a new one
		var accessControl AccessControl
		if accessControlBytes == nil {
			accessControl = AccessControl{providerID, true}
		} else {
			// If the access control record exists, update it
			err = json.Unmarshal(accessControlBytes, &accessControl)
			if err != nil {
				return fmt.Errorf("failed to unmarshal access control record: %v", err)
			}
			accessControl.authorized = true
		}

		// Write the updated access control record to the ledger
		accessControlBytes, err = json.Marshal(accessControl)
		if err != nil {
			return fmt.Errorf("failed to marshal access control record: %v", err)
		}
		err = ctx.GetStub().PutState(accessControlKey, accessControlBytes)
		if err != nil {
			return fmt.Errorf("failed to write access control record: %v", err)
		}*/

	return nil
}

/*
// define provider exists function
//to check whether the provider ID exists on the ledger or not

func (p *SmartContract) ProviderExists(ctx contractapi.TransactionContextInterface, providerID string) (bool, error) {
	providerKey, err := ctx.GetStub().CreateCompositeKey("provider", []string{providerID}) // key is created by combining the string "provider" to the provider ID
	//this key is used to retrive the record
	if err != nil {
		return false, fmt.Errorf("failed to create composite key: %v", err)
	}
	//GetState() is used to retrieve the provider record from the ledger using the composite key,record is stored in providerBytes
	providerBytes, err := ctx.GetStub().GetState(providerKey)
	if err != nil {
		return false, fmt.Errorf("failed to read provider record: %v", err)
	}
	return providerBytes != nil, nil
}*/

// define Revoke Access function
func (s *SmartContract) RevokeAccess(ctx contractapi.TransactionContextInterface, patientID string, DoctorID string) error {
	// Check if the caller is authorized to revok access
	creator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}
	if creator != patientID {
		return fmt.Errorf("caller is not authorized to grant access for this patient")
	}

	// Retrieve patient data from the ledger
	patientDataJSON, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return fmt.Errorf("failed to read patient data from world state: %v", err)
	}
	if patientDataJSON == nil {
		return fmt.Errorf("patient data with ID %s does not exist", patientID)
	}

	// Unmarshal patient data JSON into struct
	var patientData PatientData
	err = json.Unmarshal(patientDataJSON, &patientData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal patient data JSON: %v", err)
	}

	// Check if doctor has access to patient data
	if _, ok := patientData.Access[DoctorID]; !ok {
		return fmt.Errorf("Doctor with ID %s does not have access to patient data with ID %s", DoctorID, patientID)
	}

	// Revoke access for specific hospital
	//patientData.Access[DoctorID] = false
	patientData.Access = map[string]bool{DoctorID: false}

	// Marshal updated patient data into JSON format
	patientDataJSON, err = json.Marshal(patientData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient data into JSON format: %v", err)
	}

	// Update patient data in the ledger
	err = ctx.GetStub().PutState(patientID, patientDataJSON)
	if err != nil {
		return fmt.Errorf("failed to update patient data in the ledger: %v", err)
	}

	return nil
}

// define share data function

func (s *SmartContract) hasPermission(ctx contractapi.TransactionContextInterface, userID string, patientID string) (bool, error) {
	// Check if the user has permission to access the patient data
	// In this example implementation, we assume that only the patient and healthcare providers with the patient's permission have access to the data
	if userID == patientID {
		return true, nil
	}

	// Retrieve patient data from the ledger
	patientDataJSON, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return false, fmt.Errorf("failed to read patient data from world state: %v", err)
	}
	if patientDataJSON == nil {
		return false, fmt.Errorf("patient data with ID %s does not exist", patientID)
	}

	// Unmarshal patient data JSON into struct
	var patientData PatientData
	err = json.Unmarshal(patientDataJSON, &patientData)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal patient data JSON: %v", err)
	}

	access, ok := patientData.Access[userID]
	if !ok {
		return false, fmt.Errorf("%s doesn't exist in %s access list", userID, patientID)
	}
	if access == false {
		return false, fmt.Errorf("%s doesn't have access to %s", userID, patientID)
	}

	/*
		// Check if doctor has access to patient data
		if permission := patientData.Access.providerID; ok {
			return true , nil
		} else {
			return false , fmt.Errorf("Doctor with ID %s does not have access to patient data with ID %s", userID, patientID)
		}

		/*
		// Check if the user has been granted permission by the patient
		permissionKey := fmt.Sprintf("%s_%s", patientID, userID)
		permissionBytes, err := ctx.GetStub().GetState(permissionKey)
		if err != nil {
			return false
		}
		return permissionBytes != nil */
	return true, nil
}

func (s *SmartContract) ShareData(ctx contractapi.TransactionContextInterface, patientID string, recipientID string, data []byte) error {
	// Check if the caller has permission to share data
	callerID, err := ctx.GetClientIdentity().GetID()

	if err != nil {
		return err
	}

	access, err := s.hasPermission(ctx, callerID, patientID)

	if access == false && err != nil {
		return fmt.Errorf("caller does not have permission to share data for patient %s", patientID)
	}

	// Check if the recipient has permission to access the data

	permission, err := s.hasPermission(ctx, recipientID, patientID)

	if permission == false && err != nil {
		return fmt.Errorf("recipient does not have permission to access data for patient %s", patientID)
	}

	// Share the data with the recipient
	err = ctx.GetStub().PutState(fmt.Sprintf("%s_%s", patientID, recipientID), data)
	if err != nil {
		return err
	}

	return nil
}

/*func main() {
    // Create a map of patient records
    updatedRecord := []MedicalRecord{}

    // Print the patient records before updating
    fmt.Println("Patient Records before updating:")
    for name, patient := range patientRecords {
        fmt.Printf("%s: %+v\n", name, patient)
    }

    // Update Alice's record
    updatedPatient := Patient{Name: "Alice", Age: 26, PhoneNumber: "555-4321"}
    UpdatePatientRecord("Alice", updatedPatient, patientRecords)

    // Print the patient records after updating
    fmt.Println("Patient Records after updating:")
    for name, patient := range patientRecords {
        fmt.Printf("%s: %+v\n", name, patient)
    }

    // FOR DATA SHARING FUNCTION
    chaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        fmt.Printf("Error creating patient data sharing chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting patient data sharing chaincode: %s", err.Error())
    }

}
*/

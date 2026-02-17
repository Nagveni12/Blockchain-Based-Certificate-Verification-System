package main

import (
	"encoding/json" 
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Certificate struct {
	ID          string `json:"id"`
	StudentName string `json:"studentName"`
	Issuer      string `json:"issuer"`
	IssueDate   string `json:"issueDate"`
	IPFSHash    string `json:"ipfsHash"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) IssueCertificate(ctx contractapi.TransactionContextInterface, id string, studentName string, issuer string, issueDate string, ipfsHash string) error {
	
	certificate := Certificate{
		ID:          id,
		StudentName: studentName,
		Issuer:      issuer,
		IssueDate:   issueDate,
		IPFSHash:    ipfsHash,
	}

	certificateJSON, err := json.Marshal(certificate)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, certificateJSON)
}

func (s *SmartContract) GetCertificate(ctx contractapi.TransactionContextInterface, id string) (*Certificate, error) {
	certificateJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if certificateJSON == nil {
		return nil, fmt.Errorf("certificate %s does not exist", id)
	}

	var certificate Certificate
	err = json.Unmarshal(certificateJSON, &certificate)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func (s *SmartContract) GetAllCertificates(ctx contractapi.TransactionContextInterface) ([]*Certificate, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var certificates []*Certificate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var certificate Certificate
		err = json.Unmarshal(queryResponse.Value, &certificate)
		if err != nil {
			continue // Skip invalid entries
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating certificate chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting certificate chaincode: %v", err)
	}
}

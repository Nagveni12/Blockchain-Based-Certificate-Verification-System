package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type Certificate struct {
	CertificateID string `json:"certificateId"`
	StudentName   string `json:"studentName"`
	Issuer        string `json:"issuer"`
	IssueDate     string `json:"issueDate"`
	IPFSHash      string `json:"ipfsHash"`
	DocType       string `json:"docType"`
}

type CertificateContract struct {
	contractapi.Contract
}

func (cc *CertificateContract) IssueCertificate(ctx contractapi.TransactionContextInterface, certificateId, studentName, issuer, issueDate, ipfsHash string) error {
	certificate := Certificate{
		CertificateID: certificateId,
		StudentName:   studentName,
		Issuer:        issuer,
		IssueDate:     issueDate,
		IPFSHash:      ipfsHash,
		DocType:       "certificate",
	}

	certificateJSON, err := json.Marshal(certificate)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(certificateId, certificateJSON)
}

func (cc *CertificateContract) GetCertificate(ctx contractapi.TransactionContextInterface, certificateId string) (*Certificate, error) {
	certificateJSON, err := ctx.GetStub().GetState(certificateId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if certificateJSON == nil {
		return nil, fmt.Errorf("certificate %s does not exist", certificateId)
	}

	var certificate Certificate
	err = json.Unmarshal(certificateJSON, &certificate)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func (cc *CertificateContract) GetAllCertificates(ctx contractapi.TransactionContextInterface) ([]*Certificate, error) {
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
			return nil, err
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&CertificateContract{})
	if err != nil {
		log.Panicf("Error creating certificate chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting certificate chaincode: %v", err)
	}
}

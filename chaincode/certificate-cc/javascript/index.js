'use strict';

const { Contract } = require('fabric-contract-api');

class CertificateContract extends Contract {

    // Issue a new certificate
    async issueCertificate(ctx, certificateId, studentName, issuer, issueDate, ipfsHash) {
        const certificate = {
            certificateId,
            studentName,
            issuer,
            issueDate,
            ipfsHash,
            docType: 'certificate'
        };

        await ctx.stub.putState(certificateId, Buffer.from(JSON.stringify(certificate)));
        return JSON.stringify(certificate);
    }

    // Get certificate by ID
    async getCertificate(ctx, certificateId) {
        const certificateJSON = await ctx.stub.getState(certificateId);
        if (!certificateJSON || certificateJSON.length === 0) {
            throw new Error(`Certificate ${certificateId} does not exist`);
        }
        return certificateJSON.toString();
    }

    // Get all certificates
    async getAllCertificates(ctx) {
        const allResults = [];
        const iterator = await ctx.stub.getStateByRange('', '');
        let result = await iterator.next();
        while (!result.done) {
            const strValue = Buffer.from(result.value.value.toString()).toString('utf8');
            let record;
            try {
                record = JSON.parse(strValue);
            } catch (err) {
                console.log(err);
                record = strValue;
            }
            allResults.push(record);
            result = await iterator.next();
        }
        return JSON.stringify(allResults);
    }

    // Check if certificate exists
    async certificateExists(ctx, certificateId) {
        const certificateJSON = await ctx.stub.getState(certificateId);
        return certificateJSON && certificateJSON.length > 0;
    }
}

module.exports = CertificateContract;

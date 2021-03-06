/*
Copyright IBM Corp. 2017 All Rights Reserved.

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

package msp

import (
	"os"
	"reflect"
	"testing"

	"fmt"

	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/stretchr/testify/assert"
)

func TestNoopMSP(t *testing.T) {
	noopmsp := NewNoopMsp()

	id, err := noopmsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded")
		return
	}

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded")
		return
	}

	idBack, err := noopmsp.DeserializeIdentity(serializedID)
	if err != nil {
		t.Fatalf("DeserializeIdentity should have succeeded")
		return
	}

	msg := []byte("foo")
	sig, err := id.Sign(msg)
	if err != nil {
		t.Fatalf("Sign should have succeeded")
		return
	}

	err = id.Verify(msg, sig)
	if err != nil {
		t.Fatalf("The signature should be valid")
		return
	}

	err = idBack.Verify(msg, sig)
	if err != nil {
		t.Fatalf("The signature should be valid")
		return
	}
}

func TestMSPSetupBad(t *testing.T) {
	_, err := GetLocalMspConfig("barf", nil, "DEFAULT")
	if err == nil {
		t.Fatalf("Setup should have failed on an invalid config file")
		return
	}
}

func TestGetIdentities(t *testing.T) {
	_, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetDefaultSigningIdentity failed with err %s", err)
		return
	}
}

func TestSerializeIdentities(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded, got err %s", err)
		return
	}

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded, got err %s", err)
		return
	}

	idBack, err := localMsp.DeserializeIdentity(serializedID)
	if err != nil {
		t.Fatalf("DeserializeIdentity should have succeeded, got err %s", err)
		return
	}

	err = localMsp.Validate(idBack)
	if err != nil {
		t.Fatalf("The identity should be valid, got err %s", err)
		return
	}

	if !reflect.DeepEqual(id.GetPublicVersion(), idBack) {
		t.Fatalf("Identities should be equal (%s) (%s)", id, idBack)
		return
	}
}

func TestValidateCAIdentity(t *testing.T) {
	caID := getIdentity(t, cacerts)

	err := localMsp.Validate(caID)
	assert.Error(t, err)
}

func TestValidateAdminIdentity(t *testing.T) {
	caID := getIdentity(t, admincerts)

	err := localMsp.Validate(caID)
	assert.NoError(t, err)
}

func TestSerializeIdentitiesWithWrongMSP(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded, got err %s", err)
		return
	}

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded, got err %s", err)
		return
	}

	sid := &msp.SerializedIdentity{}
	err = proto.Unmarshal(serializedID, sid)
	assert.NoError(t, err)

	sid.Mspid += "BARF"

	serializedID, err = proto.Marshal(sid)
	assert.NoError(t, err)

	_, err = localMsp.DeserializeIdentity(serializedID)
	assert.Error(t, err)
}

func TestSerializeIdentitiesWithMSPManager(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded, got err %s", err)
		return
	}

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded, got err %s", err)
		return
	}

	_, err = mspMgr.DeserializeIdentity(serializedID)
	assert.NoError(t, err)

	sid := &msp.SerializedIdentity{}
	err = proto.Unmarshal(serializedID, sid)
	assert.NoError(t, err)

	sid.Mspid += "BARF"

	serializedID, err = proto.Marshal(sid)
	assert.NoError(t, err)

	_, err = mspMgr.DeserializeIdentity(serializedID)
	assert.Error(t, err)
}

func TestSignAndVerify(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded")
		return
	}

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded")
		return
	}

	idBack, err := localMsp.DeserializeIdentity(serializedID)
	if err != nil {
		t.Fatalf("DeserializeIdentity should have succeeded")
		return
	}

	msg := []byte("foo")
	sig, err := id.Sign(msg)
	if err != nil {
		t.Fatalf("Sign should have succeeded")
		return
	}

	err = id.Verify(msg, sig)
	if err != nil {
		t.Fatalf("The signature should be valid")
		return
	}

	err = idBack.Verify(msg, sig)
	if err != nil {
		t.Fatalf("The signature should be valid")
		return
	}
}

func TestSignAndVerify_longMessage(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded")
		return
	}

	serializedID, err := id.Serialize()
	if err != nil {
		t.Fatalf("Serialize should have succeeded")
		return
	}

	idBack, err := localMsp.DeserializeIdentity(serializedID)
	if err != nil {
		t.Fatalf("DeserializeIdentity should have succeeded")
		return
	}

	msg := []byte("ABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFGABCDEFG")
	sig, err := id.Sign(msg)
	if err != nil {
		t.Fatalf("Sign should have succeeded")
		return
	}

	err = id.Verify(msg, sig)
	if err != nil {
		t.Fatalf("The signature should be valid")
		return
	}

	err = idBack.Verify(msg, sig)
	if err != nil {
		t.Fatalf("The signature should be valid")
		return
	}
}

func TestGetOU(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetSigningIdentity should have succeeded")
		return
	}

	assert.Equal(t, "COP", id.GetOrganizationalUnits()[0].OrganizationalUnitIdentifier)
}

func TestCertificationIdentifierComputation(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	chain, err := localMsp.(*bccspmsp).getCertificationChain(id.GetPublicVersion())
	assert.NoError(t, err)

	// Hash the chain
	hf, err := localMsp.(*bccspmsp).bccsp.GetHash(&bccsp.SHA256Opts{})
	assert.NoError(t, err)
	for i := 0; i < len(chain); i++ {
		hf.Write(chain[i].Raw)
	}
	sum := hf.Sum(nil)

	assert.Equal(t, sum, id.GetOrganizationalUnits()[0].CertifiersIdentifier)
}

func TestOUPolicyPrincipal(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	cid, err := localMsp.(*bccspmsp).getCertificationChainIdentifier(id.GetPublicVersion())
	assert.NoError(t, err)

	ou := &msp.OrganizationUnit{
		OrganizationalUnitIdentifier: "COP",
		MspIdentifier:                "DEFAULT",
		CertifiersIdentifier:         cid,
	}
	bytes, err := proto.Marshal(ou)
	assert.NoError(t, err)

	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ORGANIZATION_UNIT,
		Principal:               bytes,
	}

	err = id.SatisfiesPrincipal(principal)
	assert.NoError(t, err)
}

func TestOUPolicyPrincipalBadPath(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	ou := &msp.OrganizationUnit{
		OrganizationalUnitIdentifier: "COP",
		MspIdentifier:                "DEFAULT",
		CertifiersIdentifier:         nil,
	}
	bytes, err := proto.Marshal(ou)
	assert.NoError(t, err)

	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ORGANIZATION_UNIT,
		Principal:               bytes,
	}

	err = id.SatisfiesPrincipal(principal)
	assert.Error(t, err)

	ou = &msp.OrganizationUnit{
		OrganizationalUnitIdentifier: "COP",
		MspIdentifier:                "DEFAULT",
		CertifiersIdentifier:         []byte{0, 1, 2, 3, 4},
	}
	bytes, err = proto.Marshal(ou)
	assert.NoError(t, err)

	principal = &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ORGANIZATION_UNIT,
		Principal:               bytes,
	}

	err = id.SatisfiesPrincipal(principal)
	assert.Error(t, err)
}

func TestAdminPolicyPrincipal(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_ADMIN, MspIdentifier: "DEFAULT"})
	assert.NoError(t, err)

	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               principalBytes}

	err = id.SatisfiesPrincipal(principal)
	assert.NoError(t, err)
}

func TestAdminPolicyPrincipalFails(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	principalBytes, err := proto.Marshal(&msp.MSPRole{Role: msp.MSPRole_ADMIN, MspIdentifier: "DEFAULT"})
	assert.NoError(t, err)

	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               principalBytes}

	// remove the admin so validation will fail
	localMsp.(*bccspmsp).admins = make([]Identity, 0)

	err = id.SatisfiesPrincipal(principal)
	assert.Error(t, err)
}

func TestIdentityPolicyPrincipal(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	idSerialized, err := id.Serialize()
	assert.NoError(t, err)

	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_IDENTITY,
		Principal:               idSerialized}

	err = id.SatisfiesPrincipal(principal)
	assert.NoError(t, err)
}

func TestMSPOus(t *testing.T) {
	// Set the OUIdentifiers
	backup := localMsp.(*bccspmsp).ouIdentifiers
	defer func() { localMsp.(*bccspmsp).ouIdentifiers = backup }()

	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	localMsp.(*bccspmsp).ouIdentifiers = []*msp.FabricOUIdentifier{
		&msp.FabricOUIdentifier{
			OrganizationalUnitIdentifier: "COP",
			CertifiersIdentifier:         id.GetOrganizationalUnits()[0].CertifiersIdentifier,
		},
	}
	assert.NoError(t, localMsp.Validate(id.GetPublicVersion()))

	localMsp.(*bccspmsp).ouIdentifiers = []*msp.FabricOUIdentifier{
		&msp.FabricOUIdentifier{
			OrganizationalUnitIdentifier: "COP2",
			CertifiersIdentifier:         id.GetOrganizationalUnits()[0].CertifiersIdentifier,
		},
	}
	assert.Error(t, localMsp.Validate(id.GetPublicVersion()))

	localMsp.(*bccspmsp).ouIdentifiers = []*msp.FabricOUIdentifier{
		&msp.FabricOUIdentifier{
			OrganizationalUnitIdentifier: "COP",
			CertifiersIdentifier:         []byte{0, 1, 2, 3, 4},
		},
	}
	assert.Error(t, localMsp.Validate(id.GetPublicVersion()))
}

const othercert = `-----BEGIN CERTIFICATE-----
MIIDAzCCAqigAwIBAgIBAjAKBggqhkjOPQQDAjBsMQswCQYDVQQGEwJHQjEQMA4G
A1UECAwHRW5nbGFuZDEOMAwGA1UECgwFQmFyMTkxDjAMBgNVBAsMBUJhcjE5MQ4w
DAYDVQQDDAVCYXIxOTEbMBkGCSqGSIb3DQEJARYMQmFyMTktY2xpZW50MB4XDTE3
MDIwOTE2MDcxMFoXDTE4MDIxOTE2MDcxMFowfDELMAkGA1UEBhMCR0IxEDAOBgNV
BAgMB0VuZ2xhbmQxEDAOBgNVBAcMB0lwc3dpY2gxDjAMBgNVBAoMBUJhcjE5MQ4w
DAYDVQQLDAVCYXIxOTEOMAwGA1UEAwwFQmFyMTkxGTAXBgkqhkiG9w0BCQEWCkJh
cjE5LXBlZXIwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQlRSnAyD+ND6qmaRV7
AS/BPJKX5dZt3gBe1v/RewOpc1zJeXQNWACAk0ae3mv5u9l0HxI6TXJIAQSwJACu
Rqsyo4IBKTCCASUwCQYDVR0TBAIwADARBglghkgBhvhCAQEEBAMCBkAwMwYJYIZI
AYb4QgENBCYWJE9wZW5TU0wgR2VuZXJhdGVkIFNlcnZlciBDZXJ0aWZpY2F0ZTAd
BgNVHQ4EFgQUwHzbLJQMaWd1cpHdkSaEFxdKB1owgYsGA1UdIwSBgzCBgIAUYxFe
+cXOD5iQ223bZNdOuKCRiTKhZaRjMGExCzAJBgNVBAYTAkdCMRAwDgYDVQQIDAdF
bmdsYW5kMRAwDgYDVQQHDAdJcHN3aWNoMQ4wDAYDVQQKDAVCYXIxOTEOMAwGA1UE
CwwFQmFyMTkxDjAMBgNVBAMMBUJhcjE5ggEBMA4GA1UdDwEB/wQEAwIFoDATBgNV
HSUEDDAKBggrBgEFBQcDATAKBggqhkjOPQQDAgNJADBGAiEAuMq65lOaie4705Ol
Ow52DjbaO2YuIxK2auBCqNIu0gECIQCDoKdUQ/sa+9Ah1mzneE6iz/f/YFVWo4EP
HeamPGiDTQ==
-----END CERTIFICATE-----
`

func TestIdentityPolicyPrincipalFails(t *testing.T) {
	id, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	sid, err := NewSerializedIdentity("DEFAULT", []byte(othercert))
	assert.NoError(t, err)

	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_IDENTITY,
		Principal:               sid}

	err = id.SatisfiesPrincipal(principal)
	assert.Error(t, err)
}

var conf *msp.MSPConfig
var localMsp MSP
var mspMgr MSPManager

func TestMain(m *testing.M) {
	var err error
	mspDir, err := config.GetDevMspDir()
	if err != nil {
		fmt.Printf("Errog getting DevMspDir: %s", err)
		os.Exit(-1)
	}

	conf, err = GetLocalMspConfig(mspDir, nil, "DEFAULT")
	if err != nil {
		fmt.Printf("Setup should have succeeded, got err %s instead", err)
		os.Exit(-1)
	}

	localMsp, err = NewBccspMsp()
	if err != nil {
		fmt.Printf("Constructor for msp should have succeeded, got err %s instead", err)
		os.Exit(-1)
	}

	err = localMsp.Setup(conf)
	if err != nil {
		fmt.Printf("Setup for msp should have succeeded, got err %s instead", err)
		os.Exit(-1)
	}

	mspMgr = NewMSPManager()
	err = mspMgr.Setup([]MSP{localMsp})
	if err != nil {
		fmt.Printf("Setup for msp manager should have succeeded, got err %s instead", err)
		os.Exit(-1)
	}

	retVal := m.Run()
	os.Exit(retVal)
}

func getIdentity(t *testing.T, path string) Identity {
	mspDir, err := config.GetDevMspDir()
	assert.NoError(t, err)

	pems, err := getPemMaterialFromDir(filepath.Join(mspDir, path))
	assert.NoError(t, err)

	id, _, err := localMsp.(*bccspmsp).getIdentityFromConf(pems[0])
	assert.NoError(t, err)

	return id
}

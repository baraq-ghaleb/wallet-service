package packers

import (
	"encoding/json"
	"fmt"
	"testing"

	core "github.com/iden3/go-iden3-core"
	"github.com/lastingasset/wallet-service/go-circuits"
	"github.com/lastingasset/wallet-service/go-jwz"
	"github.com/lastingasset/wallet-service/iden3comm/mock"
	"github.com/lastingasset/wallet-service/iden3comm/protocol"
	"github.com/stretchr/testify/assert"
)

func TestZKPPacker_Pack(t *testing.T) {

	mockedProvingMethod := &mock.ProvingMethodGroth16AuthV2{ProvingMethodAlg: jwz.ProvingMethodAlg{Alg: "groth16-mock", CircuitID: "authV2"}}

	jwz.RegisterProvingMethod(mockedProvingMethod.ProvingMethodAlg, func() jwz.ProvingMethod {
		return mockedProvingMethod
	})

	mockVerificationParam := make(map[jwz.ProvingMethodAlg]VerificationParams)
	mockVerificationParam[mockedProvingMethod.ProvingMethodAlg] = NewVerificationParams([]byte(""), mock.VerifyState)

	mockProvingParamMap := make(map[jwz.ProvingMethodAlg]ProvingParams)
	mockProvingParamMap[mockedProvingMethod.ProvingMethodAlg] = NewProvingParams(mock.PrepareAuthInputs, []byte{}, []byte{})

	p := NewZKPPacker(mockProvingParamMap, mockVerificationParam)

	msgBytes := []byte(`{"type":"https://iden3-communication.io/authorization/1.0/response","from":"did:iden3:polygon:mumbai:x4jcHP4XHTK3vX58AHZPyHE8kYjneyE6FZRfz7K29","body":{"scope":[{"type":"zeroknowledge","circuit_id":"auth","pub_signals":["1","18311560525383319719311394957064820091354976310599818797157189568621466950811","323416925264666217617288569742564703632850816035761084002720090377353297920"],"proof_data":{"pi_a":["11130843150540789299458990586020000719280246153797882843214290541980522375072","1300841912943781723022032355836893831132920783788455531838254465784605762713","1"],"pi_b":[["20615768536988438336537777909042352056392862251785722796637590212160561351656","10371144806107778890538857700855108667622042215096971747203105997454625814080"],["19598541350804478549141207835028671111063915635580679694907635914279928677812","15264553045517065669171584943964322117397645147006909167427809837929458012913"],["1","0"]],"pi_c":["16443309279825508893086251290003936935077348754097470818523558082502364822049","2984180227766048100510120407150752052334571876681304999595544138155611963273","1"],"protocol":""}}]}}`)

	identifier := "did:iden3:polygon:mumbai:x4jcHP4XHTK3vX58AHZPyHE8kYjneyE6FZRfz7K29"

	senderDID, err := core.ParseDID(identifier)
	assert.NoError(t, err)

	b, err := p.Pack(msgBytes, ZKPPackerParams{
		SenderID:         senderDID,
		ProvingMethodAlg: jwz.ProvingMethodAlg{Alg: "groth16-mock", CircuitID: "authV2"},
	})
	assert.Nil(t, err)

	token, err := jwz.Parse(string(b))
	assert.Nil(t, err)

	fmt.Println(string(b))
	outs := circuits.AuthV2PubSignals{}
	err = token.ParsePubSignals(&outs)
	assert.Nil(t, err)

	didFromToken, err := core.ParseDIDFromID(*outs.UserID)
	assert.Nil(t, err)

	assert.EqualValues(t, senderDID.String(), didFromToken.String())

}

func TestPlainMessagePacker_Unpack(t *testing.T) {
	mockedProvingMethod := &mock.ProvingMethodGroth16AuthV2{ProvingMethodAlg: jwz.ProvingMethodAlg{Alg: "groth16-mock", CircuitID: "authV2"}}

	jwz.RegisterProvingMethod(mockedProvingMethod.ProvingMethodAlg, func() jwz.ProvingMethod {
		return mockedProvingMethod
	})

	mockVerificationParam := make(map[jwz.ProvingMethodAlg]VerificationParams)
	mockVerificationParam[mockedProvingMethod.ProvingMethodAlg] = NewVerificationParams([]byte(""), mock.VerifyState)

	mockProvingParamMap := make(map[jwz.ProvingMethodAlg]ProvingParams)
	mockProvingParamMap[mockedProvingMethod.ProvingMethodAlg] = NewProvingParams(mock.PrepareAuthInputs, []byte{}, []byte{})

	jwz.RegisterProvingMethod(mockedProvingMethod.ProvingMethodAlg, func() jwz.ProvingMethod {
		return mockedProvingMethod
	})
	p := NewZKPPacker(mockProvingParamMap, mockVerificationParam)

	msgZKP := []byte(`eyJhbGciOiJncm90aDE2LW1vY2siLCJjaXJjdWl0SWQiOiJhdXRoVjIiLCJjcml0IjpbImNpcmN1aXRJZCJdLCJ0eXAiOiJhcHBsaWNhdGlvbi9pZGVuMy16a3AtanNvbiJ9.eyJ0eXBlIjoiaHR0cHM6Ly9pZGVuMy1jb21tdW5pY2F0aW9uLmlvL2F1dGhvcml6YXRpb24vMS4wL3Jlc3BvbnNlIiwiZnJvbSI6ImRpZDppZGVuMzpwb2x5Z29uOm11bWJhaTp4NGpjSFA0WEhUSzN2WDU4QUhaUHlIRThrWWpuZXlFNkZaUmZ6N0syOSIsImJvZHkiOnsic2NvcGUiOlt7InR5cGUiOiJ6ZXJva25vd2xlZGdlIiwiY2lyY3VpdF9pZCI6ImF1dGgiLCJwdWJfc2lnbmFscyI6WyIxIiwiMTgzMTE1NjA1MjUzODMzMTk3MTkzMTEzOTQ5NTcwNjQ4MjAwOTEzNTQ5NzYzMTA1OTk4MTg3OTcxNTcxODk1Njg2MjE0NjY5NTA4MTEiLCIzMjM0MTY5MjUyNjQ2NjYyMTc2MTcyODg1Njk3NDI1NjQ3MDM2MzI4NTA4MTYwMzU3NjEwODQwMDI3MjAwOTAzNzczNTMyOTc5MjAiXSwicHJvb2ZfZGF0YSI6eyJwaV9hIjpbIjExMTMwODQzMTUwNTQwNzg5Mjk5NDU4OTkwNTg2MDIwMDAwNzE5MjgwMjQ2MTUzNzk3ODgyODQzMjE0MjkwNTQxOTgwNTIyMzc1MDcyIiwiMTMwMDg0MTkxMjk0Mzc4MTcyMzAyMjAzMjM1NTgzNjg5MzgzMTEzMjkyMDc4Mzc4ODQ1NTUzMTgzODI1NDQ2NTc4NDYwNTc2MjcxMyIsIjEiXSwicGlfYiI6W1siMjA2MTU3Njg1MzY5ODg0MzgzMzY1Mzc3Nzc5MDkwNDIzNTIwNTYzOTI4NjIyNTE3ODU3MjI3OTY2Mzc1OTAyMTIxNjA1NjEzNTE2NTYiLCIxMDM3MTE0NDgwNjEwNzc3ODg5MDUzODg1NzcwMDg1NTEwODY2NzYyMjA0MjIxNTA5Njk3MTc0NzIwMzEwNTk5NzQ1NDYyNTgxNDA4MCJdLFsiMTk1OTg1NDEzNTA4MDQ0Nzg1NDkxNDEyMDc4MzUwMjg2NzExMTEwNjM5MTU2MzU1ODA2Nzk2OTQ5MDc2MzU5MTQyNzk5Mjg2Nzc4MTIiLCIxNTI2NDU1MzA0NTUxNzA2NTY2OTE3MTU4NDk0Mzk2NDMyMjExNzM5NzY0NTE0NzAwNjkwOTE2NzQyNzgwOTgzNzkyOTQ1ODAxMjkxMyJdLFsiMSIsIjAiXV0sInBpX2MiOlsiMTY0NDMzMDkyNzk4MjU1MDg4OTMwODYyNTEyOTAwMDM5MzY5MzUwNzczNDg3NTQwOTc0NzA4MTg1MjM1NTgwODI1MDIzNjQ4MjIwNDkiLCIyOTg0MTgwMjI3NzY2MDQ4MTAwNTEwMTIwNDA3MTUwNzUyMDUyMzM0NTcxODc2NjgxMzA0OTk5NTk1NTQ0MTM4MTU1NjExOTYzMjczIiwiMSJdLCJwcm90b2NvbCI6IiJ9fV19fQ.eyJwcm9vZiI6eyJwaV9hIjpudWxsLCJwaV9iIjpudWxsLCJwaV9jIjpudWxsLCJwcm90b2NvbCI6Imdyb3RoMTYifSwicHViX3NpZ25hbHMiOlsiMTkyMjkwODQ4NzM3MDQ1NTAzNTcyMzI4ODcxNDI3NzQ2MDU0NDIyOTczMzcyMjkxNzY1NzkyMjkwMTEzNDIwOTE1OTQxNzQ5NzciLCI2MTEwNTE3NzY4MjQ5NTU5MjM4MTkzNDc3NDM1NDU0NzkyMDI0NzMyMTczODY1NDg4OTAwMjcwODQ5NjI0MzI4NjUwNzY1NjkxNDk0IiwiMTI0MzkwNDcxMTQyOTk2MTg1ODc3NDIyMDY0NzYxMDcyNDI3Mzc5ODkxODQ1Nzk5MTQ4NjAzMTU2NzI0NDEwMDc2NzI1OTIzOTc0NyJdfQ`)
	iden3msg, err := p.Unpack(msgZKP)
	assert.NoError(t, err)
	msgBytes, err := json.Marshal(iden3msg)
	assert.Nil(t, err)
	var authResponse protocol.AuthorizationResponseMessage
	err = json.Unmarshal(msgBytes, &authResponse)
	assert.Nil(t, err)

	assert.Equal(t, protocol.AuthorizationResponseMessageType, authResponse.Type)
	assert.Len(t, authResponse.Body.Scope, 1)

}

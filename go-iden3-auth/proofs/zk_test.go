package proofs

import (
	"os"
	"testing"

	"github.com/iden3/go-rapidsnark/types"
	"github.com/lastingasset/wallet-service/go-circuits"
	"github.com/lastingasset/wallet-service/iden3comm/protocol"
	"github.com/stretchr/testify/assert"
)

func TestVerifyProof(t *testing.T) {

	var err error
	proofMessage := protocol.ZeroKnowledgeProofResponse{
		ZKProof: types.ZKProof{
			Proof: &types.ProofData{
				A: []string{
					"9517112492422486418344671523752691163637612305590571624363668885796911150333",
					"8855938450276251202387073646943136306720422603123854769235151758541434807968",
					"1",
				},
				B: [][]string{
					{
						"18880568320884466923930564925565727939067628655227999252296084923782755860476",
						"8724893415197458543695192455798597402395044930214471497778888748319129905479",
					},
					{
						"9807559381041464075347519433137353143151890330916363861193891037865993320923",
						"6995202980453256069532771522391679223085808426805857698209331232672383046019",
					},
					{
						"1",
						"0",
					}},
				C: []string{
					"16453660244095377174525331937765624986258178472608723119429308977591704509298",
					"7523187725705152586426891868747265746542072544935310991409893207335385519512",
					"1",
				},
				Protocol: "groth16",
			},
			PubSignals: []string{
				"1",
				"25054465935916343733470065977393556898165832783214621882239050035846517250",
				"10",
				"25054465935916343733470065977393556898165832783214621882239050035846517250",
				"7120485770008490579908343167068999806468056401802904713650068500000641772574",
				"1",
				"7120485770008490579908343167068999806468056401802904713650068500000641772574",
				"1671543597",
				"336615423900919464193075592850483704600",
				"0",
				"17002437119434618783545694633038537380726339994244684348913844923422470806844",
				"0",
				"5",
				"840",
				"120",
				"340",
				"509",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
				"0",
			},
		},
	}
	proofMessage.CircuitID = string(circuits.AtomicQueryMTPV2CircuitID)

	verificationKey, err := os.ReadFile("../testdata/credentialAtomicQueryMTPV2.json")
	assert.NoError(t, err)

	proofMessage.ID = 1

	err = VerifyProof(proofMessage, verificationKey)
	assert.Nil(t, err)
}

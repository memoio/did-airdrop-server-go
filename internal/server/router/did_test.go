package router

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	sig  = "0xbc90763fd92811bd890b153f3ab3d393a5aeebc207d2561e014c8887ec60c4b92cc5b1a98c1a3b8c0c84d172cfc44ef6844ba829a45cad940e9ab1dc335269d81b"
	hash = "0x1b2a7fe5414ff93753cbf4418a49398436b863d3f97a4c101d218c596bc9b5e3"
	msg  = "0x637265617465444944373864623561626564386432653333646234626231346632333263396563373938303361306263373437316639313462623030643665303965313161383236364563647361536563703235366b31566572696669636174696f6e4b657932303139c145a262565c746fc1596ba92b85e43f006b95660000000000000000"
)

func TestSignature(t *testing.T) {
	SigByte, err := hexutil.Decode(sig)
	if err != nil {
		t.Error(err)
		return
	}
	msgByte, err := hexutil.Decode(msg)
	if err != nil {
		t.Error(err)
		return
	}

	hashByte := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(string(msgByte)), string(msgByte))))
	SigByte[len(SigByte)-1] %= 27

	pk, err := crypto.SigToPub(hashByte, SigByte)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(crypto.PubkeyToAddress(*pk).Hex())
}

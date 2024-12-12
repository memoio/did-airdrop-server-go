package did

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	klog "github.com/go-kratos/kratos/v2/log"
)

const (
	privatekey = "9b4fc2a14cbc63a0d338377413ca72949cbb2fd5be1b08844b4b5e332597d91e"
	publickey  = "0x03ecc373891778bed36426ddcd682bf1e0b5a99a8d8534be05a000ddc4faaccea0"
	did        = "did:memo:3e237e60f5d68a5f1a73bc108dd9150a5cdf439754463acc2aa0962876ba4ce7"
	address    = "0x47D4f617A654337AFB121F455629fF7d92b670eA"
)

func TestCreateDID(t *testing.T) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.CreateDIDByAaddress(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(did.String())
}

func TestGetNonce(t *testing.T) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("endpoint: %s, ins: %s, proxy: %s, account: %s", memoDID.Controller.EndPoint(), memoDID.Controller.Instance().String(), memoDID.Controller.Proxy().String(), memoDID.Controller.Account().String())

	nonce, err := memoDID.Controller.GetNonce(did)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(nonce)
}

func TestRegisterDID(t *testing.T) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	pubKeyBytes, err := hexutil.Decode(publickey)
	if err != nil {
		t.Fatal(err)
	}

	unsig, err := memoDID.getCreateDIDHashPubkey(did, publickey, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("did:", did)

	unSigByte, err := hexutil.Decode(unsig)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		t.Fatal(err)
	}

	sig, err := crypto.Sign(unSigByte, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.RegisterDID(did, "memo", pubKeyBytes, sig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(did)
}

func TestGetDIDInfo(t *testing.T) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.GetDIDInfo(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(did)
}

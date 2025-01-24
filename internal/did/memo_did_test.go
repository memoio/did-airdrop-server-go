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
	sk1        = "cf9f8e55aaf30ab82d6daec06248cdfb1a761db68bc5ac30b230c4beaa48e3e4"
	publickey  = "0x03ecc373891778bed36426ddcd682bf1e0b5a99a8d8534be05a000ddc4faaccea0"
	address    = "0x47D4f617A654337AFB121F455629fF7d92b670eA"
	address1   = "0x594CE7BA907710f5647C6ec58db168B0a2686de4"
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

	did, err := memoDID.CreateDIDByAddress(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(did.String())
}

func TestGetNonce(t *testing.T) {
	addr := "0x9Cf73Ad845075227566EdC8503DF843D529eD3b9"
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("product", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.CreateDIDByAddress(addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("endpoint: %s, ins: %s, proxy: %s, account: %s", memoDID.Controller.EndPoint(), memoDID.Controller.Instance().String(), memoDID.Controller.Proxy().String(), memoDID.Controller.Account().String())

	nonce, err := memoDID.Controller.GetNonce(did.Identifier)
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

	did, err := memoDID.CreateDIDByAddress(address1)
	if err != nil {
		t.Fatal(err)
	}

	nonce, err := memoDID.Controller.GetNonce(did.String())
	if err != nil {
		t.Fatal(err)
	}

	unsig, err := memoDID.getCreateDIDHashPubkey(did.String(), publickey, nonce)
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

	didStr, err := memoDID.RegisterDIDByPublic(publickey, sig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(didStr)
}

func TestGetHashByAddress(t *testing.T) {
	addr := "0x2EB682d7387d65a785EbF983987E5977dc6700D4"

	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.CreateDIDByAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("did:", did.String())

	nonce, err := memoDID.Controller.GetNonce(did.Identifier)
	if err != nil {
		t.Fatal(err)
	}

	unsig, err := memoDID.getCreateDIDHashByAddress(did.Identifier, addr, nonce)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(unsig)

}

func TestRegisterDIDByAddress(t *testing.T) {
	addr := "0xc145A262565C746fc1596ba92b85E43F006b9566"

	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.CreateDIDByAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("did:", did.String())

	didStr, err := memoDID.RegisterDIDByAddressByAdmin(addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(didStr)
}

func TestGetDIDInfo(t *testing.T) {
	addr := "0x013A08061C08E3852aBb921F305B304e0C165eB2"
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, number, err := memoDID.GetDIDInfo(addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(did, number)
}

func TestGetDIDExist(t *testing.T) {
	addr := "0xc145A262565C746fc1596ba92b85E43F006b9566"
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	memoDID, err := NewMemoDID("dev", klog.NewHelper(logger))
	if err != nil {
		t.Fatal(err)
	}

	did, err := memoDID.GetDIDExist(addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(did)
}

func TestParseMfileDID(t *testing.T) {
	
}
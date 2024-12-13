package contract

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"golang.org/x/xerrors"

	com "github.com/memoio/contractsv2/common"
	"github.com/memoio/contractsv2/go_contracts/erc"
	inst "github.com/memoio/contractsv2/go_contracts/instance"
)

var (
	privatekey    = "9b4fc2a14cbc63a0d338377413ca72949cbb2fd5be1b08844b4b5e332597d91e"
	address       = "0x47D4f617A654337AFB121F455629fF7d92b670eA"
	publickey     = "0x03ecc373891778bed36426ddcd682bf1e0b5a99a8d8534be05a000ddc4faaccea0"
	sk1           = "cf9f8e55aaf30ab82d6daec06248cdfb1a761db68bc5ac30b230c4beaa48e3e4"
	walletPrivate = "7a71499718e02f0bef77a9a34cd8eca62a3c3964caf5c599c16855e6774953a1"
	message       = "0x2e91efb096c737d898be749479e99f5b4ce386127ee616dfb206a218c77304e5"
)

var (
	errorSig     = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _ = abi.NewType("string", "", nil)
)

func TestCreateSK(t *testing.T) {
	sk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(sk)
	t.Log(hexutil.Encode(privateKeyBytes))
}

func TestGetPublicKey(t *testing.T) {
	privateKey, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.CompressPubkey(publicKeyECDSA)

	t.Log(hexutil.Encode(publicKeyBytes))
}

func TestSignatureMsg(t *testing.T) {
	privateKey, err := crypto.HexToECDSA(sk1)
	if err != nil {
		t.Fatal(err)
	}

	msg := hexutil.MustDecode("0x1322499d95e9b59914a793bef45d8d9e979c5302a74e351dbe6f57a62f4cf243")
	sig, err := crypto.Sign(msg, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	t.Log(hexutil.Encode(sig))
}
func TestGetAddress(t *testing.T) {
	privateKey, err := crypto.HexToECDSA(sk1)
	if err != nil {
		t.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	t.Log(address.Hex())
}

func TestGetBalance(t *testing.T) {
	balance, balanceErc20, err := getBalance("dev", walletPrivate)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance, balanceErc20)
}

func TestTransfer(t *testing.T) {
	var pledge = big.NewInt(1000000000000000000)
	var amount = new(big.Int).Mul(pledge, big.NewInt(1000))

	ethAddr := common.HexToAddress(address)

	transferMemo("dev", walletPrivate, ethAddr, amount)
	transferEth("dev", walletPrivate, ethAddr, amount)
}

func getBalance(chain string, sk string) (*big.Int, *big.Int, error) {
	instanceAddr, endpoint := com.GetInsEndPointByChain(chain)

	client, err := ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		return nil, nil, err
	}

	instanceIns, err := inst.NewInstance(instanceAddr, client)
	if err != nil {
		return nil, nil, err
	}

	tokenAddr, err := instanceIns.Instances(&bind.CallOpts{}, com.TypeERC20)
	if err != nil {
		return nil, nil, err
	}

	tokenIns, err := erc.NewERC20(tokenAddr, client)
	if err != nil {
		return nil, nil, err
	}

	sk0, err := crypto.HexToECDSA(sk)
	if err != nil {
		return nil, nil, err
	}
	balance, err := client.BalanceAt(context.TODO(), crypto.PubkeyToAddress(sk0.PublicKey), nil)
	if err != nil {
		return nil, nil, err
	}

	balanceErc20, err := tokenIns.BalanceOf(&bind.CallOpts{}, crypto.PubkeyToAddress(sk0.PublicKey))
	if err != nil {
		return nil, nil, err
	}
	return balance, balanceErc20, nil
}

func transferMemo(chain string, fromSK string, to common.Address, amount *big.Int) error {
	instanceAddr, endpoint := com.GetInsEndPointByChain(chain)

	client, err := ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		return err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	instanceIns, err := inst.NewInstance(instanceAddr, client)
	if err != nil {
		return err
	}

	tokenAddr, err := instanceIns.Instances(&bind.CallOpts{}, com.TypeERC20)
	if err != nil {
		return err
	}

	tokenIns, err := erc.NewERC20(tokenAddr, client)
	if err != nil {
		return err
	}

	auth, err := com.MakeAuth(chainID, fromSK)
	if err != nil {
		return err
	}
	tx, err := tokenIns.Transfer(auth, to, amount)
	if err != nil {
		return err
	}
	return CheckTx(endpoint, auth.From, tx, "transfer memo")
}

func transferEth(chain string, fromSK string, to common.Address, amount *big.Int) error {
	_, endpoint := com.GetInsEndPointByChain(chain)

	client, err := ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		return err
	}

	privateKey, err := crypto.HexToECDSA(fromSK)
	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return err
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		chainID = big.NewInt(985)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return CheckTx(endpoint, crypto.PubkeyToAddress(privateKey.PublicKey), signedTx, "transfer eth")
}

// CheckTx check whether transaction is successful through receipt
func CheckTx(endPoint string, from common.Address, tx *types.Transaction, name string) error {
	var receipt *types.Receipt

	t := checkTxSleepTime
	for i := 0; i < 30; i++ {
		time.Sleep(time.Duration(t) * time.Second)
		receipt = com.GetTransactionReceipt(endPoint, tx.Hash())
		if receipt != nil {
			break
		}
		t = nextBlockTime
	}

	if receipt == nil {
		return xerrors.Errorf("%s: cann't get transaction(%s) receipt, not packaged", name, tx.Hash())
	}

	// 0 means fail
	if receipt.Status == 0 {
		if receipt.GasUsed != receipt.CumulativeGasUsed {
			return xerrors.Errorf("%s: transaction(%s) exceed gas limit", name, tx.Hash())
		}
		reason, err := getErrorReason(context.TODO(), endPoint, from, tx)
		if err != nil {
			return xerrors.Errorf("%s: transaction(%s) mined but execution failed: %s", name, tx.Hash(), err.Error())
		}
		return xerrors.Errorf("%s: transaction(%s) revert(%s)", name, tx.Hash(), reason)
	}
	return nil
}

func getErrorReason(ctx context.Context, endpoint string, from common.Address, tx *types.Transaction) (string, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return "", err
	}
	defer client.Close()

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	res, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		log.Println(res)
		return "", err
	}
	return unpackError(res)
}

func unpackError(result []byte) (string, error) {
	log.Println(string(result))
	if !bytes.Equal(result[:4], errorSig) {
		return "<tx result not Error(string)>", errors.New("TX result not of type Error(string)")
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return "<invalid tx result>", errors.Wrap(err, "unpacking revert reason")
	}
	return vs[0].(string), nil
}

func TestGetNonce(t *testing.T) {
	// logger := klog.With(klog.NewStdLogger(os.Stdout),
	// 	"ts", klog.DefaultTimestamp,
	// 	"caller", klog.DefaultCaller,
	// )
	// controller, err := NewController("dev", klog.NewHelper(logger))
	// if err != nil {
	// 	t.Fatal(err)
	// }

}

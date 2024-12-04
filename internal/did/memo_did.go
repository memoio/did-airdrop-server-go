package did

import (
	"github.com/nuts-foundation/did-ockam"
	"golang.org/x/xerrors"
)

type MemoDID struct {
	Method string

	Identifier string

	Identifiers []string
}

func ParseMemoDID(didStr string) (*MemoDID, error) {
	did, err := did.Parse(didStr)
	if err != nil {
		return nil, err
	}
	if did.IsURL() {
		return nil, xerrors.Errorf("%s is did url", didStr)
	}
	if did.Method != "memo" {
		return nil, xerrors.Errorf("unsupported method %s", did.Method)
	}
	if len(did.IDStrings) > 1 {
		// TODO: check didString[2:len(didStrings)-1] ==? {chain id}
		return nil, xerrors.Errorf("TODO: support chain id")
	}
	if isNot32ByteHex(did.IDStrings[len(did.IDStrings)-1]) {
		return nil, xerrors.Errorf("%s is not 32 byte hex string", did.IDStrings[len(did.IDStrings)-1])
	}
	return &MemoDID{
		Method:      "memo",
		Identifier:  did.ID,
		Identifiers: did.IDStrings,
	}, nil
}

func isNot32ByteHex(s string) bool {
	if len(s) != 64 {
		return true
	}

	for _, b := range s {
		if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')) {
			return true
		}
	}

	return false
}

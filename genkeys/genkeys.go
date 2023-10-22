package genkeys

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"log"
)

//goland:noinspection GoUnusedExportedFunction
func GenerateKey() (wif string, address string) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Println("generate private key error:", err.Error())
		return "", ""
	}
	if len(privateKey.Serialize()) != 32 {
		for {
			privateKey, err := btcec.NewPrivateKey()
			if err != nil {
				continue
			}
			if len(privateKey.Serialize()) == 32 {
				break
			}
		}
	}
	return hex.EncodeToString(privateKey.Serialize()), addr.PubkeyToAddress(privateKey.ToECDSA().PublicKey).String()
}

//goland:noinspection GoUnusedExportedFunction
func CreateAddressBySeed(seed []byte) (string, error) {
	if len(seed) != 32 {
		return "", fmt.Errorf("seed len=[%d] is not equal 32", len(seed))
	}

	privateKey, _ := btcec.PrivKeyFromBytes(seed)
	if privateKey == nil {
		return "", errors.New("private key is nil ptr")
	}
	return addr.PubkeyToAddress(privateKey.ToECDSA().PublicKey).String(), nil
}

//goland:noinspection GoUnusedExportedFunction
func AddressB58ToHex(b58 string) (string, error) {
	a, err := addr.Base58ToAddress(b58)
	if err != nil {
		return "", err
	}
	return a.Hex(), nil
}

//goland:noinspection GoUnusedExportedFunction
func AddressHexToB58(hexAddress string) string {
	a := addr.HexToAddress(hexAddress)
	return a.String()
}

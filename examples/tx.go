package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/SomniSom/trx-sign-go/grpcs"
	"github.com/SomniSom/trx-sign-go/sign"
	"github.com/btcsuite/btcutil/base58"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/pkg/errors"
	"math/big"
)

func TransferTrx(from, to, key string, amount int64) error {
	c, err := grpcs.NewClient("54.168.218.95:50051")
	if err != nil {
		return errors.Wrap(err, "new grpc client on 54.168.218.95:50051")
	}
	tx, err := c.Transfer(from, to, amount)
	if err != nil {
		return errors.Wrap(err, "create transfer")
	}
	signTx, err := sign.SignTransaction(tx.Transaction, key)
	if err != nil {
		return err
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		return err
	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
	return nil
}

func GetBalance(wallet string) error {
	c, err := grpcs.NewClient("3.225.171.164:50051")
	if err != nil {
		return errors.Wrap(err, "new grpc client on 3.225.171.164:50051")
	}
	acc, err := c.GetTrxBalance(wallet)
	if err != nil {
		return errors.Wrap(err, "get tron balance on wallet")
	}
	d, err := json.Marshal(acc)
	if err != nil {
		return errors.Wrap(err, "marshal data")
	}
	fmt.Println(string(d))
	fmt.Println(acc.GetBalance())
	return nil
}

func GetTrc20Balance(addr, contract string) error {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051")
	if err != nil {
		return err
	}
	amount, err := c.GetTrc20Balance(addr, contract)
	if err != nil {
		return err
	}
	fmt.Println(amount.String())
	return nil
}

func TransferTrc20(from, to, contract, key string, amt, fee int64) error {
	c, err := grpcs.NewClient("54.168.218.95:50051")
	if err != nil {
		return err
	}
	if fee == 0 {
		fee = 100000000
	}
	amount := big.NewInt(amt)
	amount = amount.Mul(amount, big.NewInt(1000000000000000000))
	tx, err := c.TransferTrc20(from, to, contract, amount, fee)
	signTx, err := sign.SignTransaction(tx.Transaction, key)
	if err != nil {
		return err
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		return err

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
	return nil
}

func TransferTrc10(fromAddr, toAddr, tokenID, key string, amount int64) error {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		return err
	}
	from, _ := addr.Base58ToAddress(fromAddr)
	to, _ := addr.Base58ToAddress(toAddr)
	//tokenID := "1000016"
	tx, err := c.GRPC.TransferAsset(from.String(), to.String(), tokenID, amount)
	signTx, err := sign.SignTransaction(tx.Transaction, key)
	if err != nil {
		return err
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		return err

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
	return nil
}

func GetTrc10Balance(addr, assetID string) error {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051")
	if err != nil {
		return err
	}
	amount, err := c.GetTrc10Balance(addr, assetID)
	if err != nil {
		return err
	}
	fmt.Println(amount)
	return nil
}

func DecodeCheck(input string) ([]byte, error) {
	decodeCheck := base58.Decode(input)
	if len(decodeCheck) == 0 {
		return nil, fmt.Errorf("b58 decode %s error", input)
	}
	if len(decodeCheck) < 4 {
		return nil, fmt.Errorf("b58 check error")
	}

	decodeData := decodeCheck[:len(decodeCheck)-4]

	h256h0 := sha256.New()
	h256h0.Write(decodeData)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	if h1[0] == decodeCheck[len(decodeData)] &&
		h1[1] == decodeCheck[len(decodeData)+1] &&
		h1[2] == decodeCheck[len(decodeData)+2] &&
		h1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}
	return nil, fmt.Errorf("b58 check error")
}

func GetBlock() error {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		return err
	}
	lb, err := c.GRPC.GetNowBlock()
	if err != nil {
		return err
	}
	fmt.Println(lb.BlockHeader.RawData.Number)
	fmt.Println(hex.EncodeToString(lb.Blockid))
	return nil
}

func GetTxByTxid(txID string) (*core.TransactionInfo, error) {
	c, err := grpcs.NewClient("grpc.trongrid.io:50051")
	if err != nil {
		return nil, err
	}
	ti, err := c.GRPC.GetTransactionInfoByID(txID)
	if err != nil {
		return nil, err
	}
	//fee := ti.Receipt.GetEnergyFee() + ti.Receipt.GetNetFee()
	//fmt.Println(fee)
	return ti, nil
}

//func GetTransaction(txID string) error {
//	c, err := grpcs.NewClient("3.225.171.164:50051")
//	if err != nil {
//		return err
//	}
//
//	txInfo, err := c.GRPC.GetTransactionByID(txID)
//	if err != nil {
//		return err
//	}
//
//	r, err := c.GRPC.GetTransactionInfoByID(txID)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func DecodeMessage(t *testing.T) error {
//	data := "CMN5oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABlSYXRlIHN0YWxlIG9yIG5vdCBhIHN5bnRoAAAAAAAAAA=="
//	d, err := base64.StdEncoding.DecodeString(data)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(hex.EncodeToString(d))
//	dd, _ := hex.DecodeString("1952617465207374616c65206f72206e6f7420612073796e746800000000000000")
//	fmt.Println(string(dd))
//}

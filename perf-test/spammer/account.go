package spammer

import (
	"crypto/ecdsa"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	PrivKey string
	Nonce   uint64
	Mutex   sync.Mutex
}

func (a *Account) PrivateKey() *ecdsa.PrivateKey {
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)

	privateKey, err := crypto.HexToECDSA(a.PrivKey)
	if err != nil {
		logger.Error("Failed to load private key: %v", err)
	}
	return privateKey
}

func (a *Account) Address() *common.Address {
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)

	publicKey := a.PrivateKey().Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Error("Error casting public key to ECDSA")
	}
	addr := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &addr
}

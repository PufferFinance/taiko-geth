package spammer

import (
	"context"
	"math/big"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/ethereum/go-ethereum/core/types"
)

type Spammer struct {
	ctx               context.Context
	cancel            context.CancelFunc
	wg                *sync.WaitGroup
	client            *EthClient
	logger            *log.Logger
	accounts          []*Account
	maxTxsPerAccount  uint64
	prefundedAccounts []*Account
}

func New(url string, chainID *big.Int, logger *log.Logger, accounts []*Account, maxTxsPerAccount uint64, prefundedAccounts []*Account) *Spammer {
	var wg = new(sync.WaitGroup)

	ctx, cancel := context.WithCancel(context.Background())

	client, err := NewEthClient(ctx, url, chainID, logger)
	if err != nil {
		logger.Error("Failed to connect to the Ethereum client", "error", err)
	}

	return &Spammer{
		ctx:               ctx,
		cancel:            cancel,
		wg:                wg,
		client:            client,
		accounts:          accounts,
		logger:            logger,
		maxTxsPerAccount:  maxTxsPerAccount,
		prefundedAccounts: prefundedAccounts,
	}
}

func (s *Spammer) Start() {

	for _, account := range s.prefundedAccounts {
		nonce, err := s.client.GetNonce(account)

		if err != nil {
			s.logger.Error("Failed to get nonce", "error", err)
			continue
		}

		account.Mutex.Lock()
		account.Nonce = nonce
		account.Mutex.Unlock()
	}

	s.sendTxs()
}

func (s *Spammer) sendTxs() {
	for _, account := range s.prefundedAccounts {
		s.wg.Add(1)

		go func(account *Account) {
			defer s.wg.Done()

			for i := uint64(0); i < s.maxTxsPerAccount; i++ {
				legacyTx := &types.LegacyTx{
					Nonce:    account.Nonce,
					To:       account.Address(),
					Value:    big.NewInt(20000000000000),
					Gas:      21000,
					GasPrice: big.NewInt(24000000000),
				}
				tx := types.NewTx(legacyTx)
				// Sign and send the tx
				_, err := s.client.SendTx(account, tx)
				if err != nil {
					s.logger.Error("Failed to send tx", "error", err)
					continue
				}
				account.Mutex.Lock()
				account.Nonce++
				account.Mutex.Unlock()
			}
		}(account)
	}

	s.wg.Wait()
}

// PrefundAccounts is a function that prefunds the accounts with ETH
func (s *Spammer) PrefundAccounts() {
	sender := s.accounts[0]
	for _, account := range s.prefundedAccounts {
		// Get nonce
		nonce, err := s.client.GetNonce(sender)
		if err != nil {
			s.logger.Error("Failed to get nonce", "error", err, "account", account.Address())
			return
		}
		s.logger.Info("Got nonce for account", "nonce", nonce, "account", sender.Address())

		legacyTx := &types.LegacyTx{
			Nonce:    nonce,
			To:       account.Address(),
			Value:    big.NewInt(10000000000000000),
			Gas:      21000,
			GasPrice: big.NewInt(24000000000),
		}

		tx := types.NewTx(legacyTx)

		// Sign and send the tx
		signedTx, err := s.client.SendTx(sender, tx)
		if err != nil {
			s.logger.Error("Failed to send tx", "error", err, "nonce", nonce)
			return
		}

		s.logger.Info("Transaction submitted",
			"hash", signedTx.Hash().Hex(),
			"nonce", nonce,
			"value", legacyTx.Value,
			"gasPrice", legacyTx.GasPrice,
			"to", account.Address())

		// Wait for transaction receipt
		receipt, err := s.client.WaitForTxReceipt(s.ctx, signedTx.Hash())
		if err != nil {
			s.logger.Error("Failed waiting for transaction receipt", "error", err)
			return
		}

		if receipt.Status == 1 {
			s.logger.Info("Transaction confirmed",
				"hash", signedTx.Hash().Hex(),
				"blockNumber", receipt.BlockNumber,
				"gasUsed", receipt.GasUsed)
		} else {
			s.logger.Error("Transaction failed",
				"hash", signedTx.Hash().Hex(),
				"blockNumber", receipt.BlockNumber)
			return
		}
	}

	// Check final balance
	balance, err := s.client.GetBalance(sender)
	if err != nil {
		s.logger.Error("Failed to get balance", "error", err, "account", sender.Address())
		return
	}
	s.logger.Info("Account balance",
		"account", sender.Address(),
		"balance", balance)
}

func (s *Spammer) GetBalances() {
	for _, account := range s.prefundedAccounts {
		balance, err := s.client.GetBalance(account)
		if err != nil {
			s.logger.Error("Failed to get balance", "error", err, "account", account.Address())
			continue
		}
		s.logger.Info("Account balance", "account", account.Address(), "balance", balance)
	}
}

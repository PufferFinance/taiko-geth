package spammer

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	ctx context.Context
	*ethclient.Client
	chainID *big.Int
	logger  *log.Logger
}

func NewEthClient(ctx context.Context, url string, chainID *big.Int, logger *log.Logger) (*EthClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthClient{ctx, client, chainID, logger}, nil
}

func (ec *EthClient) GetBalance(account *Account) (*big.Int, error) {
	return ec.BalanceAt(ec.ctx, *account.Address(), nil)
}

func (ec *EthClient) FetchAssignedSlots() ([]uint64, error) {
	var assignedSlots []uint64
	err := ec.Client.Client().CallContext(ec.ctx, &assignedSlots, "taiko_fetchAssignedSlots")
	if err != nil {
		return nil, err
	}
	return assignedSlots, nil
}

func (ec *EthClient) FetchL1GenesisTimestamp() (uint64, error) {
	var l1GenesisTimestamp uint64
	err := ec.Client.Client().CallContext(ec.ctx, &l1GenesisTimestamp, "taiko_fetchL1GenesisTimestamp")
	if err != nil {
		return 0, err
	}
	return l1GenesisTimestamp, nil
}

func (ec *EthClient) FetchCurrentSlot(now int64) (uint64, uint64, error) {
	l1GenesisTimestamp, err := ec.FetchL1GenesisTimestamp()
	if err != nil {
		return 0, 0, err
	}
	headSlot, _ := HeadSlotAndEpoch(l1GenesisTimestamp, now)
	currentSlot := headSlot + 1
	_, headSlotEndTime := HeadSlotStartEndTime(l1GenesisTimestamp, now)
	return currentSlot, headSlotEndTime, nil
}

func (ec *EthClient) GetNonce(account *Account) (uint64, error) {
	nonce, err := ec.PendingNonceAt(ec.ctx, *account.Address())
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (ec *EthClient) SendTx(account *Account, tx *types.Transaction) (*types.Transaction, error) {
	ec.logger.Info("Sending tx", "nonce", tx.Nonce(), "to", tx.To().Hex(), "value", tx.Value().String(), "gas", tx.Gas(), "gas price", tx.GasPrice().String(), "private key", account.Address())
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(ec.chainID), account.PrivateKey())
	if err != nil {
		ec.logger.Error("Failed to sign transaction", "error", err)
		return nil, err
	}

	res := ec.SendTransaction(ec.ctx, signedTx)

	if reflect.DeepEqual(res, common.Hash{}) {
		ec.logger.Error("Tx rejected", "hash", res)
	} else {
		ec.logger.Info("Submitted tx", "hash", signedTx.Hash(), "nonce", signedTx.Nonce(), "slot deadline")
	}

	return signedTx, nil
}

func (ec *EthClient) LogTx(signedTx *types.Transaction) {
	_, _, err := ec.TransactionByHash(ec.ctx, signedTx.Hash())
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			ec.logger.Error("Transaction not found", "tx hash", signedTx.Hash())
			return
		} else {
			ec.logger.Error("Failed to get transaction by hash", "error", err, "tx hash", signedTx.Hash())
		}
	}
}

func (ec *EthClient) LogReceipt(signedTx *types.Transaction) {
	txReceipt, err := ec.TransactionReceipt(ec.ctx, signedTx.Hash())
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			ec.logger.Error("Transaction receipt not found", "tx hash ", signedTx.Hash())
			return
		} else {
			ec.logger.Error("Failed to get transaction receipt", "error", err)
		}
	}
	ec.logger.Warn("Transaction receipt", "tx hash", signedTx.Hash(), "block number", txReceipt.BlockNumber, "status", txReceipt.Status, "cumulative gas used", txReceipt.CumulativeGasUsed, "effective gas price", txReceipt.EffectiveGasPrice, "gas used", txReceipt.GasUsed)
}

func (c *EthClient) WaitForTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			receipt, err := c.TransactionReceipt(ctx, txHash)
			if err == nil && receipt != nil {
				return receipt, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

var (
	EpochLength  = 32 // number of slots per epoch
	SlotDuration = 12 // duration of each slot in seconds
)

func HeadSlotAndEpoch(genesisTimestamp uint64, now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	headSlot := uint64(elapsedTime) / uint64(SlotDuration)
	headEpoch := headSlot / uint64(EpochLength)
	return headSlot, headEpoch
}

func HeadSlotStartEndTime(genesisTimestamp uint64, now int64) (uint64, uint64) {
	headSlot, _ := HeadSlotAndEpoch(genesisTimestamp, now)
	slotStartTime := genesisTimestamp + headSlot*uint64(SlotDuration)
	slotEndTime := slotStartTime + uint64(SlotDuration)
	return slotStartTime, slotEndTime
}
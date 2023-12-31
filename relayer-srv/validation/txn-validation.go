package validation

import (
	"fmt"
	"sync"

	protocols_p2p "github.com/volmexfinance/relayers/relayer-srv/chat"
	"github.com/volmexfinance/relayers/relayer-srv/db"
)

var lock sync.Mutex

// TODO: Handle issue with transaction status
// HandleTransactionMessage handles transaction message received on p2p
func HandleTransactionMessage(dbs *db.PostgresDataBase, m *protocols_p2p.GossipMessage) error {
	lock.Lock()
	defer lock.Unlock()

	txHash := m.GetTransactionMessage().TransactionMessage.TransactionHash
	txID := m.GetTransactionMessage().TransactionMessage.TxID
	txOrderID := m.GetTransactionMessage().TransactionMessage.OrderID
	txNextLeader := m.GetTransactionMessage().TransactionMessage.NextLeader
	txNextPubkey := m.GetTransactionMessage().TransactionMessage.NextPubKey
	chain := m.GetTransactionMessage().MessageData.Chain

	if txHash == "" {
		for i := 0; i < int(len(txOrderID)); i++ {
			getOrder, err := dbs.FindOrder(txOrderID[i])
			if err != nil {
				return fmt.Errorf("HandleTransactionMessage: Order ID not found in database  %w", err)
			}
			// TODO: check status to be allowed or not allowed
			if err := dbs.UpdateOrderStatus(getOrder.OrderID, []db.MatchedStatus{db.MatchedStatusInit, db.MatchedStatusFullMatchFound, db.MatchedStatusPartialMatchFound, db.MatchedStatusPartialMatchConfirmed, db.MatchedStatusSentFailed, db.MatchedStatusValidationConfirmed, db.MatchedStatusValidated},
				[]db.MatchedStatus{db.MatchedStatusSentToContract, db.MatchedStatusFailedConfirmed, db.MatchedStatusFullMatchConfirmed},
				db.MatchedStatusSentFailed); err != nil {
				return fmt.Errorf("HandleTransactionMessage: Unable to update status of order %w", err)
			}
		}
		return nil
	}

	txSentByID, err := dbs.GetTxnByID(txID, chain)
	if err != nil {
		return fmt.Errorf("HandleTransactionMessage: Unable to GetTxnByID %w", err)
	} else if txSentByID.ID == 0 {
		newTxSent := db.TransactionSent{OrderID: txOrderID, TransactionHash: txHash, TransactionStatus: db.TransactionStatusTypeInit}
		if err := dbs.CreateTxSent(&newTxSent, chain); err != nil {
			return fmt.Errorf("HandleTransactionMessage: Unable to create next leader transaction in DB %w", err)
		}
	} else {
		txSentByID.OrderID = txOrderID
		txSentByID.TransactionHash = txHash
		err := dbs.UpdateTxnStatus(txSentByID, txSentByID.TransactionStatus)
		if err != nil {
			return fmt.Errorf("HandleTransactionMessage: Unable to update transaction in DB %w", err)
		}
	}

	for i := 0; i < int(len(txOrderID)); i++ {
		getOrder, err := dbs.FindOrder(txOrderID[i])
		if err != nil {
			return fmt.Errorf("HandleTransactionMessage: Order ID not found in database  %w", err)
		}
		// TODO: check status to be allowed or not allowed
		if err := dbs.UpdateOrderStatus(getOrder.OrderID, []db.MatchedStatus{db.MatchedStatusInit, db.MatchedStatusFullMatchFound, db.MatchedStatusPartialMatchFound, db.MatchedStatusPartialMatchConfirmed, db.MatchedStatusSentFailed, db.MatchedStatusValidationConfirmed, db.MatchedStatusValidated},
			[]db.MatchedStatus{db.MatchedStatusSentToContract, db.MatchedStatusFailedConfirmed, db.MatchedStatusFullMatchConfirmed},
			db.MatchedStatusSentToContract); err != nil {
			return fmt.Errorf("HandleTransactionMessage: Unable to update status of order %w", err)
		}
	}

	/// New leader found from p2p and created in db
	newLeaderTx := db.TransactionSent{Leader: txNextLeader, LeaderPublicKey: txNextPubkey, TransactionStatus: db.TransactionStatusTypeInit}
	if err := dbs.CreateTxSent(&newLeaderTx, chain); err != nil {
		return fmt.Errorf("HandleTransactionMessage: Unable to create next leader transaction in DB %w", err)
	}
	///TODO: Nonce verification and txn nonce calculation and verification for leader randomization
	return nil
}

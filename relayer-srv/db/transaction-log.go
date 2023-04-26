package db

import (
	"time"
)

// CreateTransactionLog created a transaction log in DB
func (db *PostgresDataBase) CreateTransactionLog(tLog *TransactionLog) error {
	tLog.CreatedAt = time.Now().Unix()

	if err := db.DB.Create(&tLog); err.Error != nil {
		return err.Error
	}
	return nil
}

// CreateTransactionLog created a transaction log in DB
func (db *SQLiteDataBase) CreateTransactionLog(tLog *TransactionLog) error {
	tLog.CreatedAt = time.Now().Unix()

	if err := db.DB.Create(&tLog); err.Error != nil {
		return err.Error
	}
	return nil
}

// FindLatestTxnLog returns last transaction log stored in DB
func (db *PostgresDataBase) FindLatestTxnLog(chain string) (*TransactionLog, error) {
	var lastTx TransactionLog
	if result := db.DB.Model(TransactionLog{}).Where("chain_name = ?", chain).Order("id desc").Limit(1).First(&lastTx); result.Error != nil {
		return nil, result.Error
	}
	return &lastTx, nil
}

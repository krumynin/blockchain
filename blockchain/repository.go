package blockchain

import "database/sql"

func GetBlocks(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT block FROM blockchain`)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func createBlockchainTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE blockchain (
			id INTEGER PRIMARY KEY autoincrement,
			hash VARCHAR(44) UNIQUE,
			block text
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func insertBlockToBlockchain(db *sql.DB, currHash, block string) error {
	_, err := db.Exec(`INSERT INTO blockchain (hash, block) VALUES ($1, $2)`,
		currHash,
		block,
	)
	if err != nil {
		return err
	}

	return nil
}

func getLastBlockId(db *sql.DB) uint64 {
	var index uint64
	row := db.QueryRow(`SELECT id FROM blockchain ORDER BY id DESC`)
	row.Scan(&index)

	return index
}

func getPreviousBlocks(db *sql.DB, index uint64) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT block FROM blockchain WHERE id <= $1 ORDER BY id DESC`, index)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func getLastHash(db *sql.DB) string {
	var hash string
	row := db.QueryRow(`SELECT hash FROM blockchain ORDER BY id DESC`)
	row.Scan(&hash)

	return hash
}

func getIdByHash(db *sql.DB, hash string) uint64 {
	var id uint64
	row := db.QueryRow(`SELECT id FROM blockchain WHERE hash = $1`, hash)
	row.Scan(&id)

	return id
}

func getBlockByHash(db *sql.DB, hash string) string {
	var block string
	row := db.QueryRow(`SELECT block FROM blockchain WHERE hash = $1`, hash)
	row.Scan(&block)

	return block
}

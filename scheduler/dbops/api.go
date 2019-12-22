package dbops

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func AddDeletion(vid string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO delect_vid(id) VALUES ?")
	if err != nil {
		log.Printf("add deletionID db prepare error: %v")
		return err
	}

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("add vid db error: %v\n", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ReadDeletion(count int) ([]string, error) {
	var ids []string
	stmtOut, err := dbConn.Prepare("SELECT id FROM delect_vid LIMIT ?")
	if err != nil {
		log.Printf("read deletionID db prepare error: %v!\n", err)
		return ids, err
	}

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("read vid db error: %v\n", err)
		return ids, err
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return ids, nil

}

func DeleteDeletion(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM delect_vid WHERE id = ?")
	if err != nil {
		log.Printf("delete vid db prepare error: %v\n", err)
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("delete vid db prepare error: %v\n", err)
		return err
	}

	defer stmtDel.Close()
	return nil

}

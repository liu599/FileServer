package model

import (
	"nekohandfileserver/middleware/data"
	"fmt"
	"nekohandfileserver/middleware/func"
	"time"
)

func FindFile(fi data.NekohandFile) (error, bool) {
	var ac int
	statement := fmt.Sprintf("select count(pid) from files where fileId = '%s'", fi.FileId)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err, false
	}

	err = db.QueryRow(statement).Scan(&ac)

	if err != nil {
		return err, false
	}

	return nil, ac > 0
}

func CreateFile(fi data.NekohandFile) (error) {
	createdTime := time.Now().Unix()

	statement := fmt.Sprintf("INSERT INTO files (fileId, hashId, fileName, createdAt, modifiedAt) VALUES('%s', '%s', '%s', '%d', '%d')", fi.FileId, fi.HashId, fi.FileName, createdTime, createdTime)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err
	}

	_, err = db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}

func UpdateFile(fi data.NekohandFile) error {
	if err, nm := FindFile(fi); nm == false {
		return err
	}
	fi.ModifiedAt = time.Now().Unix()
	statement := fmt.Sprintf("UPDATE post SET hashId='%s', filename='%s', modifiedAt='%d' WHERE fileid='%s'", fi.HashId, fi.FileName, fi.ModifiedAt, fi.ModifiedAt)

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err
	}

	_, err = db.Exec(statement)

	if err != nil {
		return err
	}

	return err
}

func FetchFile(fileId string) string {
	var nfile data.NekohandFile
	statement := fmt.Sprintf("select fileId, fileName from files where fileId = '%s'", fileId)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err.Error()
	}
	err = db.QueryRow(statement).Scan(&nfile.FileId, &nfile.FileName)

	if err != nil {
		return err.Error()
	}

	return nfile.FileId + "_" + nfile.FileName
}

func FetchFileList() (error, []data.NekohandFile) {
	statement := fmt.Sprintf("select * from files")
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return err, nil
	}
	rows, err := db.Query(statement)

	if err != nil {
		return err, []data.NekohandFile{}
	}

	nfiles := []data.NekohandFile{}

	for rows.Next() {
		var nf data.NekohandFile
		if err = rows.Scan(&nf.FID, &nf.FileId, &nf.HashId, &nf.FileName, &nf.CreatedAt, &nf.ModifiedAt); err != nil {
			return err, nil
		}
		nfiles = append(nfiles, nf)
	}
	return nil, nfiles
}
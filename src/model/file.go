package model

import (
	"fmt"
	"github.com/liu599/FileServer/src/data"
	"github.com/liu599/FileServer/src/middleware/func"
	"github.com/liu599/FileServer/src/utils"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path/filepath"
	"strings"
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

	var relativePathSQL string

	relativePathSQL = filepath.ToSlash(fi.RelativePath)

	fmt.Println(relativePathSQL)

	statement := fmt.Sprintf("INSERT INTO files (fileId, hashId, fileName, createdAt, modifiedAt, relativePath) VALUES('%s', '%s', '%s', '%d', '%d', '%s')", fi.FileId, fi.HashId, fi.FileName, createdTime, createdTime, relativePathSQL)

	fmt.Println(statement)

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

func FetchFile(fileId string) (string, error) {
	var nfile data.NekohandFile
	statement := fmt.Sprintf("select fileId, fileName, relativePath from files where fileId = '%s'", fileId)
	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		return "_", err
	}
	err = db.QueryRow(statement).Scan(&nfile.FileId, &nfile.FileName, &nfile.RelativePath)

	if err != nil {
		return "_", err
	}

	return nfile.RelativePath + nfile.FileId + "_" + nfile.FileName, nil
}

func FetchFileList(fileType string) (error, []data.NekohandFile) {
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
		if err = rows.Scan(&nf.FID, &nf.FileId, &nf.HashId, &nf.FileName, &nf.CreatedAt, &nf.ModifiedAt, &nf.RelativePath); err != nil {
			return err, nil
		}
		ftype := strings.Split(nf.FileName, ".")
		if fileType == "ALL" {
			nfiles = append(nfiles, nf)
		}
		if strings.ToLower(ftype[len(ftype)-1]) == fileType && fileType != "ALL" {
			nfiles = append(nfiles, nf)
		}
	}
	return nil, nfiles
}

func FixPath(userpath string, filetype string) error {
	var mdst string
	xfiles, err  := utils.GetAllFiles(userpath, filetype)
	if err != nil {
		return err
	}
	for _, file := range xfiles {
		salt := bson.NewObjectId().Hex()
		// tmp[0] path, tmp[1] filename tmp2 relative path
		tmp := strings.Split(file,"|")
		tmp2 := strings.Split(tmp[0], userpath)
		filephyurl := strings.Join(tmp, "")
		if mdst, err = utils.HashFileMd5(filephyurl); err != nil {
			return err
		}

		//fmt.Printf("获取的文件为[%s]， salt[%s], md5[%s], url[%v]\n", file, salt, mdst, tmp2)

		if err = CreateFile(data.NekohandFile{
			FileId: salt,
			HashId: mdst,
			FileName: tmp[1],
			RelativePath: filepath.FromSlash(strings.Join(tmp2, "")),
		}); err != nil {
			return err
		}

		err = os.Rename(filephyurl, tmp[0]+salt+"_"+tmp[1])

		if err != nil {
			return err
		}
	}
	return nil
}
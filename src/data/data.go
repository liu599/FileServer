package data

type (
	Error struct {
		Code    string                 `json:"code"`    // 错误代码
		Message string                 `json:"message"` // 错误信息
		Fields  map[string]interface{} `json:"fields,omitempty"`  // 错误字段信息
	}
	Database struct {
		DB      interface{}
		Driver  string
		MaxOpen int
		MaxIdle int
		Name    string
		Source  string
	}
	NekohandFile struct {
		FID      int      `json:"fid"`
		FileName string   `json:"filename"`
		FileId   string	  `json:"fileid"`
		HashId   string   `json:"filehash"`
		CreatedAt  int64  `json:"createdAt"`
		ModifiedAt int64  `json:"modifiedAt"`
		RelativePath string `json:"relativePath"`
		Src string `json:"src"`
	}
)
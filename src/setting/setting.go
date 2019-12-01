package setting

import (
	"github.com/go-ini/ini"
	"log"
	"path/filepath"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort string
	//ReadTimeout time.Duration

	Source string
	FileRoot string
	MaxIdle int
	MaxOpen int

	PageSize int
	JwtSecret string
	JwtTimeoutMinute int


)

func init() {
	var err error
	examplePath := filepath.FromSlash("./src/conf/my-config.ini")
	Cfg, err = ini.Load(examplePath)
	if err != nil {
		log.Fatalf("Fail to parse 'conf/my.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	//LoadApp()
	LoadDataBase()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	HTTPPort = sec.Key("HTTP_PORT").MustString(":17699")
}

func LoadDataBase() {
	sec, err := Cfg.GetSection("database")
	if  err!=nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	Source = sec.Key("SOURCE").MustString("root:root@tcp(127.0.0.1:3306)/shana?charset=utf8")
	FileRoot = sec.Key("FILEROOT").MustString("/home/wwwroot/shana/data/")
	MaxIdle = sec.Key("MAX_IDLE").MustInt(2)
	MaxOpen = sec.Key("MAX_OPEN").MustInt(5)
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	JwtTimeoutMinute = sec.Key("JwtTimeoutMinute").MustInt(30)
	// PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
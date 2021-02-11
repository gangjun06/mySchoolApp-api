package conf

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	JWTSecret string
	MongoDB   string
	Server    struct {
		Release bool
		Port    int
	}
	CoolSMS struct {
		Enable   bool
		ClientID string
		Secret   string
		From     string
	}
}

var conf *Config

func Init() {
	conf = &Config{}

	flag.StringVar(&conf.JWTSecret, "jwtsecret", os.Getenv("JWTSECRET"), "JWT Secret")
	flag.BoolVar(&conf.Server.Release, "release", parseBool(os.Getenv("SERVER_RELEASE")), "start server with release")
	flag.IntVar(&conf.Server.Port, "port", parseInt(os.Getenv("SERVER_PORT")), "Web Sserver Port")
	flag.StringVar(&conf.MongoDB, "mongodb", os.Getenv("DB_MONGODB"), "MongoDB ConnectURI")

	flag.BoolVar(&conf.CoolSMS.Enable, "coolsms_enable", parseBool(os.Getenv("COOLSMS_ENABLE")), "enable send sms")
	flag.StringVar(&conf.CoolSMS.ClientID, "coolsms_clientid", os.Getenv("COOLSMS_CLIENTID"), "coolsms client id")
	flag.StringVar(&conf.CoolSMS.Secret, "coolsms_secret", os.Getenv("COOLSMS_SECRET"), "coolsms secret")
	flag.StringVar(&conf.CoolSMS.From, "coolsms_from", os.Getenv("COOLSMS_FROM"), "coolsms from")

	flag.Parse()
}

func parseBool(str string) bool {
	val, _ := strconv.ParseBool(str)
	return val
}

func parseInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

func Get() *Config {
	return conf
}

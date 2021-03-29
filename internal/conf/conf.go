package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
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
		Enable bool
		ApiKey string
		Secret string
		From   string
	}
	NeisAPIKey string
	Redis      struct {
		Addr string
		Pass string
	}
	Email struct {
		Url      string
		UserName string
		Password string
	}
}

type DiscordConf struct {
	GroupID string `json"groupID"`
	Token   string `json"token"`
	SubPost []struct {
		CategoryID       string `json:"categoryID"`
		CategoryName     string `json:"categoryName"`
		DiscordChannelID string `json:"discordChannelID"`
	} `json"subpost"`
}

var conf *Config
var confDiscord *DiscordConf

func Init() {
	conf = &Config{}

	flag.StringVar(&conf.JWTSecret, "jwtsecret", os.Getenv("JWTSECRET"), "JWT Secret")
	flag.BoolVar(&conf.Server.Release, "release", parseBool(os.Getenv("SERVER_RELEASE")), "start server with release")
	flag.IntVar(&conf.Server.Port, "port", parseInt(os.Getenv("SERVER_PORT")), "Web Sserver Port")
	flag.StringVar(&conf.MongoDB, "mongodb", os.Getenv("DB_MONGODB"), "MongoDB ConnectURI")

	flag.BoolVar(&conf.CoolSMS.Enable, "coolsms_enable", parseBool(os.Getenv("COOLSMS_ENABLE")), "enable send sms")
	flag.StringVar(&conf.CoolSMS.ApiKey, "coolsms_apikey", os.Getenv("COOLSMS_APIKEY"), "coolsms api key")
	flag.StringVar(&conf.CoolSMS.Secret, "coolsms_secret", os.Getenv("COOLSMS_SECRET"), "coolsms secret")
	flag.StringVar(&conf.CoolSMS.From, "coolsms_from", os.Getenv("COOLSMS_FROM"), "coolsms from")

	flag.StringVar(&conf.NeisAPIKey, "neis_apikey", os.Getenv("NEIS_APIKEY"), "https://open.neis.go.kr api key")

	flag.StringVar(&conf.Redis.Addr, "redis_addr", os.Getenv("REDIS_ADDR"), "address of redis")
	flag.StringVar(&conf.Redis.Pass, "redis_pass", os.Getenv("REDIS_PASS"), "password of redis")

	flag.StringVar(&conf.Email.Url, "mail_url", os.Getenv("EMAIL_URL"), "MailInABox url")
	flag.StringVar(&conf.Email.UserName, "mail_username", os.Getenv("EMAIL_USERNAME"), "username")
	flag.StringVar(&conf.Email.Password, "mail_password", os.Getenv("EMAIL_PASSWORD"), "password")

	flag.Parse()

	file, err := ioutil.ReadFile("./discordConf.json")
	if err != nil {
		return
	}
	json.Unmarshal(file, &confDiscord)

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

func Discord() *DiscordConf {
	return confDiscord
}

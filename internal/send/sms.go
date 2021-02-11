package send

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/osang-school/backend/internal/conf"
	"github.com/osang-school/backend/internal/utils"
)

func Sms(to, text string) error {
	config := conf.Get().CoolSMS

	if config.Enable {
		fmt.Println(to, " >> ", text)
		return nil
	}

	h := hmac.New(sha256.New, []byte(config.Secret))
	time := time.Now().Format("2006-01-02T15:04:05-0700")
	salt := utils.CreateRandomString(32)

	if _, err := h.Write(append([]byte(time), []byte(salt)...)); err != nil {
		return err
	}
	signature := hex.EncodeToString(h.Sum(nil))

	header := fmt.Sprintf("HMAC-SHA256 apiKey=%s, date=%s, salt=%s, signature=%s", config.Secret, time, salt, signature)

	bodyBytes, err := json.Marshal(map[string]interface{}{
		"message": map[string]interface{}{
			"to":   to,
			"from": config.From,
			"text": text,
			"type": "SMS",
		},
	})
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest("POST", "https://api.coolsms.co.kr/messages/v4/send", body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	return nil
}

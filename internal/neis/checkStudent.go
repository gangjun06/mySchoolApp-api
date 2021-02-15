package neis

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	AREA_SEOUL Area = 1 + iota
	AREA_BUSAN
	AREA_DAEGU
	AREA_INCHEON
	AREA_GWANGJU
	AREA_DAEJEON
	AREA_ULSAN
	AREA_SEJONG
	AREA_GYEONGGI
	AREA_GANGWON
	AREA_CHUNGBUK
	AREA_CHUNGNAM
	AREA_JEONBUK
	AREA_JEONNAM
	AREA_GYEONGBUK
	AREA_GYEONGNAM
	AREA_JEJ
)

type Area uint8

// Encrypt 택스트를 자가진단 사이트에 요청 보낼 수 있게 암호화 시킴
func Encrypt(text string) *string {
	keyOrigin := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA81dCnCKt0NVH7j5Oh2+SGgEU0aqi5u6sYXemouJWXOlZO3jqDsHYM1qfEjVvCOmeoMNFXYSXdNhflU7mjWP8jWUmkYIQ8o3FGqMzsMTNxr+bAp0cULWu9eYmycjJwWIxxB7vUwvpEUNicgW7v5nCwmF5HS33Hmn7yDzcfjfBs99K5xJEppHG0qc+q3YXxxPpwZNIRFn0Wtxt0Muh1U8avvWyw03uQ/wMBnzhwUC8T4G5NclLEWzOQExbQ4oDlZBv8BM/WxxuOyu0I8bDUDdutJOfREYRZBlazFHvRKNNQQD2qDfjRz484uFs7b5nykjaMB9k/EJAuHjJzGs9MMMWtQIDAQAB"

	key, err := base64.StdEncoding.DecodeString(keyOrigin)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, isRSAPublicKey := publicKeyInterface.(*rsa.PublicKey)
	if !isRSAPublicKey {
		log.Fatal("It is not RSA Public Key")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(text))
	if err != nil {
		log.Fatal(err)
	}

	result := base64.StdEncoding.EncodeToString(ciphertext)

	return &result
}

func getAreaURL(area Area) string {
	AreaURL := []string{"sen", "pen", "dge", "ice", "gen", "dje", "use", "sje", "goe", "kwe", "cbe", "cne", "jbe", "jne", "gbe", "gne", "jje"}
	return AreaURL[area-1]
}

// CheCheckStudent 학생의 학교정보, 생일, 자가진단 비밀번호를 바탕으로 학생 본인인지 확인합니다
// school: 학교 org코드
// birth: 주민등록번호 앞 6글자
func CheckStudent(area Area, school, birth, name, password string) error {
	areaUrl := getAreaURL(area)
	url := fmt.Sprintf("https://%shcs.eduro.go.kr/v2/findUser", areaUrl)
	reqBody, err := json.Marshal(map[string]interface{}{
		"name":      Encrypt(name),
		"birthday":  Encrypt(birth),
		"orgCode":   school,
		"loginType": "school",
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("bad request")
	}

	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var parse map[string]interface{}

	if err := json.Unmarshal(bodyByte, &parse); err != nil {
		return err
	}

	token, ok := parse["token"].(string)
	if !ok {
		return fmt.Errorf("error")
	}

	url = fmt.Sprintf("https://%shcs.eduro.go.kr/v2/validatePassword", areaUrl)
	reqBody2, err := json.Marshal(map[string]interface{}{
		"password":   Encrypt(password),
		"deviceUuid": "",
	})
	if err != nil {
		return err
	}

	body2 := bytes.NewBuffer(reqBody2)

	req, err := http.NewRequest(http.MethodPost, url, body2)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		return err
	} else if resp2.StatusCode != 200 {
		return fmt.Errorf("bad request")
	}

	return nil
}

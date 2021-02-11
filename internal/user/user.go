package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/muesli/cache2go"
	"github.com/osang-school/backend/internal/utils"
)

type (
	User struct {
	}
)

var (
	ErrTooManyVerifyReq = errors.New("Too Many Verify Req")
)

func PhoneVerifyCode(ip, phone string) (string, error) {
	cache := cache2go.Cache("phoneVerify")
	cnt := 0
	if res, err := cache.Value("cnt:" + ip); err == nil {
		cnt = res.Data().(int)
		if cnt > 3 {
			return "", ErrTooManyVerifyReq
		}
	}

	code := utils.CreateRandomNum(6)
	cache.Add(phone, time.Minute*5, code)

	cnt++
	cache.Add("cnt:"+ip, utils.TimeLeftUntilMidnight(), cnt)
	return code, nil
}

func PhoneVerifyCheck(phone, code string) error {
	cache := cache2go.Cache("phoneVerify")
	res, err := cache.Value(phone)
	if err != nil {
		return err
	} else if res.Data().(string) != code {
		return fmt.Errorf("Code is not correct")
	}
	return nil
}

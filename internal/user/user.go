package user

import (
	"errors"
	"time"

	"github.com/muesli/cache2go"
	"github.com/osang-school/backend/internal/send"
	"github.com/osang-school/backend/internal/utils"
)

type (
	User struct {
	}
)

var (
	ErrTooManyVerifyReq = errors.New("Too Many Verify Req")
)

// SendMax per day
const SendMax = 5

func PhoneVerifyCode(ip, phone string) error {
	cache := cache2go.Cache("phoneVerify")
	cnt := 0
	if res, err := cache.Value("cnt:" + ip); err == nil {
		cnt = res.Data().(int)
		if cnt >= SendMax {
			return ErrTooManyVerifyReq
		}
	}

	code := utils.CreateRandomNum(6)
	cache.Add(phone, time.Minute*5, code)

	if err := send.Sms(phone, "오상중학교 회원가입 본인확인 인증코드는 ["+code+"]입니다."); err != nil {
		return err
	}

	cnt++
	cache.Add("cnt:"+ip, utils.TimeLeftUntilMidnight(), cnt)
	return nil
}

func PhoneVerifyCheck(phone, code string) (bool, error) {
	cache := cache2go.Cache("phoneVerify")
	res, err := cache.Value(phone)
	if err != nil {
		return false, err
	} else if res.Data().(string) != code {
		return false, nil
	}
	return true, nil
}

package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/muesli/cache2go"
	"github.com/osang-school/backend/graph/model"
	"github.com/osang-school/backend/internal/db/mongodb"
	"github.com/osang-school/backend/internal/send"
	"github.com/osang-school/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Role   uint8
	Status uint8
	User   struct {
		ID       primitive.ObjectID `bson:"_id,omitempty"`
		Status   Status             `bson:"status,omitempty"`
		Name     string             `bson:"name,omitempty"`
		Phone    string             `bson:"phone,omitempty"`
		Password string             `bson:"password,omitempty"`
		Nickname string             `bson:"nickname,omitempty"`
		Role     Role               `bson:"role,omitempty"`
		Student  *Student           `bson:"student,omitempty"`
		Teacher  *Teacher           `bson:"teacher,omitempty"`
		Officals *Officals          `bson:"officals,omitempty"`
	}
	Student struct {
		Grade  int `bson:"grade,omitempty"`
		Class  int `bson:"class,omitempty"`
		Number int `bson:"number,omitempty"`
	}
	Teacher struct {
		Subject []string `bson:"subject,omitempty"`
	}
	Officals struct {
		Role        string `bson:"role,omitempty"`
		Description string `bson:"description,omitempty"`
	}
)

const (
	RoleStudent Role = iota + 1
	RoleTeacher
	RoleOfficals
)

const (
	StatusUser Status = iota + 1
	StatusWait
	StatusBan
)

var (
	ErrTooManyVerifyReq = errors.New("Too Many Verify Req")
)

// SendMax per day
const SendMax = 5

// PhoneVerifyCode generate and send phone verify code
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

// PhoneVerifyCheck check phone verify code is correct
func PhoneVerifyCheck(phone, code string) (string, error) {
	cache := cache2go.Cache("phoneVerify")
	res, err := cache.Value(phone)
	if err != nil {
		return "", err
	} else if res.Data().(string) != code {
		return "", fmt.Errorf("Not Valid Code")
	}

	signupCode := utils.CreateRandomString(6)
	cache.Add("signup:"+signupCode, time.Hour, phone)
	return signupCode, nil
}

// PhoneSignUpCheck load phone number from phone signup code
func PhoneSignUpCheck(code string) (string, error) {
	cache := cache2go.Cache("phoneVerify")
	res, err := cache.Value("signup:" + code)
	if err != nil {
		return "", err
	}
	return res.Data().(string), nil
}

// CheckStudentDup exits = true, not exits = false
func CheckStudentDup(grade, class, number int) (bool, error) {
	filter := bson.M{"student.grade": grade, "student.class": class, "student.number": number}
	var result User
	if err := mongodb.User.FindOne(nil, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckUserDup exits = true, not exits = false
func CheckUserDup(phone string) (bool, error) {
	filter := bson.M{"phone": phone}
	var result User
	if err := mongodb.User.FindOne(nil, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SignUp Just signup user
func SignUp(user *User) (primitive.ObjectID, error) {
	user.Status = StatusWait
	result, err := mongodb.User.InsertOne(nil, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), err
}

func DetailToUnion(d interface{}) model.ProfileDetail {
	var result model.ProfileDetail
	switch v := d.(type) {
	case *Student:
		result = model.StudentProfile{
			Grade:  v.Grade,
			Class:  v.Class,
			Number: v.Number,
		}
	case *Teacher:
		result = model.TeacherProfile{
			Subject: v.Subject,
		}
	case *Officals:
		result = model.OfficalsProfile{
			Role:        v.Role,
			Description: v.Description,
		}
	}
	return result
}

func StatusToEnum(s Status) model.UserStatus {
	switch s {
	case StatusUser:
		return model.UserStatusUser
	case StatusWait:
		return model.UserStatusWait
	default:
		return model.UserStatusBan
	}
}

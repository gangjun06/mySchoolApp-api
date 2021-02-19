package user

import (
	"strconv"
	"time"

	"github.com/osang-school/backend/graph/model"
	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/db/mongodb"
	"github.com/osang-school/backend/internal/db/redis"
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
		ID          primitive.ObjectID `bson:"_id,omitempty"`
		Permissions []string           `bson:"permissions,omitempty"`
		Status      Status             `bson:"status,omitempty"`
		Name        string             `bson:"name,omitempty"`
		Phone       string             `bson:"phone,omitempty"`
		Password    string             `bson:"password,omitempty"`
		Nickname    string             `bson:"nickname,omitempty"`
		Role        Role               `bson:"role,omitempty"`
		Student     *Student           `bson:"student,omitempty"`
		Teacher     *Teacher           `bson:"teacher,omitempty"`
		Officals    *Officals          `bson:"officals,omitempty"`
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
	Anon struct {
	}
)

const (
	RoleStudent Role = iota + 1
	RoleTeacher
	RoleOfficals
	RoleAnon
)

const (
	StatusUser Status = iota + 1
	StatusWait
	StatusBan
)

// SendMax per day
const SendMax = 5

// PhoneVerifyCode generate and send phone verify code
func PhoneVerifyCode(ip, phone string) error {
	cnt := 0
	if res, err := redis.C.Get("phone_verify:cnt:" + ip).Result(); err == nil {
		cnt, _ = strconv.Atoi(res)
		if cnt >= SendMax {
			return myerr.New(myerr.ErrTooManyReq, "request for verify beyond the limit")
		}
	}

	code := utils.CreateRandomNum(6)
	redis.C.Set("phone_verify:code:"+phone, code, time.Minute*5)

	if err := send.Sms(phone, "오상중학교 회원가입 본인확인 인증코드는 ["+code+"]입니다."); err != nil {
		return err
	}

	cnt++
	redis.C.Set("phone_verify:cnt:"+ip, cnt, 0)
	redis.C.ExpireAt("phone_verify:cnt:"+ip, utils.TodayTimeNoon())
	return nil
}

// PhoneVerifyCheck check phone verify code is correct
func PhoneVerifyCheck(phone, code string) (string, error) {

	str, err := redis.C.Get("phone_verify:code:" + phone).Result()
	if err != nil {
		if redis.IsNil(err) {
			return "", myerr.New(myerr.ErrBadRequest, "not valid code")
		}
		return "", err
	}
	if str != code {
		return "", myerr.New(myerr.ErrBadRequest, "not valid code")
	}

	signupCode := utils.CreateRandomString(6)
	redis.C.Set("phone_verify:signup:"+signupCode, phone, time.Hour)
	return signupCode, nil
}

// PhoneSignUpCheck load phone number from phone signup code
func PhoneSignUpCheck(code string) (string, error) {
	phone, err := redis.C.Get("phone_verify:signup:" + code).Result()
	if err != nil {
		if redis.IsNil(err) {
			return "", myerr.New(myerr.ErrBadRequest, "not valid code")
		}
		return "", myerr.New(myerr.ErrServer, err.Error())
	}
	return phone, nil
}

// ChutheckStudentDup exits = true, not exits = false
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

// GetUserByPhone
func GetUserByPhone(phone model.Phone) (*User, error) {
	filter := bson.M{"phone": phone}
	var result User
	if err := mongodb.User.FindOne(nil, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, myerr.New(myerr.ErrNotFound, "user not found")
		}
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	return &result, nil
}

// GetUserByID
func GetUserByID(id primitive.ObjectID) (*User, error) {
	filter := bson.M{"_id": id}
	var result User
	if err := mongodb.User.FindOne(nil, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, myerr.New(myerr.ErrNotFound, "user not found")
		}
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	return &result, nil
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
	default:
		result = model.AnonProfile{}
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

func UserToGqlType(u *User) *model.Profile {
	phone := model.Phone(u.Phone)
	profile := &model.Profile{
		ID:       model.ObjectID(u.ID),
		Name:     u.Name,
		Nickname: u.Nickname,
		Phone:    &phone,
		Status:   StatusToEnum(u.Status),
	}
	switch u.Role {
	case RoleStudent:
		profile.Detail = model.StudentProfile{
			Grade:  u.Student.Grade,
			Class:  u.Student.Class,
			Number: u.Student.Number,
		}
	case RoleOfficals:
		profile.Detail = model.OfficalsProfile{
			Role:        u.Officals.Role,
			Description: u.Officals.Description,
		}
	case RoleTeacher:
		profile.Detail = model.TeacherProfile{
			Subject: u.Teacher.Subject,
		}
	case RoleAnon:
		profile.Detail = model.AnonProfile{}
	}
	return profile
}

package db

import (
	"time"

	"github.com/gangjun06/mySchoolApp-api/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name              string `gorm:"type:varchar(255)"`
	IsTeacher         bool   `gorm:"default:false"`
	TeacherSubject    string `gorm:"type:varchar(100)"`
	StudentGrade      int    `gorm:"type:tinyint;default:0"`
	StudentClass      int    `gorm:"type:tinyint;default:0"`
	StudentNumber     int    `gorm:"type:tinyint;default:0"`
	IsVerifiedUser    bool   `gorm:"default:false"`
	Birth             *time.Time
	StatusMessage     string `gorm:"varchar(255)"`
	Avatar            string `gorm:"varchar(255)"`
	Password          string `gorm:"varchar(255)"`
	PasswordResetDate *time.Time
	PasswordResetCode string `gorm:"type:varchar(255);uniqueIndex"`
	Phone             string `gorm:"type:varchar(255);uniqueIndex"`
	PhoneVerifyCode   string `gorm:"type:varchar(100);uniqueIndex"`
	PhoneDate         *time.Time
	PhoneVerified     bool `gorm:"default:false"`
}

func CreateUser(name, password string) (*User, error) {
	user := User{
		Name:          name,
		Avatar:        "",
		StatusMessage: "",
		Password:      utils.HashAndSalt(password),
	}
	if err := utils.GetDB().Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (user *User, err error) {
	err = utils.GetDB().Find(user, id).Error
	return
}

func GetUserByPhone(phone string) (user *User, err error) {
	err = utils.GetDB().Where("phone = ?", phone).Find(user).Error
	return
}

func (u *User) GetToken() (string, error) {
	return utils.GetJwtToken(u.ID)
}

func (u *User) IsPasswordMatch(password string) bool {
	return utils.CheckPassword(password, u.Password)
}

func (u *User) UpdateAvatar(avatar string) error {
	return utils.GetDB().Model(u).Updates(&User{
		Avatar: avatar,
	}).Error
}

func (u *User) UpdateTeacher(teacherSubject string) error {
	return utils.GetDB().Model(u).Updates(&User{
		IsTeacher:      true,
		TeacherSubject: teacherSubject,
		StudentGrade:   0,
		StudentClass:   0,
		StudentNumber:  0,
	}).Error
}

func (u *User) UpdateStudent(grade, class, number int) error {
	return utils.GetDB().Model(u).Updates(&User{
		IsTeacher:      false,
		TeacherSubject: "",
		StudentGrade:   grade,
		StudentClass:   class,
		StudentNumber:  number,
	}).Error
}

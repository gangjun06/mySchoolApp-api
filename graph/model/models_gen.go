// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type ProfileDetail interface {
	IsProfileDetail()
}

type AnonProfile struct {
	Dummy *string `json:"dummy"`
}

func (AnonProfile) IsProfileDetail() {}

type Calendar struct {
	ID          ObjectID `json:"id"`
	Year        uint     `json:"year"`
	Month       uint     `json:"month"`
	Day         uint     `json:"day"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
}

type CalendarFilter struct {
	Year  uint `json:"year"`
	Month uint `json:"month"`
}

type Category struct {
	ID            ObjectID   `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	ReqPermission []string   `json:"reqPermission"`
	AnonAble      bool       `json:"anonAble"`
	ReadAbleRole  []UserRole `json:"readAbleRole"`
	WriteAbleRole []UserRole `json:"writeAbleRole"`
}

type Comment struct {
	ID       ObjectID   `json:"id"`
	Author   *Profile   `json:"author"`
	Content  string     `json:"content"`
	CreateAt Timestamp  `json:"createAt"`
	UpdateAt Timestamp  `json:"updateAt"`
	Status   PostStatus `json:"status"`
}

type CommentFilter struct {
	Limit           *int  `json:"limit"`
	Offset          *int  `json:"offset"`
	LoadOnlyComment *bool `json:"loadOnlyComment"`
}

type HomepageDetailFilter struct {
	Board HomepageBoard `json:"board"`
	ID    uint          `json:"id"`
}

type HomepageDetailType struct {
	ID        uint                `json:"id"`
	Title     string              `json:"title"`
	WrittenBy string              `json:"writtenBy"`
	CreateAt  Timestamp           `json:"createAt"`
	Content   string              `json:"content"`
	Images    []string            `json:"images"`
	Files     []*HomepageFileType `json:"files"`
}

type HomepageFileType struct {
	Name     string `json:"name"`
	Download string `json:"download"`
	Preview  string `json:"preview"`
}

type HomepageListFilter struct {
	Board HomepageBoard `json:"board"`
	Page  uint          `json:"page"`
}

type HomepageListType struct {
	ID        uint      `json:"id"`
	Number    uint      `json:"number"`
	Title     string    `json:"title"`
	WrittenBy string    `json:"writtenBy"`
	CreateAt  Timestamp `json:"createAt"`
}

type LikePostInput struct {
	Post   ObjectID `json:"post"`
	Status bool     `json:"status"`
}

type NewCalendar struct {
	Year        uint   `json:"year"`
	Month       uint   `json:"month"`
	Day         uint   `json:"day"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type NewCategory struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	ReqPermission []string   `json:"reqPermission"`
	AnonAble      bool       `json:"anonAble"`
	ReadAbleRole  []UserRole `json:"readAbleRole"`
	WriteAbleRole []UserRole `json:"writeAbleRole"`
}

type NewComment struct {
	Post    ObjectID `json:"post"`
	Content string   `json:"content"`
	Anon    *bool    `json:"anon"`
}

type NewPost struct {
	Category ObjectID `json:"category"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Anon     *bool    `json:"anon"`
}

type OfficialsProfile struct {
	Role        string `json:"role"`
	Description string `json:"description"`
}

func (OfficialsProfile) IsProfileDetail() {}

type OfficialsProfileInput struct {
	Role        string  `json:"role"`
	Description *string `json:"description"`
}

type Post struct {
	ID       ObjectID   `json:"id"`
	Category *Category  `json:"category"`
	Like     *int       `json:"like"`
	IsLike   *bool      `json:"isLike"`
	Author   *Profile   `json:"author"`
	Title    string     `json:"title"`
	Content  *string    `json:"content"`
	CreateAt Timestamp  `json:"createAt"`
	UpdateAt Timestamp  `json:"updateAt"`
	Comment  []*Comment `json:"comment"`
	Status   PostStatus `json:"status"`
}

type Profile struct {
	ID       ObjectID      `json:"id"`
	Name     string        `json:"name"`
	Nickname string        `json:"nickname"`
	Phone    *Phone        `json:"phone"`
	Detail   ProfileDetail `json:"detail"`
	Status   UserStatus    `json:"status"`
	Role     UserRole      `json:"role"`
}

type ProfileWithToken struct {
	Profile *Profile `json:"profile"`
	Token   string   `json:"token"`
}

type Schedule struct {
	Dow         uint   `json:"dow"`
	Period      uint   `json:"period"`
	Grade       uint   `json:"grade"`
	Class       uint   `json:"class"`
	Subject     string `json:"subject"`
	Teacher     string `json:"teacher"`
	Description string `json:"description"`
	ClassRoom   string `json:"classRoom"`
}

type ScheduleDelFilter struct {
	Grade  uint `json:"grade"`
	Class  uint `json:"class"`
	Dow    uint `json:"dow"`
	Period uint `json:"period"`
}

type ScheduleFilter struct {
	Grade *uint   `json:"grade"`
	Class *uint   `json:"class"`
	Dow   uint    `json:"dow"`
	Name  *string `json:"name"`
}

type SchoolMeal struct {
	Type     SchoolMealType `json:"type"`
	Calorie  string         `json:"calorie"`
	Content  string         `json:"content"`
	Nutrient string         `json:"nutrient"`
	Origin   string         `json:"origin"`
	Date     Timestamp      `json:"date"`
}

type SchoolMealFilter struct {
	DateStart *Timestamp      `json:"dateStart"`
	DateEnd   *Timestamp      `json:"dateEnd"`
	Type      *SchoolMealType `json:"type"`
}

type SignUpInput struct {
	Name     string  `json:"name"`
	Nickname *string `json:"nickname"`
	Password string  `json:"password"`
	Phone    string  `json:"phone"`
	Detail   string  `json:"detail"`
}

type StudentProfile struct {
	Grade  int `json:"grade"`
	Class  int `json:"class"`
	Number int `json:"number"`
}

func (StudentProfile) IsProfileDetail() {}

type StudentProfileInput struct {
	Grade  int `json:"grade"`
	Class  int `json:"class"`
	Number int `json:"number"`
}

type TeacherProfile struct {
	Subject []string `json:"subject"`
}

func (TeacherProfile) IsProfileDetail() {}

type TeacherProfileInput struct {
	Subject []string `json:"subject"`
}

type UpdateSchedule struct {
	Dow         uint   `json:"dow"`
	Period      uint   `json:"period"`
	Grade       uint   `json:"grade"`
	Class       uint   `json:"class"`
	Subject     string `json:"subject"`
	Teacher     string `json:"teacher"`
	Description string `json:"description"`
	ClassRoom   string `json:"classRoom"`
}

type UserNotificationID struct {
	ID string `json:"id"`
}

type HomepageBoard string

const (
	HomepageBoardNotice         HomepageBoard = "Notice"
	HomepageBoardPrints         HomepageBoard = "Prints"
	HomepageBoardRule           HomepageBoard = "Rule"
	HomepageBoardEvaluationPlan HomepageBoard = "EvaluationPlan"
	HomepageBoardAdministration HomepageBoard = "Administration"
)

var AllHomepageBoard = []HomepageBoard{
	HomepageBoardNotice,
	HomepageBoardPrints,
	HomepageBoardRule,
	HomepageBoardEvaluationPlan,
	HomepageBoardAdministration,
}

func (e HomepageBoard) IsValid() bool {
	switch e {
	case HomepageBoardNotice, HomepageBoardPrints, HomepageBoardRule, HomepageBoardEvaluationPlan, HomepageBoardAdministration:
		return true
	}
	return false
}

func (e HomepageBoard) String() string {
	return string(e)
}

func (e *HomepageBoard) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HomepageBoard(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HomepageBoard", str)
	}
	return nil
}

func (e HomepageBoard) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PostStatus string

const (
	PostStatusNormal   PostStatus = "Normal"
	PostStatusDeleted  PostStatus = "Deleted"
	PostStatusReported PostStatus = "Reported"
)

var AllPostStatus = []PostStatus{
	PostStatusNormal,
	PostStatusDeleted,
	PostStatusReported,
}

func (e PostStatus) IsValid() bool {
	switch e {
	case PostStatusNormal, PostStatusDeleted, PostStatusReported:
		return true
	}
	return false
}

func (e PostStatus) String() string {
	return string(e)
}

func (e *PostStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PostStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PostStatus", str)
	}
	return nil
}

func (e PostStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SchoolMealType string

const (
	SchoolMealTypeBreakfast SchoolMealType = "BREAKFAST"
	SchoolMealTypeLunch     SchoolMealType = "LUNCH"
	SchoolMealTypeDinner    SchoolMealType = "DINNER"
)

var AllSchoolMealType = []SchoolMealType{
	SchoolMealTypeBreakfast,
	SchoolMealTypeLunch,
	SchoolMealTypeDinner,
}

func (e SchoolMealType) IsValid() bool {
	switch e {
	case SchoolMealTypeBreakfast, SchoolMealTypeLunch, SchoolMealTypeDinner:
		return true
	}
	return false
}

func (e SchoolMealType) String() string {
	return string(e)
}

func (e *SchoolMealType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SchoolMealType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SchoolMealType", str)
	}
	return nil
}

func (e SchoolMealType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserRole string

const (
	UserRoleStudent   UserRole = "Student"
	UserRoleTeacher   UserRole = "Teacher"
	UserRoleOfficials UserRole = "Officials"
)

var AllUserRole = []UserRole{
	UserRoleStudent,
	UserRoleTeacher,
	UserRoleOfficials,
}

func (e UserRole) IsValid() bool {
	switch e {
	case UserRoleStudent, UserRoleTeacher, UserRoleOfficials:
		return true
	}
	return false
}

func (e UserRole) String() string {
	return string(e)
}

func (e *UserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserRole", str)
	}
	return nil
}

func (e UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserStatus string

const (
	UserStatusWait UserStatus = "WAIT"
	UserStatusUser UserStatus = "USER"
	UserStatusBan  UserStatus = "BAN"
)

var AllUserStatus = []UserStatus{
	UserStatusWait,
	UserStatusUser,
	UserStatusBan,
}

func (e UserStatus) IsValid() bool {
	switch e {
	case UserStatusWait, UserStatusUser, UserStatusBan:
		return true
	}
	return false
}

func (e UserStatus) String() string {
	return string(e)
}

func (e *UserStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserStatus", str)
	}
	return nil
}

func (e UserStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

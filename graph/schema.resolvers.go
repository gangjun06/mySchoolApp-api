package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/muesli/cache2go"
	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/graph/model"
	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/conf"
	"github.com/osang-school/backend/internal/discord"
	"github.com/osang-school/backend/internal/info"
	"github.com/osang-school/backend/internal/neis"
	"github.com/osang-school/backend/internal/post"
	"github.com/osang-school/backend/internal/session"
	"github.com/osang-school/backend/internal/user"
	"github.com/osang-school/backend/internal/utils"
	osangdata "github.com/osang-school/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) SignIn(ctx context.Context, phone model.Phone, password string) (*model.ProfileWithToken, error) {
	userData, err := user.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if ok := utils.CheckPassword(password, userData.Password); !ok {
		return nil, myerr.New(myerr.ErrPasswordWrong, "")
	}
	token, err := session.CreateToken(userData.ID, userData.Role, userData.Permissions)
	if err != nil {
		return nil, err
	}
	return &model.ProfileWithToken{
		Profile: user.UserToGqlType(userData),
		Token:   token,
	}, nil
}

func (r *mutationResolver) SignOut(ctx context.Context) (string, error) {
	if err := ctx.Value("data").(*session.Data).Expiry(); err != nil {
		return "", myerr.New(myerr.ErrServer, err.Error())
	}
	return "", nil
}

func (r *mutationResolver) VerifyPhone(ctx context.Context, number model.Phone) (string, error) {
	exits, err := user.CheckUserDup(string(number))
	if err != nil {
		return "", err
	} else if exits {
		return "", myerr.New(myerr.ErrDuplicate, "")
	}
	return "", user.PhoneVerifyCode(ctx.Value("ip").(string), string(number))
}

func (r *mutationResolver) CheckVerifyPhoneCode(ctx context.Context, number model.Phone, code string) (string, error) {
	return user.PhoneVerifyCheck(string(number), code)
}

func (r *mutationResolver) SetProfile(ctx context.Context, student *model.StudentProfileInput, teacher *model.TeacherProfileInput, officials *model.OfficialsProfileInput) (string, error) {
	randomStr := utils.CreateRandomString(6)
	cache := cache2go.Cache("profile")
	if student != nil {
		exits, err := user.CheckStudentDup(student.Grade, student.Class, student.Number)
		if err != nil {
			return "", err
		} else if exits {
			return "", myerr.New(myerr.ErrDuplicate, "")
		}
		cache.Add(randomStr, time.Hour, student)
	} else if teacher != nil {
		cache.Add(randomStr, time.Hour, teacher)
	} else if officials != nil {
		cache.Add(randomStr, time.Hour, officials)
	} else {
		return "", myerr.New(myerr.ErrBadRequest, "")
	}
	return randomStr, nil
}

func (r *mutationResolver) SignUp(ctx context.Context, input model.SignUpInput) (*model.ProfileWithToken, error) {
	phone, err := user.PhoneSignUpCheck(input.Phone)
	if err != nil {
		return nil, err
	}
	cache := cache2go.Cache("profile")
	res, err := cache.Value(input.Detail)
	if err != nil {
		return nil, myerr.New(myerr.ErrBadRequest, "")
	}
	detailData := res.Data()

	newUser := &user.User{
		Name:     input.Name,
		Phone:    phone,
		Status:   user.StatusWait,
		Password: utils.HashAndSalt(input.Password),
	}
	if input.Nickname != nil {
		newUser.Nickname = *input.Nickname
	}
	var resultDetail interface{}
	switch v := detailData.(type) {
	case *model.StudentProfileInput:
		newUser.Role = user.RoleStudent
		newUser.Student = &user.Student{
			Grade:  v.Grade,
			Class:  v.Class,
			Number: v.Number,
		}
		resultDetail = newUser.Student
	case *model.TeacherProfileInput:
		newUser.Role = user.RoleTeacher
		newUser.Teacher = &user.Teacher{
			Subject: v.Subject,
		}
		resultDetail = newUser.Teacher
	case *model.OfficialsProfileInput:
		newUser.Role = user.RoleOfficials
		newUser.Officials = &user.Officials{
			Role: v.Role,
		}
		if v.Description != nil {
			newUser.Officials.Description = *v.Description
		}
		resultDetail = newUser.Officials
	}

	id, err := user.SignUp(newUser)
	if err != nil {
		return nil, fmt.Errorf("Error While Signup")
	}

	phoneResult := model.Phone(newUser.Phone)
	profile := &model.Profile{
		ID:       model.ObjectID(id),
		Name:     newUser.Name,
		Nickname: newUser.Nickname,
		Phone:    &phoneResult,
		Status:   user.StatusToEnum(user.StatusWait),
		Detail:   user.DetailToUnion(resultDetail),
	}

	token, err := session.CreateToken(id, newUser.Role, newUser.Permissions)
	if err != nil {
		return nil, err
	}

	result := &model.ProfileWithToken{
		Profile: profile,
		Token:   token,
	}
	return result, nil
}

func (r *mutationResolver) CreateCategory(ctx context.Context, input model.NewCategory) (model.ObjectID, error) {
	convert := func(i model.UserRole) user.Role {
		switch i {
		case model.UserRoleStudent:
			return user.RoleStudent
		case model.UserRoleTeacher:
			return user.RoleTeacher
		case model.UserRoleOfficials:
			return user.RoleOfficials
		}
		return user.RoleOfficials
	}
	category := &post.Category{
		Name:          input.Name,
		ReqPermission: input.ReqPermission,
		AnonAble:      input.AnonAble,
		Description:   input.Description,
	}
	for _, v := range input.ReadAbleRole {
		category.ReadAbleRole = append(category.ReadAbleRole, convert(v))
	}
	for _, v := range input.WriteAbleRole {
		category.WriteAbleRole = append(category.WriteAbleRole, convert(v))
	}

	id, err := post.NewCategory(category)
	if err != nil {
		return model.ObjectID(primitive.NilObjectID), err
	}
	return model.ObjectID(id), nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (model.ObjectID, error) {
	user := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategory(primitive.ObjectID(input.Category)); err != nil {
		return model.ObjectID(primitive.NilObjectID), err
	} else {
		if ok := post.CheckUserPermission("write", category, user.Role, user.Permission); !ok {
			return model.ObjectID(primitive.NilObjectID), myerr.New(myerr.ErrPermission, "")
		}
	}

	anon := false
	if *input.Anon {
		anon = true
	}
	id, err := post.NewPost(primitive.ObjectID(input.Category), user.ID, input.Title, input.Content, anon)

	if conf.Discord() != nil {
		for _, d := range conf.Discord().SubPost {
			if primitive.ObjectID(input.Category).Hex() == d.CategoryID {
				discord.SendEmbed(d.DiscordChannelID, &discordgo.MessageEmbed{
					Title:       input.Title,
					Description: input.Content,
					Color:       0xbedbe9,
					Footer: &discordgo.MessageEmbedFooter{
						Text: d.CategoryName,
					},
					Timestamp: time.Now().Format(time.RFC3339),
				})
			}
		}
	}

	return model.ObjectID(id), err
}

func (r *mutationResolver) LikePost(ctx context.Context, input model.LikePostInput) (*string, error) {
	user := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategoryByPost(primitive.ObjectID(input.Post)); err != nil {
		return nil, err
	} else {
		if ok := post.CheckUserPermission("read", category, user.Role, user.Permission); !ok {
			return nil, myerr.New(myerr.ErrPermission, "")
		}
	}
	err := post.PostLike(primitive.ObjectID(input.Post), user.ID, input.Status)
	return nil, err
}

func (r *mutationResolver) AddComment(ctx context.Context, input model.NewComment) (model.ObjectID, error) {
	user := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategoryByPost(primitive.ObjectID(input.Post)); err != nil {
		return model.ObjectID(primitive.NilObjectID), err
	} else {
		if ok := post.CheckUserPermission("read", category, user.Role, user.Permission); !ok {
			return model.ObjectID(primitive.NilObjectID), myerr.New(myerr.ErrPermission, "")
		}
	}
	anon := false
	if input.Anon != nil {
		anon = *input.Anon
	}
	id, err := post.NewComment(primitive.ObjectID(input.Post), user.ID, input.Content, anon)
	return model.ObjectID(id), err
}

func (r *mutationResolver) DeleteComment(ctx context.Context, postID model.ObjectID, commentID model.ObjectID) (string, error) {
	return "", post.DeleteComment(primitive.ObjectID(postID), primitive.ObjectID(commentID))
}

func (r *mutationResolver) AddCalendar(ctx context.Context, input model.NewCalendar) (model.ObjectID, error) {
	objID, err := info.NewCalendar(uint(input.Year), uint(input.Month), uint(input.Day), input.Title, input.Description, input.Icon)
	return model.ObjectID(objID), err
}

func (r *mutationResolver) DeleteCalendar(ctx context.Context, target model.ObjectID) (string, error) {
	return "", info.DeleteCalendar(primitive.ObjectID(target))
}

func (r *mutationResolver) InsertSchedule(ctx context.Context, input []*model.UpdateSchedule) (string, error) {
	var data []*info.Schedule
	for _, d := range input {
		data = append(data, &info.Schedule{
			ID:          primitive.NewObjectID(),
			Grade:       d.Grade,
			Class:       d.Class,
			Dow:         d.Dow,
			Period:      d.Period,
			Subject:     d.Subject,
			Teacher:     d.Teacher,
			Description: d.Description,
			ClassRoom:   d.ClassRoom,
		})
	}

	return "", info.InsertSchedules(data)
}

func (r *mutationResolver) UpdateSchedule(ctx context.Context, input model.UpdateSchedule) (string, error) {
	data := info.UpdateScheduleInput{
		input.Grade,
		input.Class,
		input.Dow,
		input.Period,
		input.Subject,
		input.Teacher,
		input.Description,
		input.ClassRoom,
	}
	return "", info.UpdateSchedule(&data)
}

func (r *mutationResolver) DeleteSchedule(ctx context.Context, target model.ScheduleDelFilter) (string, error) {
	return "", info.DeleteSchedule(uint(target.Grade), uint(target.Class), uint(target.Dow), uint(target.Period))
}

func (r *mutationResolver) UpdateEmailAliases(ctx context.Context, input model.EmailAliasesInput) (string, error) {
	if ok := strings.HasSuffix(input.From, "@osang.xyz"); !ok {
		return "", myerr.New(myerr.ErrBadRequest, "invalid email format")
	}
	userData := ctx.Value("user").(*user.User)
	origin := userData.EmailAliases.From

	if err := user.MailAliasesUpdate(userData.ID, input.From, input.To, origin != "" && origin == input.From); err != nil {
		return "", err
	}

	if origin != "" && origin != input.From {
		if err := user.MailAliasesRemove(origin); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (r *mutationResolver) DeleteEmailAliases(ctx context.Context) (string, error) {
	userData := ctx.Value("user").(*user.User)
	if err := user.MailAliasesRemove(userData.EmailAliases.From); err != nil {
		return "", err
	}
	if err := user.MailAliasesRemoveDB(userData.ID); err != nil {
		return "", err
	}
	return "", nil
}

func (r *queryResolver) MyProfile(ctx context.Context) (*model.Profile, error) {
	userData := ctx.Value("user").(*user.User)
	return user.UserToGqlType(userData), nil
}

func (r *queryResolver) SchoolMeal(ctx context.Context, filter *model.SchoolMealFilter) ([]*model.SchoolMeal, error) {
	if filter == nil {
		filter = &model.SchoolMealFilter{}
	}

	return neis.GetSchoolMeal(filter)
}

func (r *queryResolver) Post(ctx context.Context, id model.ObjectID, comment *model.CommentFilter) (*model.Post, error) {
	userData := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategoryByPost(primitive.ObjectID(id)); err != nil {
		return nil, err
	} else {
		if ok := post.CheckUserPermission("read", category, userData.Role, userData.Permission); !ok {
			return nil, myerr.New(myerr.ErrPermission, "")
		}
	}

	loadPost := true
	offset := 0
	limit := 20
	if comment != nil {
		if comment.LoadOnlyComment != nil {
			loadPost = !*comment.LoadOnlyComment
		}
		if comment.Offset != nil {
			offset = *comment.Offset
		}
		if comment.Limit != nil {
			limit = *comment.Limit
		}
	}

	data, err := post.GetPost(primitive.ObjectID(id), userData.ID, loadPost, offset, limit)
	if err != nil {
		return nil, err
	}

	resultComment := []*model.Comment{}
	for _, v := range data.Comment {
		resultComment = append(resultComment, &model.Comment{
			ID:       model.ObjectID(v.ID),
			Author:   user.UserToGqlType(v.AuthorData),
			Content:  v.Content,
			CreateAt: model.Timestamp(v.CreateAt),
			UpdateAt: model.Timestamp(v.UpdateAt),
			Status:   post.StatusToGqlType(v.Status),
		})
	}
	likeCnt := data.LikeCnt
	isLike := data.IsLike
	content := data.Content
	return &model.Post{
		ID: model.ObjectID(data.ID),
		Category: &model.Category{
			ID:            model.ObjectID(data.CategoryData.ID),
			Name:          data.CategoryData.Name,
			ReqPermission: data.CategoryData.ReqPermission,
			AnonAble:      data.CategoryData.AnonAble,
			WriteAbleRole: user.RoleListToGql(data.CategoryData.WriteAbleRole),
			ReadAbleRole:  user.RoleListToGql(data.CategoryData.ReadAbleRole),
		},
		Like:     &likeCnt,
		IsLike:   &isLike,
		Author:   user.UserToGqlType(data.AuthorData),
		Title:    data.Title,
		Content:  &content,
		CreateAt: model.Timestamp(data.CreateAt),
		UpdateAt: model.Timestamp(data.UpdateAt),
		Comment:  resultComment,
		Status:   post.StatusToGqlType(data.Status),
	}, nil
}

func (r *queryResolver) Posts(ctx context.Context, categoryID model.ObjectID, offset *int, limit *int) ([]*model.Post, error) {
	userData := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategory(primitive.ObjectID(categoryID)); err != nil {
		return nil, err
	} else {
		if ok := post.CheckUserPermission("read", category, userData.Role, userData.Permission); !ok {
			return nil, myerr.New(myerr.ErrPermission, "")
		}
	}

	realOffset := utils.IfInt(offset != nil, offset, 0)
	realLimit := utils.IfInt(limit != nil, limit, 20)

	data, err := post.GetPosts(primitive.ObjectID(categoryID), *realOffset, *realLimit)
	if err != nil {
		return nil, err
	}

	result := []*model.Post{}
	for _, d := range data {
		result = append(result, &model.Post{
			ID: model.ObjectID(d.ID),
			Category: &model.Category{
				ID:            model.ObjectID(d.CategoryData.ID),
				Name:          d.CategoryData.Name,
				ReqPermission: d.CategoryData.ReqPermission,
				AnonAble:      d.CategoryData.AnonAble,
				WriteAbleRole: user.RoleListToGql(d.CategoryData.WriteAbleRole),
				ReadAbleRole:  user.RoleListToGql(d.CategoryData.ReadAbleRole),
			},
			Author:   user.UserToGqlType(d.AuthorData),
			Title:    d.Title,
			CreateAt: model.Timestamp(d.CreateAt),
			UpdateAt: model.Timestamp(d.UpdateAt),
		})
	}

	return result, nil
}

func (r *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
	var result []*model.Category
	result = []*model.Category{}
	data, err := post.GetAllCategory()
	if err != nil {
		return result, err
	}
	for _, d := range data {
		result = append(result, &model.Category{
			ID:            model.ObjectID(d.ID),
			Name:          d.Name,
			Description:   d.Description,
			ReqPermission: d.ReqPermission,
			AnonAble:      d.AnonAble,
			WriteAbleRole: user.RoleListToGql(d.WriteAbleRole),
			ReadAbleRole:  user.RoleListToGql(d.ReadAbleRole),
		})
	}
	return result, nil
}

func (r *queryResolver) Calendar(ctx context.Context, filter model.CalendarFilter) ([]*model.Calendar, error) {
	data, err := info.FindCalendar(filter.Year, filter.Month)
	if err != nil {
		return nil, err
	}
	var result []*model.Calendar
	for _, d := range data {
		result = append(result, &model.Calendar{
			ID:          model.ObjectID(d.ID),
			Year:        d.Year,
			Month:       d.Month,
			Day:         d.Day,
			Title:       d.Title,
			Description: d.Description,
			Icon:        d.Icon,
		})
	}
	return result, nil
}

func (r *queryResolver) Schedule(ctx context.Context, filter model.ScheduleFilter) ([]*model.Schedule, error) {
	data, err := info.FindSchedule(filter.Grade, filter.Class, filter.Dow)
	if err != nil {
		return nil, err
	}
	var result []*model.Schedule
	for _, d := range data {
		result = append(result, &model.Schedule{
			Dow:         d.Dow,
			Period:      d.Period,
			Grade:       d.Grade,
			Class:       d.Class,
			Subject:     d.Subject,
			Teacher:     d.Teacher,
			Description: d.Description,
			ClassRoom:   d.ClassRoom,
		})
	}
	return result, nil
}

func (r *queryResolver) EmailAliases(ctx context.Context) (*model.EmailAliases, error) {
	userData := ctx.Value("user").(*user.User)
	return &model.EmailAliases{
		From: userData.EmailAliases.From,
		To:   userData.EmailAliases.To,
	}, nil
}

func (r *queryResolver) HomepageList(ctx context.Context, filter *model.HomepageListFilter) ([]*model.HomepageListType, error) {
	var url osangdata.Url
	switch filter.Board {
	case model.HomepageBoardNotice:
		url = osangdata.UrlNotice
	case model.HomepageBoardAdministration:
		url = osangdata.UrlAdministration
	case model.HomepageBoardEvaluationPlan:
		url = osangdata.UrlEvaluationPlan
	case model.HomepageBoardPrints:
		url = osangdata.UrlPrints
	case model.HomepageBoardRule:
		url = osangdata.UrlRule
	}

	data, err := osangdata.CrawlList(url, filter.Page)
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	var result []*model.HomepageListType
	for _, d := range data {
		result = append(result, &model.HomepageListType{
			ID:        d.ID,
			Number:    d.Number,
			Title:     d.Title,
			WrittenBy: d.WrittenBy,
			CreateAt:  model.Timestamp(d.CreateAt),
		})
	}
	return result, nil
}

func (r *queryResolver) HomepageDetail(ctx context.Context, filter *model.HomepageDetailFilter) (*model.HomepageDetailType, error) {
	var url osangdata.Url
	switch filter.Board {
	case model.HomepageBoardNotice:
		url = osangdata.UrlNotice
	case model.HomepageBoardAdministration:
		url = osangdata.UrlAdministration
	case model.HomepageBoardEvaluationPlan:
		url = osangdata.UrlEvaluationPlan
	case model.HomepageBoardPrints:
		url = osangdata.UrlPrints
	case model.HomepageBoardRule:
		url = osangdata.UrlRule
	}

	data, err := osangdata.CrawlPage(url, filter.ID)
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	var files []*model.HomepageFileType
	for _, d := range data.Files {
		files = append(files, &model.HomepageFileType{
			Name:     d.Name,
			Download: d.Download,
			Preview:  d.Preview,
		})
	}
	return &model.HomepageDetailType{
		ID:        data.ID,
		Title:     data.Title,
		WrittenBy: data.WrittenBy,
		Content:   data.Content,
		Images:    data.Images,
		Files:     files,
		CreateAt:  model.Timestamp(data.CreateAt),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

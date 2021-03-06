package neis

import (
	"fmt"
	"time"

	"github.com/imroc/req"
	"github.com/osang-school/backend/graph/model"
	"github.com/osang-school/backend/internal/conf"
)

type (
	Cafeteria struct {
		Type     string
		Calorie  string
		Content  string
		Nutrient string
		Origin   string
		Date     string
	}
)

func GetSchoolMeal(filter *model.SchoolMealFilter) ([]*model.SchoolMeal, error) {
	emptyList := []*model.SchoolMeal{}
	param := req.Param{
		"KEY":                conf.Get().NeisAPIKey,
		"Type":               "json",
		"pIndex":             1,
		"pSize":              100,
		"ATPT_OFCDC_SC_CODE": "R10",
		"SD_SCHUL_CODE":      "8750198",
	}
	if filter.DateStart != nil && filter.DateEnd != nil {
		param["MLSV_FROM_YMD"] = time.Time(*filter.DateStart).Format("20060102")
		param["MLSV_TO_YMD"] = time.Time(*filter.DateEnd).Format("20060102")
	} else if filter.DateStart != nil {
		param["MLSV_YMD"] = time.Time(*filter.DateStart).Format("20060102")
	} else {
		param["MLSV_YMD"] = time.Now().Format("20060102")
	}

	if filter.Type != nil {
		switch *filter.Type {
		case model.SchoolMealTypeBreakfast:
			param["MMEAL_SC_CODE"] = 1
		case model.SchoolMealTypeLunch:
			param["MMEAL_SC_CODE"] = 2
		case model.SchoolMealTypeDinner:
			param["MMEAL_SC_CODE"] = 3
		}
	}

	res, err := req.Get("https://open.neis.go.kr/hub/mealServiceDietInfo", param)
	if err != nil {
		return emptyList, err
	}
	var data map[string]interface{}
	res.ToJSON(&data)
	if d, ok := data["RESULT"].(map[string]interface{}); ok {
		if _, ok := d["MESSAGE"].(string); ok {
			return emptyList, nil
		}
	}

	defer func() {
		if v := recover(); v != nil {
			fmt.Println(v)
		}
	}()

	target, _ := data["mealServiceDietInfo"].([]interface{})[1].(map[string]interface{})["row"].([]interface{})

	var result []*model.SchoolMeal
	for _, d := range target {
		v := d.(map[string]interface{})
		date, err := time.Parse("20060102", v["MLSV_FROM_YMD"].(string))
		if err != nil {
			date = time.Now()
		}
		newField := model.SchoolMeal{
			Calorie:  v["CAL_INFO"].(string),
			Content:  v["DDISH_NM"].(string),
			Nutrient: v["NTR_INFO"].(string),
			Origin:   v["ORPLC_INFO"].(string),
			Date:     model.Timestamp(date),
		}
		switch v["MMEAL_SC_CODE"].(string) {
		case "1":
			newField.Type = model.SchoolMealTypeBreakfast
		case "2":
			newField.Type = model.SchoolMealTypeLunch
		case "3":
			newField.Type = model.SchoolMealTypeDinner
		}
		result = append(result, &newField)
	}

	return result, nil
}

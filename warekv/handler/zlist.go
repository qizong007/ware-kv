package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

type SetZListParam struct {
	Key        string           `json:"k"`
	Val        []util.SlElement `json:"v"`
	ExpireTime int64            `json:"expire_time" binding:"-"`
}

type AddZListParam struct {
	Element  *util.SlElement   `json:"e"`
	Elements *[]util.SlElement `json:"elements"`
}

type RemoveZListByScoreParam struct {
	Scores *[]float64
	Score  *float64
	Min    *float64
	Max    *float64
}

func SetZList(c *gin.Context) {
	param := SetZListParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	key := storage.MakeKey(param.Key)

	newVal := ds.MakeZList(param.Val)

	set(key, newVal, param.ExpireTime)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetZListLen(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2ZList(val).GetLen(),
	})
}

func GetZListByPos(c *gin.Context) {
	posStr := c.Param("pos")

	if posStr == "" {
		log.Println("pos is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "pos is <nil>!",
		})
		return
	}

	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}

	zList := ds.Value2ZList(val)

	pos, err := util.Str2Int(posStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	res, err := zList.GetElementAt(pos)
	if err != nil {
		log.Println("GetElementAt fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  res,
	})
}

func GetZListBetween(c *gin.Context) {
	leftStr := c.Param("left")
	rightStr := c.Param("right")

	if leftStr == "" && rightStr == "" || leftStr != "" && rightStr == "" || leftStr == "" && rightStr != "" {
		log.Println("GetZListBetween left or right is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "left or right is <nil>!",
		})
		return
	}

	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}

	zList := ds.Value2ZList(val)

	left, err := util.Str2Int(leftStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	right, err := util.Str2Int(rightStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	res, err := zList.GetListBetween(left, right)
	if err != nil {
		log.Println("GetListBetween fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  res,
	})
}

func GetZListStartAt(c *gin.Context) {
	leftStr := c.Param("left")

	if leftStr == "" {
		log.Println("GetZListStartAt left is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "left is <nil>!",
		})
		return
	}

	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}

	zList := ds.Value2ZList(val)

	left, err := util.Str2Int(leftStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	res, err := zList.GetListStartWith(left)
	if err != nil {
		log.Println("GetZListStartAt fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  res,
	})
}

func GetZListEndAt(c *gin.Context) {
	rightStr := c.Param("right")

	if rightStr == "" {
		log.Println("GetZListEndAt right is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "right is <nil>!",
		})
		return
	}

	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}

	zList := ds.Value2ZList(val)

	right, err := util.Str2Int(rightStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	res, err := zList.GetListEndAt(right)
	if err != nil {
		log.Println("GetZListEndAt fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  res,
	})
}

func AddZList(c *gin.Context) {
	param := AddZListParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if param.Element == nil && param.Elements == nil {
		log.Println("AddZListParam all <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param all <nil>!",
		})
		return
	}

	key, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}

	zList := ds.Value2ZList(val)

	if param.Element != nil {
		zList.Add([]util.SlElement{*param.Element})
		setNotify(key, zList)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	zList.Add(*param.Elements)
	setNotify(key, zList)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func RemoveZListByScore(c *gin.Context) {
	param := RemoveZListByScoreParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if param.Score == nil && param.Scores == nil && param.Min == nil && param.Max == nil {
		log.Println("RemoveZListByScoreParam all <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param all <nil>!",
		})
		return
	}

	if param.Min != nil && param.Max == nil || param.Min == nil && param.Max != nil {
		log.Println("RemoveZListByScoreParam should have min and max at same time!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Check your min and max!",
		})
		return
	}

	key, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}

	zList := ds.Value2ZList(val)

	if param.Score != nil {
		zList.Remove(*param.Score)
		setNotify(key, zList)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	if param.Scores != nil {
		zList.RemoveScores(*param.Scores)
		setNotify(key, zList)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	// param.Min != nil && param.Max != nil
	if err = zList.RemoveInScore(*param.Min, *param.Max); err != nil {
		log.Println("RemoveInScore Fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ScopeError,
		})
		return
	}

	setNotify(key, zList)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

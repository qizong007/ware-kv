package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	zlist "github.com/qizong007/ware-kv/warekv/util"
	"log"
	"time"
)

type SetZListParam struct {
	Key        string            `json:"k"`
	Val        []zlist.SlElement `json:"v"`
	ExpireTime int64             `json:"expire_time" binding:"-"`
}

type AddZListParam struct {
	Element  *zlist.SlElement   `json:"e"`
	Elements *[]zlist.SlElement `json:"elements"`
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

	cmd := tracker.NewCreateCommand(param.Key, tracker.ZListStruct, param.Val, newVal.CreateTime, param.ExpireTime)
	set(key, newVal, param.ExpireTime, cmd)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetZListLen(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  storage.Value2ZList(val).GetLen(),
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

	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	pos, err := util.Str2Int(posStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
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

	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	left, err := util.Str2Int(leftStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}
	right, err := util.Str2Int(rightStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
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

	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	left, err := util.Str2Int(leftStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
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

	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	right, err := util.Str2Int(rightStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
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

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	if param.Element != nil {
		slElement := []zlist.SlElement{*param.Element}
		wal(tracker.NewModifyCommand(key.GetKey(), tracker.ZListAdd, time.Now().Unix(), slElement))
		zList.Add(slElement)
		setNotify(key, zList)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	wal(tracker.NewModifyCommand(key.GetKey(), tracker.ZListAdd, time.Now().Unix(), *param.Elements))
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

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	if param.Score != nil {
		wal(tracker.NewModifyCommand(key.GetKey(), tracker.ZListRemoveScore, time.Now().Unix(), *param.Score))
		zList.RemoveScore(*param.Score)
		setNotify(key, zList)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	if param.Scores != nil {
		wal(tracker.NewModifyCommand(key.GetKey(), tracker.ZListRemoveScores, time.Now().Unix(), *param.Scores))
		zList.RemoveScores(*param.Scores)
		setNotify(key, zList)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	// param.Min != nil && param.Max != nil
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.ZListRemoveInScore, time.Now().Unix(), *param.Min, *param.Max))
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

func GetZListBetweenScores(c *gin.Context) {
	minStr := c.Param("min")
	maxStr := c.Param("max")

	if minStr == "" && maxStr == "" || minStr != "" && maxStr == "" || minStr == "" && maxStr != "" {
		log.Println("GetZListBetweenScores min or max is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "min or max is <nil>!",
		})
		return
	}

	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, zlist.ZListDS) {
		return
	}

	zList := storage.Value2ZList(val)

	min, err := util.Str2Float64(minStr)
	if err != nil {
		log.Println("Str2Float64 fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}
	max, err := util.Str2Float64(maxStr)
	if err != nil {
		log.Println("Str2Float64 fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}
	res, err := zList.GetListInScore(min, max)
	if err != nil {
		log.Println("GetListInScore fail", err)
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

package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/util"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
)

type SetListParam struct {
	Key        string        `json:"k"`
	Val        []interface{} `json:"v"`
	ExpireTime int64         `json:"expire_time" binding:"-"`
}

type AddListParam struct {
	Element  interface{}   `json:"e"`
	Elements []interface{} `json:"elements"`
}

type RemoveListElementParam struct {
	Val interface{}
	Pos *int
}

func SetList(c *gin.Context) {
	param := SetListParam{}
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
	newVal := ds.MakeList(param.Val)

	set(key, newVal, param.ExpireTime)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetListLen(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2List(val).GetLen(),
	})
}

func GetListByPos(c *gin.Context) {
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

	list := ds.Value2List(val)

	pos, err := util.Str2Int(posStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}
	res, err := list.GetElementAt(pos)
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

func GetListBetween(c *gin.Context) {
	leftStr := c.Param("left")
	rightStr := c.Param("right")

	if leftStr == "" && rightStr == "" || leftStr != "" && rightStr == "" || leftStr == "" && rightStr != "" {
		log.Println("GetListBetween left or right is <nil>!")
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

	list := ds.Value2List(val)

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
	res, err := list.GetListBetween(left, right)
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

func GetListStartAt(c *gin.Context) {
	leftStr := c.Param("left")

	if leftStr == "" {
		log.Println("GetListStartAt left is <nil>!")
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

	list := ds.Value2List(val)

	left, err := util.Str2Int(leftStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}
	res, err := list.GetListStartWith(left)
	if err != nil {
		log.Println("GetListStartAt fail", err)
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

func GetListEndAt(c *gin.Context) {
	rightStr := c.Param("right")

	if rightStr == "" {
		log.Println("GetListEndAt right is <nil>!")
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

	list := ds.Value2List(val)

	right, err := util.Str2Int(rightStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}
	res, err := list.GetListEndAt(right)
	if err != nil {
		log.Println("GetListEndAt fail", err)
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

func AddList(c *gin.Context) {
	param := AddListParam{}
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
		log.Println("AddListParam all <nil>!")
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

	list := ds.Value2List(val)

	if param.Element != nil {
		list.Append([]interface{}{param.Element})
		setNotify(key, list)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	list.Append(param.Elements)
	setNotify(key, list)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func RemoveListElement(c *gin.Context) {
	param := RemoveListElementParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if param.Pos == nil && param.Val == nil {
		log.Println("RemoveListElement all <nil>!")
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

	list := ds.Value2List(val)

	if param.Val != nil {
		list.RemoveVal(param.Val)
		setNotify(key, list)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.Success,
		})
		return
	}

	// param.Pos != nil
	list.RemoveAt(*param.Pos)
	setNotify(key, list)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
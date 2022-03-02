package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/util"
)

type SetBitmapParam struct {
	ExpireTime int64 `json:"expire_time" binding:"-"`
}

func SetBitmap(c *gin.Context) {
	param := SetListParam{}
	_ = c.ShouldBindJSON(&param)

	numStr := c.Param("num")
	if numStr == "" {
		paramNull(c, "num")
		return
	}
	num, err := util.Str2Int(numStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}

	var bitmap *ds.Bitmap

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}

	if isValNil(val) {
		// new
		bitmap = ds.MakeBitmap()
		bitmap.SetBit(num)
		set(key, bitmap, param.ExpireTime)
	} else {
		// update
		if !isKVEffective(c, val) {
			return
		}
		bitmap = ds.Value2Bitmap(val)
		bitmap.SetBit(num)
		setNotify(key, bitmap)
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetBitmapLen(c *gin.Context) {
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
		Val:  ds.Value2Bitmap(val).GetLen(),
	})
}

func GetBitmapBit(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	numStr := c.Param("num")
	if numStr == "" {
		paramNull(c, "num")
		return
	}
	num, err := util.Str2Int(numStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2Bitmap(val).GetBit(num),
	})
}

func GetBitCount(c *gin.Context) {
	leftStr := c.Param("left")
	rightStr := c.Param("right")

	if leftStr == "" && rightStr == "" || leftStr != "" && rightStr == "" || leftStr == "" && rightStr != "" {
		log.Println("GetBitCount left or right is <nil>!")
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

	bitmap := ds.Value2Bitmap(val)

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

	res, err := bitmap.GetBitCount(left, right)
	if err != nil {
		log.Println("GetBitCount fail", err)
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

func ClearBitmap(c *gin.Context) {
	numStr := c.Param("num")
	if numStr == "" {
		paramNull(c, "num")
		return
	}
	num, err := util.Str2Int(numStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
		})
		return
	}

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}

	bitmap := ds.Value2Bitmap(val)
	bitmap.ClearBit(num)

	setNotify(key, bitmap)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

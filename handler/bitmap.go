package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	dstype "github.com/qizong007/ware-kv/warekv/util"
	"log"
	"time"
)

type SetBitmapParam struct {
	ExpireTime int64 `json:"expire_time" binding:"-"`
}

func SetBitmap(c *gin.Context) {
	param := SetBitmapParam{}
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
		cmd := tracker.NewCreateCommand(key.GetKey(), tracker.BitmapStruct, num, bitmap.CreateTime, param.ExpireTime)
		set(key, bitmap, param.ExpireTime, cmd)
	} else {
		// update
		if !isKVEffective(c, val) {
			return
		}
		if !isKVTypeCorrect(c, val, dstype.BitmapDS) {
			return
		}

		bitmap = storage.Value2Bitmap(val)
		wal(tracker.NewModifyCommand(key.GetKey(), tracker.BitmapSet, time.Now().Unix(), num))
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
	if !isKVTypeCorrect(c, val, dstype.BitmapDS) {
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  storage.Value2Bitmap(val).GetLen(),
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
	if !isKVTypeCorrect(c, val, dstype.BitmapDS) {
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
		Val:  storage.Value2Bitmap(val).GetBit(num),
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
	if !isKVTypeCorrect(c, val, dstype.BitmapDS) {
		return
	}

	bitmap := storage.Value2Bitmap(val)

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
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.BitmapDS) {
		return
	}

	bitmap := storage.Value2Bitmap(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.BitmapClear, time.Now().Unix(), num))
	bitmap.ClearBit(num)

	setNotify(key, bitmap)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

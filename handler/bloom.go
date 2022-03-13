package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/util"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
)

type BloomSpecificOption struct {
	M uint64
	K uint64
}

type SetBloomSpecificParam struct {
	Key        string              `json:"k"`
	Val        BloomSpecificOption `json:"v"`
	ExpireTime int64               `json:"expire_time" binding:"-"`
}

type BloomFuzzyOption struct {
	N  uint
	Fp float64
}

type SetBloomFuzzyParam struct {
	Key        string           `json:"k"`
	Val        BloomFuzzyOption `json:"v"`
	ExpireTime int64            `json:"expire_time" binding:"-"`
}

type AddBloomParam struct {
	Key string
}

func SetBloomSpecific(c *gin.Context) {
	param := SetBloomSpecificParam{}
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
	option := ds.BloomFilterSpecificOption{
		M: param.Val.M,
		K: param.Val.K,
	}
	newVal := ds.MakeBloomFilterSpecific(option)

	set(key, newVal, param.ExpireTime)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func SetBloomFuzzy(c *gin.Context) {
	param := SetBloomFuzzyParam{}
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
	option := ds.BloomFilterFuzzyOption{
		N:  param.Val.N,
		Fp: param.Val.Fp,
	}
	newVal := ds.MakeBloomFilterFuzzy(option)

	set(key, newVal, param.ExpireTime)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetBloomSize(c *gin.Context) {
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
		Val:  ds.Value2BloomFilter(val).GetSize(),
	})
}

func ClearBloom(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	filter := ds.Value2BloomFilter(val)

	filter.ClearAll()
	setNotify(key, filter)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func TestBloom(c *gin.Context) {
	data := c.Param("data")
	if data == "" {
		paramNull(c, "data")
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

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2BloomFilter(val).Test(data),
	})
}

func GetBloomFalseRate(c *gin.Context) {
	nStr := c.Param("n")
	if nStr == "" {
		log.Println("n is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "n is <nil>!",
		})
		return
	}

	n, err := util.Str2Uint(nStr)
	if err != nil {
		log.Println("Str2Int fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.TypeTransformError,
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
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2BloomFilter(val).EstimateFalsePositiveRate(n),
	})
}

func AddBloom(c *gin.Context) {
	param := AddBloomParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	data := param.Key

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	filter := ds.Value2BloomFilter(val)
	filter.Add(data)
	setNotify(key, filter)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

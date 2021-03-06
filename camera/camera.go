package camera

import (
	"fmt"
	tool "github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv/manager"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/util"
	"golang.org/x/crypto/blake2b"
	"io/ioutil"
	"log"
	"math"
	"sync"
	"time"
)

// The 'dump' file looks like:
// | magic number (6 bytes) | -> 'warekv'
// |    version (3 bytes)   |
// |     meta (64 bytes)    |
// |   table flag (1 byte)  | -> to prove that the things below is about 'KVTable'
// |    keys num (4 bytes)  |
// |    table kv pairs      |
// | subscribe center flag  | -> to prove that the things below is about 'Subscribe Center'
// |    keys num (4 bytes)  |
// |     sub kv pairs       |
// |          ...           |
// |          ...           |
// |   check sum (64 bytes) |
//
// The 'meta' be like:
// | magic switch (1 byte) | create time (8 byte) | keep bytes |
// [ magic switch ]: 0 0 0 0 0 0 0 is_zip
//
// The 'kv pairs' be like:
// |    type (1 byte)   |
// |  key len (4 bytes) |   key (key len bytes)   |
// | base len (4 bytes) |  base (base len bytes)  |
// | value len (4 bytes)| value (value len bytes) |
// ( the 'key' is just a 'string' )

type Camera struct {
	filePath   string
	lock       sync.Mutex
	ticker     *time.Ticker
	closer     chan bool
	isZip      bool
	isOpen     bool
	createTime int64
}

var (
	camera *Camera
)

const (
	magicHead           = "warekv"
	magicHeadLen        = len(magicHead)
	wareKVVersion       = tool.WareKVVersionForCamera
	wareKVVersionLen    = len(wareKVVersion)
	totalHeadLen        = magicHeadLen + wareKVVersionLen + metaDataLen
	metaDataLen         = 64
	checkSumLen         = blake2b.Size
	zipFlag             = 1 << 0
	defaultCameraPath   = "./photo"
	defaultTickInterval = 15                 // minutes
	tickIntervalMin     = 5                  // minutes
	tickIntervalMax     = int(math.MaxInt32) // minutes
)

type CameraOption struct {
	Open             bool   `yaml:"Open"`
	IsZip            bool   `yaml:"IsZip"`
	FilePath         string `yaml:"FilePath"`
	SaveTickInterval uint   `yaml:"SaveTickInterval"`
}

func DefaultOption() *CameraOption {
	return &CameraOption{
		Open:             true,
		IsZip:            false,
		FilePath:         defaultCameraPath,
		SaveTickInterval: defaultTickInterval,
	}
}

func NewCamera(option *CameraOption) *Camera {
	filePath := defaultCameraPath
	tickInterval := uint(defaultTickInterval)
	isZip := false
	if option != nil {
		if !option.Open {
			camera = &Camera{isOpen: false}
			return camera
		}
		filePath = option.FilePath
		isZip = option.IsZip
		tickInterval = option.SaveTickInterval
		tickInterval = uint(util.SetIfHitLimit(int(tickInterval), tickIntervalMin, tickIntervalMax))
	}
	camera = &Camera{
		filePath: filePath,
		isOpen:   true,
		isZip:    isZip,
	}
	camera.closer = make(chan bool)
	camera.ticker = time.NewTicker(time.Duration(tickInterval) * time.Minute)
	camera.start()
	return camera
}

func GetCamera() *Camera {
	return camera
}

func (c *Camera) start() {
	go c.scheduledSave()
	log.Println("Camera's Save worker starts working...")
}

func (c *Camera) Close() {
	if !c.isOpen {
		return
	}
	c.closer <- true
}

func (c *Camera) scheduledSave() {
	for {
		select {
		case <-c.ticker.C:
			p := []storage.Photographer{storage.GlobalTable, manager.GetSubscribeCenter()}
			c.TakePhotos(p, c.isZip)
		case <-c.closer:
			c.ticker.Stop()
			close(c.closer)
			log.Println("Camera's Save worker starts working...")
			return
		}
	}
}

// TakePhotos encoding to bin file, just like 'take photos'
func (c *Camera) TakePhotos(p []storage.Photographer, needZip bool) {
	if !c.isOpen {
		log.Println("You don't have camera, can't take photos...")
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	data := make([]byte, 0)

	// magic number (6 bytes)
	magicHeadBytes := []byte(magicHead)

	// version 0.0.1 (3 bytes)
	versionBytes := []byte(wareKVVersion)

	// meta (64 bytes)
	metaBytes := make([]byte, metaDataLen)
	// fill the magic switch
	metaBytes[0] = c.generateMagicSwitch()
	// fill the time
	createTimeBytes := util.Int64ToBytes(time.Now().Unix())
	copy(metaBytes[1:], createTimeBytes)

	data = append(data, magicHeadBytes...)
	data = append(data, versionBytes...)
	data = append(data, metaBytes...)

	content := make([]byte, 0)
	for _, kvPairs := range p {
		view := kvPairs.View()
		content = append(content, view...)
	}
	if needZip {
		tool.ZipBytes(content)
	}

	data = append(data, content...)

	// check sum (64 bytes)
	checkSum := blake2b.Sum512(data)
	data = append(data, checkSum[:]...)

	// save the bin file
	c.save(data)

	return
}

func (c *Camera) generateMagicSwitch() byte {
	// [ magic switch ]: 0 0 0 0 0 0 0 is_zip
	magicSwitch := uint8(0)
	if c.isZip {
		magicSwitch = magicSwitch | zipFlag
	}
	return magicSwitch
}

func (c *Camera) save(data []byte) {
	err := ioutil.WriteFile(c.filePath, data, 0666)
	if err != nil {
		panic(fmt.Sprintf("TakePhotos Fail: %v", err))
		return
	}
}

// DevelopPhotos decoding the bin file, load the data
func (c *Camera) DevelopPhotos() {
	if !c.isOpen {
		return
	}

	start := time.Now()
	data, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		log.Println("DevelopPhotos ReadFile Failed:", err)
		return
	}
	if !ifCheckHeadOK(data) {
		log.Println("DevelopPhotos check head Fail!")
		return
	}
	if !ifCheckSumOK(data) {
		log.Println("DevelopPhotos check sum Fail!")
		return
	}

	meta := reduceMetaInfo(data)
	c.createTime = meta.CreateTime

	content := data[totalHeadLen : len(data)-checkSumLen]
	if meta.IsZip {
		content = tool.UnzipBytes(content)
	}

	reduceContent(content)
	log.Printf("Camera finished loading in %s...\n", time.Since(start).String())
}

func ifCheckHeadOK(data []byte) bool {
	magicHeadBytes := data[:magicHeadLen]
	versionBytes := data[magicHeadLen : magicHeadLen+wareKVVersionLen]
	if string(magicHeadBytes) != magicHead {
		return false
	}
	if string(versionBytes) != wareKVVersion {
		return false
	}
	return true
}

func ifCheckSumOK(data []byte) bool {
	n := len(data)

	dataNeedCheck := data[:n-checkSumLen]
	sum := blake2b.Sum512(dataNeedCheck)
	checkSum := data[n-checkSumLen:]

	for i := 0; i < checkSumLen; i++ {
		if sum[i] != checkSum[i] {
			return false
		}
	}
	return true
}

func (c *Camera) GetCreateTime() int64 {
	return c.createTime
}

func (c *Camera) IsActive() bool {
	return c.isOpen
}

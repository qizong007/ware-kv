package storage

type Key struct {
	val string
}

func (k *Key) Hashcode() int {
	h := 0
	if len(k.val) > 0 {
		for _, c := range k.val {
			// (h << 5) - h --> h*31
			// 31是素数，相乘得到的结果比其他方式更容易产生唯一性
			// 也就是说产生 hash 值重复的概率比较小 --> 降低冲突概率
			h = ((h << 5) - h) + int(c)
		}
	}
	return h
}

func (k *Key) GetKey() string {
	return k.val
}

func (k *Key) SetKey(val string) {
	k.val = val
}

func MakeKey(key string) *Key {
	return &Key{
		val: key,
	}
}

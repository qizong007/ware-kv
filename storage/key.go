package storage

type Key struct {
	Val  string
	Hash int
}

func (k Key) Hashcode() int {
	h := k.Hash
	if h == 0 && len(k.Val) > 0 {
		for _, c := range k.Val {
			// (h << 5) - h --> h*31
			// 31是素数，相乘得到的结果比其他方式更容易产生唯一性
			// 也就是说产生 hash 值重复的概率比较小 --> 降低冲突概率
			h = ((h << 5) - h) + int(c)
		}
		k.Hash = h
	}
	return h
}

func MakeKey(key string) *Key {
	return &Key{
		Val:  key,
		Hash: 0,
	}
}

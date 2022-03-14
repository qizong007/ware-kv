package storage

type Key struct {
	val string
}

func (k *Key) Hashcode() int {
	h := 0
	if len(k.val) > 0 {
		for _, c := range k.val {
			// (h << 5) - h --> h*31, the '31' is a prime number
			// The result is more likely to be unique when multiplied than otherwise,
			// and the probability of creating hash value duplicates is low,
			// so it [ Reduces the probability of Conflicts ]
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

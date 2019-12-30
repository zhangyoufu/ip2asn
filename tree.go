package ip2asn

import (
	"errors"
	"math/bits"
)

type node struct {
	Child [2]*node
	AsNum uint32
}

type tree struct {
	Root node
}

func (t *tree) Insert(subnet uint32, mask uint, asNum uint32) error {
	if mask < 1 || mask > 32 {
		return errors.New("invalid mask")
	}
	n := &t.Root
	for k := bits.Reverse32(subnet); mask > 0; mask-- {
		c := n.Child[k&1]
		if c == nil {
			c = new(node)
			n.Child[k&1] = c
		}
		n = c
		k >>= 1
	}
	if n.AsNum != 0 && n.AsNum != asNum {
		return errors.New("conflict AS number")
	}
	n.AsNum = asNum
	return nil
}

func (t *tree) QueryOne(ip uint32) uint32 {
	k := bits.Reverse32(ip)
	n := &t.Root
	for n != nil {
		if n.AsNum != 0 {
			return n.AsNum
		}
		n = n.Child[k&1]
		k >>= 1
	}
	return 0
}

// TODO: QueryAll

package ip2asn

import "fmt"

type DataSet struct {
	tree  tree
	asMap map[uint32]string
}

func newDataSet() *DataSet {
	return &DataSet{
		asMap: map[uint32]string{},
	}
}

func (ds *DataSet) GetAsNum(ip string) (asNum uint32) {
	addr, err := parseIpv4(ip)
	if err != nil {
		return
	}
	asNum = ds.tree.QueryOne(addr)
	return
}

func (ds *DataSet) GetAsName(asNum uint32) string {
	return ds.asMap[asNum]
}

func parseIpv4(s string) (addr uint32, err error) {
	var ip [4]byte
	if _, err = fmt.Sscanf(s, "%d.%d.%d.%d", &ip[0], &ip[1], &ip[2], &ip[3]); err != nil {
		return
	}
	addr = uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	return
}

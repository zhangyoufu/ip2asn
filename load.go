package ip2asn

import "fmt"

type Record struct {
	CIDR   string
	AsNum  uint32
	AsName string
}

// Load records from a given data source. Load returns after data source is
// drained (closed).
func Load(dataSource <-chan Record, errc chan<- error) (ds *DataSet) {
	ds = newDataSet()

	for record := range dataSource {
		subnet, mask := parseCIDR(record.CIDR)
		if mask == 0 {
			reportError(errc, "invalid CIDR %q", record.CIDR)
			continue
		}

		if record.AsNum == 0 {
			reportError(errc, "AS number missing for %#v", record)
			continue
		}

		if record.AsName != "" {
			name, existed := ds.asMap[record.AsNum]
			if existed {
				if name != record.AsName {
					reportError(errc, "AS%d name conflict, %q vs. %q", record.AsNum, name, record.AsName)
				}
			} else {
				ds.asMap[record.AsNum] = record.AsName
			}
		}

		if err := ds.tree.Insert(subnet, mask, record.AsNum); err != nil {
			reportError(errc, "%w for %#v", err, record)
			continue
		}
	}

	return
}

// parseCIDR returns mask == 0 if parsing failed.
func parseCIDR(cidr string) (subnet uint32, mask uint) {
	var ip [4]byte
	if _, err := fmt.Sscanf(cidr, "%d.%d.%d.%d/%d", &ip[0], &ip[1], &ip[2], &ip[3], &mask); err != nil {
		goto fail
	}
	if mask < 1 || mask > 32 {
		goto fail
	}
	subnet = uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	if subnet<<mask != 0 {
		goto fail
	}
	return

fail:
	mask = 0
	return
}

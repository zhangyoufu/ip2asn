# introduction

A golang package to lookup AS number for IPv4 address.
This package does not come with a database, you have to download one yourself.

# supported data source

1. [IP2Location LITE IP-ASN Database](https://lite.ip2location.com/database/ip-asn) (recommended)
2. [MaxMind GeoLite2 ASN CSV Database](https://dev.maxmind.com/geoip/geoip2/geolite2-asn-csv-database/)

# example

```go
package main

import (
    "bufio"
    "fmt"
    "os"

    "github.com/zhangyoufu/ip2asn"
)

func main() {
    fmt.Println("loading...")
    ds, err := ip2asn.IP2LocationDataSource.Load("IP2LOCATION-LITE-ASN.CSV.ZIP", nil)
    if err != nil {
        panic(err)
    }
    fmt.Println("loaded")

    fmt.Println("input any IPv4 address and press Enter to query")
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        ipAddr := s.Text()
        asNum := ds.GetAsNum(ipAddr)
        if asNum == 0 {
            fmt.Println("not found")
        } else {
            fmt.Printf("AS%d %s\n", asNum, ds.GetAsName(asNum))
        }
    }
}
```

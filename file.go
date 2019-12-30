package ip2asn

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
)

type ZippedCsvDataSource struct {
	CsvFilename  string
	SkipHeader   bool
	NumColumn    int
	CidrColumn   int
	AsNumColumn  int
	AsNameColumn int
}

func (config *ZippedCsvDataSource) Load(zipFilename string, errc chan<- error) (ds *DataSet, err error) {
	zipReadCloser, err := zip.OpenReader(zipFilename)
	if err != nil {
		err = fmt.Errorf("cannot open zip archive, %w", err)
		return
	}
	defer zipReadCloser.Close()

	var csvReadCloser io.ReadCloser
	for _, zipFile := range zipReadCloser.File {
		if path.Base(zipFile.Name) != config.CsvFilename {
			continue
		}
		csvReadCloser, err = zipFile.Open()
		if err != nil {
			err = fmt.Errorf("cannot open csv inside zip archive, %w", err)
			return
		}
		defer csvReadCloser.Close()
		break
	}
	if csvReadCloser == nil {
		err = errors.New("cannot find csv inside zip archive")
		return
	}

	csvReader := csv.NewReader(csvReadCloser)
	csvReader.FieldsPerRecord = config.NumColumn
	csvReader.ReuseRecord = true

	if config.SkipHeader {
		if _, err = csvReader.Read(); err != nil {
			err = fmt.Errorf("cannot read from csv, %w", err)
			return
		}
	}

	ch := make(chan Record, 128)
	go func() {
		defer close(ch)

		// read rows
		for {
			record, err := csvReader.Read()
			switch {
			case err == io.EOF:
				return
			case err == nil:
				asNumText := record[config.AsNumColumn]
				asNum, err := strconv.ParseUint(asNumText, 10, 32)
				if err != nil {
					reportError(errc, "cannot parse AS number %q", asNumText)
					continue
				}

				ch <- Record{
					CIDR:   record[config.CidrColumn],
					AsNum:  uint32(asNum),
					AsName: record[config.AsNameColumn],
				}
			default:
				// TODO: this may spam errc when SkipHeader == false and
				//       number of columns mismatch for every record
				reportError(errc, "cannot read from csv, %w", err)
			}
		}
	}()

	ds = Load(ch, errc)
	// when Load returns, ch is guaranteed to be closed (goroutine finished)
	// we're safe to close opened files here

	return
}

var (
	MaxMindDataSource = ZippedCsvDataSource{
		CsvFilename:  "GeoLite2-ASN-Blocks-IPv4.csv",
		SkipHeader:   true,
		NumColumn:    3,
		CidrColumn:   0,
		AsNumColumn:  1,
		AsNameColumn: 2,
	}
	IP2LocationDataSource = ZippedCsvDataSource{
		CsvFilename:  "IP2LOCATION-LITE-ASN.CSV",
		SkipHeader:   false,
		NumColumn:    5,
		CidrColumn:   2,
		AsNumColumn:  3,
		AsNameColumn: 4,
	}
)

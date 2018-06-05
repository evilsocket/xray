package units

import (
	"fmt"
)

type DataType int

const (
	DataTypeDomain DataType = iota
	DataTypeIP
	DataTypeMAC
	DataTypeOpenPort
)

type Data struct {
	Type  DataType
	Data  string
	Extra interface{}
}

func (d Data) String() string {
	switch d.Type {
	case DataTypeDomain:
		return fmt.Sprintf("%s %v", d.Data, d.Extra)
	}
	return fmt.Sprintf("%v", d)
}

func (d Data) explodeDomain() []Data {
	exploded := make([]Data, 0)

	exploded = append(exploded, Data{Type: DataTypeDomain, Data: d.Data})
	if d.Extra != nil {
		if ips, ok := d.Extra.([]string); ok {
			for _, ip := range ips {
				exploded = append(exploded, Data{Type: DataTypeIP, Data: ip})
			}
		}
	}

	return exploded
}

func (d Data) Explode() []Data {
	switch d.Type {
	case DataTypeDomain:
		return d.explodeDomain()
	}
	return []Data{d}
}

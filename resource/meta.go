/*
 * Copyright 2025 Secberus, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/netip"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"

	v1 "github.com/secberus/go-push-api/types/v1"
)

func columnsFor[T any](pk string) []*v1.Column {
	t := reflect.TypeFor[T]()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	n := t.NumField()
	columns := make([]*v1.Column, n)

	for i := 0; i < n; i++ {
		f := t.Field(i)

		columns[i] = columnFor(f)
		if f.Name == pk || columns[i].Name == pk {
			columns[i].PrimaryKey = true
		}
	}

	return columns
}

var (
	_Text        = &v1.DataType_Text{Text: &v1.Text{}}
	_Boolean     = &v1.DataType_Boolean{Boolean: &v1.Boolean{}}
	_Integer     = &v1.DataType_Integer{Integer: &v1.Integer{}}
	_Smallint    = &v1.DataType_Smallint{Smallint: &v1.Smallint{}}
	_Bigint      = &v1.DataType_Bigint{Bigint: &v1.Bigint{}}
	_Real        = &v1.DataType_Real{Real: &v1.Real{}}
	_Double      = &v1.DataType_Double{Double: &v1.Double{}}
	_Bytea       = &v1.DataType_Bytea{Bytea: &v1.Bytea{}}
	_Timestamptz = &v1.DataType_Timestamptz{Timestamptz: &v1.Timestamptz{}}
	_Jsonb       = &v1.DataType_Jsonb{Jsonb: &v1.Jsonb{}}
	_Inet        = &v1.DataType_Inet{Inet: &v1.Inet{}}
	_Cidr        = &v1.DataType_Cidr{Cidr: &v1.Cidr{}}
	_Macaddr     = &v1.DataType_Macaddr{Macaddr: &v1.Macaddr{}}
)

func columnFor(f reflect.StructField) *v1.Column {
	c := v1.Column{
		Name:     snakecase(f.Name),
		DataType: &v1.DataType{},
	}

	t := f.Type
	if t.Kind() == reflect.Pointer {
		c.Nillable = true
		t = t.Elem()
	}

	u := &c.DataType.Union

	switch t.Kind() {
	case reflect.String:
		*u = _Text
	case reflect.Bool:
		*u = _Boolean
	case reflect.Int, reflect.Uint16, reflect.Int32:
		*u = _Integer
	case reflect.Int8, reflect.Uint8, reflect.Int16:
		*u = _Smallint
	case reflect.Int64, reflect.Uint32:
		*u = _Bigint
	case reflect.Float32:
		*u = _Real
	case reflect.Float64:
		*u = _Double
	case reflect.Array, reflect.Slice:
		if t == reflect.TypeFor[net.IP]() {
			*u = _Inet
		} else if t == reflect.TypeFor[net.HardwareAddr]() {
			*u = _Macaddr
		} else {
			switch t.Elem().Kind() {
			case reflect.Int8, reflect.Uint8:
				*u = _Bytea
			default:
				*u = _Jsonb
			}
		}
	case reflect.Map:
		if t.Key().Kind() == reflect.String {
			*u = _Jsonb
		} else {
			// TODO better handle non-string-keyed maps
			log.Printf("unhandled type for struct field %q: map with %s keys\n", f.Name, t.Key())
		}
	case reflect.Struct:
		switch t {
		case reflect.TypeFor[bytes.Buffer]():
			*u = _Bytea
		case reflect.TypeFor[time.Time]():
			*u = _Timestamptz
		case reflect.TypeFor[net.IPAddr](), reflect.TypeFor[netip.Addr]():
			*u = _Inet
		case reflect.TypeFor[net.IPNet](), reflect.TypeFor[netip.Prefix]():
			*u = _Cidr
		default:
			*u = _Jsonb
		}
	default:
		log.Printf("unhandled type for struct field %q: %s\n", f.Name, f.Type)
		*u = _Jsonb
	}
	return &c
}

func RecordFor(t *v1.Table, v any) (*v1.Record, error) {
	cs, err := columnValuesFor(t, v)
	if err != nil {
		return nil, fmt.Errorf("failed to create Record for table %q: %w", t.Name, err)
	}
	return &v1.Record{
		TableName: t.Name,
		Columns:   cs,
	}, nil
}

func copyColumn(c *v1.Column) *v1.Column {
	return &v1.Column{
		Name:     c.Name,
		DataType: proto.Clone(c.DataType).(*v1.DataType),
	}
}

const pgTstzFmt = "2006-01-02 15:04:05.999999999Z07:00"

func columnValuesFor(t *v1.Table, v any) ([]*v1.Column, error) {
	rv := reflect.ValueOf(v)

	row := make([]*v1.Column, len(t.Columns))
	for i, c := range t.Columns {
		if c == nil {
			continue
		}
		row[i] = copyColumn(c)

		rcv := rv.Field(i)
		cv := rcv.Interface()
		dt := row[i].DataType

		var err error
		if tt := dt.GetText(); tt != nil {
			tt.Value, err = box[string](cv)
		} else if bt := dt.GetBoolean(); bt != nil {
			bt.Value, err = box[bool](cv)
		} else if it := dt.GetInteger(); it != nil {
			it.Value, err = tryAll(cv, box[int32], boxWith(ii32), boxWith(u16i32))
		} else if st := dt.GetSmallint(); st != nil {
			st.Value, err = tryAll(cv, boxWith(i8i32), boxWith(u8i32), boxWith(i16i32))
		} else if bgt := dt.GetBigint(); bgt != nil {
			bgt.Value, err = tryAll(cv, box[int64], boxWith(u32i64))
		} else if ft := dt.GetReal(); ft != nil {
			ft.Value, err = box[float32](cv)
		} else if dbt := dt.GetDouble(); dbt != nil {
			dbt.Value, err = box[float64](cv)
		} else if bat := dt.GetBytea(); bat != nil {
			switch u := cv.(type) {
			case []uint8:
				bat.Value = u
			case []int8:
				bat.Value = converts(u, i8u8)
			case bytes.Buffer:
				bat.Value = u.Bytes()
			case *bytes.Buffer:
				bat.Value = u.Bytes()
			default:
				if rcv.Kind() == reflect.Array {
					switch rcv.Elem().Kind() {
					case reflect.Uint8:
						bat.Value = rcv.Bytes()
					case reflect.Int8:
						bat.Value = converts(rcv.Slice(0, rcv.Len()).Interface().([]int8), i8u8)
					}
				}
			}
		} else if tmt := dt.GetTimestamptz(); tmt != nil {
			switch u := cv.(type) {
			case time.Time:
				tmt.Value = ptr(u.Truncate(time.Microsecond).Format(pgTstzFmt))
			case *time.Time:
				tmt.Value = ptr(u.Truncate(time.Microsecond).Format(pgTstzFmt))
			}
		} else if intt := dt.GetInet(); intt != nil {
			switch u := cv.(type) {
			case net.IP:
				intt.Value = ptr(u.String())
			case net.IPAddr:
				intt.Value = ptr(u.IP.String())
			case netip.Addr:
				intt.Value = ptr(u.String())
			}
		} else if mt := dt.GetMacaddr(); mt != nil {
			switch u := cv.(type) {
			case net.HardwareAddr:
				mt.Value = ptr(u.String())
			}
		} else if cdt := dt.GetCidr(); cdt != nil {
			switch u := cv.(type) {
			case net.IPNet:
				cdt.Value = ptr(u.String())
			case netip.Prefix:
				cdt.Value = ptr(u.String())
			}
		} else if jt := dt.GetJsonb(); jt != nil {
			var enc []byte
			if enc, err = json.Marshal(cv); err == nil {
				jt.Value = ptr(string(enc))
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to encode column %q: %w", c.Name, err)
		}
	}

	return row, nil
}

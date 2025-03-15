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
	"errors"
	"fmt"
	"unicode"
)

func snakecase(s string) string {
	var out bytes.Buffer
	caps := true
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && !caps {
				out.WriteByte('_')
			}
			caps = true
			r = unicode.ToLower(r)
		} else {
			caps = false
		}
		out.WriteRune(r)
	}
	return out.String()
}

func ptr[T any](v T) *T {
	return &v
}

func zero[T any]() (t T) {
	return
}

func box[T any](v any) (*T, error) {
	if p, ok := v.(*T); ok {
		return p, nil
	} else if u, ok := v.(T); ok {
		return &u, nil
	}
	return nil, fmt.Errorf("%+v is not (*)%T", v, zero[T]())
}

func boxWith[U, V any](conv func(U) V) func(any) (*V, error) {
	return func(v any) (*V, error) {
		if pu, ok := v.(*U); ok {
			if pu == nil {
				return nil, nil
			}
			return ptr(conv(*pu)), nil
		} else if u, ok := v.(U); ok {
			return ptr(conv(u)), nil
		}
		return nil, fmt.Errorf("%+v is not (*)%T", v, zero[U]())
	}
}

func tryAll[U, V any](u U, fs ...func(U) (V, error)) (V, error) {
	errs := make([]error, len(fs))
	for i, f := range fs {
		if v, err := f(u); err == nil {
			return v, nil
		} else {
			errs[i] = err
		}
	}
	return zero[V](), errors.Join(errs...)
}

func converts[U, V any](us []U, conv func(U) V) []V {
	vs := make([]V, len(us))
	for i, u := range us {
		vs[i] = conv(u)
	}
	return vs
}

func ii32(i int) int32      { return int32(i) }
func i8u8(i int8) uint8     { return uint8(i) }
func i8i32(i int8) int32    { return int32(i) }
func u8i32(i uint8) int32   { return int32(i) }
func i16i32(i int16) int32  { return int32(i) }
func u16i32(i uint16) int32 { return int32(i) }
func u32i64(i uint32) int64 { return int64(i) }

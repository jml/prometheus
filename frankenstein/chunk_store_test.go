// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frankenstein

import (
	"reflect"
	"testing"

	"github.com/prometheus/prometheus/frankenstein/wire"
)

func c(id string) wire.Chunk {
	return wire.Chunk{ID: id}
}

func TestIntersect(t *testing.T) {
	for _, tc := range []struct {
		in   [][]wire.Chunk
		want []wire.Chunk
	}{
		{nil, []wire.Chunk{}},
		{[][]wire.Chunk{{c("a"), c("b"), c("c")}}, []wire.Chunk{c("a"), c("b"), c("c")}},
		{[][]wire.Chunk{{c("a"), c("b"), c("c")}, {c("a"), c("c")}}, []wire.Chunk{c("a"), c("c")}},
		{[][]wire.Chunk{{c("a"), c("b"), c("c")}, {c("a"), c("c")}, {c("b")}}, []wire.Chunk{}},
		{[][]wire.Chunk{{c("a"), c("b"), c("c")}, {c("a"), c("c")}, {c("a")}}, []wire.Chunk{c("a")}},
	} {
		have := nWayIntersect(tc.in)
		if !reflect.DeepEqual(have, tc.want) {
			t.Errorf("%v != %v", have, tc.want)
		}
	}
}

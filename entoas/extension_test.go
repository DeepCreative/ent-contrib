// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entoas

import (
	"os"
	"testing"

	"entgo.io/contrib/entoas/spec"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

func TestExtension(t *testing.T) {
	ex, err := NewExtension(
		DefaultPolicy(PolicyExpose),
		Mutations(func(graph *gen.Graph, spec *spec.Spec) error { return nil }),
		SpecTitle("Spec Title"),
		SpecDescription("Spec Description"),
		SpecVersion("Spec Version"),
		WriteTo(os.Stdout),
	)
	require.NoError(t, err)
	require.Equal(t, ex.config.DefaultPolicy, PolicyExpose)
	require.Len(t, ex.mutations, 4)
	require.Equal(t, ex.out, os.Stdout)
}
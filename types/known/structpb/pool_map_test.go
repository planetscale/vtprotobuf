package structpb

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestStructPooledFieldsRoundTrip(t *testing.T) {
	src, err := structpb.NewStruct(map[string]any{
		"a": 1.0, "b": "s", "c": []any{1.0, "x"}, "d": map[string]any{"k": true},
	})
	require.NoError(t, err)
	data, err := proto.Marshal(src)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		got := StructFromVTPool()
		require.NoError(t, got.UnmarshalVT(data))
		require.True(t, proto.Equal(src, (*structpb.Struct)(got)))
		got.ReturnToVTPool()
	}
}

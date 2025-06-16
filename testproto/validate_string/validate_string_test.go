package validate_string

import (
	"testing"
	utf8 "unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalVTInvalidUTF8(t *testing.T) {
	payload := []byte{0x12, 0x34, 0x56, 0x78, 0x9a}
	require.False(t, utf8.Valid(payload))
	tests := []struct {
		name string
		msg  *ValidateStringTest
	}{
		{
			"string field",
			&ValidateStringTest{
				Sub: &ValidateStringTest_Sub1_{
					&ValidateStringTest_Sub1{
						S: string(payload),
					},
				},
			},
		},
		{
			"repeated string field",
			&ValidateStringTest{
				Sub: &ValidateStringTest_Sub2_{
					&ValidateStringTest_Sub2{
						S: []string{string(payload)},
					},
				},
			},
		},
		{
			"map string field invalid key",
			&ValidateStringTest{
				Sub: &ValidateStringTest_Sub3_{
					&ValidateStringTest_Sub3{
						Foo: map[string]string{
							string(payload): "value",
						},
					},
				},
			},
		},
		{
			"map string field invalid value",
			&ValidateStringTest{
				Sub: &ValidateStringTest_Sub3_{
					&ValidateStringTest_Sub3{
						Foo: map[string]string{
							"key": string(payload),
						},
					},
				},
			},
		},
		{
			"oneof string field",
			&ValidateStringTest{
				Sub: &ValidateStringTest_Sub4_{
					&ValidateStringTest_Sub4{
						Foo: &ValidateStringTest_Sub4_S{
							S: string(payload),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.msg.MarshalVT()
			require.NoError(t, err)
			var msg ValidateStringTest
			err = msg.UnmarshalVT(data)
			require.Equal(t, "proto: invalid UTF-8 string", err.Error())
		})
	}
}

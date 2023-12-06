package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerICS20Denom(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		want        string
		expectedErr bool
	}{
		{
			name:        "HexAddress",
			in:          "0x4639F884305273E856dBa51AF60c10a5b5E0F482",
			want:        "0x4639f884305273e856dba51af60c10a5b5e0f482",
			expectedErr: false,
		},
		{
			name:        "WithPrefix",
			in:          "Port/channel-0/0x4639F884305273E856dBa51AF60c10a5b5E0F482",
			want:        "Port/channel-0/0x4639f884305273e856dba51af60c10a5b5e0f482",
			expectedErr: false,
		},
		{
			name:        "WithMultiplePrefix",
			in:          "Port-0/channel-0/Port-1/channel-1/0x4639F884305273E856dBa51AF60c10a5b5E0F482",
			want:        "Port-0/channel-0/Port-1/channel-1/0x4639f884305273e856dba51af60c10a5b5e0f482",
			expectedErr: false,
		},
		{
			name:        "Invalid Address Format",
			in:          "invalid address format",
			expectedErr: true,
		},
		{
			name:        "Invalid Prefix Format",
			in:          "Portchannel-0/0x4639F884305273E856dBa51AF60c10a5b5E0F482",
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ret, err := ToCanonicalICS20Denom(c.in)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.want, ret)
			}
		})
	}
}

//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package ip

import (
	"fmt"
	"testing"
)

func TestParseNetworkSegment(t *testing.T) {
	type args struct {
		cidr string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "parseNetworkSegmentTest1",
			args: args{
				cidr: "1.0.1.0/24",
			},
			want:  "1.0.1.1",
			want1: "1.0.1.255",
		},
		{
			name: "parseNetworkSegmentTest2",
			args: args{
				cidr: "1.0.2.0/23",
			},
			want:  "",
			want1: "",
		},
		{
			name: "parseNetworkSegmentTest3",
			args: args{
				cidr: "1.0.8.0/21",
			},
			want:  "",
			want1: "",
		},
		{
			name: "parseNetworkSegmentTest4",
			args: args{
				cidr: "43.224.208.0/21",
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCidr(tt.args.cidr)

			fmt.Println(got)
			fmt.Println("主机数量", getCidrHostNum(21))
		})
	}
}

func Test_parseIpSegSlice(t *testing.T) {
	type args struct {
		ipSegSlice []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "解析网段转化为二进制",
			args: args{ipSegSlice: []string{"192", "168", "32", "0"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseIpSegSlice(tt.args.ipSegSlice)
		})
	}
}

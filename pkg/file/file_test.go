//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package file

import "testing"

const testPath = "../../assets/originSource.txt"

func TestOpenFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "openSource_test1", args: args{path: testPath}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OpenFile(tt.args.path)
		})
	}
}

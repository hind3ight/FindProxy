//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package file

import (
	"io/ioutil"
	"os"
)

//
//  OpenFile
//  @Description open file and print
//  @param path string
//  @author hind3ight
//  @createdtime 2022-07-31 16:49:02
//  @updatedtime 2022-07-31 16:49:02
//
func OpenFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(file)

	return content, err
}

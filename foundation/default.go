// +build !asset

package foundation

import "github.com/deadcheat/cashew/setting/loader/file"

func init() {
	loader = newLoader(file.New())
}

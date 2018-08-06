// +build !asset

package foundation

import "github.com/deadcheat/cashew/setting/file"

func init() {
	loader = newLoader(file.New())
}

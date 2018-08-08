// +build !asset

package foundation

import "github.com/deadcheat/cashew/setting/loader/file"

func initializeApp() {
	loader = newLoader(file.New())
}

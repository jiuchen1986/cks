/*
Copyright © 2020 Xin Chen <devops.chen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"fmt"
	"time"
)

const (
	timeLayout string = "Mon Jan 2 15:04:05.0000 MST 2006"
)

// Printf is a wapper that prints msg with a prefix
// which should be used in case log system isn't initiated completely
func Printf(format string, a ...interface{}) (int, error) {
	msg := fmt.Sprintf("%s * * * * * * %s\n", time.Now().Local().Format(timeLayout), format)
	return fmt.Printf(msg, a...)
}

// Println is a wapper that prints msg with a prefix
// which should be used in case log system isn't initiated completely
func Println(a ...interface{}) (int, error) {
	msg := fmt.Sprintf("%s * * * * * *", time.Now().Local().Format(timeLayout))
	as := []interface{}{msg}
	as = append(as, a...)
	return fmt.Println(as...)
}

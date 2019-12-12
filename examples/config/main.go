// Copyright 2019 The magicdb Authors
//
// Licensed under the Apache Licence, Version 2.0(the "License");
// You may not use the file except in compliance with the Licence.
// You may obtain a copy of the Licence at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distrubuted under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	for {
		viper.WatchConfig()
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		appName := viper.Get("appName")
		fmt.Println(appName)

		time.Sleep(2 * time.Second)
	}

}

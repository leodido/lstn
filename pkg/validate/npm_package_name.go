// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2023 The listen.dev team <engineering@garnet.ai>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package validate

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/leodido/go-npmpackagename"
)

func isNpmPackageName(fl validator.FieldLevel) bool {
	field := fl.Field()
	// Do you want strict validation or not?
	strict := false
	if fl.Param() == "strict" {
		strict = true
	}

	if field.Kind() == reflect.String {
		valid, warnings, err := npmpackagename.Validate([]byte(field.String()))
		if strict {
			if err == nil && len(warnings) == 0 {
				return valid
			}
		} else {
			if err == nil {
				return valid
			}
		}

		return false
	}

	panic(fmt.Sprintf("bad field type: %T", field.Interface()))
}

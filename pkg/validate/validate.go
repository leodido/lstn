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
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationErrors = validator.ValidationErrors

// Singleton is the validator singleton instance.
//
// This way it caches the structs info.
var Singleton *validator.Validate

// Translator is the universal translator for validation errors.
var Translator ut.Translator

func init() {
	Singleton = validator.New()

	// Register a function to get the field name from "name" tags.
	Singleton.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	if err := Singleton.RegisterValidation("endpoint", isEndpoint); err != nil {
		panic(err)
	}
	if err := Singleton.RegisterValidation("jq", jqQueryCompiles); err != nil {
		panic(err)
	}
	if err := Singleton.RegisterValidation("npm_package_name", isNpmPackageName); err != nil {
		panic(err)
	}
	if err := Singleton.RegisterValidation("version_constraint", isVersionConstraint); err != nil {
		panic(err)
	}

	Singleton.RegisterAlias("shasum", "len=40")
	Singleton.RegisterAlias("mandatory", "required")

	eng := en.New()
	Translator, _ = (ut.New(eng, eng)).GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(Singleton, Translator); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"mandatory",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("mandatory", "{0} is mandatory", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mandatory", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"excluded_without",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("excluded_without", "cannot use --{0} without specifying --{1}", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			// NOTE > Assuming that the flag is the lowercase of the struct field name we are depending on
			dependingOn := strings.ToLower(fe.Param())
			t, _ := ut.T("excluded_without", fe.Field(), dependingOn)

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"endpoint",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("endpoint", "{0} must be a valid listen.dev endpoint", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("endpoint", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"jq",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("jq", "{0} must be a valid JQ query", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("jq", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"npm_package_name",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("npm_package_name", "{0} is not a valid npm package name", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("npm_package_name", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"version_constraint",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("version_constraint", "{0} is not a valid version constraint", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("version_constraint", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"reporter",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("reporter", "{0} is not a valid reporter", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("reporter", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}
}

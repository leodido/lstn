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
package templates

import (
	"io"
	"text/template"

	"github.com/listendev/lstn/pkg/listen"
)

const singleProblemsTpl = `
{{ if gt (len .Problems) 0 }}
### <b><a href="https://verdicts.listen.dev/npm/{{ .Name }}/{{ .Version }}">{{ .Name }}@{{ .Version }}</a></b><br>

{{ range .Problems }}
{{ $title := .Title}}
{{ $url := .Type }}
- {{ $title }} (<a href="{{ $url }}">learn more :link:</a>)
{{ end }}
{{ end }}
`

func RenderSingleProblemsPackage(w io.Writer, p listen.Package) error {
	ct := template.Must(template.New("single_problem").Parse(singleProblemsTpl))
	err := ct.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}

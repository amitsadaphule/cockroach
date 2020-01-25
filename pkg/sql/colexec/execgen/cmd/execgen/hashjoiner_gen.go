// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import (
	"io"
	"io/ioutil"
	"text/template"
)

func genHashJoiner(wr io.Writer) error {
	t, err := ioutil.ReadFile("pkg/sql/colexec/hashjoiner_tmpl.go")
	if err != nil {
		return err
	}

	s := string(t)

	distinctCollectRightOuter := makeFunctionRegex("_DISTINCT_COLLECT_RIGHT_OUTER", 3)
	s = distinctCollectRightOuter.ReplaceAllString(s, `{{template "distinctCollectRightOuter" buildDict "Global" . "UseSel" $3}}`)

	distinctCollectNoOuter := makeFunctionRegex("_DISTINCT_COLLECT_NO_OUTER", 4)
	s = distinctCollectNoOuter.ReplaceAllString(s, `{{template "distinctCollectNoOuter" buildDict "Global" . "UseSel" $4}}`)

	collectRightOuter := makeFunctionRegex("_COLLECT_RIGHT_OUTER", 5)
	s = collectRightOuter.ReplaceAllString(s, `{{template "collectRightOuter" buildDict "Global" . "UseSel" $5}}`)

	collectNoOuter := makeFunctionRegex("_COLLECT_NO_OUTER", 5)
	s = collectNoOuter.ReplaceAllString(s, `{{template "collectNoOuter" buildDict "Global" . "UseSel" $5}}`)

	tmpl, err := template.New("hashjoiner_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, struct{}{})
}

func init() {
	registerGenerator(genHashJoiner, "hashjoiner.eg.go")
}

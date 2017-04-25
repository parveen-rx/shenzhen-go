// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package view

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/google/shenzhen-go/model"
)

const graphEditorTemplateSrc = `<html>
<head>
	<title>{{$.Graph.Name}}</title>
	<link type="text/css" rel="stylesheet" href="/.static/fonts.css">
	<link type="text/css" rel="stylesheet" href="/.static/main.css">
	<script>
		var graphPath = '{{$.Graph.URLPath}}';
		var apiURL = '/.api';
		var GraphJSON = "{{$.GraphJSON}}";
	</script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/ace.js" type="text/javascript" charset="utf-8"></script>
</head>
<body>
	<div class="head">
		<a href="?up" title="Go up to the files in the current directory">Up</a> |
		<a id="graph-save" href="javascript:void(0)" title="Save current changes to disk">Save</a> | 
		<a href="?reload" class="destructive" title="Revert to last saved file">Revert</a> |
		{{if $.Graph.IsCommand -}}
		<a href="?install" target="_blank" title="Export the graph to a Go package and 'go install' it">Install</a> | 
		{{else -}}
		<a href="?build" target="_blank" title="Export the graph to a Go package and 'go build' it">Build</a> | 
		{{end -}}
		<a href="?run" target="_blank" title="Export the graph to a Go package and execute it">Run</a> | 
		<span class="dropdown">
			<a href="javascript:void(0)">New goroutine</a> 
			<div class="dropdown-content">
				<ul>
				{{range $t, $null := $.PartTypes -}}
					<li><a href="?node=new&type={{$t}}">{{$t}}</a></li>
				{{- end}}
				</ul>
			</div>
		</span> |
		View as: <a href="?go" target="_blank">Go</a> <a href="?json" target="_blank">JSON</a>
	</div>
	<div class="box">
		<div class="container" id="diagram-container">
			<svg id="diagram" width="1600" height="1600" viewBox="0 0 1600 1600" />
		</div>
		<div class="container" id="panels-container">
			<div id="graph-properties">
				<h3>Graph Properties</h3>
				<a id="graph-properties-save" href="javascript:void(0)">Save</a>
				<div class="formfield">
				    <label for="Name">Name</label>
					<input id="graph-prop-name" name="Name" type="text" required value="{{$.Graph.Name}}">
				</div>
				<div class="formfield">
				    <label for="PackagePath">Package path</label>
					<input id="graph-prop-package-path" name="PackagePath" type="text" required value="{{$.Graph.PackagePath}}">
				</div>
				<div class="formfield">
				    <label for="IsCommand">Is a command?</label>
					<input id="graph-prop-is-command" name="IsCommand" type="checkbox" {{if $.Graph.IsCommand}}checked{{end}} title="Selecting this means the generated package line will be 'package main' instead of 'package [packagename]', which allows your package to run as a standalone command and be installed with 'go install'. De-selecting this causes the package to be usable as a library.">
				</div>
			</div>
			<div id="node-properties" style="display:none">
				<div id="node-actions" class="head">
					<a id="node-save-link" href="javascript:void(0);" title="Save changes to this goroutine.">Save</a> |
					<a id="node-clone-link" href="javascript:void(0);" title="Make a copy of this goroutine.">Clone</a> | 
					<a id="node-convert-link" href="javascript:void(0);" class="destructive" title="Change this goroutine into a Code goroutine; it cannot be converted back.">Convert to Code</a> |
					<a id="node-delete-link" href="javascript:void(0);" class="destructive" title="Delete this goroutine">Delete</a> | 
				</div>
				<div id="node-panels" class="head">
					<a id="node-metadata-link" href="javascript:void(0);">Properties</a> 
					{{range $tk, $type := $.PartTypes}}
					<span id="node-{{$tk}}-links" style="display:none">
					{{range $type.Panels }}
					| <a id="node-{{$tk}}-{{.Name}}-link" href="javascript:void(0);">{{.Name}}</a>
					{{end}}
					</span>
					{{end}}
				</div>
				<div id="node-metadata-panel">
					<div class="formfield">
						<label for="Name">Name</label>
						<input id="node-name" name="Name" type="text" required value="{.Name}">
					</div>
					<div class="formfield">
						<label for="Enabled">Enabled</label>
						<input id="node-enabled" name="Enabled" type="checkbox" checked>
					</div>
					<div class="formfield">
						<label for="Multiplicity">Multiplicity</label>
						<input id="node-multiplicity" name="Multiplicity" type="number" required pattern="^[1-9][0-9]*$" title="Must be a whole number, at least 1." value="1">
					</div>
					<div class="formfield">
						<label for="Wait">Wait for this to finish</label>
						<input id="node-wait" name="Wait" type="checkbox" checked>
					</div>
				</div>
				{{range $tk, $type := $.PartTypes}}
				{{range $type.Panels}}
				<div class="node-panel" id="node-{{$tk}}-{{.Name}}-panel" style="display:none">
					{{.Editor}}
				</div>
				{{end}}
				{{end}}
			</div>
		</div>
	</div>
	<script src="/.static/svg.js"></script>
</body>
</html>`

var graphEditorTemplate = template.Must(template.New("graphEditor").Parse(graphEditorTemplateSrc))

type editorInput struct {
	Graph     *model.Graph
	GraphJSON string
	PartTypes map[string]*model.PartType
}

// Graph displays a graph.
func Graph(w http.ResponseWriter, g *model.Graph) {
	gj, err := json.Marshal(g)
	if err != nil {
		log.Printf("Could not execute graph editor template: %v", err)
		http.Error(w, "Could not execute graph editor template", http.StatusInternalServerError)
		return
	}

	d := &editorInput{
		Graph:     g,
		GraphJSON: string(gj),
		PartTypes: model.PartTypes,
	}
	if err := graphEditorTemplate.Execute(w, d); err != nil {
		log.Printf("Could not execute graph editor template: %v", err)
		http.Error(w, "Could not execute graph editor template", http.StatusInternalServerError)
	}
}

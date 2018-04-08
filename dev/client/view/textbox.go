// Copyright 2017 Google Inc.
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

import "github.com/google/shenzhen-go/dev/dom"

// TextBox is an SVG group containing a filled rectangle and text.
type TextBox struct {
	// Children Rectangle and Text, and Text has child TextNode.
	Group
	Rectangle dom.Element
	Text      dom.Element
	TextNode  dom.Element

	MinWidth    float64
	Margin      float64
	TextOffsetY float64
}

// MakeElements creates the DOM elements, organises them,
// and sets default attributes. Note the default is to create hidden.
// The return is the main group.
func (b *TextBox) MakeElements(doc dom.Document) *TextBox {
	b.Group = NewGroup(doc)
	b.Group.Hide()
	b.Rectangle = doc.MakeSVGElement("rect")
	b.Text = doc.MakeSVGElement("text")
	b.TextNode = doc.MakeTextNode("")
	b.Group.
		AddChildren(b.Rectangle, b.Text)
	b.Text.
		SetAttribute("text-anchor", "middle").
		SetAttribute("unselectable", "on").
		AddChildren(b.TextNode)
	return b
}

// SetHeight sets the textbox height.
func (b *TextBox) SetHeight(height float64) *TextBox {
	b.Rectangle.SetAttribute("height", height)
	b.Text.SetAttribute("y", height/2+b.TextOffsetY)
	return b
}

// SetRectangleStyle sets the style of the rectangle.
func (b *TextBox) SetRectangleStyle(style string) *TextBox {
	b.Rectangle.SetAttribute("style", style)
	return b
}

// SetText sets te text in the textbox.
func (b *TextBox) SetText(text string) *TextBox {
	b.TextNode.Set("nodeValue", text)
	return b
}

// SetTextStyle sets the style attribute of the text.
func (b *TextBox) SetTextStyle(style string) *TextBox {
	b.Text.SetAttribute("style", style)
	return b
}

// SetWidth sets the width of the textbox, unless the width is less than the MinWidth,
// in which case MinWidth is used instead.
func (b *TextBox) SetWidth(w float64) *TextBox {
	if w < b.MinWidth {
		w = b.MinWidth
	}
	b.Rectangle.SetAttribute("width", w)
	b.Text.SetAttribute("x", w/2)
	return b
}

// Width returns the current width.
func (b *TextBox) Width() float64 {
	return b.Rectangle.GetAttribute("width").Float()
}

// RecomputeWidth resizes the textbox to fit all text (plus a margin).
func (b *TextBox) RecomputeWidth() *TextBox {
	return b.SetWidth(b.Text.Call("getComputedTextLength").Float() + 2*b.Margin)
}

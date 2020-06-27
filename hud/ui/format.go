package ui

import (
	"regexp"
	"strings"
)

// NewFormatPane creates container with components positioned by formatted text.
//
// The idea is that you will define layout as the text with multiple lines.
// On the places where you want to have components you will use some placeholders.
// Then for each placeholder you define how to create component on it's position.
//
// Created container will always contain at least one Text component.
// It will also contain additional components depending on given components builders.
//
// See ComponentBuilder for more info.
func NewFormatPane(text string, builders []*ComponentBuilder) *BaseContainer {
	p := NewBaseContainer()
	components := []Component{}
	lines := strings.Split(text, "\n")

	// process all component builders
	for _, cb := range builders {
		ci := 0
		pattern := regexp.MustCompile(cb.Pattern)

		// check all lines for pattern matches
		for line, text := range lines {
			locs := pattern.FindAllStringIndex(text, -1)
			matches := pattern.FindAllString(text, -1)

			// create component for each match and remove matched text from the source
			for mi, match := range matches {
				// remove match is done to support components with smaller size than placeholder size
				lines[line] = lines[line][:locs[mi][0]] + strings.Repeat(" ", len(matches[mi])) + lines[line][locs[mi][1]:]

				// new component with location by match position
				component := cb.Build(ci, match)
				component.Position().X = locs[mi][0]
				component.Position().Y = line
				ci++

				if !cb.Skip {
					components = append(components, component)
				}
			}
		}
	}

	// adding text before other components to be in background
	// using lines with replaced matches
	p.Add(NewTextFromLines(lines))
	p.Add(components...)

	return p
}

// ComponentBuilder defines how to create one or more components for each match of given Pattern
type ComponentBuilder struct {
	// Pattern is regex pattern used to locate placeholder
	Pattern string
	// Build is called for each match of the Pattern.
	// It can be called multiple times - once per each match.
	// It will get index of the match (starting from zero) and the match itself as parameters.
	// It should return some component.
	// Position of returned component will be changed to the position of corresponding match.
	// Each matched text will be removed from source text.
	Build func(i int, s string) Component
	// Skip if true will cause components created by this builder not be added to the container.
	// But underlying matched text will be still removed from source text.
	Skip bool
}

/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"fmt"
	"regexp"
	"strings"
)

// getReferenceString returns a reference to the source.
func (s *SourceRecord) GetReferenceString() string {

	var refs []string
	if s.Author != "" {
		refs = append(refs, s.Author)
	}
	if s.Title != "" {
		refs = append(refs, fmt.Sprintf("\"%s\"", s.Title))
	}

	var pubFacts, pubName, pubDate string
	pubParts := strings.Split(s.Publication, ";")
	r, _ := regexp.Compile("^ *([^:]+): (.+)$")
	for _, p := range pubParts {
		match := r.FindStringSubmatch(p)
		if len(match) > 0 {
			switch label := match[1]; label {
			case "Location":
				pubFacts = match[2]
			case "Name":
				pubName = match[2]
			case "Date":
				pubDate = match[2]
			}
		}
	}
	if pubName != "" {
		if pubFacts != "" {
			pubFacts += ", "
		}
		pubFacts += pubName
	}
	if pubDate != "" {
		if pubFacts != "" {
			pubFacts += ", "
		}
		pubFacts += pubDate
	}
	if pubFacts != "" {
		refs = append(refs, "("+pubFacts+")")
	}

	return strings.Join(refs, ", ")
}

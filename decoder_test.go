/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	data []byte
)

var g *Gedcom

func init() {
	var err error
	data, err = ioutil.ReadFile("testdata/allged.ged")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {

	d := NewDecoder(bytes.NewReader(data))
	d.SetUnrecTagFunc(func(l int, t, v, x string) {
		fmt.Printf("Unrecognized: %d %s %s", l, t, v)
		if x != "" {
			fmt.Printf(" (%s)", x)
		}
		fmt.Println("")
	})

	var err error
	g, err = d.Decode()
	if err != nil {
		log.Fatal("Result of decoding gedcom gave error, expected no error")
	}

	retCode := m.Run()
	os.Exit(retCode)
}

func TestStructuresAreInitialized(t *testing.T) {

	if g == nil {
		t.Fatalf("Result of decoding gedcom was nil, expected valid object")
	}

	if g.Header == nil {
		t.Fatalf("Header record was nil, expected valid record")
	}

	if g.Individual == nil {
		t.Fatalf("Individual list was nil, expected valid slice")
	}

	if g.Family == nil {
		t.Fatalf("Family list was nil, expected valid slice")
	}

	if g.Media == nil {
		t.Fatalf("Media list was nil, expected valid slice")
	}

	if g.Repository == nil {
		t.Fatalf("Repository list was nil, expected valid slice")
	}

	if g.Source == nil {
		t.Fatalf("Source list was nil, expected valid slice")
	}

	if g.Submitter == nil {
		t.Fatalf("Submitter list was nil, expected valid slice")
	}

}

func TestHeader(t *testing.T) {

	h := g.Header

	testCases := []struct {
		tested   string
		expected string
		actual   string
	}{
		{"Gedcom File Name", "ALLGED.GED", h.File},
		{"Copyright", "(C) 1997-2000 by H. Eichmann. You can use and distribute this file freely as long as you do not charge for it", h.Copyright},
		{"Language", "language", h.Language},
		{"Destination", "Destination of transmission", h.Destination},
		{"Encoding", "ASCII Version number of ASCII (whatever it means)", fmt.Sprintf("%s %s", h.Encoding.Name, h.Encoding.Version)},
		{"Timestamp", "1 JAN 1998 13:57:24.80", fmt.Sprintf("%s %s", h.Timestamp.Date, h.Timestamp.Time)},
		{"Source", "APPROVED_SOURCE_NAME", h.Source.Source},
		{"Source version", "Version number of source-program", h.Source.Version},
		{"Source name", "Name of source-program", h.Source.Name},
		{"Source corp name", "Corporation name", h.Source.Corporation.Name},
		{"Source Address Line 0", "Corporation address line 1\nCorporation address line 2\nCorporation address line 3\nCorporation address line 4", h.Source.Corporation.Address.Full},
		{"Source Address Line 1", "Corporation address line 1", h.Source.Corporation.Address.Line1},
		{"Source Address Line 2", "Corporation address line 2", h.Source.Corporation.Address.Line2},
		{"Source Address City", "Corporation address city", h.Source.Corporation.Address.City},
		{"Source Address State", "Corporation address state", h.Source.Corporation.Address.State},
		{"Source Address Postal Code", "Corporation address ZIP code", h.Source.Corporation.Address.PostalCode},
		{"Source Address Country", "Corporation address country", h.Source.Corporation.Address.Country},
		{"Source Phone 3", "Corporation phone number 3 (last one!)", h.Source.Corporation.Phone[2]},
	}

	for _, tc := range testCases {
		if tc.expected != tc.actual {
			t.Errorf("%s was [%s], expected [%s]", tc.tested, tc.actual, tc.expected)
		}
	}
}

func TestIndividual(t *testing.T) {

	name1 := &NameRecord{
		Name: "given name /surname/jr.",
		Citation: []*CitationRecord{
			&CitationRecord{
				Source: &SourceRecord{
					Xref:   "SOURCE1",
					Title:  "Title of source\nTitle continued here. The word TEST should not be broken!",
					Author: "Author of source\nAuthor continued here. The word TEST should not be broken!",
				},

				Page: "42",
				Quay: "0",
				Data: DataRecord{
					Date: "BEF 1 JAN 1900",
					Text: []string{
						"a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: []*NoteRecord{
					&NoteRecord{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: []*NoteRecord{
			&NoteRecord{
				Note: "Personal Name note\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	event1 := &EventRecord{
		Tag:  "BIRT",
		Date: "31 DEC 1997",
		Place: PlaceRecord{
			Name: "The place",
		},
		Note: []*NoteRecord{
			&NoteRecord{
				Note: "BIRTH event note (the event of entering into life)\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	att1 := &EventRecord{
		Tag:   "CAST",
		Value: "Cast name",
		Date:  "31 DEC 1997",
		Place: PlaceRecord{
			Name: "The place",
		},
		Note: []*NoteRecord{
			&NoteRecord{
				Note: "CASTE event note (the name of an individual's rank or status in society, based   on racial or religious differences, or differences in wealth, inherited   rank, profession, occupation, etc)\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	i1 := g.Individual[0]

	iTestCases := []struct {
		format   string
		expected int
		actual   int
	}{
		{"Individual list length was [%d]", 8, len(g.Individual)},
		{"Individual 0 had [%d] names", 2, len(i1.Name)},
		{"Individual 0 had [%d] events", 24, len(i1.Event)},
		{"Individual 0 had [%d] attributes", 14, len(i1.Attribute)},
		{"Individual 0 had [%d] parent families", 2, len(i1.Parents)},
	}

	sTestCases := []struct {
		tested   string
		expected string
		actual   string
	}{
		{"Individual 0 xref", "PERSON1", i1.Xref},
		{"Individual 0 sex", "M", i1.Sex},
		{"Individual 0 Name", name1.Name, i1.Name[0].Name},
		{"Individual 0 Name citation source", name1.Citation[0].Source.Xref, i1.Name[0].Citation[0].Source.Xref},
		{"Individual 0 Name citation Title", name1.Citation[0].Source.Title, i1.Name[0].Citation[0].Source.Title},
		{"Individual 0 Name citation Author", name1.Citation[0].Source.Author, i1.Name[0].Citation[0].Source.Author},
		{"Individual 0 Note 0", name1.Note[0].Note, i1.Name[0].Note[0].Note},
		{"Individual 0 Event 0 Tag", event1.Tag, i1.Event[0].Tag},
		{"Individual 0 Event 0 Date", event1.Date, i1.Event[0].Date},
		{"Individual 0 Event 0 Place Name", event1.Place.Name, i1.Event[0].Place.Name},
		{"Individual 0 Event 0 Note", event1.Note[0].Note, i1.Event[0].Note[0].Note},
		{"Individual 0 Attribute 0 Tag", att1.Tag, i1.Attribute[0].Tag},
		{"Individual 0 Attribute 0 Value", att1.Value, i1.Attribute[0].Value},
		{"Individual 0 Attribute 0 Date", att1.Date, i1.Attribute[0].Date},
		{"Individual 0 Attribute 0 Place Name", att1.Place.Name, i1.Attribute[0].Place.Name},
		{"Individual 0 Attribute 0 Note", att1.Note[0].Note, i1.Attribute[0].Note[0].Note},
	}

	for _, tc := range iTestCases {
		if tc.expected != tc.actual {
			t.Fatalf(tc.format+", expected [%d]", tc.actual, tc.expected)
		}
	}

	for _, tc := range sTestCases {
		if tc.expected != tc.actual {
			t.Errorf("%s was [%s], expected [%s]", tc.tested, tc.actual, tc.expected)
		}
	}

}

func TestSubmitter(t *testing.T) {

	iTestCases := []struct {
		format   string
		expected int
		actual   int
	}{
		{"Submitter list length was [%d]", 1, len(g.Submitter)},
	}

	for _, tc := range iTestCases {
		if tc.expected != tc.actual {
			t.Fatalf(tc.format+", expected [%d]", tc.actual, tc.expected)
		}
	}
}

func TestFamily(t *testing.T) {

	iTestCases := []struct {
		format   string
		expected int
		actual   int
	}{
		{"Family list length was [%d]", 4, len(g.Family)},
	}

	for _, tc := range iTestCases {
		if tc.expected != tc.actual {
			t.Fatalf(tc.format+", expected [%d]", tc.actual, tc.expected)
		}
	}
}

func TestSource(t *testing.T) {

	s := g.Source[0]

	iTestCases := []struct {
		format   string
		expected int
		actual   int
	}{
		{"Source list length was [%d]", 1, len(g.Source)},
		{"Source file list length was [%d]", 2, len(s.File)},
	}

	sTestCases := []struct {
		tested   string
		expected string
		actual   string
	}{
		{"Second source file name", "file2", s.File[1]},
		{"Source title", "Title of source\nTitle continued here. The word TEST should not be broken!", s.Title},
		{"Source submitter name", "A submitter", s.Submitter[0]},
		{"Source repository", "A repository", s.Repository[0]},
		{"Source periodical name", "A periodical name", s.Periodical},
		{"Source volume", "1", s.Volume},
		{"Source page", "3", s.Page[0]},
		{"Source film reference", "at 11", s.Film[0]},
	}

	for _, tc := range iTestCases {
		if tc.expected != tc.actual {
			t.Fatalf(tc.format+", expected [%d]", tc.actual, tc.expected)
		}
	}

	for _, tc := range sTestCases {
		if tc.expected != tc.actual {
			t.Errorf("%s was [%s], expected [%s]", tc.tested, tc.actual, tc.expected)
		}
	}

}

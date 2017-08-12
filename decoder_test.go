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

func TestIndividual(t *testing.T) {

	if len(g.Individual) != 8 {
		t.Fatalf("Individual list length was %d, expected 8", len(g.Individual))
	}

	i1 := g.Individual[0]

	if i1.Xref != "PERSON1" {
		t.Errorf(`Individual 0 xref was "%s", expected @PERSON1@`, i1.Xref)
	}

	if i1.Sex != "M" {
		t.Errorf(`Individual 0 sex "%s" names, expected "M"`, i1.Sex)
	}

	if len(i1.Name) != 2 {
		t.Fatalf(`Individual 0 had %d names, expected 2`, len(i1.Name))
	}

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

	if i1.Name[0].Name != name1.Name {
		t.Errorf(`Individual 0 Name is "%s" names, expected "%s"`, i1.Name[0].Name, name1.Name)
	}
	if i1.Name[0].Citation[0].Source.Xref != name1.Citation[0].Source.Xref {
		t.Errorf(`Individual 0 Name citation source is "%s", expected "%s"`, i1.Name[0].Citation[0].Source.Xref, name1.Citation[0].Source.Xref)
	}
	if i1.Name[0].Citation[0].Source.Title != name1.Citation[0].Source.Title {
		t.Errorf(`Individual 0 Name citation Title is "%s", expected "%s"`, i1.Name[0].Citation[0].Source.Title, name1.Citation[0].Source.Title)
	}
	if i1.Name[0].Citation[0].Source.Author != name1.Citation[0].Source.Author {
		t.Errorf(`Individual 0 Name citation Author is "%s", expected "%s"`, i1.Name[0].Citation[0].Source.Author, name1.Citation[0].Source.Author)
	}
	if i1.Name[0].Note[0].Note != name1.Note[0].Note {
		t.Errorf(`Individual 0 Note 0 is "%s", expected "%s"`, i1.Name[0].Note[0].Note, name1.Note[0].Note)
	}

	if len(i1.Event) != 24 {
		t.Fatalf(`Individual 0 had %d events, expected 24`, len(i1.Event))
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

	if i1.Event[0].Tag != event1.Tag {
		t.Errorf(`Individual 0 Event 0 Tag is "%s", expected "%s"`, i1.Event[0].Tag, event1.Tag)
	}
	if i1.Event[0].Date != event1.Date {
		t.Errorf(`Individual 0 Event 0 Date is "%s", expected "%s"`, i1.Event[0].Date, event1.Date)
	}
	if i1.Event[0].Place.Name != event1.Place.Name {
		t.Errorf(`Individual 0 Event 0 Place Name is "%s", expected "%s"`, i1.Event[0].Place.Name, event1.Place.Name)
	}
	if i1.Event[0].Note[0].Note != event1.Note[0].Note {
		t.Errorf(`Individual 0 Event 0 Note is "%s", expected "%s"`, i1.Event[0].Note[0].Note, event1.Place.Note[0].Note)
	}

	if len(i1.Attribute) != 14 {
		t.Fatalf(`Individual 0 had %d attributes, expected 15`, len(i1.Attribute))
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

	if i1.Attribute[0].Tag != att1.Tag {
		t.Errorf(`Individual 0 Attribute 0 Tag is "%s" names, expected "%s"`, i1.Attribute[0].Tag, att1.Tag)
	}
	if i1.Attribute[0].Value != att1.Value {
		t.Errorf(`Individual 0 Attribute 0 Value is "%s" names, expected "%s"`, i1.Attribute[0].Value, att1.Value)
	}
	if i1.Attribute[0].Date != att1.Date {
		t.Errorf(`Individual 0 Attribute 0 Date is "%s", expected "%s"`, i1.Attribute[0].Date, att1.Date)
	}
	if i1.Attribute[0].Place.Name != att1.Place.Name {
		t.Errorf(`Individual 0 Attribute 0 Place Name is "%s", expected "%s"`, i1.Attribute[0].Place.Name, att1.Place.Name)
	}
	if i1.Attribute[0].Note[0].Note != att1.Note[0].Note {
		t.Errorf(`Individual 0 Attribute 0 Note is "%s", expected "%s"`, i1.Attribute[0].Note[0].Note, att1.Place.Note[0].Note)
	}

	if len(i1.Parents) != 2 {
		t.Fatalf(`Individual 0 had %d parent families, expected 2`, len(i1.Parents))
	}

}

func TestSubmitter(t *testing.T) {

	if len(g.Submitter) != 1 {
		t.Fatalf("Submitter list length was %d, expected 1", len(g.Submitter))
	}

}

func TestFamily(t *testing.T) {

	if len(g.Family) != 4 {
		t.Fatalf("Family list length was %d, expected 4", len(g.Family))
	}

}

func TestSource(t *testing.T) {

	if len(g.Source) != 1 {
		t.Fatalf("Source list length was %d, expected 1", len(g.Source))
	}

	s := g.Source[0]
	if len(s.File) != 2 {
		t.Fatalf("Source file list length was %d, expected 2", len(s.File))
	}
	if s.File[1] != "file2" {
		t.Fatalf("Second source file name was [%s], expected [2]", s.File[1])
	}
	if s.Title != "Title of source\nTitle continued here. The word TEST should not be broken!" {
		t.Fatalf("Source title was [%s], expected [Title of source\nTitle continued here. The word TEST should not be broken!]", s.Title)
	}
	if s.Submitter[0] != "A submitter" {
		t.Fatalf("Source submitter name was [%s], expected [A submitter]", s.Submitter[0])
	}
	if s.Repository[0] != "A repository" {
		t.Fatalf("Source repository name was [%s], expected [A repository]", s.Repository[0])
	}
	if s.Periodical != "A periodical name" {
		t.Fatalf("Source periodical name was [%s], expected [A periodical name]", s.Periodical)
	}
	if s.Volume != "1" {
		t.Fatalf("Source volume was [%s], expected [1]", s.Volume)
	}
	if s.Page[0] != "3" {
		t.Fatalf("Source page was [%s], expected [3]", s.Page[0])
	}
	if s.Film[0] != "at 11" {
		t.Fatalf("Source film reference was [%s], expected [at 11]", s.Film[0])
	}
}

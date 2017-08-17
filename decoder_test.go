/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

import (
	"fmt"
	"io/ioutil"
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

	stringTestCases{
		{"Gedcom File Name", "ALLGED.GED", h.File},
		{"Copyright", "(C) 1997-2000 by H. Eichmann. You can use and distribute this file freely as long as you do not charge for it", h.Copyright},
		{"Language", "language", h.Language},
		{"Destination", "Destination of transmission", h.Destination},
		{"Encoding", "ASCII Version number of ASCII (whatever it means)", fmt.Sprintf("%s %s", h.Encoding.Name, h.Encoding.Version)},
		{"Timestamp", "1 JAN 1998 13:57:24.80", fmt.Sprintf("%s %s", h.Timestamp.Date, h.Timestamp.Time)},
		{"Source", "APPROVED_SOURCE_NAME", h.Source.Source},
		{"Source version", "Version number of source-program", h.Source.Version},
		{"Source name", "Name of source-program", h.Source.Name},
		{"Source form", "LINEAGE-LINKED", h.Source.Form},
		{"Source corp name", "Corporation name", h.Source.Corporation.Name},
		{"Source Address Line 0", "Corporation address line 1\nCorporation address line 2\nCorporation address line 3\nCorporation address line 4", h.Source.Corporation.Address.Full},
		{"Source Address Line 1", "Corporation address line 1", h.Source.Corporation.Address.Line1},
		{"Source Address Line 2", "Corporation address line 2", h.Source.Corporation.Address.Line2},
		{"Source Address City", "Corporation address city", h.Source.Corporation.Address.City},
		{"Source Address State", "Corporation address state", h.Source.Corporation.Address.State},
		{"Source Address Postal Code", "Corporation address ZIP code", h.Source.Corporation.Address.PostalCode},
		{"Source Address Country", "Corporation address country", h.Source.Corporation.Address.Country},
		{"Source Phone 3", "Corporation phone number 3 (last one!)", h.Source.Corporation.Phone[2]},
		{"Source Data Name", "Name of source data", h.Source.Data.Name},
		{"Source Data Date", "1 JAN 1998", h.Source.Data.Date},
		{"Source Data Copyright", "Copyright of source data", h.Source.Data.Copyright},
		{"Submitter name", "/Submitter-Name/", h.Submitter.Name},
		{"Submission family file", "NameOfFamilyFile", h.Submission.FamilyFile},
		{"Gedcom version", "5.5", h.Info.Version},
		{"Gedcom form", "LINEAGE-LINKED", h.Info.Form},
		{"Note", "A general note", h.Note.Note[:14]},
	}.run(t)
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

				Page:    "42",
				Quality: "0",
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

	intTestCases{
		{"Individual list length was [%d]", 8, len(g.Individual)},
		{"Individual 0 had [%d] names", 2, len(i1.Name)},
		{"Individual 0 had [%d] events", 24, len(i1.Event)},
		{"Individual 0 had [%d] attributes", 14, len(i1.Attribute)},
		{"Individual 0 had [%d] parent families", 2, len(i1.Parents)},
	}.run(t)

	stringTestCases{
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
		{"Individual 0 birth father name ", "/Father/", i1.Event[0].Parents[0].Family.Husband.Name[0].Name},
		{"Individual 0 adopted by", "BOTH", i1.Event[6].Parents[0].AdoptedBy},
		{"Individual 0 citation page", "42", i1.Citation[0].Page},
		{"Individual 0 citation source author", "Author of source", i1.Citation[0].Source.Author[:16]},
		{"Individual 0 object title", "A gif picture", i1.Object[0].Title},
		{"Individual 0 Note 0", "A note about the individual\nNote continued here. The word TEST should not be broken!", i1.Note[0].Note},
		{"Individual 0 change date", "1 APR 1998", i1.Changed.Stamp.Date},
		{"Individual 0 change time", "12:34:56.789", i1.Changed.Stamp.Time},
	}.run(t)
}

func TestSubmitter(t *testing.T) {

	r := g.Submitter

	intTestCases{
		{"Submitter list length was [%d]", 1, len(r)},
	}.run(t)

	stringTestCases{
		{"Submitter xref", "SUBMITTER", r[0].Xref},
		{"Submitter name", "/Submitter-Name/", r[0].Name},
		{"Submitter address country", "Submitter address country", r[0].Address.Country},
		{"Submitter phone nuumber 2", "Submitter phone number 2", r[0].Phone[1]},
		{"Submitter language", "English", r[0].Language},
		{"Submitter change date", "19 JUN 2000", r[0].Changed.Stamp.Date},
		{"Submitter change time", "12:34:56.789", r[0].Changed.Stamp.Time},
		{"Submitter note", "A note\nNote continued here. The word TEST should not be broken!", r[0].Changed.Note[0].Note},
	}.run(t)
}

func TestSubmission(t *testing.T) {

	r := g.Submission

	stringTestCases{
		{"Submission xref", "SUBMISSION", r.Xref},
		{"Submission family file", "NameOfFamilyFile", r.FamilyFile},
		{"Submission temple code", "Abreviated temple code", r.TempleCode},
		{"Submission ancestors", "1", r.Ancestors},
		{"Submission descendants", "1", r.Descendants},
		{"Submission ordinance", "yes", r.Ordinance},
		{"Submission submitter name", "/Submitter-Name/", r.Submitter.Name},
		{"Submission submitter language", "English", r.Submitter.Language},
	}.run(t)
}

func TestFamily(t *testing.T) {

	f := g.Family

	intTestCases{
		{"Family list length was [%d]", 4, len(f)},
		{"Family 0 event list length was [%d]", 11, len(f[0].Event)},
	}.run(t)

	stringTestCases{
		{"Husband's age at annulment", "42y", f[0].Event[0].SpouseInfo[0].Age},
		{"Wife's age at divorce filing", "CHILD", f[0].Event[3].SpouseInfo[1].Age},
		{"Wife's spouse type at divorce filing", "WIFE", f[0].Event[3].SpouseInfo[1].Spouse},
		{"Number of children", "42", f[0].NumberOfChildren.Value},
		{"Family citation quality", "0", f[0].Citation[0].Quality},
		{"Family citation first file", "file1", f[0].Citation[0].Source.File[0]},
		{"Family object title", "A bmp picture", f[0].Object[0].Title},
		{"Family note 0", "A note about the family\nNote continued here. The word TEST should not be broken!", f[0].Note[0].Note},
		{"Family change date", "1 APR 1998", f[0].Changed.Stamp.Date},
		{"Family change time", "12:34:56.789", f[0].Changed.Stamp.Time},
	}.run(t)
}

func TestSource(t *testing.T) {

	s := g.Source[0]

	intTestCases{
		{"Source list length was [%d]", 1, len(g.Source)},
		{"Source file list length was [%d]", 2, len(s.File)},
	}.run(t)

	stringTestCases{
		{"Second source file name", "file2", s.File[1]},
		{"Source title", "Title of source\nTitle continued here. The word TEST should not be broken!", s.Title},
		{"Source submitter name", "A submitter", s.Submitter[0]},
		{"Source repository", "A repository", s.Repository[0]},
		{"Source periodical name", "A periodical name", s.Periodical},
		{"Source volume", "1", s.Volume},
		{"Source page", "3", s.Page[0]},
		{"Source film reference", "at 11", s.Film[0]},
		{"Source object 0 title", "A bmp picture", s.Object[0].Title},
		{"Source event data responsible agency", "Responsible agency", s.EventData.Agency},
		{"Source event data note", "A note about whatever\nNote continued here. The word TEST should not be broken!", s.EventData.Note[0].Note},
		{"Source birth and christening event tags", "BIRT, CHR", s.EventData.Event[0].Value},
		{"Source death event place", "Another place", s.EventData.Event[1].Place.Name},
	}.run(t)
}

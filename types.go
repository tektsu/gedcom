/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

// Gedcom is the top level structure.
type Gedcom struct {
	Header           *Header
	SubmissionRecord *SubmissionRecord
	Family           []*FamilyRecord
	Individual       []*IndividualRecord
	Media            []*MediaRecord
	Repository       []*RepositoryRecord
	Source           []*SourceRecord
	Submitter        []*SubmitterRecord
	Trailer          *Trailer
}

// Header is the heading of the Gedcom file.
type Header struct {
	SourceSystem SystemRecord
}

// SystemRecord is the administrative information about the data.
type SystemRecord struct {
	ID              string
	Version         string
	ProductName     string
	BusinessName    string
	Address         AddressRecord
	SourceName      string
	SourceDate      string
	SourceCopyright string
}

// SubmissionRecord ...
type SubmissionRecord struct {
	Xref string
}

// The Trailer doesn't do anything but mark the end of the file.
type Trailer struct {
}

// FamilyRecord describes a family unit.
type FamilyRecord struct {
	Xref    string
	Husband *IndividualRecord
	Wife    *IndividualRecord
	Child   []*IndividualRecord
	Event   []*EventRecord
}

// IndividualRecord describes a single person.
type IndividualRecord struct {
	Xref      string
	Name      []*NameRecord
	Sex       string
	Event     []*EventRecord
	Attribute []*EventRecord
	Parents   []*FamilyLinkRecord
	Family    []*FamilyLinkRecord
}

// MediaRecord is currently not implemented.
type MediaRecord struct {
}

// RepositoryRecord is currently not implemented.
type RepositoryRecord struct {
}

// SourceRecord describes a single source document.
type SourceRecord struct {
	Xref        string
	Title       string
	Author      string
	Abbr        string
	Publication string
	Text        string
	Media       []*MediaRecord
	Note        []*NoteRecord
}

// CitationRecord links another record and a source.
type CitationRecord struct {
	Source *SourceRecord
	Page   string
	Data   DataRecord
	Quay   string
	Media  []*MediaRecord
	Note   []*NoteRecord
}

// SubmitterRecord is not yet implemented.
type SubmitterRecord struct {
}

// NameRecord describes a person's name.
type NameRecord struct {
	Name     string
	Citation []*CitationRecord
	Note     []*NoteRecord
}

// DataRecord ...
type DataRecord struct {
	Date string
	Text []string
}

// EventRecord describes a life event.
type EventRecord struct {
	Tag      string
	Value    string
	Type     string
	Date     string
	Place    PlaceRecord
	Address  AddressRecord
	Age      string
	Agency   string
	Cause    string
	Citation []*CitationRecord
	Media    []*MediaRecord
	Note     []*NoteRecord
}

// NoteRecord describes a text note.
type NoteRecord struct {
	Note     string
	Citation []*CitationRecord
}

// PlaceRecord describes a location.
type PlaceRecord struct {
	Name     string
	Citation []*CitationRecord
	Note     []*NoteRecord
}

// FamilyLinkRecord ...
type FamilyLinkRecord struct {
	Family *FamilyRecord
	Type   string
	Note   []*NoteRecord
}

// AddressRecord describes and address.
type AddressRecord struct {
	Full       string
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
	Phone      string
}

// ObjectRecord describes a source object.
type ObjectRecord struct {
	Title       string
	Form        string
	File        string
	Type        string
	Place       string
	Date        string
	Media       string
	DocLocation string
	Note        []*NoteRecord
}

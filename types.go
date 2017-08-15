/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import "io"

// A Decoder reads and decodes GEDCOM objects from an input stream.
type Decoder struct {
	r                 io.Reader
	parsers           []parser
	refs              map[string]interface{}
	cbUnrecognizedTag func(int, string, string, string)
}

// Gedcom is the top level structure.
type Gedcom struct {
	Header     *HeaderRecord
	Submission *SubmissionRecord
	Submitter  []*SubmitterRecord
	Family     []*FamilyRecord
	Individual []*IndividualRecord
	Media      []*MediaRecord
	Repository []*RepositoryRecord
	Source     []*SourceRecord
	Trailer    *Trailer
}

// HeaderInfoRecord contains information about the Gedcom version
type HeaderInfoRecord struct {
	Version string
	Form    string
}

// HeaderRecord is the heading of the Gedcom file.
type HeaderRecord struct {
	Date        string
	Destination string
	File        string
	Copyright   string
	Language    string
	Timestamp   *TimestampRecord
	Encoding    *EncodingRecord
	Source      *HeaderSourceRecord
	Submitter   *SubmitterRecord
	Submission  *SubmissionRecord
	Info        *HeaderInfoRecord
	Note        *NoteRecord
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
	Xref        string
	FamilyFile  string
	TempleCode  string
	Ancestors   string
	Descendants string
	Ordinance   string
	Submitter   *SubmitterRecord
}

// SubmitterRecord describes a submitter.
type SubmitterRecord struct {
	Xref     string
	Name     string
	Language string
	Phone    []string
	Address  *AddressRecord
	Changed  *ChangedRecord
}

// ChangedRecord describes a document change.
type ChangedRecord struct {
	Stamp *TimestampRecord
	Note  []*NoteRecord
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

// HeaderSourceRecord ...
type HeaderSourceRecord struct {
	Source      string
	Version     string
	Name        string
	Form        string
	Corporation *CorpRecord
	Data        *HeaderDataRecord
}

// HeaderDataRecord ...
type HeaderDataRecord struct {
	Name      string
	Date      string
	Copyright string
}

// SourceRecord describes a single source document.
type SourceRecord struct {
	Xref        string
	Author      string
	Title       string
	Abbr        string
	Publication string
	Type        string
	Text        string
	MediaType   string
	Periodical  string
	Volume      string
	Page        []string
	Film        []string
	File        []string
	FileNumber  []string
	Place       []string
	Date        []string
	DateViewed  []string
	URL         []string
	DocLocation []string
	Repository  []string
	Submitter   []string
	LastChanged DataRecord
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
	Tag        string
	Value      string
	Type       string
	Date       string
	Place      PlaceRecord
	Address    AddressRecord
	Age        string
	Agency     string
	Cause      string
	Citation   []*CitationRecord
	Media      []*MediaRecord
	Note       []*NoteRecord
	Parents    []*FamilyLinkRecord
	SpouseInfo []*SpouseInfoRecord
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
	Family   *FamilyRecord
	Pedigree string
	Note     []*NoteRecord
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
	Title string
	Form  string
	File  string
	Note  []*NoteRecord
}

// EncodingRecord describes a character encoding.
type EncodingRecord struct {
	Name    string
	Version string
}

// TimestampRecord describes a timestamp.
type TimestampRecord struct {
	Date string
	Time string
}

// CorpRecord describes the corporation producing the software that generated the Gedcom.
type CorpRecord struct {
	Name    string
	Address *AddressRecord
	Phone   []string
}

// SpouseInfoRecord describes information about a spouse referenced in a family event.
type SpouseInfoRecord struct {
	Spouse string
	Age    string
}

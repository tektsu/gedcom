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
	Object     []*ObjectRecord
	Repository []*RepositoryRecord
	Source     []*SourceRecord
	Note       []*NoteRecord
	Trailer    *Trailer
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

// ChangedRecord describes a document change.
type ChangedRecord struct {
	Stamp *TimestampRecord
	Note  []*NoteRecord
}

// ChildRecord describes a child within a family.
type ChildRecord struct {
	FatherRelation string
	MotherRelation string
	Person         *IndividualRecord
}

// CitationRecord links another record and a source.
type CitationRecord struct {
	Source  *SourceRecord
	Page    string
	Data    DataRecord
	Quality string
	Object  []*ObjectRecord
	Note    []*NoteRecord
}

// CorpRecord describes the corporation producing the software that generated the Gedcom.
type CorpRecord struct {
	Name    string
	Address *AddressRecord
	Phone   []string
}

// DataRecord ...
type DataRecord struct {
	Date string
	Text []string
}

// EncodingRecord describes a character encoding.
type EncodingRecord struct {
	Name    string
	Version string
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
	Citation   []*CitationRecord
	Object     []*ObjectRecord
	Note       []*NoteRecord
	Cause      []*NoteRecord
	Parents    []*FamilyLinkRecord
	SpouseInfo []*SpouseInfoRecord
}

// FamilyLinkRecord ...
type FamilyLinkRecord struct {
	Family    *FamilyRecord
	Pedigree  string
	AdoptedBy string
	Note      []*NoteRecord
}

// FamilyRecord describes a family unit.
type FamilyRecord struct {
	Xref             string
	Husband          *IndividualRecord
	Wife             *IndividualRecord
	NumberOfChildren *EventRecord
	Changed          *ChangedRecord
	Child            []*ChildRecord
	Event            []*EventRecord
	Citation         []*CitationRecord
	Object           []*ObjectRecord
	Note             []*NoteRecord
}

// FileRecord ...
type FileRecord struct {
	Name  string
	Title string
	Form  string
}

// HeaderDataRecord ...
type HeaderDataRecord struct {
	Name      string
	Date      string
	Copyright string
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

// HeaderSourceRecord ...
type HeaderSourceRecord struct {
	Source      string
	Version     string
	Name        string
	Form        string
	Corporation *CorpRecord
	Data        *HeaderDataRecord
}

// IndividualRecord describes a single person.
type IndividualRecord struct {
	Xref      string
	Sex       string
	Changed   *ChangedRecord
	Photo     *ObjectRecord
	Name      []*NameRecord
	Event     []*EventRecord
	Attribute []*EventRecord
	Parents   []*FamilyLinkRecord
	Family    []*FamilyLinkRecord
	Citation  []*CitationRecord
	Object    []*ObjectRecord
	Note      []*NoteRecord
}

// NameRecord describes a person's name.
type NameRecord struct {
	Name     string
	Prefix   string
	Suffix   string
	Citation []*CitationRecord
	Note     []*NoteRecord
}

// NoteRecord describes a text note.
type NoteRecord struct {
	Xref     string
	Note     string
	Citation []*CitationRecord
}

// ObjectRecord describes a source object.
type ObjectRecord struct {
	Xref string
	File *FileRecord
}

// PlaceRecord describes a location.
type PlaceRecord struct {
	Name      string
	Latitude  string
	Longitude string
	Citation  []*CitationRecord
	Note      []*NoteRecord
}

// RepositoryRecord is currently not implemented.
type RepositoryRecord struct {
	Xref string
}

// SourceDataRecord describes events pertaining to this source
type SourceDataRecord struct {
	Agency string
	Event  []*EventRecord
	Note   []*NoteRecord
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
	Changed     *ChangedRecord
	EventData   *SourceDataRecord
	Note        []*NoteRecord
	Object      []*ObjectRecord
}

// SpouseInfoRecord describes information about a spouse referenced in a family event.
type SpouseInfoRecord struct {
	Spouse string
	Age    string
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

// TimestampRecord describes a timestamp.
type TimestampRecord struct {
	Date string
	Time string
}

// The Trailer doesn't do anything but mark the end of the file.
type Trailer struct {
}

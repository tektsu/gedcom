/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"io"
	"strings"
)

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:                 r,
		cbUnrecognizedTag: func(l int, t, v, x string) {},
	}
}

// Set the Unrecognized Tag Callback
func (d *Decoder) SetUnrecTagFunc(f func(int, string, string, string)) {
	if f == nil {
		return
	}
	d.cbUnrecognizedTag = f
}

// Decode reads the next GEDCOM-encoded value from its
// input and stores it in the value pointed to by v.
func (d *Decoder) Decode() (*Gedcom, error) {

	g := &Gedcom{
		Header:     &HeaderRecord{},
		Family:     make([]*FamilyRecord, 0),
		Individual: make([]*IndividualRecord, 0),
		Media:      make([]*MediaRecord, 0),
		Repository: make([]*RepositoryRecord, 0),
		Source:     make([]*SourceRecord, 0),
		Submitter:  make([]*SubmitterRecord, 0),
	}

	d.refs = make(map[string]interface{})
	d.parsers = []parser{makeRootParser(d, g)}
	d.scan(g)

	return g, nil
}

func (d *Decoder) scan(g *Gedcom) {
	s := &scanner{}
	buf := make([]byte, 512)

	n, err := d.r.Read(buf)
	if err != nil {
		// TODO
	}

	for n > 0 {
		pos := 0

		for {
			s.reset()
			offset, err := s.nextTag(buf[pos:n])
			pos += offset
			if err != nil {
				if err != io.EOF {
					println(err.Error())
					return
				}
				break
			}

			d.parsers[len(d.parsers)-1](s.level, string(s.tag), string(s.value), string(s.xref))

		}

		// shift unparsed bytes to start of buffer
		rest := copy(buf, buf[pos:])

		// top up buffer
		num, err := d.r.Read(buf[rest:len(buf)])
		if err != nil {
			break
		}

		n = rest + num - 1

	}

}

type parser func(level int, tag string, value string, xref string) error

func (d *Decoder) pushParser(p parser) {
	d.parsers = append(d.parsers, p)
}

func (d *Decoder) popParser(level int, tag string, value string, xref string) error {
	n := len(d.parsers) - 1
	if n < 1 {
		panic("MASSIVE ERROR") // TODO
	}
	d.parsers = d.parsers[0:n]

	return d.parsers[len(d.parsers)-1](level, tag, value, xref)
}

func (d *Decoder) individual(xref string) *IndividualRecord {
	if xref == "" {
		return &IndividualRecord{}
	}

	ref, found := d.refs[xref].(*IndividualRecord)
	if !found {
		rec := &IndividualRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

func (d *Decoder) family(xref string) *FamilyRecord {
	if xref == "" {
		return &FamilyRecord{}
	}

	ref, found := d.refs[xref].(*FamilyRecord)
	if !found {
		rec := &FamilyRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

func (d *Decoder) source(xref string) *SourceRecord {
	if xref == "" {
		return &SourceRecord{}
	}

	ref, found := d.refs[xref].(*SourceRecord)
	if !found {
		rec := &SourceRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

func (d *Decoder) submitter(xref string) *SubmitterRecord {
	if xref == "" {
		return &SubmitterRecord{}
	}

	ref, found := d.refs[xref].(*SubmitterRecord)
	if !found {
		rec := &SubmitterRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

func makeRootParser(d *Decoder, g *Gedcom) parser {
	return func(level int, tag string, value string, xref string) error {
		if level == 0 {
			switch tag {
			case "HEAD":
				obj := g.Header
				d.pushParser(makeHeaderParser(d, obj, level))
			case "INDI":
				obj := d.individual(xref)
				g.Individual = append(g.Individual, obj)
				d.pushParser(makeIndividualParser(d, obj, level))
			case "SUBM":
				obj := d.submitter(xref)
				g.Submitter = append(g.Submitter, obj)
				d.pushParser(makeSubmitterParser(d, obj, level))
			case "FAM":
				obj := d.family(xref)
				g.Family = append(g.Family, obj)
				d.pushParser(makeFamilyParser(d, obj, level))
			case "SOUR":
				obj := d.source(xref)
				g.Source = append(g.Source, obj)
				d.pushParser(makeSourceParser(d, obj, level))

			default:
				d.cbUnrecognizedTag(level, tag, value, xref)
				d.pushParser(makeSlurkParser(d, level))
			}
		}
		return nil
	}
}

func makeHeaderParser(d *Decoder, h *HeaderRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "COPR":
			h.Copyright = value
		case "FILE":
			h.File = value
		case "LANG":
			h.Language = value
		case "DEST":
			h.Destination = value
		case "CHAR":
			h.Encoding = &EncodingRecord{Name: value}
			d.pushParser(makeEncodingParser(d, h.Encoding, level))
		case "DATE":
			h.Timestamp = &TimestampRecord{Date: value}
			d.pushParser(makeTimestampParser(d, h.Timestamp, level))
		case "SOUR":
			h.Source = &GedcomSourceRecord{Source: value}
			d.pushParser(makeGedcomSourceParser(d, h.Source, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeIndividualParser(d *Decoder, i *IndividualRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "NAME":
			n := &NameRecord{Name: value}
			i.Name = append(i.Name, n)
			d.pushParser(makeNameParser(d, n, level))
		case "SEX":
			i.Sex = value
		case "BIRT", "CHR", "DEAT", "BURI", "CREM", "ADOP", "BAPM", "BARM", "BASM", "BLES", "CHRA", "CONF", "FCOM", "ORDN", "NATU", "EMIG", "IMMI", "CENS", "PROB", "WILL", "GRAD", "RETI", "EVEN":
			e := &EventRecord{Tag: tag, Value: value}
			i.Event = append(i.Event, e)
			d.pushParser(makeEventParser(d, e, level))
		case "CAST", "DSCR", "EDUC", "IDNO", "NATI", "NCHI", "NMR", "OCCU", "PROP", "RELI", "RESI", "SSN", "TITL", "FACT":
			e := &EventRecord{Tag: tag, Value: value}
			i.Attribute = append(i.Attribute, e)
			d.pushParser(makeEventParser(d, e, level))
		case "FAMC":
			family := d.family(stripXref(value))
			f := &FamilyLinkRecord{Family: family}
			i.Parents = append(i.Parents, f)
			d.pushParser(makeFamilyLinkParser(d, f, level))
		case "FAMS":
			family := d.family(stripXref(value))
			f := &FamilyLinkRecord{Family: family}
			i.Family = append(i.Family, f)
			d.pushParser(makeFamilyLinkParser(d, f, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeNameParser(d *Decoder, n *NameRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			n.Citation = append(n.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			n.Note = append(n.Note, r)
			d.pushParser(makeNoteParser(d, r, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeSourceParser(d *Decoder, s *SourceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "AUTH": // {0:1}
			s.Author = value
			d.pushParser(makeTextParser(d, &s.Author, level))
		case "TITL": // {0:1}
			s.Title = value
			d.pushParser(makeTextParser(d, &s.Title, level))
		case "ABBR": // {0:1}
			s.Abbr = value
			d.pushParser(makeTextParser(d, &s.Abbr, level))
		case "PUBL": // {0:1}
			s.Publication = value
			d.pushParser(makeTextParser(d, &s.Publication, level))
		case "TEXT": // {0:1}
			s.Text = value
			d.pushParser(makeTextParser(d, &s.Text, level))
		case "TYPE": // {0:1}
			s.Type = value
			d.pushParser(makeTextParser(d, &s.Type, level))
		case "CHAN": // {0:1}
			d.pushParser(makeDataParser(d, &s.LastChanged, level))
		case "NOTE": // {0:M}
			r := &NoteRecord{Note: value}
			s.Note = append(s.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		case "OBJE": // {0:M}
			o := &ObjectRecord{}
			d.pushParser(makeObjectParser(d, o, level))

		// Non-standard tags
		case "PERI": // {0:1}
			s.Periodical = value
			d.pushParser(makeTextParser(d, &s.Periodical, level))
		case "VOL": // {0:1}
			s.Volume = value
			d.pushParser(makeTextParser(d, &s.Volume, level))
		case "MEDI": // {0:1}
			s.MediaType = value
			d.pushParser(makeTextParser(d, &s.MediaType, level))
		case "PAGE": // {0:M}
			r := value
			s.Page = append(s.Page, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "FILM": // {0:M}
			r := value
			s.Film = append(s.Film, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "FILE": // {0:M}
			r := value
			s.File = append(s.File, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "FILN": // {0:M}
			r := value
			s.FileNumber = append(s.FileNumber, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "DATE": // {0:M}
			r := value
			s.Date = append(s.Date, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "PLAC": // {0:M}
			r := value
			s.Place = append(s.Place, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "DATV": // {0:M}
			r := value
			s.DateViewed = append(s.DateViewed, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "URL": // {0:M}
			r := value
			s.URL = append(s.URL, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "LOCA": // {0:M}
			r := value
			s.DocLocation = append(s.DocLocation, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "REPO": // {0:M}
			r := value
			s.Repository = append(s.Repository, r)
			d.pushParser(makeTextParser(d, &r, level))
		case "SUBM": // {0:M}
			r := value
			s.Submitter = append(s.Submitter, r)
			d.pushParser(makeTextParser(d, &r, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeCitationParser(d *Decoder, c *CitationRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "PAGE":
			c.Page = value
		case "QUAY":
			c.Quay = value
		case "NOTE":
			r := &NoteRecord{Note: value}
			c.Note = append(c.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		case "DATA":
			d.pushParser(makeDataParser(d, &c.Data, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeNoteParser(d *Decoder, n *NoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "CONT":
			n.Note = n.Note + "\n" + value
		case "CONC":
			n.Note = n.Note + value
		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			n.Citation = append(n.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeTextParser(d *Decoder, s *string, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "CONT":
			*s = *s + "\n" + value
		case "CONC":
			*s = *s + value

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeDataParser(d *Decoder, r *DataRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "DATE":
			r.Date = value
		case "TEXT":
			r.Text = append(r.Text, value)
			d.pushParser(makeTextParser(d, &r.Text[len(r.Text)-1], level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeEventParser(d *Decoder, e *EventRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "TYPE":
			e.Type = value
		case "DATE":
			e.Date = value
		case "PLAC":
			e.Place.Name = value
			d.pushParser(makePlaceParser(d, &e.Place, level))
		case "ADDR":
			e.Address.Full = value
			d.pushParser(makeAddressParser(d, &e.Address, level))
		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			e.Citation = append(e.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			e.Note = append(e.Note, r)
			d.pushParser(makeNoteParser(d, r, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makePlaceParser(d *Decoder, p *PlaceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			p.Citation = append(p.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			p.Note = append(p.Note, r)
			d.pushParser(makeNoteParser(d, r, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeFamilyLinkParser(d *Decoder, f *FamilyLinkRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "PEDI":
			f.Type = value
		case "NOTE":
			r := &NoteRecord{Note: value}
			f.Note = append(f.Note, r)
			d.pushParser(makeNoteParser(d, r, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeFamilyParser(d *Decoder, f *FamilyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "HUSB":
			f.Husband = d.individual(stripXref(value))
		case "WIFE":
			f.Wife = d.individual(stripXref(value))
		case "CHIL":
			f.Child = append(f.Child, d.individual(stripXref(value)))
		case "ANUL", "CENS", "DIV", "DIVF", "ENGA", "MARR", "MARB", "MARC", "MARL", "MARS", "EVEN":
			e := &EventRecord{Tag: tag, Value: value}
			f.Event = append(f.Event, e)
			d.pushParser(makeEventParser(d, e, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeAddressParser(d *Decoder, a *AddressRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "CONT":
			a.Full = a.Full + "\n" + value
		case "ADR1":
			a.Line1 = value
		case "ADR2":
			a.Line2 = value
		case "CITY":
			a.City = value
		case "STAE":
			a.State = value
		case "POST":
			a.PostalCode = value
		case "CTRY":
			a.Country = value
		case "PHON":
			a.Phone = value

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}

		return nil
	}
}

func makeObjectParser(d *Decoder, o *ObjectRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "TITL":
			o.Title = value
			d.pushParser(makeTextParser(d, &o.Title, level))
		case "FORM":
			o.Form = value
			d.pushParser(makeTextParser(d, &o.Form, level))
		case "FILE":
			o.File = value
			d.pushParser(makeTextParser(d, &o.File, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			o.Note = append(o.Note, r)
			d.pushParser(makeNoteParser(d, r, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeEncodingParser(d *Decoder, e *EncodingRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "VERS":
			e.Version = value

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeTimestampParser(d *Decoder, t *TimestampRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "TIME":
			t.Time = value

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeGedcomSourceParser(d *Decoder, r *GedcomSourceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "VERS":
			r.Version = value
		case "NAME":
			r.Name = value
		case "FORM":
			r.Form = value
		case "CORP":
			r.Corporation = &CorpRecord{Name: value}
			d.pushParser(makeCorpParser(d, r.Corporation, level))
		case "DATA":
			r.Data = &GedcomDataRecord{Name: value}
			d.pushParser(makeGedcomDataParser(d, r.Data, level))

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeGedcomDataParser(d *Decoder, r *GedcomDataRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "DATE":
			r.Date = value
		case "COPR":
			r.Copyright = value

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeCorpParser(d *Decoder, r *CorpRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "ADDR":
			r.Address = &AddressRecord{Full: value}
			d.pushParser(makeAddressParser(d, r.Address, level))
		case "PHON":
			p := value
			r.Phone = append(r.Phone, p)

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

func makeSubmitterParser(d *Decoder, r *SubmitterRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "NAME":
			r.Name = value
		case "LANG":
			r.Language = value
		case "ADDR":
			r.Address = &AddressRecord{Full: value}
			d.pushParser(makeAddressParser(d, r.Address, level))
		case "PHON":
			p := value
			r.Phone = append(r.Phone, p)

		default:
			d.cbUnrecognizedTag(level, tag, value, xref)
			d.pushParser(makeSlurkParser(d, level))
		}
		return nil
	}
}

// Just bypass this tag and all its children.
func makeSlurkParser(d *Decoder, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		return nil
	}
}

func stripXref(value string) string {
	return strings.Trim(value, "@")
}

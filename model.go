package musicbrainz

import (
	"regexp"
	"strconv"
	"time"
)

type Area struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Disambiguation string `json:"disambiguation"`
}

type Tag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type LifeSpan struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
	Ended bool   `json:"ended"`
}

type Artist struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Disambiguation string         `json:"disambiguation"`
	Area           Area           `json:"area"`
	Type           string         `json:"type"`
	TypeID         string         `json:"type-id"`
	BeginArea      Area           `json:"begin-area"`
	EndArea        Area           `json:"end-area"`
	Gender         string         `json:"gender"`
	Country        string         `json:"country"`
	LifeSpan       LifeSpan       `json:"life-span"`
	Tags           []Tag          `json:"tags"`
	Releases       []Release      `json:"releases"`
	ReleaseGroups  []ReleaseGroup `json:"release-groups"`
}

type ReleaseGroup struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	FirstReleaseDate string         `json:"first-release-date"`
	PrimaryType      string         `json:"primary-type"`
	PrimaryTypeID    string         `json:"primary-type-id"`
	SecondaryTypes   []string       `json:"secondary-types"`
	Disambiguation   string         `json:"disambiguation"`
	Tags             []Tag          `json:"tags"`
	Releases         []Release      `json:"releases"`
	ArtistCredit     []ArtistCredit `json:"artist-credit"`
	Genres           []Genre        `json:"genres"`
}

type Genre struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Disambiguation string `json:"disambiguation"`
	Count          int    `json:"count"`
}

type TextRepresentation struct {
	Language string `json:"language"`
	Script   string `json:"script"`
}

type ArtistCredit struct {
	Artist     Artist `json:"artist"`
	JoinPhrase string `json:"join-phrase"`
	Name       string `json:"name"`
}

type Recording struct {
	ID               string         `json:"id"`
	Disambiguation   string         `json:"disambiguation"`
	ArtistCredit     []ArtistCredit `json:"artist-credit"`
	Length           int            `json:"length"`
	Title            string         `json:"title"`
	FirstReleaseDate string         `json:"first-release-date"`
	Video            bool           `json:"video"`
}

type Track struct {
	ID           string         `json:"id"`
	Position     int            `json:"position"`
	Length       int            `json:"length"`
	Title        string         `json:"title"`
	ArtistCredit []ArtistCredit `json:"artist-credit"`
	Recording    Recording      `json:"recording"`
}

type Media struct {
	Format      string  `json:"format"`
	Title       string  `json:"title"`
	Position    int     `json:"position"`
	TrackCount  int     `json:"track-count"`
	TrackOffset int     `json:"track-offset"`
	FormatID    string  `json:"format-id"`
	Pregap      Track   `json:"pregap"`
	Tracks      []Track `json:"tracks"`
}

type Release struct {
	ID                 string             `json:"id"`
	Title              string             `json:"title"`
	Disambiguation     string             `json:"disambiguation"`
	TextRepresentation TextRepresentation `json:"text-representation"`
	Date               string             `json:"date"`
	Packaging          string             `json:"packaging"`
	PackagingID        string             `json:"packaging-id"`
	Barcode            string             `json:"barcode"`
	Quality            string             `json:"quality"`
	Country            string             `json:"country"`
	ReleaseGroup       ReleaseGroup       `json:"release-group"`
	ArtistCredit       []ArtistCredit     `json:"artist-credit"`
	Media              []Media            `json:"media"`
	LabelInfo          []LabelInfo        `json:"label-info"`
	Relations          []Relation         `json:"relations"`
}

var discogsIDRegex = regexp.MustCompile(`discogs\.com/release/(\d+)`)

func (r Release) DiscogsReleaseIDs() []int {
	var result []int

	for _, rel := range r.Relations {
		if rel.TargetType == "url" && rel.Type == "discogs" {
			matches := discogsIDRegex.FindStringSubmatch(rel.URL.Resource)
			if len(matches) > 0 {
				id, _ := strconv.Atoi(matches[1])
				result = append(result, id)
			}
		}
	}

	return result
}

type LabelInfo struct {
	CatalogNumber string `json:"catalog-number"`
	Label         Label  `json:"label"`
}

type LabelAlias struct {
	Name string `json:"name"`
}

type RelationDetail struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	TypeID         string `json:"type-id"`
	Type           string `json:"type_name"`
	Resource       string `json:"resource"`
	Disambiguation string `json:"disambiguation"`
}

type Relation struct {
	Label        RelationDetail    `json:"label"`
	Series       RelationDetail    `json:"series"`
	URL          RelationDetail    `json:"url"`
	TargetType   string            `json:"target-type"`
	Type         string            `json:"type"`
	TypeID       string            `json:"type-id"`
	Direction    string            `json:"direction"`
	TargetCredit string            `json:"target-credit"`
	SourceCredit string            `json:"source-credit"`
	OrderingKey  int               `json:"ordering-key"`
	AttributeIDs map[string]string `json:"attribute-ids"`
	Attributes   []string          `json:"attributes"`
}

var parentLabelRelationTypes = map[string]struct{}{
	"label ownership": {},
}

var parentLabelTypes = map[string]struct{}{
	"Imprint":             {},
	"Original Production": {},
}

func (r Relation) IsParentLabel() bool {
	if r.Label.ID == "" {
		return false
	}
	if r.Direction != "backward" {
		return false
	}
	if _, ok := parentLabelRelationTypes[r.Type]; !ok {
		return false
	}
	if _, ok := parentLabelTypes[r.Label.Type]; !ok {
		return false
	}
	return true
}

type Label struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Country   string     `json:"country"`
	LifeSpan  LifeSpan   `json:"life-span"`
	TypeID    string     `json:"type-id"`
	Relations []Relation `json:"relations"`
	Area      Area       `json:"area"`
}

func (a Label) ParentLabels() []Relation {
	var result []Relation
	for _, rel := range a.Relations {
		if rel.IsParentLabel() {
			result = append(result, rel)
		}
	}
	return result
}

type ReleaseType string

const (
	ReleaseTypeAlbum     ReleaseType = "Album"
	ReleaseTypeSingle    ReleaseType = "Single"
	ReleaseTypeEP        ReleaseType = "EP"
	ReleaseTypeBroadcast ReleaseType = "Broadcast"
	ReleaseTypeOther     ReleaseType = "Other"
)

type SecondaryReleaseType string

const (
	SecondaryReleaseTypeCompilation    SecondaryReleaseType = "Compilation"
	SecondaryReleaseTypeSoundtrack     SecondaryReleaseType = "Soundtrack"
	SecondaryReleaseTypeSpokenword     SecondaryReleaseType = "Spokenword"
	SecondaryReleaseTypeInterview      SecondaryReleaseType = "Interview"
	SecondaryReleaseTypeAudiobook      SecondaryReleaseType = "Audiobook"
	SecondaryReleaseTypeAudioDrama     SecondaryReleaseType = "Audio drama"
	SecondaryReleaseTypeLive           SecondaryReleaseType = "Live"
	SecondaryReleaseTypeRemix          SecondaryReleaseType = "Remix"
	SecondaryReleaseTypeDJMix          SecondaryReleaseType = "DJ-mix"
	SecondaryReleaseTypeMixtapeStreet  SecondaryReleaseType = "Mixtape/Street"
	SecondaryReleaseTypeDemo           SecondaryReleaseType = "Demo"
	SecondaryReleaseTypeFieldRecording SecondaryReleaseType = "Field recording"
)

type Record[Data any] struct {
	Date time.Time `json:"date"`
	Data Data      `json:"data"`
}

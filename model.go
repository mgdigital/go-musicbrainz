package musicbrainz

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	Score            int            `json:"score"`
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
	DiscCount   int     `json:"disc-count"`
	TrackCount  int     `json:"track-count"`
	TrackOffset int     `json:"track-offset"`
	FormatID    string  `json:"format-id"`
	Pregap      Track   `json:"pregap"`
	Tracks      []Track `json:"tracks"`
}

type Release struct {
	ID                 string             `json:"id"`
	Title              string             `json:"title"`
	Status             string             `json:"status"`
	Disambiguation     string             `json:"disambiguation"`
	TextRepresentation TextRepresentation `json:"text-representation"`
	Date               string             `json:"date"`
	Packaging          string             `json:"packaging"`
	PackagingID        string             `json:"packaging-id"`
	Barcode            string             `json:"barcode"`
	Quality            string             `json:"quality"`
	Country            string             `json:"country"`
	TrackCount         int                `json:"track-count"`
	Count              int                `json:"count"`
	Score              int                `json:"score"`
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

type SearchReleaseGroupRequest struct {
	Raw              string
	ArtistName       string
	ReleaseName      string
	FirstReleaseDate string
	Reid             string
	Rgid             string
}

func (r SearchReleaseGroupRequest) Query() string {
	var (
		orParts  []string
		andParts []string
	)
	if r.ArtistName != "" {
		orParts = append(orParts, queryPart("artistname", r.ArtistName)+"~")
	}
	if r.ReleaseName != "" {
		orParts = append(orParts, queryPart("release", r.ReleaseName)+"~")
	}
	if r.FirstReleaseDate != "" {
		orParts = append(orParts, queryPart("firstreleasedate", r.FirstReleaseDate))
	}
	if r.Raw != "" {
		andParts = append(andParts, r.Raw)
	}
	if len(orParts) > 0 {
		andParts = append(andParts, "("+strings.Join(orParts, " OR ")+")")
	}
	if r.Rgid != "" {
		andParts = append(andParts, queryPart("rgid", r.Rgid))
	}
	return strings.Join(andParts, " AND ")
}

type SearchReleaseRequest struct {
	Raw           string
	ArtistName    string
	ReleaseName   string
	ReleaseDate   string
	Format        string
	CatalogNumber string
	Tracks        int
	Reid          string
	Rgid          string
}

func queryPart(field, value string) string {
	return fmt.Sprintf("%s:\"%s\"", field, strings.ReplaceAll(value, "\"", "\\\""))
}

func (r SearchReleaseRequest) Query() string {
	var (
		orParts  []string
		andParts []string
	)
	if r.ArtistName != "" {
		orParts = append(orParts, queryPart("artistname", r.ArtistName)+"~")
	}
	if r.ReleaseName != "" {
		orParts = append(orParts, queryPart("release", r.ReleaseName)+"~")
	}
	if r.ReleaseDate != "" {
		orParts = append(orParts, queryPart("date", r.ReleaseDate))
	}
	if r.Format != "" {
		orParts = append(orParts, queryPart("format", r.Format))
	}
	if r.CatalogNumber != "" {
		orParts = append(orParts, queryPart("catno", r.CatalogNumber))
	}
	if r.Tracks > 0 {
		orParts = append(orParts, queryPart("tracks", strconv.Itoa(r.Tracks)))
	}
	if len(orParts) > 0 {
		andParts = append(andParts, "("+strings.Join(orParts, " OR ")+")")
	}
	if r.Raw != "" {
		andParts = append(andParts, r.Raw)
	}
	if r.Reid != "" {
		andParts = append(andParts, queryPart("reid", r.Reid))
	}
	if r.Rgid != "" {
		andParts = append(andParts, queryPart("rgid", r.Rgid))
	}
	return strings.Join(andParts, " AND ")
}

type SearchReleaseResult struct {
	Created  time.Time `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []Release `json:"releases"`
}

type SearchReleaseGroupResult struct {
	Created       time.Time      `json:"created"`
	Count         int            `json:"count"`
	Offset        int            `json:"offset"`
	ReleaseGroups []ReleaseGroup `json:"release-groups"`
}

type IDName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (g ReleaseGroup) String() string {
	var artists []string
	for _, credit := range g.ArtistCredit {
		artists = append(artists, credit.Name)
	}
	return fmt.Sprintf("[%s] %s - %s (%s)", g.ID, strings.Join(artists, "; "), g.Title, g.FirstReleaseDate)
}

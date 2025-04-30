package activitypub

import (
	"encoding/json"
	"time"
)

// ActivityPub const
const (
	// Activity Types: https://www.w3.org/TR/activitystreams-vocabulary/#activity-types
	ActivityTypeCreate   = "Create"
	ActivityTypeFollow   = "Follow"
	ActivityTypeAccept   = "Accept"
	ActivityTypeReject   = "Reject"
	ActivityTypeDelete   = "Delete"
	ActivityTypeAnnounce = "Announce"
	ActivityTypeLike     = "Like"
	ActivityTypeUpdate   = "Update"
	ActivityTypeUndo     = "Undo"

	// Object Types: https://www.w3.org/TR/activitystreams-vocabulary/#object-types
	ObjectTypeNote         = "Note"
	ObjectTypePerson       = "Person"
	ObjectTypeImage        = "Image"
	ObjectTypePlace        = "Place"
	ObjectTypeGroup        = "Group"
	ObjectTypeRelationship = "Relationship"
	ObjectTypeActivity     = "Activity"
	ObjectTypeTombstone    = "Tombstone"

	// Actor Types: https://www.w3.org/TR/activitystreams-vocabulary/#actor-types
	ActorTypeApplication  = "Application"
	ActorTypeGroup        = "Group"
	ActorTypeOrganization = "Organization"
	ActorTypePerson       = "Person"
	ActorTypeService      = "Service"
)

// Context: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-context
type Context []interface{}

// Object: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-object
type Object struct {
	Context      Context   `json:"@context,omitempty"`
	ID           string    `json:"id,omitempty"`
	Type         string    `json:"type"`
	AttributedTo string    `json:"attributedTo,omitempty"`
	Name         string    `json:"name,omitempty"`
	Summary      string    `json:"summary,omitempty"`
	Content      string    `json:"content,omitempty"`
	URL          string    `json:"url,omitempty"`
	MediaType    string    `json:"mediaType,omitempty"`
	Published    time.Time `json:"published,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
	Icon         *Image    `json:"icon,omitempty"`
	Image        *Image    `json:"image,omitempty"`
	Location     *Place    `json:"location,omitempty"`
	Tag          []Object  `json:"tag,omitempty"`
	Attachment   []Object  `json:"attachment,omitempty"`
	InReplyTo    string    `json:"inReplyTo,omitempty"`
	To           []string  `json:"to,omitempty"`
	Cc           []string  `json:"cc,omitempty"`
	Bto          []string  `json:"bto,omitempty"`
	Bcc          []string  `json:"bcc,omitempty"`
	Generator    *Object   `json:"generator,omitempty"`
}

// Link: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-link
type Link struct {
	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"` // relation
	MediaType string `json:"mediaType,omitempty"`
	Name      string `json:"name,omitempty"`
	HrefLang  string `json:"hrefLang,omitempty"`
	Height    int    `json:"height,omitempty"`
	Width     int    `json:"width,omitempty"`
}

// Activity: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-activity
type Activity struct {
	Context   Context     `json:"@context,omitempty"`
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Actor     string      `json:"actor"`
	Object    interface{} `json:"object"`
	Target    string      `json:"target,omitempty"`
	Result    interface{} `json:"result,omitempty"`
	Origin    string      `json:"origin,omitempty"`
	To        []string    `json:"to,omitempty"`
	Cc        []string    `json:"cc,omitempty"`  // Carbon Copy
	Bto       []string    `json:"bto,omitempty"` // Blind To
	Bcc       []string    `json:"bcc,omitempty"` // Blind Carbon Copy
	Published time.Time   `json:"published,omitempty"`
	Updated   time.Time   `json:"updated,omitempty"`
}

// Collection: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-collection
type Collection struct {
	Context    Context     `json:"@context,omitempty"`
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	TotalItems int         `json:"totalItems"`
	Items      interface{} `json:"items,omitempty"`
	First      string      `json:"first,omitempty"`
	Last       string      `json:"last,omitempty"`
	Current    string      `json:"current,omitempty"`
}

// OrderedCollection: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollection
type OrderedCollection struct {
	Context      Context     `json:"@context,omitempty"`
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	TotalItems   int         `json:"totalItems"`
	OrderedItems interface{} `json:"orderedItems,omitempty"`
	First        string      `json:"first,omitempty"`
	Last         string      `json:"last,omitempty"`
	Current      string      `json:"current,omitempty"`
}

// CollectionPage: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-collectionpage
type CollectionPage struct {
	Context Context     `json:"@context,omitempty"`
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	PartOf  string      `json:"partOf"`
	Items   interface{} `json:"items"`
	Next    string      `json:"next,omitempty"`
	Prev    string      `json:"prev,omitempty"`
}

// OrderedCollectionPage: Core Types, https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollectionpage
type OrderedCollectionPage struct {
	Context      Context     `json:"@context,omitempty"`
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	PartOf       string      `json:"partOf"`
	OrderedItems interface{} `json:"orderedItems"`
	Next         string      `json:"next,omitempty"`
	Prev         string      `json:"prev,omitempty"`
	StartIndex   int         `json:"startIndex,omitempty"`
}

// Image: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-image
type Image struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Name      string `json:"name,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

// Place: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-place
type Place struct {
	Type      string    `json:"type"`
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Accuracy  float64   `json:"accuracy,omitempty"`
	Altitude  float64   `json:"altitude,omitempty"`
	Radius    float64   `json:"radius,omitempty"`
	Units     string    `json:"units,omitempty"`
	Published time.Time `json:"published,omitempty"`
	Updated   time.Time `json:"updated,omitempty"`
}

// Person: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-person
type Person struct {
	Context           Context   `json:"@context,omitempty"`
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Name              string    `json:"name,omitempty"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	Outbox            string    `json:"outbox"`
	Following         string    `json:"following,omitempty"`
	Followers         string    `json:"followers,omitempty"`
	Liked             string    `json:"liked,omitempty"`
	URL               string    `json:"url,omitempty"`
	PublicKey         PublicKey `json:"publicKey,omitempty"`
	Icon              *Image    `json:"icon,omitempty"`
	Image             *Image    `json:"image,omitempty"`
	Tag               []Object  `json:"tag,omitempty"`
	Attachment        []Object  `json:"attachment,omitempty"`
	Published         time.Time `json:"published,omitempty"`
	Updated           time.Time `json:"updated,omitempty"`
}

// PublicKey:
type PublicKey struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

func DefaultContext() Context {
	return Context{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		map[string]interface{}{
			"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
		},
	}
}

// ToJSON convert input value to JSON.
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON convert JSON format string
func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

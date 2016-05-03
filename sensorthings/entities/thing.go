package entities

import (
	"encoding/json"
)

// Thing in SensorThings represents a physical object in the real world. A Thing is a good starting
// point in which to start creating the SensorThings model structure. A Thing has a Location and one or
// more Datastreams to collect Observations. A minimal Thing can be created without a Location and Datastream
// and there are options to create a Things with a nested linked Location and Datastream.
type Thing struct {
	ID                     string                `json:"@iot.id"`
	NavSelf                string                `json:"@iot.selfLink"`
	Description            string                `json:"description"`
	Properties             map[string]string     `json:"properties,omitempty"`
	NavLocations           string                `json:"Locations@iot.navigationLink,omitempty"`
	NavDatastreams         string                `json:"Datastreams@iot.navigationLink,omitempty"`
	NavHistoricalLocations string                `json:"HistoricalLocations@iot.navigationLink,omitempty"`
	Locations              []*Location           `json:"Locations,omitempty"`
	Datastreams            []*Datastream         `json:"Datastreams,omitempty"`
	HistoricalLocations    []*HistoricalLocation `json:"HistoricalLocations,omitempty"`
}

// GetEntityType returns the EntityType for Thing
func (t *Thing) GetEntityType() EntityType {
	return EntityTypeThing
}

// ParseEntity tries to parse the given json byte array into the current entity
func (t *Thing) ParseEntity(data []byte) error {
	thing := &t
	err := json.Unmarshal(data, thing)
	if err != nil {
		return err
	}

	return nil
}

// ContainsMandatoryParams checks if all mandatory params for Thing are available before posting.
func (t *Thing) ContainsMandatoryParams() (bool, []error) {
	err := []error{}
	CheckMandatoryParam(&err, t.Description, t.GetEntityType(), "description")

	if len(err) != 0 {
		return false, err
	}

	return true, nil
}

// SetLinks sets the entity specific navigation links if needed
func (t *Thing) SetLinks(externalURL string) {
	t.NavSelf = CreateEntitySefLink(externalURL, EntityLinkThings.ToString(), t.ID)
	t.NavLocations = CreateEntityLink(t.Locations == nil, EntityLinkThings.ToString(), EntityLinkLocations.ToString(), t.ID)
	t.NavDatastreams = CreateEntityLink(t.Datastreams == nil, EntityLinkThings.ToString(), EntityLinkDatastreams.ToString(), t.ID)
	t.NavHistoricalLocations = CreateEntityLink(t.HistoricalLocations == nil, EntityLinkThings.ToString(), EntityLinkHistoricalLocations.ToString(), t.ID)
}
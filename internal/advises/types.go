package advises

type Advise struct {
	Places []Place `json:"places" jsonschema_description:"places associated with a given place type"`
}

type Place struct {
	Name        string `json:"name" jsonschema_description:"The name of the place to go"`
	Link        string `json:"link" jsonschema_description:"The link to the place on the tripadvisor"`
	Description string `json:"description" jsonschema_description:"The description of the place to go, with no links"`
}

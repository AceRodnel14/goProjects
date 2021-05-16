package main

type Body struct {
	Graphql Graphql `json:"graphql"`
}

type Graphql struct {
	ShortcodeMedia ShortcodeMedia `json:"shortcode_media"`
}

type ShortcodeMedia struct {
	Typename         string             `json:"__typename"`
	Shortcode        string             `json:"shortcode"`
	Dimensions       Dimensions         `json:"dimensions"`
	DisplayResources []DisplayResources `json:"display_resources"`
	Caption          Caption            `json:"edge_media_to_caption"`
	Sidecar          Sidecar            `json:"edge_sidecar_to_children,omitempty"`
	VideoURL         string             `json:"video_url,omitempty"`
	Owner            Owner              `json:"owner,omitempty"`
	ProductType      string             `json:"product_type,omitempty"`
}

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Caption struct {
	Edges []Edges `json:"edges"`
}

type Edges struct {
	Node Node `json:"node"`
}

type Node struct {
	Text             string             `json:"text,omitempty"`
	NodeTypename     string             `json:"__typename,omitempty"`
	NodeDimensions   Dimensions         `json:"dimensions,omitempty"`
	DisplayResources []DisplayResources `json:"display_resources,omitempty"`
	VideoURL         string             `json:"video_url,omitempty"`
}

type Owner struct {
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"full_name,omitempty"`
}

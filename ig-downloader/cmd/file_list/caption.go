package main

type CaptionObject struct {
	OwnerInfo   CaptionOwnerInfo `json:"owner"`
	UploadDate  string           `json:"upload_date"`
	PostCaption string           `json:"caption"`
	PostType    string           `json:"type"`
	Count       int              `json:"count"`
}

type CaptionOwnerInfo struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

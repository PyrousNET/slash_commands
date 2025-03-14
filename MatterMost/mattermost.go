package MatterMost

type Response struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

type Attachment struct {
	Text     string   `json:"text"`
	ImageUrl string   `json:"image_url"`
	Actions  []Action `json:"actions"`
}

type Action struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Integration Integration `json:"integration"`
}

type Integration struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Context map[string]string `json:"context"`
	Params  map[string]string `json:"params"`
}

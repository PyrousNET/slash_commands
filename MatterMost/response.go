package MatterMost

type HamCallSign struct {
	CallSign string
	Name     string
	City     string
	Last3    string
	Class    string
	Status   string
}

func (hCS *HamCallSign) GetResponseString() string {
	return "| Data | Value |\n| :------ | :-------|\n| Callsign | " + hCS.CallSign +
		" |\n| Name | " + hCS.Name +
		" |\n| City | " + hCS.City +
		" |\n| Last3 | " + hCS.Last3 +
		" |\n| Class | " + hCS.Class +
		" |\n| Status | " + hCS.Status + " |"
}

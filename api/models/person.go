package models

type Person struct {
	Age           int    `json:"age"`
	Workclass     string `json:"workclass"`
	Fnlwgt        int    `json:"fnlwgt"`
	Education     string `json:"education"`
	EducationNum  int    `json:"education-num"`
	MaritalStatus string `json:"marital-status"`
	Occupation    string `json:"occupation"`
	Relationship  string `json:"relationship"`
	Race          string `json:"race"`
	Sex           string `json:"sex"`
	CapitalGain   int    `json:"capital-gain"`
	CapitalLoss   int    `json:"capital-loss"`
	HoursPerWeek  int    `json:"hours-per-week"`
	NativeCountry string `json:"native-country"`
	Income        string `json:"income"`
}

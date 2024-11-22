package models

type UserConfig struct {
	UserID        int    `json:"userID" gorm:"column:user_id"`
	Age           string `json:"age" gorm:"column:age"`
	Workclass     string `json:"workclass" gorm:"column:workclass"`
	Fnlwgt        string `json:"fnlwgt" gorm:"column:fnlwgt"`
	Education     string `json:"education" gorm:"column:education"`
	EducationNum  string `json:"educationNum" gorm:"column:education_num"`
	MaritalStatus string `json:"maritalStatus" gorm:"column:marital_status"`
	Occupation    string `json:"occupation" gorm:"column:occupation"`
	Relationship  string `json:"relationship" gorm:"column:relationship"`
	Race          string `json:"race" gorm:"column:race"`
	Sex           string `json:"sex" gorm:"column:sex"`
	CapitalGain   string `json:"capitalGain" gorm:"column:capital_gain"`
	CapitalLoss   string `json:"capitalLoss" gorm:"column:capital_loss"`
	HoursPerWeek  string `json:"hoursPerWeek" gorm:"column:hours_per_week"`
	NativeCountry string `json:"nativeCountry" gorm:"column:native_country"`
	Income        string `json:"income" gorm:"column:income"`
	SortBy        string `json:"sortBy" gorm:"column:sort_by"`
	SortOrder     string `json:"sortOrder" gorm:"column:sort_order"`
	Paginator     string `json:"paginator" gorm:"column:paginator"`
}

package shared

type Image struct {
	Id        int    `json:"id" gorm:"column:id"`
	Url       string `json:"url" gorm:"column:url"`
	Width     int    `json:"width" gorm:"column:width"`
	Height    int    `json:"height" gorm:"column:height"`
	HostName  string `json:"hostName" gorm:"-"`
	Extension string `json:"extension" gorm:"-"`
}

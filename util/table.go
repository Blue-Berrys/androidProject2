package util

//套一个壳上去，这样返回的json结构就是接口要求的

type User struct {
	Id              uint   `json:"id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	BackGroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	WorkCount       int64  `json:"work_count"`
}

type FriendsChat struct {
	Id            uint `json:"id"`
	User          `json:"user"`
	ImageUrl      string `json:"image_url"`
	FavoriteCount int64  `json:"favorite_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Content       string `json:"content"`
	CreateDate    string `json:"create_date"`
}

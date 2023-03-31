package util

//套一个壳上去，这样返回的json结构就是接口要求的

type User struct {
	Id              uint   `json:"id"`
	UserId          int64  `json:"user_id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar,omitempty"`
	BackGroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	WorkCount       int64  `json:"work_count"`
}

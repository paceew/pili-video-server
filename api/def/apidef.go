package def

const (
	SCHEDULER_ADDR = "localhost:9001/"
	PAGE_NUM       = 10
)

//request
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type UserModifyPwd struct {
	PTPwd string `json:"pass_time_pwd"`
	CPwd  string `json:"current_pwd"`
}

type UserModifyInfo struct {
	Name string `json:"user_name"`
}

type NewVideo struct {
	AuthorName string `json:"author_name"`
	VideoName  string `json:"video_name"`
}

type NewComment struct {
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
}

//response
type SignedUp struct {
	Success   bool   `json:"message"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Id   int    `json:"user_id"`
	Pwd  string `json:"user_pwd"`
	Name string `json:"user_name"`
}

type VideosList struct {
	Videos []*VideoInfo `json:"videos_list"`
}

type CommentsList struct {
	Comments []*CommentInfo `json:"comments_list"`
}

type MessageList struct {
	Messages []*Message `json:"messages_list"`
}

type LikeNumber struct {
	Count int `json:"like_num"`
}

type LikeStatus struct {
	Like bool `json:"like"`
}

//Data model
type User struct {
	Id       int
	Username string
	Pwd      string
}

type SimpleSession struct {
	Id       string `json:"session_id"`
	Username string `json:"user_name"`
}

type VideoInfo struct {
	Id           string `json:"video_id"`
	Name         string `json:"video_name"`
	AuthorId     int    `json:"video_aid"`
	DisplayCtime string `json:"video_ct"`
	Modular      string `json:"video_mod"`
	LikeNum      int    `json:"like_num"`
	CollectNum   int    `json:"collect_num"`
	CommentNum   int    `json:"comment_num"`
}

type CommentInfo struct {
	Id         string `json:"comment_id"`
	AuthorName string `json :"author_name"`
	VideoName  string `json:"vidoe_id"`
	Content    string `json:"content"`
}

type Message struct {
	FriendName   string         `json:"message_firend_name"`
	Num          int            `json:"message_number"`
	FristMessage *MessageDetail `json:"first_message"`
}

type MessageDetail struct {
	SendName string `json:"message_send_name"`
	Content  string `json:"message_content"`
	SendTime string `json:"message_send_time"`
	Status   int    `json:"message_status"`
}

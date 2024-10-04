package annotalk

/**
Messages examples
	42["initChat",{"gender":"man","partner_gender":"woman","captchaID":""}]
	42["onSearchingPartner",{}]
	42["onChatStart",{"partner_gender":"man","chatId":"EJG1MLuVO9pcGZgwZOBM"}]
	42["onChatStart",{"partner_gender":"man","chatId":"EJG1MLuVO9pcGZgwZOBM"}]
	42["onTyping",{}]
	42["onMessage",{"message":"szia kor","isYou":0}]
	42["onDoneTyping",{}]
	42["typing",{}]
	42["sendMessage",{"message":"Szi 27 "}]
	42["doneTyping",{}]
	42["onMessage",{"message":"Szi 27 ","isYou":1}]
	42["onChatEnd",{}]
	42["lookForPartner",{}]
	42["onSearchingPartner",{}]
*/

const (
	OnStatistics       = "onStatistics"
	InitChat           = "initChat"
	OnChatStart        = "onChatStart"
	OnTyping           = "onTyping"
	OnMessage          = "onMessage"
	OnDoneTyping       = "onDoneTyping"
	Typing             = "typing"
	SendMessage        = "sendMessage"
	DoneTyping         = "doneTyping"
	OnChatEnd          = "onChatEnd"
	LookForPartner     = "lookForPartner"
	OnSearchingPartner = "onSearchingPartner"
)

type MessageType int

const (
	MessageFromPartner MessageType = 0
	MessageFromYou                 = 1
)

type StatisticsData struct {
	Man struct {
		WantWithWoman    int `json:"wantWithWoman"`
		WantWithMan      int `json:"wantWithMan"`
		WantWithWhatever int `json:"wantWithWhatever"`
	} `json:"man"`
	Woman struct {
		WantWithWoman    int `json:"wantWithWoman"`
		WantWithMan      int `json:"wantWithMan"`
		WantWithWhatever int `json:"wantWithWhatever"`
	} `json:"woman"`
}

type InitChatData struct {
	Gender        PersonGender `json:"gender"`
	PartnerGender PersonGender `json:"partner_gender"`
	CaptchaID     string       `json:"captchaID"`
}

type OnChatStartData struct {
	PartnerGender string `json:"partner"`
	ChatID        string `json:"chatId"`
}

type OnMessageData struct {
	Message string      `json:"message"`
	IsYou   MessageType `json:"isYou"` // 0 - partner, 1 - you
}

type SendMessageData struct {
	Message string `json:"message"`
}

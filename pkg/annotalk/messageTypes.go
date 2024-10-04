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

type MessageType string

// Incoming messages
const (
	OnStatistics       MessageType = "onStatistics"
	OnChatStart        MessageType = "onChatStart"
	OnTyping           MessageType = "onTyping"
	OnMessage          MessageType = "onMessage"
	OnDoneTyping       MessageType = "onDoneTyping"
	OnChatEnd          MessageType = "onChatEnd"
	OnSearchingPartner MessageType = "onSearchingPartner"
)

// Outgoing messages types
const (
	InitChat       MessageType = "initChat"
	Typing         MessageType = "typing"
	DoneTyping     MessageType = "doneTyping"
	SendMessage    MessageType = "sendMessage"
	LookForPartner MessageType = "lookForPartner"
	LeaveChat      MessageType = "leaveChat"
)

type Sender int

const (
	MessageFromPartner Sender = 0
	MessageFromYou     Sender = 1
)

type OnStatisticsData struct {
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

func NewOnStatisticsData(data map[string]interface{}) OnStatisticsData {
	return OnStatisticsData{
		Man: struct {
			WantWithWoman    int `json:"wantWithWoman"`
			WantWithMan      int `json:"wantWithMan"`
			WantWithWhatever int `json:"wantWithWhatever"`
		}{
			int((data["man"]).(map[string]interface{})["wantWithWoman"].(float64)),
			int((data["man"]).(map[string]interface{})["wantWithMan"].(float64)),
			int((data["man"]).(map[string]interface{})["wantWithWhatever"].(float64)),
		},
		Woman: struct {
			WantWithWoman    int `json:"wantWithWoman"`
			WantWithMan      int `json:"wantWithMan"`
			WantWithWhatever int `json:"wantWithWhatever"`
		}{
			int((data["woman"]).(map[string]interface{})["wantWithWoman"].(float64)),
			int((data["woman"]).(map[string]interface{})["wantWithMan"].(float64)),
			int((data["woman"]).(map[string]interface{})["wantWithWhatever"].(float64)),
		},
	}
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

func NewOnChatStartData(data map[string]interface{}) OnChatStartData {
	return OnChatStartData{
		PartnerGender: data["partner_gender"].(string),
		ChatID:        data["chatId"].(string),
	}
}

type OnMessageData struct {
	Message string `json:"message"`
	IsYou   Sender `json:"isYou"` // 0 - partner, 1 - you
}

func NewOnMessageData(data map[string]interface{}) OnMessageData {
	return OnMessageData{
		Message: data["message"].(string),
		IsYou:   Sender(data["isYou"].(float64)),
	}
}

type SendMessageData struct {
	Message string `json:"message"`
}

func NewSendMessageData(data map[string]interface{}) SendMessageData {
	return SendMessageData{
		Message: data["message"].(string),
	}
}

type MessageEvents struct {
	Stats            chan OnStatisticsData
	ChatStart        chan OnChatStartData
	Message          chan OnMessageData
	ChatEnd          chan struct{}
	SearchingPartner chan struct{}
}

func NewMessageEvents() *MessageEvents {
	return &MessageEvents{
		Stats:            make(chan OnStatisticsData),
		ChatStart:        make(chan OnChatStartData),
		Message:          make(chan OnMessageData),
		ChatEnd:          make(chan struct{}),
		SearchingPartner: make(chan struct{}),
	}
}

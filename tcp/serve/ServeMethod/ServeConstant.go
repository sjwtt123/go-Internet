package ServeMethod

// 登录选项
const (
	ServeCommandLogin    = "1"
	ServeCommandRegister = "2"
)

// 消息类型
const (
	TypeOnline    = "Online"
	TypeUnderLine = "UnderLine"

	TypeList = "List"

	TypePrivate = "Private"

	TypeRadio = "Every"

	TypeHeart = "Ping"
)

// 传入响应数据
const (
	ReceiveUserExists     = "isCreate"
	ReceiveUserExistsOrLo = "CreateOrLo"
	ReceiveSuccess        = "success"
	ReceiveTrue           = "true"
	ReceiveFalse          = "false"
)

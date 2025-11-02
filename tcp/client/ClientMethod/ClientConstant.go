package ClientMethod

// 登录选择
const (
	ClientCommandLogin    = "1"
	ClientCommandRegister = "2"
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

// 服务器响应消息
const (
	ResponseUserExists     = "isCreate"
	ResponseUserExistsOrLo = "CreateOrLo"
	ResponseLoginFailed    = "false"
)

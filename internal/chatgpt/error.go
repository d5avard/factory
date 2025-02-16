package chatgpt

type Error struct {
	Event_id string `json:"event_id"`
	Type     string `json:"type"`
	Error    struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

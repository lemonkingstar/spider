package response

// D 通用返回数据结构
type D struct {
	Status 	int			`json:"status"`
	Code   	int         `json:"code"`
	Message	string		`json:"message,omitempty"`
	Data	interface{} `json:"data,omitempty"`
}

/*
数据结构字段说明：

Status：状态码，可以用来判断返回结果是否正确，非0为异常（有些系统使用 Result bool，来判断，比如cmdb，无关大体）
Code: 错误码，根据实际情况返回
Message: 消息字段
Data: 响应数据，可以是map，也可以是数组。列表数据一般会将 Count信息添加到其中（有些系统会将 Count拉到最外层，根据实际需要选择）

下面是CMDB系统的响应数据结构示例，可以参考：

type D struct {
	Result 	bool 		`json:"result"`
	Code   	int         `json:"code"`
	Message	string		`json:"message,omitempty"`
	Data	interface{} `json:"data,omitempty"`
}

type D struct {
	Result 	bool 		`json:"result"`
	Code   	int         `json:"code"`
	Message	string		`json:"message,omitempty"`
	Data	struct{
		Count 	int 						`json:"count"`
		Info	[]map[string]interface{}	`json:"info"`
	} `json:"data,omitempty"`
}

 */

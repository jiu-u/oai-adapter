package v1

// https://docs.siliconflow.cn/cn/api-reference/videos/videos_submit

type (
	VideoRequest struct {
		Model          string `json:"model" binding:"required"`
		Prompt         string `json:"prompt" binding:"required"`
		ImageSize      string `json:"image_size" binding:"required"`
		NegativePrompt string `json:"negative_prompt,omitempty"`
		Image          string `json:"image" binding:"required"`
		Seed           uint64 `json:"seed,omitempty"`
		// exp
		Style string `json:"style,omitempty"`
	}
	VideoResponse struct {
		RequestId string `json:"requestId"`
	}
	VideoStatusRequest struct {
		RequestId string `json:"requestId"`
	}
	VideoStatusResponse struct {
		Status       string        `json:"status"` // Succeed, InQueue, InProgress, Failed
		RawRequestId string        `json:"rawRequestId"`
		Reason       string        `json:"reason"`
		Results      []VideoResult `json:"results"`
		// 附加的错误信息
		Error any `json:"error"`
	}
	VideoResult struct {
		Videos []VideoItem `json:"videos"`
		Timing VideoTiming `json:"timing"`
		Seed   int         `json:"seed"`
	}
	VideoTiming struct {
		Inference int `json:"inference"`
	}
	VideoItem struct {
		Url string `json:"url"`
	}
)

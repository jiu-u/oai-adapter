package volcengine

type (
	CreateVideoTaskRequest struct {
		Model       string                   `json:"model" binding:"required"`
		Content     []CreateVideoTaskContent `json:"content" binding:"required"`
		Resolution  string                   `json:"resolution,omitempty"`
		Ratio       string                   `json:"ratio,omitempty"`
		Duration    int                      `json:"duration,omitempty"`
		Watermark   bool                     `json:"watermark,omitempty"`
		Seed        int64                    `json:"seed,omitempty"`
		CameraFixed bool                     `json:"camera_fixed,omitempty"`
	}
	CreateVideoTaskContent struct {
		Type     string                `json:"type" binding:"required"`
		Text     string                `json:"text"`
		ImageUrl *CreateVideoTaskImage `json:"image_url,omitempty"`
	}
	CreateVideoTaskImage struct {
		Url string `json:"url" binding:"required"`
	}
	CreateVideoTaskResponse struct {
		Id string `json:"id"`
	}
	VideoTaskStatusRequest struct {
		Id string `json:"id" binding:"required"`
	}
	VideoTaskStatusResponse struct {
		Id        string        `json:"id"`
		Model     string        `json:"model"`
		Status    string        `json:"status"`
		Error     any           `json:"error"`
		CreatedAt int64         `json:"created_at"`
		UpdatedAt int64         `json:"updated_at"`
		Content   *VideoContent `json:"content"`
		Usage     any           `json:"usage"`
	}
	VideoUsage struct {
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
	VideoContent struct {
		VideoUrl string `json:"video_url"`
	}
)

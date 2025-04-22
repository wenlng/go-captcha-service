/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package adapt

type CaptData struct {
	CaptchaKey        string `json:"captcha_key,omitempty"`
	MasterImageBase64 string `json:"master_image_base64,omitempty"`
	ThumbImageBase64  string `json:"thumb_image_base64,omitempty"`
	MasterWidth       int32  `json:"master_width,omitempty"`
	MasterHeight      int32  `json:"master_height,omitempty"`
	ThumbWidth        int32  `json:"thumb_width,omitempty"`
	ThumbHeight       int32  `json:"thumb_height,omitempty"`
	ThumbSize         int32  `json:"thumb_size,omitempty"`
	DisplayX          int32  `json:"display_x,omitempty"`
	DisplayY          int32  `json:"display_y,omitempty"`
	Id                string `json:"id,omitempty"`
}

type CaptNormalDataResponse struct {
	Code    int32       `json:"code" default:"200"`
	Message string      `json:"message" default:""`
	Data    interface{} `json:"data"`
}

type CaptStatusInfo struct {
	Info   interface{} `json:"info"`
	Status int         `json:"status"`
}

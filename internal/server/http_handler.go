package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wenlng/go-captcha-service/internal/adapt"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/logic"
	"github.com/wenlng/go-captcha-service/internal/middleware"
	"go.uber.org/zap"
)

// HTTPHandlers manages HTTP request handlers
type HTTPHandlers struct {
	config *config.Config
	logger *zap.Logger

	// Initialize logic
	clickCaptLogic  *logic.ClickCaptLogic
	slideCaptLogic  *logic.SlideCaptLogic
	rotateCaptLogic *logic.RotateCaptLogic
	commonLogic     *logic.CommonLogic
}

// NewHTTPHandlers creates a new HTTP handlers instance
func NewHTTPHandlers(svcCtx *common.SvcContext) *HTTPHandlers {
	return &HTTPHandlers{
		config:          svcCtx.Config,
		logger:          svcCtx.Logger,
		clickCaptLogic:  logic.NewClickCaptLogic(svcCtx),
		slideCaptLogic:  logic.NewSlideCaptLogic(svcCtx),
		rotateCaptLogic: logic.NewRotateCaptLogic(svcCtx),
		commonLogic:     logic.NewCommonLogic(svcCtx),
	}
}

// GetDataHandler .
func (h *HTTPHandlers) GetDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptDataResponse{Code: http.StatusOK, Message: ""}

	if r.Method != http.MethodGet {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()

	typeStr := query.Get("type")
	ctype, err := strconv.Atoi(typeStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "missing type parameter")
		return
	}

	themeStr := query.Get("theme")
	var theme int
	if themeStr != "" {
		theme, err = strconv.Atoi(themeStr)
		if err != nil {
			middleware.WriteError(w, http.StatusBadRequest, "missing theme parameter")
			return
		}
	}

	langStr := query.Get("lang")
	var lang int
	if langStr != "" {
		lang, err = strconv.Atoi(langStr)
		if err != nil {
			middleware.WriteError(w, http.StatusBadRequest, "missing lang parameter")
			return
		}
	}

	var data *adapt.CaptData
	switch ctype {
	case common.GoCaptchaTypeClick:
		data, err = h.clickCaptLogic.GetData(r.Context(), ctype, theme, lang)
		break
	case common.GoCaptchaTypeClickShape:
		data, err = h.clickCaptLogic.GetData(r.Context(), ctype, theme, lang)
		break
	case common.GoCaptchaTypeSlide:
		data, err = h.slideCaptLogic.GetData(r.Context(), ctype, theme, lang)
		break
	case common.GoCaptchaTypeDrag:
		data, err = h.slideCaptLogic.GetData(r.Context(), ctype, theme, lang)
		break
	case common.GoCaptchaTypeRotate:
		data, err = h.rotateCaptLogic.GetData(r.Context(), ctype, theme, lang)
		break
	default:
		//...
	}

	if err != nil || data == nil {
		h.logger.Error("failed to get captcha data, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusNotFound, "v")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "success"
	resp.Type = int32(ctype)

	resp.CaptchaKey = data.CaptchaKey
	resp.MasterImageBase64 = data.MasterImageBase64
	resp.ThumbImageBase64 = data.ThumbImageBase64
	resp.MasterImageWidth = data.MasterImageWidth
	resp.MasterImageHeight = data.MasterImageHeight
	resp.ThumbImageWidth = data.ThumbImageWidth
	resp.ThumbImageHeight = data.ThumbImageHeight
	resp.ThumbImageSize = data.ThumbImageSize
	resp.DisplayX = data.DisplayX
	resp.DisplayY = data.DisplayY

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// CheckDataHandler .
func (h *HTTPHandlers) CheckDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: ""}

	if r.Method != http.MethodPost {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Type       int32  `json:"type"`
		CaptchaKey string `json:"captchaKey"`
		Value      string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CaptchaKey == "" || req.Value == "" {
		middleware.WriteError(w, http.StatusBadRequest, "captchaKey and value are required")
		return
	}

	var err error
	var ok bool
	switch req.Type {
	case common.GoCaptchaTypeClick:
		ok, err = h.clickCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case common.GoCaptchaTypeClickShape:
		ok, err = h.clickCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case common.GoCaptchaTypeSlide:
		ok, err = h.slideCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case common.GoCaptchaTypeDrag:
		ok, err = h.slideCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case common.GoCaptchaTypeRotate:
		var angle int64
		angle, err = strconv.ParseInt(req.Value, 10, 64)
		if err == nil {
			ok, err = h.rotateCaptLogic.CheckData(r.Context(), req.CaptchaKey, int(angle))
		}
		break
	default:
		//...
	}

	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "failed to check captcha data")
		return
	}

	if ok {
		resp.Data = "ok"
	} else {
		resp.Data = "failure"
	}
	resp.Code = http.StatusOK

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// CheckStatusHandler .
func (h *HTTPHandlers) CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}

	if r.Method != http.MethodGet {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	captchaKey := query.Get("captchaKey")
	if captchaKey == "" {
		middleware.WriteError(w, http.StatusBadRequest, "captchaKey is required")
		return
	}

	data, err := h.commonLogic.GetStatusInfo(r.Context(), captchaKey)
	if err != nil {
		h.logger.Error("failed to check status, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "failed to check status")
		return
	}

	if data != nil && data.Status == 1 {
		resp.Code = http.StatusOK
		resp.Data = "ok"
	} else {
		resp.Data = "failure"
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// GetStatusInfoHandler .
func (h *HTTPHandlers) GetStatusInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}
	if r.Method != http.MethodGet {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	captchaKey := query.Get("captchaKey")
	if captchaKey == "" {
		middleware.WriteError(w, http.StatusBadRequest, "captchaKey is required")
		return
	}

	data, err := h.commonLogic.GetStatusInfo(r.Context(), captchaKey)
	if err != nil {
		h.logger.Error("failed to get status info, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusNotFound, "not found status info")
		return
	}

	resp.Code = http.StatusOK
	if data != nil {
		resp.Data = data
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// DelStatusInfoHandler .
func (h *HTTPHandlers) DelStatusInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}
	if r.Method != http.MethodDelete {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	captchaKey := query.Get("captchaKey")
	if captchaKey == "" {
		middleware.WriteError(w, http.StatusBadRequest, "captchaKey is required")
		return
	}

	ret, err := h.commonLogic.DelStatusInfo(r.Context(), captchaKey)
	if err != nil {
		h.logger.Error("failed to del status data, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "not found status info")
		return
	}

	if ret {
		resp.Data = "ok"
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

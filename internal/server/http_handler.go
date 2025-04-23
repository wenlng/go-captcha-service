/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wenlng/go-captcha-service/internal/adapt"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/consts"
	"github.com/wenlng/go-captcha-service/internal/helper"
	"github.com/wenlng/go-captcha-service/internal/logic"
	"github.com/wenlng/go-captcha-service/internal/middleware"
	config2 "github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"
	"go.uber.org/zap"
)

const maxUploadSize = int64(10 << 20) // 10MB

// HTTPHandlers manages HTTP request handlers
type HTTPHandlers struct {
	svcCtx     *common.SvcContext
	dynamicCfg *config.DynamicConfig
	logger     *zap.Logger

	// Initialize logic
	clickCaptLogic  *logic.ClickCaptLogic
	slideCaptLogic  *logic.SlideCaptLogic
	rotateCaptLogic *logic.RotateCaptLogic
	commonLogic     *logic.CommonLogic
	resourceLogic   *logic.ResourceLogic
}

// NewHTTPHandlers creates a new HTTP handlers instance
func NewHTTPHandlers(svcCtx *common.SvcContext) *HTTPHandlers {
	return &HTTPHandlers{
		svcCtx:          svcCtx,
		dynamicCfg:      svcCtx.DynamicConfig,
		logger:          svcCtx.Logger,
		clickCaptLogic:  logic.NewClickCaptLogic(svcCtx),
		slideCaptLogic:  logic.NewSlideCaptLogic(svcCtx),
		rotateCaptLogic: logic.NewRotateCaptLogic(svcCtx),
		commonLogic:     logic.NewCommonLogic(svcCtx),
		resourceLogic:   logic.NewResourceLogic(svcCtx),
	}
}

// HealthStatusHandler .
func (h *HTTPHandlers) HealthStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}
	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// GetDataHandler .
func (h *HTTPHandlers) GetDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: ""}

	if r.Method != http.MethodGet {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()

	id := query.Get("id")
	if id == "" {
		middleware.WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	var data *adapt.CaptData
	var err error
	switch h.svcCtx.Captcha.GetCaptTypeWithKey(id) {
	case consts.GoCaptchaTypeClick:
		data, err = h.clickCaptLogic.GetData(r.Context(), id)
		break
	case consts.GoCaptchaTypeClickShape:
		data, err = h.clickCaptLogic.GetData(r.Context(), id)
		break
	case consts.GoCaptchaTypeSlide:
		data, err = h.slideCaptLogic.GetData(r.Context(), id)
		break
	case consts.GoCaptchaTypeDrag:
		data, err = h.slideCaptLogic.GetData(r.Context(), id)
		break
	case consts.GoCaptchaTypeRotate:
		data, err = h.rotateCaptLogic.GetData(r.Context(), id)
		break
	default:
		//...
	}

	if err != nil || data == nil {
		h.logger.Warn("[HttpHandler] Failed to get captcha data, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusNotFound, "captcha type not found")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "success"

	resp.Data = &adapt.CaptData{
		Id:                id,
		CaptchaKey:        data.CaptchaKey,
		MasterImageBase64: data.MasterImageBase64,
		ThumbImageBase64:  data.ThumbImageBase64,
		MasterWidth:       data.MasterWidth,
		MasterHeight:      data.MasterHeight,
		ThumbWidth:        data.ThumbWidth,
		ThumbHeight:       data.ThumbHeight,
		ThumbSize:         data.ThumbSize,
		DisplayX:          data.DisplayX,
		DisplayY:          data.DisplayY,
	}

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
		Id         string `json:"id"`
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

	if req.Id == "" {
		middleware.WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	var ok bool
	var err error

	switch h.svcCtx.Captcha.GetCaptTypeWithKey(req.Id) {
	case consts.GoCaptchaTypeClick:
		ok, err = h.clickCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case consts.GoCaptchaTypeClickShape:
		ok, err = h.clickCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case consts.GoCaptchaTypeSlide:
		ok, err = h.slideCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case consts.GoCaptchaTypeDrag:
		ok, err = h.slideCaptLogic.CheckData(r.Context(), req.CaptchaKey, req.Value)
		break
	case consts.GoCaptchaTypeRotate:
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
		h.logger.Warn("[HttpHandler] Failed to check data, err: ", zap.Error(err))
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
		h.logger.Warn("[HttpHandler] Failed to check status, err: ", zap.Error(err))
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
		h.logger.Warn("[HttpHandler] Failed to get status info, err: ", zap.Error(err))
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
		h.logger.Warn("[HttpHandler] Failed to del status data, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "not found status info")
		return
	}

	if ret {
		resp.Data = "ok"
	} else {
		resp.Data = "no-ops"
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// UploadResourceHandler .
func (h *HTTPHandlers) UploadResourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}

	if r.Method != http.MethodPost {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	dirname := r.FormValue("dirname")
	if dirname == "" {
		middleware.WriteError(w, http.StatusBadRequest, "dirname is required")
		return
	}

	if !helper.IsValidDirName(dirname) {
		middleware.WriteError(w, http.StatusBadRequest, "invalid directory name")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse multipart/form-data
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.logger.Warn("[HttpHandler] Failed to parse form: %v ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "parse form fail")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		middleware.WriteError(w, http.StatusBadRequest, "no files uploaded")
		return
	}

	ret, allDone, err := h.resourceLogic.SaveResource(r.Context(), dirname, files)
	if !ret && err != nil {
		h.logger.Warn("[HttpHandler] Failed to save resource, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "save resource fail")
		return
	}

	if ret {
		resp.Data = "ok"
	}

	if !allDone {
		resp.Data = "some-files-ok"
		resp.Message = "some files failed to be uploaded. check if they already exist"
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// GetResourceListHandler .
func (h *HTTPHandlers) GetResourceListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}
	if r.Method != http.MethodGet {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	resourcePath := query.Get("path")
	if resourcePath == "" {
		middleware.WriteError(w, http.StatusBadRequest, "path is required")
		return
	}

	fileList, err := h.resourceLogic.GetResourceList(r.Context(), resourcePath)
	if err != nil {
		h.logger.Warn("[HttpHandler] Failed to get resource, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "get resource fail")
		return
	}

	if fileList != nil {
		resp.Data = fileList
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// DeleteResourceHandler .
func (h *HTTPHandlers) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}
	if r.Method != http.MethodDelete {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	resourcePath := query.Get("path")
	if resourcePath == "" {
		middleware.WriteError(w, http.StatusBadRequest, "path is required")
		return
	}

	ret, err := h.resourceLogic.DelResource(r.Context(), resourcePath)
	if err != nil {
		h.logger.Warn("[HttpHandler] Failed to delete resource, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "delete resource fail")
		return
	}

	if ret {
		resp.Data = "ok"
	} else {
		resp.Data = "no-ops"
	}

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// GetGoCaptchaConfigHandler .
func (h *HTTPHandlers) GetGoCaptchaConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: "success"}
	if r.Method != http.MethodGet {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	resp.Data = h.svcCtx.Captcha.DynamicCnf.Get()
	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

// UpdateHotGoCaptchaConfigHandler .
func (h *HTTPHandlers) UpdateHotGoCaptchaConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := &adapt.CaptNormalDataResponse{Code: http.StatusOK, Message: ""}

	if r.Method != http.MethodPost {
		middleware.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var conf config2.CaptchaConfig
	if err := json.NewDecoder(r.Body).Decode(&conf); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.svcCtx.Captcha.DynamicCnf.HotUpdate(conf)
	if err != nil {
		h.logger.Warn("[HttpHandler] Failed to hot update config, err: ", zap.Error(err))
		middleware.WriteError(w, http.StatusBadRequest, "hot update config fail")
		return
	}

	resp.Data = "ok"
	resp.Code = http.StatusOK

	json.NewEncoder(w).Encode(helper.Marshal(resp))
}

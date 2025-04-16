package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wenlng/go-captcha-service/internal/adapt"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/logic"
	"github.com/wenlng/go-captcha-service/proto"
	"go.uber.org/zap"
)

// GrpcServer implements the gRPC cache service
type GrpcServer struct {
	proto.UnimplementedGoCaptchaServiceServer
	config *config.Config
	logger *zap.Logger

	// Initialize logic
	clickCaptLogic  *logic.ClickCaptLogic
	slideCaptLogic  *logic.SlideCaptLogic
	rotateCaptLogic *logic.RotateCaptLogic
	commonLogic     *logic.CommonLogic
}

// NewGoCaptchaServer creates a new gRPC cache server
func NewGoCaptchaServer(svcCtx *common.SvcContext) *GrpcServer {
	return &GrpcServer{
		config:          svcCtx.Config,
		logger:          svcCtx.Logger,
		clickCaptLogic:  logic.NewClickCaptLogic(svcCtx),
		slideCaptLogic:  logic.NewSlideCaptLogic(svcCtx),
		rotateCaptLogic: logic.NewRotateCaptLogic(svcCtx),
		commonLogic:     logic.NewCommonLogic(svcCtx),
	}
}

// GetData handle
func (s *GrpcServer) GetData(ctx context.Context, req *proto.GetDataRequest) (*proto.GetDataResponse, error) {
	resp := &proto.GetDataResponse{Code: 0}
	var err error

	var data = &adapt.CaptData{}

	ctype := int(req.GetType())
	theme := int(req.GetTheme())
	lang := int(req.GetLang())

	switch req.GetType() {
	case proto.GoCaptchaType_GoCaptchaTypeClick:
		data, err = s.clickCaptLogic.GetData(ctx, ctype, theme, lang)
		break
	case proto.GoCaptchaType_GoCaptchaTypeClickShape:
		data, err = s.clickCaptLogic.GetData(ctx, ctype, theme, lang)
		break
	case proto.GoCaptchaType_GoCaptchaTypeSlide:
		data, err = s.slideCaptLogic.GetData(ctx, ctype, theme, lang)
		break
	case proto.GoCaptchaType_GoCaptchaTypeDrag:
		data, err = s.slideCaptLogic.GetData(ctx, ctype, theme, lang)
		break
	case proto.GoCaptchaType_GoCaptchaTypeRotate:
		data, err = s.rotateCaptLogic.GetData(ctx, ctype, theme, lang)
		break
	default:
		//
	}

	if err != nil || data == nil {
		s.logger.Error("failed to get captcha data, err: ", zap.Error(err))
		return &proto.GetDataResponse{Code: 0, Message: "failed to get captcha data"}, nil
	}

	resp.Type = req.GetType()
	return resp, nil
}

// CheckData handle
func (s *GrpcServer) CheckData(ctx context.Context, req *proto.CheckDataRequest) (*proto.CheckDataResponse, error) {
	resp := &proto.CheckDataResponse{Code: 0}

	if req.GetCaptchaKey() == "" || req.GetValue() == "" {
		return &proto.CheckDataResponse{Code: 1, Message: "captchaKey and value are required"}, nil
	}

	var err error
	var ok bool
	switch req.GetType() {
	case proto.GoCaptchaType_GoCaptchaTypeClick:
		ok, err = s.clickCaptLogic.CheckData(ctx, req.GetCaptchaKey(), req.GetValue())
		break
	case proto.GoCaptchaType_GoCaptchaTypeClickShape:
		ok, err = s.clickCaptLogic.CheckData(ctx, req.GetCaptchaKey(), req.GetValue())
		break
	case proto.GoCaptchaType_GoCaptchaTypeSlide:
		ok, err = s.slideCaptLogic.CheckData(ctx, req.GetCaptchaKey(), req.GetValue())
		break
	case proto.GoCaptchaType_GoCaptchaTypeDrag:
		ok, err = s.slideCaptLogic.CheckData(ctx, req.GetCaptchaKey(), req.GetValue())
		break
	case proto.GoCaptchaType_GoCaptchaTypeRotate:
		var angle int64
		angle, err = strconv.ParseInt(req.GetValue(), 10, 64)
		if err == nil {
			ok, err = s.rotateCaptLogic.CheckData(ctx, req.GetCaptchaKey(), int(angle))
		}
		break
	default:
		//...
	}

	if err != nil {
		s.logger.Error("failed to check captcha data, err: ", zap.Error(err))
		return &proto.CheckDataResponse{Code: 1, Message: "failed to check captcha data"}, nil
	}

	if ok {
		resp.Data = "ok"
	} else {
		resp.Data = "failure"
	}

	return resp, nil
}

// CheckStatus handle
func (s *GrpcServer) CheckStatus(ctx context.Context, req *proto.StatusInfoRequest) (*proto.StatusInfoResponse, error) {
	resp := &proto.StatusInfoResponse{Code: 0}

	if req.GetCaptchaKey() == "" {
		return &proto.StatusInfoResponse{Code: 1, Message: "captchaKey is required"}, nil
	}

	data, err := s.commonLogic.GetStatusInfo(ctx, req.GetCaptchaKey())
	if err != nil {
		s.logger.Error("failed to check status, err: ", zap.Error(err))
		return &proto.StatusInfoResponse{Code: 1}, nil
	}

	if data != nil && data.Status == 1 {
		resp.Data = "ok"
	} else {
		resp.Data = "failure"
	}

	return resp, nil
}

// GetStatusInfo handle
func (s *GrpcServer) GetStatusInfo(ctx context.Context, req *proto.StatusInfoRequest) (*proto.StatusInfoResponse, error) {
	resp := &proto.StatusInfoResponse{Code: 0}

	if req.CaptchaKey == "" {
		return &proto.StatusInfoResponse{Code: 1, Message: "captchaKey is required"}, nil
	}

	data, err := s.commonLogic.GetStatusInfo(ctx, req.GetCaptchaKey())
	if err != nil {
		s.logger.Error("failed to check status, err: ", zap.Error(err))
		return &proto.StatusInfoResponse{Code: 1}, nil
	}

	if data != nil && data.Status == 1 {
		dataByte, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to json marshal: %v", err)
		}

		resp.Data = string(dataByte)
	}

	return resp, nil
}

// DelStatusInfo handle
func (s *GrpcServer) DelStatusInfo(ctx context.Context, req *proto.StatusInfoRequest) (*proto.StatusInfoResponse, error) {
	resp := &proto.StatusInfoResponse{Code: 0}

	if req.CaptchaKey == "" {
		return &proto.StatusInfoResponse{Code: 1, Message: "captchaKey is required"}, nil
	}

	ret, err := s.commonLogic.DelStatusInfo(ctx, req.GetCaptchaKey())
	if err != nil {
		s.logger.Error("failed to delete status info, err: ", zap.Error(err))
		return &proto.StatusInfoResponse{Code: 1}, nil
	}

	if ret {
		resp.Data = "ok"
	}

	return resp, nil
}

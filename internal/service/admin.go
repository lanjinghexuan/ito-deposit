package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ito-deposit/internal/biz"
	data2 "ito-deposit/internal/data"
	"time"

	pb "ito-deposit/api/helloworld/v1"
)

type AdminService struct {
	pb.UnimplementedAdminServer
	data *data2.Data
	biz  *biz.AdminUsecase
}

func NewAdminService(data *data2.Data, bizdataa *biz.AdminUsecase) *AdminService {
	return &AdminService{
		data: data,
		biz:  bizdataa,
	}
}

func (s *AdminService) SetPriceRule(ctx context.Context, req *pb.SetPriceRuleReq) (*pb.SetPriceRuleRes, error) {

	if len(req.Rules) == 0 {
		return nil, status.Error(codes.InvalidArgument, "至少需要一条价格规则")
	}

	var newRule []*biz.LockerPricingRules
	// 3. 创建新规则
	for _, rule := range req.Rules {

		newRule = append(newRule, &biz.LockerPricingRules{
			NetworkId:        req.NetworkId,
			RuleName:         rule.RuleName,
			FeeType:          convertToFeeType(rule.FeeType),
			LockerType:       convertToLockerType(rule.LockerType),
			FreeDuration:     convertToDecimal(rule.FreeDuration),
			IsDepositEnabled: boolToInt(rule.IsDepositEnabled),
			IsAdvancePay:     boolToInt(rule.IsAdvancePay),
			HourlyRate:       convertToDecimal(rule.HourlyRate),
			DailyCap:         convertToDecimal(rule.DailyCap),
			DailyRate:        convertToDecimal(rule.DailyRate),
			AdvanceAmount:    convertToDecimal(rule.AdvanceAmount),
			DepositAmount:    convertToDecimal(rule.DepositAmount),
			Status:           1,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		})
	}

	err := s.biz.SetPriceRule(ctx, int32(req.NetworkId), newRule)
	if err != nil {
		return nil, err
	}

	return &pb.SetPriceRuleRes{
		Code: 200,
		Msg:  "规则更新成功",
	}, nil
}

func (s *AdminService) GetPriceRule(ctx context.Context, req *pb.GetPriceRuleReq) (*pb.GetPriceRuleRes, error) {
	var rules []*data2.LockerPricingRules

	// 只查询生效状态的规则
	err := s.data.DB.Where("network_id = ? AND status = 1", req.NetworkId).Find(&rules).Error
	if err != nil {
		return nil, err
	}

	pbRules := make([]*pb.LockerPriceRule, 0, len(rules))
	for _, r := range rules {
		pbRules = append(pbRules, &pb.LockerPriceRule{
			Id:               r.Id,
			RuleName:         r.RuleName,
			FeeType:          int32(r.FeeType),
			LockerType:       int32(r.LockerType),
			FreeDuration:     float32(r.FreeDuration),
			HourlyRate:       float32(r.HourlyRate),
			DailyCap:         float32(r.DailyCap),
			DailyRate:        float32(r.DailyRate),
			AdvanceAmount:    float32(r.AdvanceAmount),
			DepositAmount:    float32(r.DepositAmount),
			IsDepositEnabled: intToBool(int(r.IsDepositEnabled)),
			IsAdvancePay:     intToBool(int(r.IsAdvancePay)),
		})
	}

	return &pb.GetPriceRuleRes{Rules: pbRules}, nil
}

// 辅助函数

func convertToDecimal(value float32) float64 {
	// 直接转为float64，保留原始精度
	return float64(value)
}

func boolToInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func convertToFeeType(feeType int32) int8 {
	// 添加默认值逻辑
	if feeType == 0 {
		return 1 // 默认计时收费
	}
	return int8(feeType)
}

func convertToLockerType(lockerType int32) int8 {
	return int8(lockerType)
}

func validatePriceRule(rule *pb.LockerPriceRule) error {
	if rule.LockerType < 1 || rule.LockerType > 3 {
		return status.Error(codes.InvalidArgument, "柜型必须是1(小柜)或2(大柜)")
	}
	if rule.FeeType < 1 || rule.FeeType > 2 {
		return status.Error(codes.InvalidArgument, "收费类型必须是1(计时)或2(按日)")
	}
	if rule.FreeDuration < 0 {
		return status.Error(codes.InvalidArgument, "免费时长不能为负数")
	}
	return nil
}

func intToBool(i int) bool {
	return i != 0
}

package service

import (
	"context"
	"fmt"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	minio1 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ito-deposit/internal/biz"
	"ito-deposit/internal/conf"
	data2 "ito-deposit/internal/data"
	"mime/multipart"
	"time"

	pb "ito-deposit/api/helloworld/v1"
)

type AdminService struct {
	pb.UnimplementedAdminServer
	data *data2.Data
	biz  *biz.AdminUsecase
	conf *conf.Data
}

func NewAdminService(data *data2.Data, bizdataa *biz.AdminUsecase, conf *conf.Data) *AdminService {
	return &AdminService{
		data: data,
		biz:  bizdataa,
		conf: conf,
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

func (s *AdminService) DownloadFile(ctx kratosHttp.Context) error {
	req := ctx.Request()

	// 获取上传文件
	_, file, err := req.FormFile("file")
	if err != nil {
		return err
	}

	url, err := UploadFile(file.Filename, file, s.conf)
	if err != nil {
		return err
	}

	return ctx.Result(200, map[string]string{"url": url})
}

func UploadFile(objectName string, fileHeader *multipart.FileHeader, c *conf.Data) (string, error) {
	// 初始化 MinIO 客户端
	minioClient, err := minio1.New(c.Minio.Endpoint, &minio1.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyId, c.Minio.AccessKeySecret, ""),
		Secure: c.Minio.UseSsl,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	// ✅ 只打开一次
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// ✅ 构造最终对象名
	objectName = fmt.Sprintf("%s/%s", time.Now().Format("2006-01-02"), objectName)

	// ✅ 上传
	_, err = minioClient.PutObject(
		context.Background(),
		c.Minio.BucketName,
		objectName,
		src,
		fileHeader.Size,
		minio1.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	// ✅ 返回可访问地址
	return fmt.Sprintf("%s/%s/%s", c.Minio.Endpoint, c.Minio.BucketName, objectName), nil
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

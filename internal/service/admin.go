package service

import (
	"context"
	"fmt"
	jwt1 "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ito-deposit/internal/biz"
	"ito-deposit/internal/conf"
	data2 "ito-deposit/internal/data"
	"mime/multipart"
	"path/filepath"
	"time"

	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	minio1 "github.com/minio/minio-go/v7"
	pb "ito-deposit/api/helloworld/v1"
)

type AdminService struct {
	pb.UnimplementedAdminServer
	data   *data2.Data
	conf   *conf.Data
	server *conf.Server
	repo   *biz.AdminUsecase
}

func NewAdminService(data *data2.Data, conf *conf.Data, server *conf.Server, repo *biz.AdminUsecase) *AdminService {
	return &AdminService{
		data:   data,
		conf:   conf,
		server: server,
		repo:   repo,
	}
}
func (s *AdminService) AdminLogin(ctx context.Context, req *pb.AdminLoginReq) (*pb.AdminLoginRes, error) {
	get := s.data.Redis.Get(context.Background(), "sendSms"+req.Mobile+"admin_login")
	if get.Val() != req.SmsCode {
		return &pb.AdminLoginRes{
			Code: 500,
			Msg:  "验证码错误",
		}, nil
	}
	var admin data2.Admin
	err := s.data.DB.Debug().Where("mobile = ?", req.Mobile).Find(&admin).Error
	if err != nil {
		return &pb.AdminLoginRes{
			Code: 500,
			Msg:  "查询失败",
		}, nil
	}
	if admin.Id == 0 {
		return &pb.AdminLoginRes{
			Code: 500,
			Msg:  "用户不存在",
		}, nil
	}
	if req.Password != admin.Password {
		return &pb.AdminLoginRes{
			Code: 500,
			Msg:  "密码错误",
		}, nil
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 根据您的需求设置 JWT 中的声明
		"your_custom_claim": "your_custom_value",
		"id":                "123",
	})

	signedString, err := claims.SignedString([]byte(s.server.Jwt.Authkey))
	if err != nil {
		return nil, err
	}
	return &pb.AdminLoginRes{
		Code:  200,
		Msg:   "登陆成功",
		Id:    int64(admin.Id),
		Token: signedString,
	}, nil
}
func (s *AdminService) SetPriceRule(ctx context.Context, req *pb.SetPriceRuleReq) (*pb.SetPriceRuleRes, error) {
	// 0. 参数校验
	if req.NetworkId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "网点ID必须大于0")
	}
	if len(req.Rules) == 0 {
		return nil, status.Error(codes.InvalidArgument, "至少需要一条价格规则")
	}

	// 1. 开启事务
	tx := s.data.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 2. 停用旧规则（软删除）
	if err := tx.Model(&data2.LockerPricingRules{}).
		Where("network_id = ? AND status = 1", req.NetworkId).
		Update("status", 0).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "停用旧规则失败: %v", err)
	}

	// 3. 创建新规则
	for _, rule := range req.Rules {
		// 规则校验
		if err := validatePriceRule(rule); err != nil {
			tx.Rollback()
			return nil, err
		}

		newRule := &data2.LockerPricingRules{
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
		}

		if err := tx.Create(newRule).Error; err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建规则失败: %v", err)
		}
	}

	// 4. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
	}

	// 5. 清理缓存（如有）
	//go s.clearPriceRuleCache(req.NetworkId)

	return &pb.SetPriceRuleRes{
		Code: 200,
		Msg:  "规则更新成功",
	}, nil
}

func (s *AdminService) GetPriceRule(ctx context.Context, req *pb.GetPriceRuleReq) (*pb.GetPriceRuleRes, error) {
	var rules []*data2.LockerPricingRules
	fmt.Println(1)
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
func (s *AdminService) DownloadFile(ctx kratosHttp.Context) error {
	req := ctx.Request()

	// 获取上传文件
	_, file, err := req.FormFile("file")
	if err != nil {
		return err
	}
	// 文件格式验证
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" {
		return ctx.Result(500, map[string]interface{}{
			"code":    500,
			"message": "格式只能是png或者jpg",
			"data":    nil,
		})
	}

	// 文件大小验证（200KB限制）
	if file.Size >= 200*1024*1024 {
		return ctx.Result(500, map[string]interface{}{
			"code":    500,
			"message": "大小不能超过200MB",
			"data":    nil,
		})
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
func (s *AdminService) PointList(ctx context.Context, req *pb.PointListReq) (*pb.PointListRes, error) {
	var point []data2.LockerPoint
	err := s.data.DB.Debug().Find(&point).Error
	if err != nil {
		return &pb.PointListRes{
			Code: 500,
			Msg:  "查询失败",
		}, nil
	}
	var lists []*pb.PointList
	for _, l := range point {
		list := pb.PointList{
			Name:            l.Name,
			Address:         l.Address,
			AvailableLarge:  int64(l.AvailableLarge),
			AvailableMedium: int64(l.AvailableMedium),
			AvailableSmall:  int64(l.AvailableSmall),
		}
		lists = append(lists, &list)
	}
	return &pb.PointListRes{
		Code: 200,
		Msg:  "查询成功",
		List: lists,
	}, nil
}

func (s *AdminService) PointInfo(ctx context.Context, req *pb.PointInfoReq) (*pb.PointInfoRes, error) {
	var point data2.LockerPoint
	fmt.Println(req.Id)
	err := s.data.DB.Debug().Where("id = ?", req.Id).Find(&point).Error
	if err != nil {
		return &pb.PointInfoRes{
			Code: 500,
			Msg:  "查询失败",
		}, nil
	}
	return &pb.PointInfoRes{
		Code:            200,
		Msg:             "查询成功",
		Name:            point.Name,
		Address:         point.Address,
		PointType:       point.PointType,
		AvailableLarge:  int64(point.AvailableLarge),
		AvailableMedium: int64(point.AvailableMedium),
		AvailableSmall:  int64(point.AvailableSmall),
		OpenTime:        point.OpenTime,
		Staus:           point.Status,
		PointImage:      point.PointImage,
	}, nil
}

func (s *AdminService) AddPoint(ctx context.Context, req *pb.AddPointReq) (*pb.AddPointRes, error) {
	kratosToken, ok := jwt1.FromContext(ctx)
	if !ok {
		return &pb.AddPointRes{
			Code: 401,
			Msg:  "token不正确或者未传",
		}, nil
	}

	mapClaims, ok := kratosToken.(*jwt.MapClaims)

	userId := (*mapClaims)["id"].(string)
	return s.repo.AddPointAddPoint(ctx, req, userId)
}

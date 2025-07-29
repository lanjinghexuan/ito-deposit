package geo

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"ito-deposit/internal/conf"
	"strconv"
)

// GeoService 地理位置服务
type GeoService struct {
	redis *redis.Client
}

// NewGeoService 创建地理位置服务实例
func NewGeoService(c *conf.Data) *GeoService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	return &GeoService{
		redis: rdb,
	}
}

// SaveLockerPoint 保存寄存点位置信息到Redis GEO
func (s *GeoService) SaveLockerPoint(ctx context.Context, id int32, name string, longitude, latitude float64) error {
	// 使用Redis GEO命令保存位置信息
	// GEOADD key longitude latitude member
	key := "locker_points_geo"

	// 将寄存点ID转换为字符串作为成员名
	member := fmt.Sprintf("locker:%d", id)

	// 执行GEOADD命令
	_, err := s.redis.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      member,
		Longitude: longitude,
		Latitude:  latitude,
	}).Result()

	if err != nil {
		return fmt.Errorf("保存寄存点位置信息失败: %w", err)
	}

	// 同时保存寄存点名称，方便后续查询
	infoKey := fmt.Sprintf("locker_info:%d", id)
	err = s.redis.HSet(ctx, infoKey, map[string]interface{}{
		"name":      name,
		"longitude": longitude,
		"latitude":  latitude,
	}).Err()

	if err != nil {
		return fmt.Errorf("保存寄存点信息失败: %w", err)
	}

	return nil
}

// FindNearbyLockerPoints 查找附近的寄存点
func (s *GeoService) FindNearbyLockerPoints(ctx context.Context, longitude, latitude float64, radius float64, unit string, limit int64) ([]LockerPointLocation, error) {
	// 使用Redis GEORADIUS命令查询附近的寄存点
	key := "locker_points_geo"

	// 执行GEORADIUS命令
	// GEORADIUS key longitude latitude radius m|km|ft|mi WITHDIST WITHCOORD COUNT count
	res, err := s.redis.GeoRadius(ctx, key, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        unit,       // 直接使用字符串，如 "km", "m", "ft", "mi"
		WithCoord:   true,       // 返回坐标
		WithDist:    true,       // 返回距离
		WithGeoHash: false,      // 不返回geohash
		Count:       int(limit), // 限制返回数量，转换为int类型
		Sort:        "ASC",      // 按距离升序排序
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("查询附近寄存点失败: %w", err)
	}

	// 转换结果
	result := make([]LockerPointLocation, 0, len(res))
	for _, loc := range res {
		// 从成员名中提取寄存点ID
		var id int32
		_, err := fmt.Sscanf(loc.Name, "locker:%d", &id)
		if err != nil {
			continue
		}

		// 获取寄存点信息
		infoKey := fmt.Sprintf("locker_info:%d", id)
		info, err := s.redis.HGetAll(ctx, infoKey).Result()
		if err != nil {
			continue
		}

		name := info["name"]

		// 创建寄存点位置对象
		lockerLoc := LockerPointLocation{
			ID:        id,
			Name:      name,
			Distance:  loc.Dist,
			Longitude: loc.Longitude,
			Latitude:  loc.Latitude,
		}

		result = append(result, lockerLoc)
	}

	return result, nil
}

// SaveAllLockerPoints 批量保存所有寄存点位置信息
func (s *GeoService) SaveAllLockerPoints(ctx context.Context, lockers []LockerPointInfo) error {
	// 如果没有寄存点，直接返回
	if len(lockers) == 0 {
		return nil
	}

	// 使用管道批量执行命令，提高性能
	pipe := s.redis.Pipeline()

	// 先清除旧数据
	pipe.Del(ctx, "locker_points_geo")

	// 准备批量添加的位置数据
	locations := make([]*redis.GeoLocation, 0, len(lockers))
	for _, locker := range lockers {
		member := fmt.Sprintf("locker:%d", locker.ID)
		locations = append(locations, &redis.GeoLocation{
			Name:      member,
			Longitude: locker.Longitude,
			Latitude:  locker.Latitude,
		})

		// 同时保存寄存点信息
		infoKey := fmt.Sprintf("locker_info:%d", locker.ID)
		pipe.HSet(ctx, infoKey, map[string]interface{}{
			"name":      locker.Name,
			"longitude": locker.Longitude,
			"latitude":  locker.Latitude,
			"address":   locker.Address,
		})
	}

	// 批量添加位置数据
	pipe.GeoAdd(ctx, "locker_points_geo", locations...)

	// 执行管道命令
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("批量保存寄存点位置信息失败: %w", err)
	}

	return nil
}

// GetLockerPointInfo 获取寄存点信息
func (s *GeoService) GetLockerPointInfo(ctx context.Context, id int32) (*LockerPointInfo, error) {
	infoKey := fmt.Sprintf("locker_info:%d", id)
	info, err := s.redis.HGetAll(ctx, infoKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取寄存点信息失败: %w", err)
	}

	// 如果没有找到信息，返回nil
	if len(info) == 0 {
		return nil, nil
	}

	// 解析经纬度
	longitude, _ := strconv.ParseFloat(info["longitude"], 64)
	latitude, _ := strconv.ParseFloat(info["latitude"], 64)

	return &LockerPointInfo{
		ID:        id,
		Name:      info["name"],
		Address:   info["address"],
		Longitude: longitude,
		Latitude:  latitude,
	}, nil
}

// SaveUserLocation 保存用户实时位置信息到Redis
func (s *GeoService) SaveUserLocation(ctx context.Context, userID string, longitude, latitude float64) error {
	// 使用Redis GEO命令保存用户位置信息
	key := "user_locations_geo"

	// 执行GEOADD命令
	_, err := s.redis.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      userID,
		Longitude: longitude,
		Latitude:  latitude,
	}).Result()

	if err != nil {
		return fmt.Errorf("保存用户位置信息失败: %w", err)
	}

	return nil
}

// GetUserLocation 获取用户位置信息
func (s *GeoService) GetUserLocation(ctx context.Context, userID string) (float64, float64, error) {
	// 使用Redis GEOPOS命令获取用户位置信息
	key := "user_locations_geo"

	// 执行GEOPOS命令
	res, err := s.redis.GeoPos(ctx, key, userID).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("获取用户位置信息失败: %w", err)
	}

	// 检查结果
	if len(res) == 0 || res[0] == nil {
		return 0, 0, fmt.Errorf("用户位置信息不存在")
	}

	return res[0].Longitude, res[0].Latitude, nil
}

// LockerPointLocation 寄存点位置信息
type LockerPointLocation struct {
	ID        int32   `json:"id"`        // 寄存点ID
	Name      string  `json:"name"`      // 寄存点名称
	Distance  float64 `json:"distance"`  // 距离（米或公里，取决于查询时使用的单位）
	Longitude float64 `json:"longitude"` // 经度
	Latitude  float64 `json:"latitude"`  // 纬度
}

// LockerPointInfo 寄存点基本信息
type LockerPointInfo struct {
	ID        int32   `json:"id"`        // 寄存点ID
	Name      string  `json:"name"`      // 寄存点名称
	Address   string  `json:"address"`   // 地址
	Longitude float64 `json:"longitude"` // 经度
	Latitude  float64 `json:"latitude"`  // 纬度
}

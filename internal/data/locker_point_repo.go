package data

import (
	"context"
)

type lockerPointRepo struct {
	data *Data
}

func NewlockerPointRepo(data *Data) *lockerPointRepo {
	return &lockerPointRepo{
		data: data,
	}
}

// GetTypes 小表全量缓存
func (r *lockerPointRepo) GetTypes(ctx context.Context) ([]*LockerType, error) {

	var types []*LockerType
	if err := r.data.DB.WithContext(ctx).Find(&types).Error; err != nil {
		return nil, err
	}

	return types, nil
}

// CountAvailableByType 只扫 locker 单表
func (r *lockerPointRepo) CountAvailableByType(ctx context.Context, pointID int64) (map[int32]int32, error) {
	type result struct {
		TypeID int32
		Num    int32
	}
	var list []result
	err := r.data.DB.WithContext(ctx).
		Model(&Locker{}).
		Select("type_id AS type_id, COUNT(*) AS num").
		Where("locker_point_id = ? AND status = 1", pointID).
		Group("type_id").
		Scan(&list).Error
	m := make(map[int32]int32)
	for _, v := range list {
		m[v.TypeID] = v.Num
	}
	return m, err
}

// GetPricingRule 单条索引查询
func (r *lockerPointRepo) GetPricingRule(ctx context.Context, pointID, typeID int32) (*LockerPricingRules, error) {
	var rule LockerPricingRules
	err := r.data.DB.WithContext(ctx).
		Where("network_id = ? AND locker_type = ? AND status = 1", pointID, typeID).
		Order("effective_time DESC").
		First(&rule).Error
	return &rule, err
}

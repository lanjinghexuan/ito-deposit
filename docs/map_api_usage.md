# 城市寄存点分布图API使用说明

## API接口

### 获取城市寄存点分布图数据

**接口地址：** `GET /api/nearby/city/map`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| city_name | string | 是 | 城市名称，如"北京" |
| north_lat | double | 是 | 北纬度（地图边界） |
| south_lat | double | 是 | 南纬度（地图边界） |
| east_lng | double | 是 | 东经度（地图边界） |
| west_lng | double | 是 | 西经度（地图边界） |
| zoom_level | int32 | 是 | 地图缩放级别（1-20） |
| enable_cluster | bool | 否 | 是否启用聚合（默认true） |

**请求示例：**
```
GET /api/nearby/city/map?city_name=北京&north_lat=40.0&south_lat=39.8&east_lng=116.5&west_lng=116.3&zoom_level=12&enable_cluster=true
```

**响应格式：**

```json
{
  "points": [
    {
      "id": 1,
      "name": "北京西站寄存点",
      "address": "北京市西城区莲花池东路118号",
      "longitude": 116.321614,
      "latitude": 39.894217,
      "total_available": 15,
      "status": "available"
    }
  ],
  "clusters": [
    {
      "longitude": 116.4,
      "latitude": 39.9,
      "count": 5,
      "total_available": 50,
      "point_ids": [1, 2, 3, 4, 5]
    }
  ],
  "total_count": 100,
  "zoom_level": 12,
  "is_clustered": true
}
```

## 聚合逻辑说明

### 聚合触发条件

系统会根据以下条件决定是否进行聚合：

1. **缩放级别 <= 10**：点位数量 > 10 时聚合
2. **缩放级别 11-15**：点位数量 > 50 时聚合
3. **缩放级别 > 15**：点位数量 > 100 时聚合

### 聚合距离

根据缩放级别确定聚合距离：

- **缩放级别 <= 5**：0.1度（约11公里）
- **缩放级别 6-10**：0.05度（约5.5公里）
- **缩放级别 11-15**：0.01度（约1.1公里）
- **缩放级别 > 15**：0.005度（约550米）

### 状态说明

寄存点状态根据可用柜数量确定：

- **available**：可用柜数量 > 5
- **busy**：可用柜数量 1-5
- **full**：可用柜数量 = 0

## 错误处理

### 常见错误码

- **400 Bad Request**：参数验证失败
  - 城市名称为空
  - 纬度/经度边界不正确
  - 缩放级别不在1-20范围内

- **500 Internal Server Error**：服务器内部错误
  - 数据库查询失败
  - 城市不存在

### 错误响应示例

```json
{
  "code": 400,
  "message": "城市名称不能为空",
  "reason": "NEARBY_BAD_REQUEST"
}
```

## 前端集成示例

### JavaScript示例

```javascript
// 获取地图数据
async function getMapData(cityName, bounds, zoomLevel) {
  const params = new URLSearchParams({
    city_name: cityName,
    north_lat: bounds.north,
    south_lat: bounds.south,
    east_lng: bounds.east,
    west_lng: bounds.west,
    zoom_level: zoomLevel,
    enable_cluster: true
  });
  
  try {
    const response = await fetch(`/api/nearby/city/map?${params}`);
    const data = await response.json();
    
    if (data.is_clustered) {
      // 渲染聚合点
      renderClusters(data.clusters);
    } else {
      // 渲染详细点位
      renderPoints(data.points);
    }
  } catch (error) {
    console.error('获取地图数据失败:', error);
  }
}

// 渲染聚合点
function renderClusters(clusters) {
  clusters.forEach(cluster => {
    // 在地图上添加聚合标记
    addClusterMarker({
      position: [cluster.longitude, cluster.latitude],
      count: cluster.count,
      totalAvailable: cluster.total_available
    });
  });
}

// 渲染详细点位
function renderPoints(points) {
  points.forEach(point => {
    // 在地图上添加寄存点标记
    addPointMarker({
      id: point.id,
      position: [point.longitude, point.latitude],
      name: point.name,
      address: point.address,
      status: point.status,
      totalAvailable: point.total_available
    });
  });
}
```

## 性能优化建议

1. **合理设置边界范围**：避免查询过大的地理范围
2. **启用聚合功能**：在低缩放级别时使用聚合减少数据量
3. **缓存策略**：对相同参数的请求进行缓存
4. **分页加载**：对于大量数据可以考虑分页加载
5. **数据库索引**：确保经纬度字段有适当的索引

## 数据库索引建议

```sql
-- 为寄存点表添加地理位置索引
CREATE INDEX idx_locker_point_location ON locker_point(latitude, longitude);

-- 为城市关联添加索引
CREATE INDEX idx_locker_point_location_id ON locker_point(location_id);

-- 复合索引优化边界查询
CREATE INDEX idx_locker_point_bounds ON locker_point(location_id, latitude, longitude);
```
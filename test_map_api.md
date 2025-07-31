# 测试地图API

## 1. 确保服务正常启动

首先确保你的服务已经正常启动，没有编译错误。

## 2. 测试API

### 使用curl测试

```bash
curl -X GET "http://localhost:8000/api/nearby/city/map?city_name=北京&north_lat=40.0&south_lat=39.8&east_lng=116.5&west_lng=116.3&zoom_level=12&enable_cluster=true"
```

### 使用APIPost测试

**请求方法**: GET
**URL**: `http://localhost:8000/api/nearby/city/map`

**查询参数**:
- city_name: 北京
- north_lat: 40.0
- south_lat: 39.8
- east_lng: 116.5
- west_lng: 116.3
- zoom_level: 12
- enable_cluster: true

## 3. 检查服务日志

如果API返回错误，请检查服务日志，看看是否有错误信息。

## 4. 常见问题排查

### 问题1: "未找到api"
- 确保URL路径正确：`/api/nearby/city/map`
- 确保服务已经重新启动
- 检查服务端口是否正确（默认8000）

### 问题2: 参数验证错误
- 确保所有必填参数都已提供
- 确保参数类型正确（数字类型不要加引号）
- 确保纬度和经度范围合理

### 问题3: 城市不存在
- 确保数据库中有对应的城市数据
- 城市名称要完全匹配数据库中的记录

## 5. 预期响应格式

成功响应示例：
```json
{
  "points": [
    {
      "id": 1,
      "name": "测试寄存点",
      "address": "测试地址",
      "longitude": 116.4,
      "latitude": 39.9,
      "total_available": 10,
      "status": "available"
    }
  ],
  "clusters": [],
  "total_count": 1,
  "zoom_level": 12,
  "is_clustered": false
}
```

错误响应示例：
```json
{
  "code": 400,
  "message": "城市名称不能为空",
  "reason": "NEARBY_BAD_REQUEST"
}
```
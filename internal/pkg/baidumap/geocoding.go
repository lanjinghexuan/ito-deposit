package baidumap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// BaiduMapClient 百度地图API客户端
type BaiduMapClient struct {
	AK string // 百度地图API密钥
}

// NewBaiduMapClient 创建一个新的百度地图API客户端
func NewBaiduMapClient(ak string) *BaiduMapClient {
	return &BaiduMapClient{
		AK: ak,
	}
}

// 使用更灵活的结构来处理百度地图API的响应
type GeocodeResponse struct {
	Status  int    `json:"status"`  // 状态码，0表示成功
	Message string `json:"message"` // 错误信息

	// 使用map来存储结果，避免结构体字段不匹配的问题
	RawResult map[string]interface{} `json:"-"`

	// 提取后的关键信息
	Longitude float64 // 经度
	Latitude  float64 // 纬度
	AdCode    string  // 行政区划代码（6位数字）
	CityCode  string  // 城市编码（字母缩写）
}

// Geocode 根据城市名称获取地理编码信息
func (c *BaiduMapClient) Geocode(cityName string) (*GeocodeResponse, error) {
	// 构建API请求URL
	apiURL := "http://api.map.baidu.com/geocoding/v3/"

	// 设置请求参数
	params := url.Values{}
	params.Set("address", cityName)
	params.Set("output", "json")
	params.Set("ak", c.AK)
	params.Set("ret_coordtype", "wgs84ll") // 使用WGS84坐标系，与前端地图一致
	params.Set("extensions_town", "true")  // 返回乡镇信息
	params.Set("extensions_poi", "0")      // 不返回POI信息

	// 发送HTTP GET请求
	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析为通用的map结构
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 创建结果对象
	result := &GeocodeResponse{
		RawResult: rawResponse,
	}

	// 检查API响应状态
	if status, ok := rawResponse["status"].(float64); ok {
		result.Status = int(status)
	}
	if message, ok := rawResponse["message"].(string); ok {
		result.Message = message
	}

	// 如果状态不为0，表示API调用失败
	if result.Status != 0 {
		return nil, fmt.Errorf("百度地图API错误: %s", result.Message)
	}

	// 从原始响应中提取关键信息
	if resultData, ok := rawResponse["result"].(map[string]interface{}); ok {
		// 提取经纬度
		if location, ok := resultData["location"].(map[string]interface{}); ok {
			if lng, ok := location["lng"].(float64); ok {
				result.Longitude = lng
			}
			if lat, ok := location["lat"].(float64); ok {
				result.Latitude = lat
			}
		}

		// 提取城市编码
		// 1. 尝试从addressComponent中提取
		if addressComponent, ok := resultData["addressComponent"].(map[string]interface{}); ok {
			// 尝试获取adcode（6位数字行政区划编码）
			if adcode, ok := addressComponent["adcode"].(string); ok && adcode != "" {
				result.AdCode = adcode
			}

			// 尝试获取citycode（字母缩写）
			if citycode, ok := addressComponent["citycode"].(string); ok && citycode != "" {
				result.CityCode = citycode
			}
		}

		// 2. 尝试从顶层result中提取
		if adcode, ok := resultData["adcode"].(string); ok && adcode != "" && result.AdCode == "" {
			result.AdCode = adcode
		}

		if citycode, ok := resultData["cityCode"].(string); ok && citycode != "" && result.CityCode == "" {
			result.CityCode = citycode
		}
	}

	// 如果没有提取到任何编码，但API调用成功，则使用城市名称作为编码
	if result.AdCode == "" && result.CityCode == "" {
		// 这里我们直接使用城市名称作为编码，确保不返回空
		result.CityCode = cityName
	}

	return result, nil
}

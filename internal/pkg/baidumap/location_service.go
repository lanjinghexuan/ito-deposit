package baidumap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// LocationResponse 百度地图IP定位API响应
type LocationResponse struct {
	Status  int    `json:"status"`  // 状态码，0表示成功
	Message string `json:"message"` // 错误信息
	
	// 位置信息
	Content struct {
		Address       string `json:"address"`        // 详细地址
		AddressDetail struct {
			City     string `json:"city"`     // 城市
			CityCode int    `json:"city_code"` // 城市编码
			District string `json:"district"` // 区县
			Province string `json:"province"` // 省份
			Street   string `json:"street"`   // 街道
		} `json:"address_detail"`
		Point struct {
			X float64 `json:"x"` // 经度
			Y float64 `json:"y"` // 纬度
		} `json:"point"`
	} `json:"content"`
}

// RealtimeLocationResponse 百度地图实时定位API响应
type RealtimeLocationResponse struct {
	Status  int    `json:"status"`  // 状态码，0表示成功
	Message string `json:"message"` // 错误信息
	
	// 位置信息
	Result struct {
		Location struct {
			Lng float64 `json:"lng"` // 经度
			Lat float64 `json:"lat"` // 纬度
		} `json:"location"`
		Accuracy float64 `json:"accuracy"` // 精确度
		Address  string  `json:"address"`  // 地址
		City     string  `json:"city"`     // 城市
		CityCode string  `json:"city_code"` // 城市编码
	} `json:"result"`
}

// GetLocation 根据IP地址获取位置信息
func (c *BaiduMapClient) GetLocation(ip string) (*LocationResponse, error) {
	// 构建API请求URL
	apiURL := "http://api.map.baidu.com/location/ip"

	// 设置请求参数
	params := url.Values{}
	params.Set("ak", c.AK)
	params.Set("coor", "bd09ll") // 百度坐标系
	
	// 如果提供了IP地址，则使用该IP
	if ip != "" {
		params.Set("ip", ip)
	}

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

	// 解析响应
	var result LocationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查API响应状态
	if result.Status != 0 {
		return nil, fmt.Errorf("百度地图API错误: %s", result.Message)
	}

	return &result, nil
}

// GetRealtimeLocation 获取实时位置信息
// 参数说明：
// - cityCode: 城市编码，用于限定定位范围
// - ip: 客户端IP地址，可选
// - coor: 坐标系类型，默认为bd09ll（百度坐标系）
func (c *BaiduMapClient) GetRealtimeLocation(cityCode, ip, coor string) (*RealtimeLocationResponse, error) {
	// 构建API请求URL
	apiURL := "http://api.map.baidu.com/location/api/v2/sdk"

	// 设置请求参数
	params := url.Values{}
	params.Set("ak", c.AK)
	
	// 设置坐标系，默认为百度坐标系
	if coor == "" {
		coor = "bd09ll"
	}
	params.Set("coor", coor)
	
	// 如果提供了城市编码，则使用该城市编码
	if cityCode != "" {
		params.Set("city_code", cityCode)
	}
	
	// 如果提供了IP地址，则使用该IP
	if ip != "" {
		params.Set("ip", ip)
	}

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

	// 解析响应
	var result RealtimeLocationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查API响应状态
	if result.Status != 0 {
		return nil, fmt.Errorf("百度地图API错误: %s", result.Message)
	}

	return &result, nil
}
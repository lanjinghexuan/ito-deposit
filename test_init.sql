-- 创建测试数据库和表
CREATE DATABASE IF NOT EXISTS ito CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE ito;

-- 创建城市表
CREATE TABLE IF NOT EXISTS `city` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '城市ID',
  `name` varchar(50) NOT NULL COMMENT '城市名称',
  `code` varchar(20) NOT NULL COMMENT '城市编码',
  `latitude` decimal(11,6) NOT NULL COMMENT '纬度',
  `longitude` decimal(11,6) NOT NULL COMMENT '经度',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态(1:启用,0:禁用)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建寄存点表
CREATE TABLE IF NOT EXISTS `locker_point` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '寄存点ID',
  `location_id` int NOT NULL COMMENT '所属地点ID',
  `name` varchar(30) NOT NULL COMMENT '寄存点名称',
  `address` varchar(50) DEFAULT NULL COMMENT '详细地址',
  `latitude` decimal(11,6) NOT NULL COMMENT '纬度',
  `longitude` decimal(11,6) NOT NULL COMMENT '经度',
  `available_large` int DEFAULT '0' COMMENT '可用大柜数量',
  `available_medium` int DEFAULT '0' COMMENT '可用中柜数量',
  `available_small` int DEFAULT '0' COMMENT '可用小柜数量',
  `open_time` varchar(30) DEFAULT NULL COMMENT '营业时间',
  `mobile` varchar(20) DEFAULT NULL COMMENT '联系电话',
  PRIMARY KEY (`id`),
  KEY `location_id` (`location_id`),
  FOREIGN KEY (`location_id`) REFERENCES `city` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试城市数据
INSERT INTO `city` (`name`, `code`, `latitude`, `longitude`, `status`) VALUES
('郑州市', '410100', 34.746611, 113.625328, 1),
('北京市', '110100', 39.904200, 116.407396, 1),
('上海市', '310100', 31.230416, 121.473701, 1),
('广州市', '440100', 23.125178, 113.280637, 1),
('深圳市', '440300', 22.547, 114.085947, 1);

-- 插入测试寄存点数据
INSERT INTO `locker_point` (`location_id`, `name`, `address`, `latitude`, `longitude`, `available_large`, `available_medium`, `available_small`, `open_time`, `mobile`) VALUES
-- 郑州市寄存点
(1, '郑州火车站', '郑州市二七区二马路82号', 34.746611, 113.625328, 5, 10, 15, '06:00-24:00', '0371-12345678'),
(1, '郑州东站', '郑州市金水区心怡路1号', 34.717299, 113.744398, 8, 12, 20, '05:30-24:00', '0371-87654321'),
(1, '二七广场', '郑州市二七区二七路230号', 34.754364, 113.625328, 3, 8, 12, '08:00-22:00', '0371-11111111'),
(1, '中交锦兰荟', '郑州市金水区农业路与花园路交叉口', 34.760000, 113.650000, 6, 10, 18, '24小时', '0371-22222222'),

-- 北京市寄存点
(2, '北京站', '北京市东城区毛家湾胡同甲13号', 39.902200, 116.427300, 10, 15, 25, '05:00-24:00', '010-12345678'),
(2, '北京南站', '北京市丰台区永外大街车站路12号', 39.865000, 116.378600, 12, 18, 30, '05:00-24:00', '010-87654321'),
(2, '天安门广场', '北京市东城区东长安街', 39.903000, 116.397128, 5, 8, 15, '06:00-22:00', '010-11111111'),

-- 上海市寄存点
(3, '上海火车站', '上海市静安区秣陵路303号', 31.249162, 121.455890, 8, 12, 20, '05:30-24:00', '021-12345678'),
(3, '上海虹桥站', '上海市闵行区申贵路1500号', 31.197646, 121.327170, 15, 20, 35, '05:00-24:00', '021-87654321'),
(3, '外滩', '上海市黄浦区中山东一路', 31.239663, 121.490317, 4, 6, 10, '08:00-22:00', '021-11111111');
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `author_info`;
CREATE TABLE `author_info` (
  `id` int(8) NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `unique_id` varchar(255) NOT NULL COMMENT '唯一 ID',
  `platform` varchar(20) NOT NULL COMMENT '博主平台',
  `nickname` varchar(255) NOT NULL COMMENT '博主昵称',
  `signature` varchar(255) DEFAULT NULL COMMENT '签名',
  `avatar_url` text COMMENT '头像url',
  `ip` varchar(20) DEFAULT NULL COMMENT '博主 IP',
  `refer` varchar(255) NOT NULL COMMENT '来源',
  `follower_count` bigint(20) NOT NULL COMMENT '粉丝数量',
  `following_count` int(8) NOT NULL COMMENT '关注数量',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作者信息表';

DROP TABLE IF EXISTS `author_info_history`;
CREATE TABLE `author_info_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `author_id` int(8) NOT NULL COMMENT '作者 ID',
  `last_follower_count` bigint(20) NOT NULL COMMENT '上次粉丝数量',
  `last_following_count` int(8) NOT NULL COMMENT '上次关注数量',
  `num` int(8) NOT NULL COMMENT '涨粉数量',
  `day` varchar(20) NOT NULL COMMENT '记录日期',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_author_day` (`author_id`,`day`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作者信息历史记录表';

INSERT INTO `sys_menu` VALUES (308, '博主管理', 3, 8, 'author', 'get/author/manage/list', 'live/author/index', '', 1, 0, 'C', '0', '0', 'live:author:list', 'user', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_menu` VALUES (30801, '新增', 308, 1, '', 'post/author/manage', '', '', 1, 0, 'F', '0', '0', 'live:author:add', '', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_menu` VALUES (30802, '详情', 308, 2, '', 'get/author/manage/{id}', '', '', 1, 0, 'F', '0', '0', 'live:author:get', '', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_menu` VALUES (30803, '删除', 308, 3, '', 'delete/author/manage/{id}', '', '', 1, 0, 'F', '0', '0', 'live:author:delete', '', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_menu` VALUES (30804, '趋势', 308, 4, '', 'get/author/manage/trend', '', '', 1, 0, 'F', '0', '0', 'live:author:trend', '', 'admin', sysdate(), 'admin', sysdate(), '');

SET FOREIGN_KEY_CHECKS = 1;
SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `media_parse`
ADD COLUMN `video_data` json NULL COMMENT '视频数据' AFTER `video_cover_url`,
ADD COLUMN `referer` varchar(1000) NULL COMMENT '来源' AFTER `platform`;


DROP TABLE IF EXISTS `download_record`;
CREATE TABLE `download_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `title` text COMMENT '任务名称',
  `task_id` varchar(50) NOT NULL COMMENT '任务 ID',
  `status` varchar(20) NOT NULL COMMENT '任务状态',
  `output` varchar(255) DEFAULT NULL COMMENT '输出路径',
  `error_msg` text COMMENT '错误信息',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `update_time` datetime NOT NULL COMMENT '结束时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='下载记录表';

INSERT INTO `sys_menu` VALUES (30604, '下载', 306, 4, '', 'get/media/parse/download', '', '', 1, 0, 'F', '0', '0', 'live:parse:download', '', 'admin', sysdate(), 'admin', null, '');
INSERT INTO `sys_menu` VALUES (307, '下载中心', 3, 7, 'download', 'get/media/download/list', 'live/download/index', '', 1, 0, 'C', '0', '0', 'live:download:list', 'download', 'admin', sysdate(), 'admin', null, '');
INSERT INTO `sys_menu` VALUES (30701, '删除', 307, 1, '', 'delete/media/download/{id}', '', '', 1, 0, 'F', '0', '0', 'live:download:delete', '', 'admin', sysdate(), 'admin', null, '');
UPDATE `sys_menu` SET icon='clock-history' WHERE menu_id=304;
UPDATE `sys_menu` SET icon='cookie' WHERE menu_id=305;

SET FOREIGN_KEY_CHECKS = 1;
SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `media_parse`
ADD COLUMN `video_data` json NULL COMMENT '视频数据' AFTER `video_cover_url`,
ADD COLUMN `referer` varchar(1000) NULL COMMENT '来源' AFTER `platform`;

SET FOREIGN_KEY_CHECKS = 1;
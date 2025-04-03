SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `sys_dict_type` (`dict_name`, `dict_type`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES ('下载状态', 'download_status', '0', 'admin', sysdate(), 'admin', sysdate(), '');

INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (1, '待下载', 'pending', 'download_status', '', 'default', '', '0', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (2, '下载中', 'running', 'download_status', '', 'default', '', '0', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (3, '转换中', 'converting', 'download_status', '', 'default', '', '0', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (4, '已完成', 'completed', 'download_status', '', 'default', '', '0', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (5, '部分完成', 'partSucceed', 'download_status', '', 'default', '', '0', 'admin', sysdate(), 'admin', sysdate(), '');
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `status`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (6, '下载失败', 'error', 'download_status', '', 'default', '', '0', 'admin', sysdate(), 'admin', sysdate(), '');

SET FOREIGN_KEY_CHECKS = 1;
package consts

// SpecialApiPath contains a list of API paths that have special handling requirements
var SpecialApiPath = map[string]bool{
	//user
	"get/getInfo":                       true,
	"get/getRouters":                    true,
	"get/system/user/authRole/{userId}": true,
	"get/system/user/profile":           true,
	"post/logout":                       true,
	"post/system/user/profile/avatar":   true,
	"post/user/sign-up":                 true,
	"put/system/user/profile":           true,
	"put/system/user/profile/updatePwd": true,

	//系统工具
	"get/system/menu/roleMenuTreeselect/{roleId}": true,

	//菜单
	"get/system/menu/treeselect": true,

	//字典
	"get/system/dict/type":              true,
	"get/system/dict/type/optionselect": true,

	// 播放
	"get/file/manage/play": true,
}

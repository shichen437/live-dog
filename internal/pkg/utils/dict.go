package utils

type SystemDictData struct {
	DictType  string `json:"dictType"`
	DictLabel string `json:"dictLabel"`
	DictValue string `json:"dictValue"`
	DictSort  int    `json:"dictSort"`
}

// 默认字典类型
var (
	AssistLivePlatform = "sys_internal_assist_live_platform"
	JobInvoke          = "sys_internal_job_invoke"
	internalDictType   = []string{AssistLivePlatform, JobInvoke}
	m                  = make(map[string]([]SystemDictData), 2)
)

func init() {
	addDictData(AssistLivePlatform, []SystemDictData{
		{DictType: AssistLivePlatform, DictLabel: "抖音", DictValue: "douyin", DictSort: 1},
		{DictType: AssistLivePlatform, DictLabel: "B站", DictValue: "bilibili", DictSort: 2},
	})
	addDictData(JobInvoke, []SystemDictData{
		{DictType: JobInvoke, DictLabel: "空间预警", DictValue: "storageWarning", DictSort: 1},
	})
}

func IsInternalDictData(dictType string) bool {
	return InSliceString(dictType, &internalDictType)
}

func GetDictDataByType(dictType string) *[]SystemDictData {
	data, ok := m[dictType]
	if !ok {
		return &[]SystemDictData{}
	}
	return &data
}

func addDictData(dictType string, dictData []SystemDictData) {
	data, ok := m[dictType]
	if !ok {
		data = []SystemDictData{}
	}
	m[dictType] = append(data, dictData...)
}

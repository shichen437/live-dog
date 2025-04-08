<template>
    <div class="app-container">
        <el-form :model="queryParams" ref="queryRef" v-show="showSearch" :inline="true" label-width="68px">
            <el-form-item label="作者名称" prop="nickname">
                <el-input v-model="queryParams.nickname" placeholder="请输入作者名称" clearable style="width: 180px"
                    @keyup.enter="handleQuery" />
            </el-form-item>
            <el-form-item>
                <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
                <el-button icon="Refresh" @click="resetQuery">重置</el-button>
            </el-form-item>
        </el-form>
        <el-row :gutter="10" class="mb8">
            <el-col :span="1.5">
                <el-button type="primary" plain icon="Plus" @click="handleAdd"
                    v-hasPermi="['live:author:add']">添加</el-button>
            </el-col>
            <el-col :span="1.5">
                <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="handleDelete"
                    v-hasPermi="['live:author:delete']">删除</el-button>
            </el-col>
            <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
        </el-row>

        <!-- 表格数据 -->
        <el-table v-loading="loading" :data="authorList" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" align="center" />
            <el-table-column label="数据ID" align="center" prop="id" width="80" />
            <el-table-column label="平台" align="center" prop="platform" :show-overflow-tooltip="true" width="80">
                <template #default="scope">
                    <dict-tag :options="sys_internal_assist_live_platform" :value="scope.row.platform" />
                </template>
            </el-table-column>
            <el-table-column label="作者" align="center" prop="nickname" :show-overflow-tooltip="true" width="140" />
            <el-table-column label="签名" align="center" prop="signature" :show-overflow-tooltip="true" width="240" />
            <el-table-column label="关注" align="center" prop="followingCount" :show-overflow-tooltip="true" width="160" />
            <el-table-column label="粉丝" align="center" prop="followerCount" :show-overflow-tooltip="true" width="160" />
            <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
                <template #default="scope">
                    <el-tooltip content="详情" placement="top" v-if="scope.row.roleId !== 1">
                        <el-button link type="primary" icon="View" @click="handleDetail(scope.row)"
                            v-hasPermi="['live:author:get']"></el-button>
                    </el-tooltip>
                    <el-tooltip content="统计" placement="top">
                        <el-button link type="primary" icon="Download" @click="handleDownload(scope.row)"
                            v-hasPermi="['live:author:trend']"></el-button>
                    </el-tooltip>
                    <el-tooltip content="删除" placement="top">
                        <el-button link type="primary" icon="Delete" @click="handleDelete(scope.row)"
                            v-hasPermi="['live:author:delete']"></el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
        </el-table>

        <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum"
            v-model:limit="queryParams.pageSize" @pagination="getList" />

        <el-dialog title="主页解析" v-model="parse" align-center close-on-press-escape>
            <template #title>
                媒体解析
                <el-popover placement="right" title="支持平台及类型" :width="300" effect="dark" trigger="hover">
                    <div slot="content">
                        抖音web版分享链接(视频/图集)
                    </div>
                    <template #reference>
                        <el-icon class="m-2">
                            <InfoFilled />
                        </el-icon>
                    </template>
                </el-popover>
            </template>
            <el-form ref="authorRef" :model="form" :rules="rules" label-width="100px">
                <el-form-item label="主页链接" prop="url" class="centered-form-item">
                    <el-input v-model="form.url" placeholder="请输入博主主页链接" clearable size="large">
                        <template #append>
                            <el-button type="primary" @click="submitForm" style="height: 40px; padding: 0 20px;">
                                添加
                            </el-button>
                        </template>
                    </el-input>
                </el-form-item>
            </el-form>
        </el-dialog>

        <el-dialog title="博主详情" v-model="detailDialog" width="60%">
            <el-descriptions :column="2" border v-if="detailData">
                <el-descriptions-item label="作者" :span="1">{{ detailData.nickname }}</el-descriptions-item>
                <el-descriptions-item label="签名" :span="1">{{ detailData.signature }}</el-descriptions-item>
                <el-descriptions-item label="作品描述" :span="2">{{ detailData.desc }}</el-descriptions-item>
            </el-descriptions>
        </el-dialog>

    </div>
</template>

<script setup name="Author">
import {
    listAuthor,
    getAuthor,
    addAuthor,
    delAuthor,
    getTrend,
} from "@/api/live/author";

const { proxy } = getCurrentInstance();
const { sys_internal_assist_live_platform } = proxy.useDict("sys_internal_assist_live_platform");

const authorList = ref([]);
const parse = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const multiple = ref(true);
const total = ref(0);
const detailDialog = ref(false);
const detailData = ref(null);

const data = reactive({
    form: {},
    queryParams: {
        pageNum: 1,
        pageSize: 10,
        nickname: undefined,
    },
    rules: {
        url: [{ required: true, message: "主页链接不能为空", trigger: "blur" }],
    },
});

const { queryParams, form, rules } = toRefs(data);

/** 查询每日统计列表 */
function getList() {
    loading.value = true;
    listAuthor(queryParams.value).then(
        (response) => {
            authorList.value = response.data.rows;
            total.value = response.data.total;
            loading.value = false;
        }
    );
}

/** 重置新增的表单以及其他数据  */
function reset() {
    form.value = {
        id: undefined,
        url: undefined,
    };
    proxy.resetForm("authorRef");
}
/** 添加记录 */
function handleAdd() {
    reset();
    parse.value = true;
}

function handleDetail(row) {
    getAuthor(row.id).then(response => {
        detailData.value = response.data;
        detailDialog.value = true;
    });
}

/** 提交按钮 */
function submitForm() {
    proxy.$refs["authorRef"].validate((valid) => {
        if (valid) {
            addAuthor(form.value).then((response) => {
                proxy.$modal.msgSuccess("添加成功");
                parse.value = false;
                getList();
            });
        }
    });
}

/** 搜索按钮操作 */
function handleQuery() {
    queryParams.value.pageNum = 1;
    getList();
}
/** 重置按钮操作 */
function resetQuery() {
    proxy.resetForm("queryRef");
    handleQuery();
}
/** 删除按钮操作 */
function handleDelete(row) {
    const dataIds = row.id || ids.value;
    proxy.$modal
        .confirm('是否确认删除数据编号为"' + dataIds + '"的数据项?')
        .then(function () {
            return delAuthor(dataIds);
        })
        .then(() => {
            getList();
            proxy.$modal.msgSuccess("删除成功");
        })
        .catch(() => { });
}

getList();
</script>

<style scoped>
.link-container {
    width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
}

.link {
    display: inline-block;
    max-width: 80ch;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
}
</style>
<template>
    <div class="app-container">
        <el-form :model="queryParams" ref="queryRef" v-show="showSearch" :inline="true" label-width="68px">
            <el-form-item label="下载状态" prop="status">
                <el-select v-model="queryParams.status" placeholder="下载状态" clearable style="width: 120px">
                    <el-option v-for="type in download_status" :key="type.value" :label="type.label"
                        :value="type.value" />
                </el-select>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
                <el-button icon="Refresh" @click="resetQuery">重置</el-button>
            </el-form-item>
        </el-form>
        <el-row :gutter="10" class="mb8">
            <el-col :span="1.5">
                <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="handleDelete"
                    v-hasPermi="['live:download:delete']">删除</el-button>
            </el-col>
            <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
        </el-row>

        <!-- 表格数据 -->
        <el-table v-loading="loading" :data="recordList" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" align="center" />
            <el-table-column label="数据编号" align="center" prop="id" width="100" />
            <el-table-column label="源标题" align="center" prop="title" :show-overflow-tooltip="true" width="240" />
            <el-table-column label="下载状态" align="center" width="220">
                <template #default="scope">
                    <el-progress :percentage="getProgressPercentage(scope.row.status)"
                        :status="getProgressStatus(scope.row.status)"
                        :indeterminate="isIndeterminate(scope.row.status)" />
                </template>
            </el-table-column>
            <el-table-column label="错误信息" align="center" prop="errorMsg" :show-overflow-tooltip="true" width="200" />
            <el-table-column label="开始时间" align="center" prop="startTime" :show-overflow-tooltip="true" width="180" />
            <el-table-column label="结束时间" align="center" prop="updateTime" :show-overflow-tooltip="true" width="180" />
            <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
                <template #default="scope">
                    <el-tooltip content="删除" placement="top">
                        <el-button link type="primary" icon="Delete" @click="handleDelete(scope.row)"
                            v-hasPermi="['live:download:delete']"></el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
        </el-table>

        <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum"
            v-model:limit="queryParams.pageSize" @pagination="getList" />

    </div>
</template>

<script setup name="Live">
import {
    listDownload,
    delRecord,
} from "@/api/live/download";

const { proxy } = getCurrentInstance();
const { download_status } = proxy.useDict("download_status");

const recordList = ref([]);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const dateRange = ref([]);

const data = reactive({
    form: {},
    queryParams: {
        pageNum: 1,
        pageSize: 10,
        anchor: undefined,
    },
});

const { queryParams } = toRefs(data);

/** 查询下载记录列表 */
function getList() {
    loading.value = true;
    listDownload(queryParams).then(
        (response) => {
            recordList.value = response.data.rows;
            total.value = response.data.total;
            loading.value = false;
        }
    );
}

/** 搜索按钮操作 */
function handleQuery() {
    queryParams.value.pageNum = 1;
    getList();
}
/** 重置按钮操作 */
function resetQuery() {
    dateRange.value = [];
    proxy.resetForm("queryRef");
    handleQuery();
}
/** 删除按钮操作 */
function handleDelete(row) {
    const dataIds = row.id || ids.value;
    proxy.$modal
        .confirm('是否确认删除数据编号为"' + dataIds + '"的数据项?')
        .then(function () {
            return delRecord(dataIds);
        })
        .then(() => {
            getList();
            proxy.$modal.msgSuccess("删除成功");
        })
        .catch(() => { });
}

/** 多选框选中数据 */
function handleSelectionChange(selection) {
    ids.value = selection.map((item) => item.roleId);
    single.value = selection.length != 1;
    multiple.value = !selection.length;
}

function getProgressPercentage(status) {
    switch (status) {
        case 'pending':
            return 1;
        case 'running':
            return 60;
        case 'converting':
            return 80;
        case 'completed':
            return 100;
        case 'partSucceed':
            return 100;
        case 'error':
            return 10;
        default:
            return 5;
    }
}

/** 根据状态获取进度条状态 */
function getProgressStatus(status) {
    switch (status) {
        case 'completed':
            return 'success';
        case 'partSucceed':
            return 'warning';
        case 'error':
            return 'exception';
        default:
            return 'active';
    }
}

function isIndeterminate(status) {
    return status === 'running' || status === 'converting';
}
getList();
</script>
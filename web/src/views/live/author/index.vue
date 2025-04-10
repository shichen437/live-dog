<template>
    <div class="app-container">
        <el-form :model="queryParams" ref="queryRef" v-show="showSearch" :inline="true" label-width="68px">
            <el-form-item label="博主名称" prop="nickname">
                <el-input v-model="queryParams.nickname" placeholder="请输入博主名称" clearable style="width: 180px"
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
        <el-table v-loading="loading" :data="authorList">
            <el-table-column type="selection" width="55" align="center" />
            <el-table-column label="数据ID" align="center" prop="id" width="80" />
            <el-table-column label="平台" align="center" prop="platform" :show-overflow-tooltip="true" width="80">
                <template #default="scope">
                    <dict-tag :options="sys_internal_assist_live_platform" :value="scope.row.platform" />
                </template>
            </el-table-column>
            <el-table-column label="头像" align="center" width="100">
                <template #default="scope">
                    <el-image style="width: 40px; height: 40px; border-radius: 50%;" :src="scope.row.avatarUrl"
                        :preview-src-list="[scope.row.avatarUrl]" fit="cover" hide-on-click-modal
                        :referrerpolicy="scope.row.platform === 'bilibili' ? 'no-referrer' : 'origin'" />
                </template>
            </el-table-column>
            <el-table-column label="博主" align="center" prop="nickname" :show-overflow-tooltip="true" width="160" />
            <el-table-column label="签名" align="center" prop="signature" :show-overflow-tooltip="true" width="360" />
            <el-table-column label="关注" align="center" prop="followingCount" :show-overflow-tooltip="true"
                width="100" />
            <el-table-column label="粉丝" align="center" prop="followerCount" :show-overflow-tooltip="true" width="100" />
            <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
                <template #default="scope">
                    <el-tooltip content="详情" placement="top">
                        <el-button link type="primary" icon="View" @click="handleDetail(scope.row)"
                            v-hasPermi="['live:author:get']"></el-button>
                    </el-tooltip>
                    <el-tooltip content="统计" placement="top">
                        <el-button link type="primary" icon="TrendCharts" @click="handleTrend(scope.row)"
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
                <el-descriptions-item align="center" label="博主" :span="1">{{ detailData.nickname
                    }}</el-descriptions-item>
                <el-descriptions-item align="center" label="IP" :span="1">{{ detailData.ip }}</el-descriptions-item>
                <el-descriptions-item align="center" label="签名" :span="2">{{ detailData.signature
                    }}</el-descriptions-item>
                <el-descriptions-item align="center" label="关注" :span="1">{{ detailData.followingCount
                    }}</el-descriptions-item>
                <el-descriptions-item align="center" label="粉丝" :span="1">{{ detailData.followerCount
                    }}</el-descriptions-item>
            </el-descriptions>
        </el-dialog>

        <el-dialog title="粉丝趋势" v-model="trendDialog" width="70%" @closed="handleTrendDialogClosed">
            <template #header>
                <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                    <div style="display: flex; align-items: center;">
                        <span style="font-size: 18px; font-weight: 500; margin-right: 12px;">粉丝趋势</span>
                        <el-tag type="primary" size="default" style="margin-right: 5px;">
                            <dict-tag v-if="currentAuthorPlatform" :options="sys_internal_assist_live_platform"
                                :value="currentAuthorPlatform" />
                        </el-tag>
                        <el-tag type="danger" size="default" style="margin-right: 5px;">{{ currentAuthorName }}</el-tag>

                    </div>
                </div>
            </template>
            <div style="height: 400px;">
                <el-empty v-if="trendData.days.length === 0" description="暂无数据" />
                <div v-else>
                    <div ref="chartRef" style="width: 100%; height: 350px;"></div>
                    <div
                        style="text-align: center; margin-top: 12px; display: flex; justify-content: center; align-items: center; gap: 30px;">
                        <el-radio-group v-if="trendData.days.length >= 7" v-model="trendRange" @change="getTrendData"
                            size="default">
                            <el-radio-button :label="7">近7天</el-radio-button>
                            <el-radio-button :label="30">近30天</el-radio-button>
                        </el-radio-group>

                        <el-radio-group v-model="showIncrement" @change="toggleIncrement" size="default">
                            <el-radio-button :label="false">总量</el-radio-button>
                            <el-radio-button :label="true">增量</el-radio-button>
                        </el-radio-group>
                    </div>
                </div>
            </div>
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

import * as echarts from 'echarts';
import { ref, watch, nextTick } from 'vue';

const trendDialog = ref(false);
const trendData = ref({
    days: [],
    counts: [],
    nums: []
});
const trendRange = ref(7);
const currentAuthorId = ref(null);
const currentAuthorName = ref('');
const currentAuthorPlatform = ref(''); // 新增：保存当前博主平台
const chartRef = ref(null);
const showIncrement = ref(false);
let chartInstance = null;

function handleTrend(row) {
    currentAuthorId.value = row.id;
    currentAuthorName.value = row.nickname;
    currentAuthorPlatform.value = row.platform; // 保存平台信息
    getTrendData();
}

function handleTrendDialogClosed() {
    trendRange.value = 7;
    showIncrement.value = false;
    currentAuthorPlatform.value = ''; // 重置平台信息
}

function getTrendData() {
    if (!currentAuthorId.value) return;

    getTrend({
        id: currentAuthorId.value,
        range: trendRange.value
    }).then(response => {
        trendData.value = {
            days: response.data.days,
            counts: response.data.counts,
            nums: response.data.nums
        };
        trendDialog.value = true;
    });
}

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
            proxy.$modal.loading("正在解析，请稍候...");
            addAuthor(form.value).then((response) => {
                proxy.$modal.closeLoading();
                proxy.$modal.msgSuccess("添加成功");
                parse.value = false;
                getList();
            }).catch(() => {
                proxy.$modal.closeLoading();
            });
        }
    });
}

function toggleIncrement() {
    initChart();
}

function initChart() {
    if (!chartRef.value) {
        console.error('Chart container element not found');
        return;
    }

    // Dispose previous chart instance if exists
    if (chartInstance) {
        chartInstance.dispose();
    }

    try {
        chartInstance = echarts.init(chartRef.value);

        const data = showIncrement.value ? trendData.value.nums : trendData.value.counts;
        const title = showIncrement.value ? '粉丝增量' : '粉丝数';

        // 计算纵坐标范围，确保包含负数
        const minValue = Math.min(...data);
        const maxValue = Math.max(...data);
        const range = maxValue - minValue;
        const interval = Math.ceil(range / 5);

        const option = {
            tooltip: {
                trigger: 'axis',
                formatter: '{b}<br/>' + title + ': {c}'
            },
            xAxis: {
                type: 'category',
                data: trendData.value.days,
                axisLabel: {
                    rotate: 45,
                    interval: function (index, value) {
                        // 当数据量大于14个时，只显示首尾和部分中间的标签
                        if (trendData.value.days.length > 14) {
                            return index === 0 ||
                                index === trendData.value.days.length - 1 ||
                                index % Math.ceil(trendData.value.days.length / 7) === 0;
                        }
                        // 数据量较少时全部显示
                        return true;
                    },
                    formatter: function (value) {
                        // 如果是日期格式，可以进一步简化显示
                        if (value.includes('-')) {
                            return value.split('-').slice(1).join('-'); // 只显示月-日
                        }
                        return value;
                    },
                    margin: 8
                }
            },
            yAxis: {
                type: 'value',
                min: minValue - interval, // 确保负数有显示空间
                max: maxValue + interval,
                interval: interval,
                axisLabel: {
                    formatter: function (value) {
                        return value.toLocaleString();
                    }
                }
            },
            // 其余配置保持不变
            series: [{
                name: title,
                data: data,
                type: 'line',
                smooth: true,
                areaStyle: {
                    color: {
                        type: 'linear',
                        x: 0,
                        y: 0,
                        x2: 0,
                        y2: 1,
                        colorStops: [{
                            offset: 0,
                            color: showIncrement.value ?
                                (data[0] >= 0 ? 'rgba(58, 77, 233, 0.8)' : 'rgba(233, 58, 58, 0.8)') :
                                'rgba(58, 77, 233, 0.8)'
                        }, {
                            offset: 1,
                            color: showIncrement.value ?
                                (data[0] >= 0 ? 'rgba(58, 77, 233, 0.1)' : 'rgba(233, 58, 58, 0.1)') :
                                'rgba(58, 77, 233, 0.1)'
                        }]
                    }
                },
                itemStyle: {
                    color: showIncrement.value ?
                        (data[0] >= 0 ? '#3a4de9' : '#e93a3a') :
                        '#3a4de9'
                }
            }]
        };

        chartInstance.setOption(option);

        // Add resize listener
        window.addEventListener('resize', function () {
            chartInstance && chartInstance.resize();
        });
    } catch (error) {
        console.error('Failed to initialize chart:', error);
    }
}

watch(trendData, () => {
    if (trendData.value.days.length > 0 && trendDialog.value) {
        nextTick(() => {
            initChart();
        });
    }
}, { deep: true });

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
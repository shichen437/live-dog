<template>
    <div class="job-doc-container">
        <h2>定时任务说明</h2>

        <el-card class="job-card">
            <h3>空间预警</h3>
            <p>任务参数说明：</p>
            <pre class="code-block">
{
  "threshold": 90,     // 空间使用率阈值，格式为 1-99 正整数， 默认 90
}</pre>
        </el-card>

        <el-card class="job-card">
            <h3>粉丝趋势</h3>
            <p>Cron 表达式说明</p>
            <pre class="code-block">0 0 3 * * ? // 推荐每日一次</pre>
        </el-card>
    </div>
</template>

<script setup>
import { ElCard } from "element-plus";
</script>

<style scoped>
.job-doc-container {
    padding: 20px;
}

h1,
h2 {
    color: #333;
}

h3 {
    color: #7f8c8d;
    font-size: 20px;
    margin-top: 15px;
}

.job-card {
    margin-bottom: 20px;
    padding: 15px;
    background-color: #ffffff;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

p {
    margin-bottom: 10px;
}

ul {
    list-style-type: disc;
    margin-left: 20px;
}

.code-block {
    background-color: #f8f9fa;
    padding: 15px;
    border-radius: 4px;
    font-family: "Monaco", "Menlo", "Consolas", monospace;
    font-size: 14px;
    color: #2c3e50;
    white-space: pre;
    overflow-x: auto;
}
</style>

<template>
  <div class="media-doc-container">
    <h2>媒体解析说明</h2>

    <el-card class="media-card">
      <h3>支持类型</h3>
      <ul>
        <li>抖音图集及视频（Web分享链接和视频链接）</li>
        <li>B站视频（1080P及以上需要配置SESSDATA）</li>
      </ul>
    </el-card>

    <el-card class="media-card">
      <h3>过期说明</h3>
      <p>各平台解析的链接有不同的过期时间，解析后请尽快下载。</p>
    </el-card>

  </div>
</template>

<script setup>
import { ElCard } from 'element-plus';
</script>

<style scoped>
.media-doc-container {
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

.media-card {
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
</style>
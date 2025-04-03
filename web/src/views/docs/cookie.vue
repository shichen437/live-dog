<template>
  <div class="cookie-doc-container">
    <h1>Cookie 获取说明</h1>

    <el-card class="cookie-card">
      <h2>通用 Cookie 获取</h2>
      <h3>步骤一：打开开发者工具</h3>
      <p>在浏览器中按下 <b>F12</b> 或在浏览器选项中打开开发者工具。</p>
      
      <h3>步骤二：导航到“控制台”标签</h3>
      <p>在开发者工具中，选择“控制台”标签，输入 <b>document.cookie</b> ，回车后进行复制。</p>
      
      <h3>步骤三：设置Cookie</h3>
      <p>进入“Cookie管理”标签页，设置对应平台的 Cookie。</p>
    </el-card>

    <el-card class="cookie-card">
      <h2>Bilibili Cookie</h2>
      <p>B站需要额外获取SESSDATA</p>
      <h3>步骤一：打开开发者工具</h3>
      <p>在浏览器中按下 <b>F12</b> 或在浏览器选项中打开开发者工具。</p>
      
      <h3>步骤二：导航到“应用”标签</h3>
      <p>在开发者工具中，选择“应用”标签，在左侧导航中找打“www.bilibili.com”开头的 Cookie，点击后复制<b> SESSDATA </b>的值。</p>
      
      <h3>步骤三：组装Cookie</h3>
      <p>在文本编辑器中进行组装，格式：<b> SESSDATA=(步骤二中的值);(B站Cookie) </b>。（未设置将影响直播录制和视频解析）</p>
    </el-card>
  </div>
</template>

<script setup>
import { ElCard } from 'element-plus';
</script>

<style scoped>
.cookie-doc-container {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
}

h1 {
  color: #2c3e50;
  font-size: 28px;
  margin-bottom: 20px;
}

h2 {
  color: #34495e;
  font-size: 24px;
  margin-top: 20px;
}

h3 {
  color: #7f8c8d;
  font-size: 20px;
  margin-top: 15px;
}

p {
  margin-bottom: 10px;
  line-height: 1.6;
}

ul {
  list-style-type: disc;
  margin-left: 20px;
}

.cookie-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
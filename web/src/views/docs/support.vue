<template>
  <div class="support-doc-container">
    <h1>技术支持</h1>

    <el-card class="support-card">
      <h3>联系开发者</h3>
      <ul>
        <li>shichen437 <span class="email">(shichen437@126.com)</span> </li>
      </ul>
    </el-card>
  </div>
</template>

<script setup>
import { ElCard } from "element-plus";
</script>

<style scoped>
.support-doc-container {
  padding: 20px;
}

h1,
h2 {
  color: #333;
}

p {
  margin-bottom: 10px;
}

h3 {
  color: #7f8c8d;
  font-size: 20px;
  margin-top: 15px;
}

ul {
  list-style-type: disc;
  margin-left: 20px;
}

.email {
  text-decoration: underline;
  color: #3498db;
}

.support-card {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>

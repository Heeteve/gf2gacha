<script setup lang="ts">
import {CaptureClose, CaptureStart} from "../../wailsjs/go/main/App";
import {onMounted} from "vue";
import {EventsOn} from "../../wailsjs/runtime";
import {ElMessage} from "element-plus";

const visible = defineModel({type: Boolean, required: true})

onMounted(() => {
  EventsOn("captureSuccess", () => {
    visible.value = false;
    ElMessage({message: "捕获成功", type: 'success', plain: true, showClose: true, duration: 2000})
  })
})

</script>

<template>
  <el-dialog v-model="visible" width="600" destroy-on-close @open="CaptureStart()" @close="CaptureClose()">
    <template #header>
      <div class="text-xl font-bold">通过抓包捕获用户信息</div>
    </template>
    <div class="flex flex-col gap-1">
      <p class="text-orange-500">首次打开这个界面会弹出安装CA证书</p>
      <p class="text-orange-500">若程序目录下无CA证书则会自动生成私有CA证书，可使用自己的CA证书替换</p>
      <p class="text-red-500 mb-2">不建议使用公用CA证书，存在安全隐患</p>
      <p>1.保持当前界面打开状态</p>
      <p>2.进入游戏打开<span class="font-bold text-blue-500">招募</span>界面</p>
      <p>3.点击任意卡池的<span class="font-bold text-blue-500">访问详情</span></p>
      <p class="mt-2">若成功获取信息，当前界面会自动关闭</p>
    </div>
  </el-dialog>
</template>
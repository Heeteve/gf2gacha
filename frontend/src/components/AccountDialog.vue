<script setup lang="ts">
import {watch, ref} from "vue";
import {GetAccountName, SaveAccountName, SaveAccountPasswd} from "../../wailsjs/go/main/App";
import {ElMessage} from "element-plus";

const visible = defineModel<boolean>({ required: true })
const currentUid = defineModel<string>('currentUid', { required: true })
const accountName = ref("");
const passwd = ref("");

async function handleSubmit() {
  try {
    await SaveAccountName(currentUid.value, accountName.value)
    await SaveAccountPasswd(currentUid.value, passwd.value)
    ElMessage.success({message:'保存成功', plain: true, duration: 1500})
    visible.value = false
  } catch (err) {
    ElMessage.error({message:'保存失败：' + (err), plain: true, duration: 1500})
  }
}

watch(visible, async (val) => {
  if (val) {
    try {
      accountName.value = await GetAccountName(currentUid.value);
    } catch (err) {
    }
    passwd.value = '';
  }
});

</script>

<template>
  <el-dialog v-model="visible" width="400" destroy-on-close>
    <template #header>
      <div class="text-xl font-bold">账号设置</div>
    </template>
    <div class="flex flex-col gap-4">
      <el-alert title="保存账号密码，用于社区Token失效后自动登录。" type="info" show-icon :closable="false"></el-alert>
      <div class="flex items-center gap-2">
        <div class="w-12 shrink-0">Uid：</div>
        <el-input class="grow" readonly v-model="currentUid"></el-input>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-12 shrink-0">账号</div>
        <el-input class="grow" v-model="accountName"></el-input>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-12 shrink-0">密码</div>
        <el-input class="grow" v-model="passwd" type="password"></el-input>
      </div>
      <div class="flex justify-end">
        <el-button type="primary" @click="handleSubmit">提交</el-button>
      </div>
    </div>
  </el-dialog>
</template>
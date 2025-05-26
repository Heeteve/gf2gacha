<script lang="ts" setup>
import { onBeforeMount, onMounted, ref } from "vue";
import { ExportMccExcel, ExportRawJson, GetCurrentVersion, GetLatestVersion, GetLogInfo, GetPoolInfo, GetSettingFont, GetUserList, HandleCommunityTasks, ImportMccExcel, ImportRawJson, MergeEreRecord, StartLogCapture, StopLogCapture, UpdatePoolInfo, UpdateTo } from "../wailsjs/go/main/App";
import PoolCard from "./components/PoolCard.vue";
import { model } from "../wailsjs/go/models";
import 'element-plus/es/components/message/style/css'
import { ElLoading, ElMessage, ElMessageBox } from "element-plus";
import { Connection, CopyDocument, Setting as SettingIcon } from "@element-plus/icons-vue";
import { ClipboardSetText } from "../wailsjs/runtime";
import SettingDialog from "./components/SettingDialog.vue";
import Pool = model.Pool;
import LogInfo = model.LogInfo;
import { useLayoutStore } from "./stores/layout.ts";

const version = ref('')
const currentUid = ref("");
const uidList = ref<string[]>([]);
const poolList = ref<Pool[]>([]);
const logInfo = ref<LogInfo>(new LogInfo)

const loading = ref(false);
const isCapturing = ref(false);
const dialogInfoVisible = ref(false)
const dialogSettingVisible = ref(false)
const dialogCaptureGuideVisible = ref(false)

const getUidList = async () => {
    await GetUserList().then(result => {
        if (result) {
            uidList.value = result
        }
    })
}

const getPoolInfo = async (poolType: number) => {
    await GetPoolInfo(currentUid.value, poolType).then(result => {
        let list = poolList.value
        list.push(result)
        poolList.value = list
    })
}

const getAllPoolInfo = async () => {
    poolList.value = []
    await getPoolInfo(3)
    await getPoolInfo(4)
    await getPoolInfo(6)
    await getPoolInfo(7)
    await getPoolInfo(1)
    await getPoolInfo(5)
    await getPoolInfo(8)
}

const updatePoolInfo = async (isFull: boolean) => {
    loading.value = true
    await UpdatePoolInfo(isFull).then(result => {
        let uid = result[0]
        if (!uidList.value.includes(uid)) {
            uidList.value.push(uid)
        }
        currentUid.value = uid
        ElMessage({ message: result.join("<br/>"), type: 'success', plain: true, showClose: true, duration: 0, dangerouslyUseHTMLString: true })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 0 })
    })
    await getAllPoolInfo()
    loading.value = false
}

const startCapture = async () => {
    isCapturing.value = true
    ElMessage({ message: '开始抓取，首次使用请安装证书(见抓取说明)', type: 'success', plain: true, showClose: true, duration: 0 })
    await StartLogCapture().then(result => {
        if (result) {
            logInfo.value = result
            ElMessage({ message: '抓取成功', type: 'success', plain: true, showClose: true, duration: 0 })
            // dialogInfoVisible.value = true
            currentUid.value = logInfo.value.uid
            stopCapture()
            isCapturing.value = false
        }
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 0 })
        isCapturing.value = false
    })
}

const stopCapture = async () => {
    await StopLogCapture().then(() => {
        ElMessage({ message: '已停止抓取', type: 'success', plain: true, showClose: true, duration: 2000 })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 0 })
    })
    isCapturing.value = false
}

const openInfoDialog = async () => {
    await GetLogInfo().then(result => {
        logInfo.value = result
    })
    dialogInfoVisible.value = true
}

const mergeEreRecord = async (typ: string) => {
    loading.value = true
    await MergeEreRecord(currentUid.value, typ).then(result => {
        ElMessage({ message: result, type: 'success', plain: true, showClose: true, duration: 2000 })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 0 })
    })
    await getAllPoolInfo()
    loading.value = false
}

const importRawJson = async (isReverse: boolean) => {
    loading.value = true
    await ImportRawJson(currentUid.value, isReverse).then(result => {
        ElMessage({ message: result, type: 'success', plain: true, showClose: true, duration: 2000 })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 0 })
    })
    await getAllPoolInfo()
    loading.value = false
}

const exportRawJson = () => {
    ExportRawJson(currentUid.value).then(result => {
        ElMessage({ message: result, type: 'success', plain: true, showClose: true, duration: 2000 })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 2000 })
    })
}

const importMccExcel = async () => {
    loading.value = true
    await ImportMccExcel(currentUid.value).then(result => {
        ElMessage({ message: result, type: 'success', plain: true, showClose: true, duration: 2000 })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 2000 })
    })
    await getAllPoolInfo()
    loading.value = false
}

const exportMccExcel = () => {
    ExportMccExcel(currentUid.value).then(result => {
        ElMessage({ message: result, type: 'success', plain: true, showClose: true, duration: 2000 })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 2000 })
    })
}

const copyUid = () => {
    ClipboardSetText(logInfo.value.uid)
    ElMessage({ message: 'UID已复制', type: 'success', plain: true, showClose: true, duration: 1000 })
}

const copyGachaUrl = () => {
    ClipboardSetText(logInfo.value.gachaUrl)
    ElMessage({ message: '抽卡链接已复制', type: 'success', plain: true, showClose: true, duration: 1000 })
}

const copyAccessToken = () => {
    ClipboardSetText(logInfo.value.accessToken)
    ElMessage({ message: 'AccessToken已复制', type: 'success', plain: true, showClose: true, duration: 1000 })
}

const handleCommunityTasks = () => {
    HandleCommunityTasks().then(result => {
        ElMessage({ message: result.join("<br/>"), type: 'success', plain: true, showClose: true, duration: 0, dangerouslyUseHTMLString: true })
    }).catch(err => {
        ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 0 })
    })
}

const checkUpdate = () => {
    GetLatestVersion().then(latestVersion => {
        if (latestVersion != version.value) {
            ElMessageBox.confirm("有可用的新版本，是否升级", latestVersion, { confirmButtonText: "是", cancelButtonText: "否", type: 'info' }).then(() => {
                const loading = ElLoading.service({ lock: true, text: `Update to ${latestVersion}...`, background: 'rgba(0, 0, 0, 0.7)' })
                UpdateTo(latestVersion).catch(err => {
                    ElMessage({ message: err, type: 'error', plain: true, showClose: true, duration: 2000 })
                }).finally(() => {
                    loading.close()
                })
            })
        }
    })
}

const layoutStore = useLayoutStore()

onBeforeMount(async () => {
    let body = document?.querySelector('body')
    if (body) {
        await GetSettingFont().then(result => {
            if (result) {
                body.style.fontFamily = result
            } else {
                body.style.fontFamily = 'ChillRoundM'
            }
        })
    }
})

onMounted(async () => {
    await getUidList()
    if (uidList.value.length > 0) {
        currentUid.value = uidList.value[0]
        await getAllPoolInfo()
    }

    await GetCurrentVersion().then(res => {
        if (res) {
            version.value = res;
    }
  })

  checkUpdate()
})

</script>

<template>
    <div class="h-dvh w-full flex flex-col p-4 gap-4" v-loading="loading" element-loading-text="Loading...">
        <div class="flex">
            <div class="grow">
                <el-button type="success" class="font-bold" @click="updatePoolInfo(false)">增量更新</el-button>
                <el-button type="primary" class="font-bold" @click="updatePoolInfo(true)">全量更新</el-button>
                <el-dropdown class="ml-3">
                    <el-button type="danger" class="font-bold">导入导出</el-button>
                    <template #dropdown>
                        <el-dropdown-menu :disabled="!currentUid">
                            <el-dropdown-item @click="mergeEreRecord('json')">导入EreJson</el-dropdown-item>
                            <el-dropdown-item @click="mergeEreRecord('excel')">导入EreExcel</el-dropdown-item>
                            <el-dropdown-item @click="importRawJson(false)">
                                <el-tooltip effect="dark" content="同一个RawJson里时间顺序一定要保持一致"
                                    placement="right">导入RawJson</el-tooltip>
                            </el-dropdown-item>
                            <el-dropdown-item @click="importRawJson(true)">
                                <el-tooltip effect="dark" content="同一个RawJson里时间顺序一定要保持一致"
                                    placement="right">导入RawJson(时间倒序)</el-tooltip>
                            </el-dropdown-item>
                            <el-dropdown-item @click="importMccExcel">导入MccExcel</el-dropdown-item>
                            <el-dropdown-item divided @click="exportRawJson">导出RawJson</el-dropdown-item>
                            <el-dropdown-item @click="exportMccExcel">导出MccExcel</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
            <div class="flex items-center gap-2">
                <el-button type="primary" class="font-bold" @click="handleCommunityTasks">一键社区</el-button>
                <div>UID:</div>
                <el-select v-model="currentUid" class="!w-28" @change="getAllPoolInfo">
                    <el-option v-for="uid in uidList" :key="uid" :label="uid" :value="uid" />
                </el-select>
                <el-button :type="isCapturing ? 'danger' : 'warning'" class="font-bold"
                    @click="isCapturing ? stopCapture() : startCapture()">
                    {{ isCapturing ? '停止抓取' : '开始抓取' }}
                </el-button>
                <!-- <el-dropdown trigger="click"> -->
                <el-button type="info" class="!ml-0" @click="dialogCaptureGuideVisible = true">抓取说明</el-button>
                <!-- <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="dialogCaptureGuideVisible = true">查看说明</el-dropdown-item>
                        </el-dropdown-menu>
                    </template> -->
                <!-- </el-dropdown> -->
                <el-button text :icon="Connection" circle @click="openInfoDialog" />
                <el-button class="!ml-0" text :icon="SettingIcon" circle @click="dialogSettingVisible = true" />
            </div>
        </div>
        <div class="w-full flex flex-wrap"
            :class="{ 'gap-16': layoutStore.layoutType===0,'gap-4':layoutStore.layoutType===1}">
            <PoolCard v-for="pool in poolList" :pool="pool"></PoolCard>
        </div>
        <el-dialog v-model="dialogInfoVisible" width="600">
            <template #header>
                <div class="text-xl font-bold">关于</div>
            </template>
            <div class="flex flex-col gap-4">
                <div class="flex items-center gap-2">
                    <div class="w-24 shrink-0">项目地址</div>
                    <div class="grow text-blue-500">https://github.com/MatchaCabin/gf2gacha</div>
                </div>
                <div class="flex items-center gap-2">
                    <div class="w-24 shrink-0">UID</div>
                    <el-input class="grow" readonly v-model="logInfo.uid" />
                    <el-button text :icon="CopyDocument" circle @click="copyUid" />
                </div>
                <div class="flex items-center gap-2">
                    <div class="w-24 shrink-0">抽卡链接</div>
                    <el-input class="grow" readonly v-model="logInfo.gachaUrl" />
                    <el-button text :icon="CopyDocument" circle @click="copyGachaUrl" />
                </div>
                <div class="flex items-center gap-2">
                    <div class="w-24 shrink-0">AccessToken</div>
                    <el-input class="grow" readonly type="password" v-model="logInfo.accessToken" />
                    <el-button text :icon="CopyDocument" circle @click="copyAccessToken" />
                </div>
                <el-alert title="AccessToken是您的临时登录凭证，请自行把控风险，切勿随意泄露" type="warning" show-icon
                    :closable="false"></el-alert>
            </div>
        </el-dialog>
        <SettingDialog v-model="dialogSettingVisible" />
        <!-- 说明对话框 -->
        <el-dialog v-model="dialogCaptureGuideVisible" title="抓取说明" width="800">
            <div class="whitespace-pre-wrap">
                <p>抓取功能由Heeteve基于<el-link type="primary" href="https://github.com/lqqyt2423/go-mitmproxy/blob/main/README_CN.md"
                        target="_blank">go-mitmproxy</el-link>使用Claude4魔改。修改后的<el-link
                        type="primary" href="https://github.com/Heeteve/gf2gacha"
                        target="_blank">源码地址</el-link>。
                   AI编写，仅自用，大概会有bug，能用就行( ，注意备份抽卡数据库gf2gacha.db</p>
                <br>
                <p><strong>抓取使用说明：</strong></p>
                <ol>
                    <li>1. 在游戏未登录的状态下，点击 <strong>开始抓取</strong></li>
                    <li>2. <strong>首次抓取时</strong>
                        需要安装cert文件夹内的mitmproxy-ca-cert.cer证书到当前用户&gt;将所有的证书放入下列存储-受信任的根证书颁发机构。（叠甲：该证书由<el-link
                            type="primary" href="https://github.com/lqqyt2423/go-mitmproxy/blob/main/README_CN.md"
                            target="_blank">go-mitmproxy</el-link>通过<el-link
                            type="primary" href="https://docs.mitmproxy.org/stable/concepts/certificates/"
                            target="_blank">mitmproxy</el-link>签发，本人不保证安全性，若不再使用请及时<el-link
                            type="primary" href="https://jingyan.baidu.com/article/6c67b1d659f8612786bb1e50.html"
                            target="_blank">删除</el-link>)</li>
                    <li>3. 打开游戏，登录，正常情况应提示抓取成功</li>
                </ol>
                <br>
                <p><strong>常见问题：</strong></p>
                <ul>
                    <li><strong>用户中心初始化失败/你的连接不是专用连接</strong>：未安装证书或证书过期，检查是否完成第2步</li>
                    <li><strong>token无效或已过期</strong>：从第1步开始重新抓取</li>
                    <li><strong>抓取完上不了网了?</strong>：1. 确保关闭软件前先停止抓取。2.
                        请检查Windows代理是否关闭：Windows设置&gt;网络和Internet&gt;代理&gt;手动设置代理&gt;关。
                    </li>
                </ul>
            </div>
        </el-dialog>
    </div>
</template>

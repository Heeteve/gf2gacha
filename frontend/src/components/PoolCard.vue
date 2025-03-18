<script setup lang="ts">
import {model} from "../../wailsjs/go/models";
import {use} from 'echarts/core';
import {PieChart} from 'echarts/charts';
import {LegendComponent, TitleComponent, TooltipComponent} from 'echarts/components';
import {CanvasRenderer} from 'echarts/renderers';
import VChart from 'vue-echarts';
import {ref} from "vue";
import {Share} from "@element-plus/icons-vue";
import html2canvas from "html2canvas-pro";
import Pool = model.Pool;

use([PieChart, TitleComponent, TooltipComponent, LegendComponent, CanvasRenderer]);

const props = defineProps<{ pool: Pool }>()
const title = (): string => {
  switch (props.pool.poolType) {
    case 1:
      return '常驻池'
    case 3:
      return '角色池'
    case 4:
      return '武器池'
    case 5:
      return '新手池'
    case 6:
      return '自选角色池'
    case 7:
      return '自选武器池'
    case 8:
      return '神秘箱'
    default:
      return `未知PoolType: ${props.pool.poolType}`
  }
}
// const tagType = ['primary', 'success', 'warning', 'danger']

const option = {
  tooltip: {
    trigger: 'item',
  },
  legend: {
    orient: 'horizontal',
    left: 'center',
  },
  series: [
    {
      type: 'pie',
      radius: '80%',
      data: [
        {value: props.pool.rank5Count, name: '五星', itemStyle: {color: '#fdcb51'}},
        {value: props.pool.rank4Count, name: '四星', itemStyle: {color: '#ddb0e2'}},
        {value: props.pool.rank3Count, name: '三星', itemStyle: {color: '#409EFF'}}],
      itemStyle: {
        borderRadius: 4, // 让每个扇形区域有圆角
        borderColor: '#fff',
        borderWidth: 2,
      },
      label: {
        show: false,
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 20,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)',
        },
      },
    },
  ],
}

const isDesc = ref(true)
const getOrderedRecordList = () => {
  if (isDesc.value && props.pool.recordList) {
    return props.pool.recordList.slice().reverse()
  } else {
    return props.pool.recordList
  }
}

const screenshotContainer = ref<HTMLElement | null>(null);
const screenshot = async () => {
  if (!screenshotContainer.value) {
    return;
  }

  const canvas = await html2canvas(screenshotContainer.value, {
    backgroundColor: '#fff',
    scrollX: 0,
    scrollY: 0,
    windowWidth: screenshotContainer.value.scrollWidth,
    windowHeight: screenshotContainer.value.scrollHeight,
    scale: 2,
    removeContainer: true,
  });

  // 转换为Blob
  const blob = await new Promise<Blob>((resolve, reject) => {
    canvas.toBlob((blob) => {
      if (blob) {
        resolve(blob)
      } else {
        reject(new Error('Canvas 转换失败'))
      }
    }, 'image/png')
  })

  // 复制到剪贴板
  await navigator.clipboard.write([
    new ClipboardItem({'image/png': blob})
  ])
}

</script>

<template>
  <div ref="screenshotContainer" class="relative w-xl shrink-0 grow-0 flex flex-col items-center gap-2 my-2 p-4 shadow rounded-xl" v-if="pool.gachaCount">
    <div class="absolute right-2 top-2">
      <el-button text :icon="Share" size="large" circle @click="screenshot"/>
    </div>
    <div class="font-bold text-xl">{{ title() }}</div>
    <div class="h-64 w-64">
      <v-chart class="h-64" :option="option"></v-chart>
    </div>
    <div class="flex flex-col w-full">
      <div class="w-full text-sm">一共 <span class="text-red-600">{{ pool.gachaCount }}</span> 抽，已垫 <span class="text-red-600">{{ pool.storedCount }}</span> 抽</div>
      <div class="w-full text-sm text-amber-600">五星: {{ pool.rank5Count }} [{{ pool.gachaCount > 0 ? Math.round(pool.rank5Count * 10000 / pool.gachaCount) / 100 + '%' : '0%' }}]</div>
      <div class="w-full text-sm text-purple-600">四星: {{ pool.rank4Count }} [{{ pool.gachaCount > 0 ? Math.round(pool.rank4Count * 10000 / pool.gachaCount) / 100 + '%' : '0%' }}]</div>
      <div class="w-full text-sm text-blue-600">三星: {{ pool.rank3Count }} [{{ pool.gachaCount > 0 ? Math.round(pool.rank3Count * 10000 / pool.gachaCount) / 100 + '%' : '0%' }}]</div>
      <div class="w-full text-sm text-green-600">平均出金抽数：{{ pool.rank5Count > 0 ? ((pool.gachaCount - pool.storedCount) / pool.rank5Count).toFixed(1) : '无' }}</div>
      <div class="w-full text-sm text-pink-400">平均出Up抽数：{{ pool.rank5Count - pool.loseCount > 0 ? ((pool.gachaCount - pool.storedCount) / (pool.rank5Count - pool.loseCount)).toFixed(1) : '无' }}</div>
      <div class="w-full text-sm text-red-600" v-if="pool.poolType==3||pool.poolType==4">歪率: {{ pool.rank5Count > 0 ? Math.round(pool.loseCount * 10000 / (pool.rank5Count - pool.guaranteesCount)) / 100 + '%' : '0%' }}</div>
    </div>
    <div class="w-full text-sm text-gray-400 flex flex-row justify-between" v-if="pool.recordList">
      <div>五星抽卡记录：</div>
      <div class="cursor-pointer mr-2 text-blue-400" @click="isDesc=!isDesc">{{ isDesc ? '倒序' : '正序' }}</div>
    </div>
    <div class="w-full flex flex-wrap gap-2">
      <el-tag class="relative text-sm w-32" v-for="record in getOrderedRecordList()" effect="dark" :type="record.count<=40?'success':record.count<=60?'primary':record.count<=65?'warning':'danger'">
        <span>{{ record.name }}</span>
        <span class="font-bold">「{{ record.count }}」</span>
        <span v-if="record.lose" class="text-purple-600 font-bold absolute right-1">歪</span>
      </el-tag>
    </div>
  </div>
</template>
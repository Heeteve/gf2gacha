import { defineStore } from 'pinia'
import {ref} from "vue";

export const useLayoutStore = defineStore('counter', () => {
    const layoutType = ref(0)

    return { layoutType }
})
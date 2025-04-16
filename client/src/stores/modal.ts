import {defineStore} from "pinia";
import {ref} from "vue";

export const useModal = defineStore('modal', () => {
    let isOpen = ref(false)
    let buildingContext = ref('')

    function open(building: string): void {
        isOpen.value = true
        buildingContext.value = building;

        console.log(isOpen, buildingContext)
    }

    function close(): void {
        isOpen.value = false;
    }

    return {open, close, isOpen, buildingContext}
})

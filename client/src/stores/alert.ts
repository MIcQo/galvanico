import {defineStore} from "pinia";
import {ref} from "vue";

export enum AlertType {
  success = "success",
  warning = "warning",
  danger = "danger",
}

export const useAlert = defineStore('alert', () => {
  let isOpen = ref(false)
  let message = ref('');
  let alertType = ref(AlertType.success);

  let timeoutHandler = null;

  function open(msg: string, type: AlertType = AlertType.success, timeout: number = 5): void {
    isOpen.value = true
    message.value = msg;
    alertType.value = type;

    if (timeout > 0) {
      if (timeoutHandler != null) {
        window.clearTimeout(timeoutHandler);
      }
      timeoutHandler = setTimeout(() => {
        close()
      }, timeout * 1000)
    }
  }

  function close(): void {
    isOpen.value = false
    message.value = ''
    alertType.value = AlertType.success;
  }

  return {isOpen, open, alertType, message, close}
})

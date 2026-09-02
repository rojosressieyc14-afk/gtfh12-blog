import { ref } from "vue";

export function useFlashMessage(timeout = 2200) {
  const flash = ref("");

  function say(message) {
    flash.value = message;
    window.setTimeout(() => {
      if (flash.value === message) {
        flash.value = "";
      }
    }, timeout);
  }

  return { flash, say };
}

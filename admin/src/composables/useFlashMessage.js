import { ref } from "vue";

export function useFlashMessage(timeout = 2200) {
  const flash = ref("");

  function say(msg) {
    flash.value = msg;
    window.setTimeout(() => {
      if (flash.value === msg) flash.value = "";
    }, timeout);
  }

  return { flash, say };
}

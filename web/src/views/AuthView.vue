<template>
  <section
    class="auth-stage"
    :class="[
      `auth-stage--${mode}`,
      `auth-stage--${introPhase}`,
      {
        'auth-stage--password': focusedField === 'password',
        'auth-stage--switching': switching,
        'auth-stage--celebrating': celebrating,
        'auth-stage--motion-ready': motionReady
      }
    ]"
  >
    <div v-if="celebrating" class="auth-brand-flash">
      <div class="brand-flash__core">
        <span class="brand-badge brand-badge--flash">P</span>
        <strong>PulseBlog</strong>
      </div>
    </div>

    <div class="auth-lane">
      <div
        ref="sceneRef"
        class="auth-scene panel-card"
        @mousemove="handleSceneMove"
        @mouseleave="resetSceneMove"
      >
        <div class="auth-figure-wrap">
          <div class="auth-figure-glow"></div>
          <svg class="auth-figure" viewBox="0 0 620 500" fill="none" xmlns="http://www.w3.org/2000/svg">
            <g class="ambient ambient--one">
              <path d="M72 132C102 86 164 74 208 104" />
            </g>
            <g class="ambient ambient--two">
              <path d="M402 86C454 58 526 78 554 126" />
            </g>
            <g class="ambient ambient--three">
              <path d="M138 392C202 426 286 424 332 390" />
            </g>

            <g class="spark spark--one"><circle cx="132" cy="98" r="5" /></g>
            <g class="spark spark--two"><circle cx="508" cy="106" r="6" /></g>
            <g class="spark spark--three"><circle cx="460" cy="364" r="4" /></g>

            <g class="fragment fragment--one">
              <path d="M92 164L108 146L124 162L108 178Z" />
            </g>
            <g class="fragment fragment--two">
              <circle cx="504" cy="152" r="12" />
            </g>
            <g class="fragment fragment--three">
              <path d="M178 404L196 388L214 404L196 422Z" />
            </g>
            <g class="fragment fragment--four">
              <path d="M442 360C454 350 470 352 478 364C486 376 482 392 470 398C458 404 444 400 438 388C432 378 433 366 442 360Z" />
            </g>

            <g class="desk-group" :style="deskStyle">
              <rect x="146" y="332" width="332" height="24" rx="12" class="desk-line" />
              <path d="M184 356V424" class="desk-leg" />
              <path d="M440 356V424" class="desk-leg" />
            </g>

            <g class="shape-team" :style="teamStyle">
              <g class="shape-slot" :style="triangleSlotStyle">
                <g class="shape shape--triangle" :style="triangleInnerStyle">
                  <path class="shape-fill triangle-fill" d="M180 278L118 178L244 176L180 278Z" />
                  <path class="shape-line" d="M180 278L118 178L244 176L180 278Z" />
                  <g class="eyes">
                    <ellipse cx="165" cy="212" rx="15" :ry="triangleEyeRy" class="eye-shell" />
                    <ellipse cx="200" cy="212" rx="15" :ry="triangleEyeRy" class="eye-shell" />
                    <circle :cx="165 + pupilX" :cy="212 + pupilY" :r="trianglePupilR" class="pupil" />
                    <circle :cx="200 + pupilX" :cy="212 + pupilY" :r="trianglePupilR" class="pupil" />
                    <circle :cx="168 + pupilX * 0.5" :cy="209 + pupilY * 0.35" r="2.4" class="pupil-glint" />
                    <circle :cx="203 + pupilX * 0.5" :cy="209 + pupilY * 0.35" r="2.4" class="pupil-glint" />
                  </g>
                </g>
              </g>

              <g class="shape-slot" :style="squareSlotStyle">
                <g class="shape shape--square" :style="squareInnerStyle">
                  <rect x="226" y="154" width="132" height="132" rx="34" class="shape-fill square-fill" />
                  <rect x="226" y="154" width="132" height="132" rx="34" class="shape-line" />
                  <g class="eyes">
                    <ellipse cx="268" cy="208" rx="18" :ry="mainEyeRy" class="eye-shell" />
                    <ellipse cx="316" cy="208" rx="18" :ry="mainEyeRy" class="eye-shell" />
                    <circle :cx="268 + pupilX" :cy="208 + pupilY" :r="mainPupilR" class="pupil" />
                    <circle :cx="316 + pupilX" :cy="208 + pupilY" :r="mainPupilR" class="pupil" />
                    <circle :cx="271 + pupilX * 0.5" :cy="204 + pupilY * 0.35" r="2.8" class="pupil-glint" />
                    <circle :cx="319 + pupilX * 0.5" :cy="204 + pupilY * 0.35" r="2.8" class="pupil-glint" />
                  </g>
                  <path v-if="focusedField !== 'password'" class="mouth-line" d="M270 242C282 250 302 250 314 242" />
                  <path v-else class="mouth-line" d="M274 246C286 238 301 238 312 246" />
                </g>
              </g>

              <g class="shape-slot" :style="circleSlotStyle">
                <g class="shape shape--circle" :style="circleInnerStyle">
                  <circle cx="418" cy="236" r="58" class="shape-fill circle-fill" />
                  <circle cx="418" cy="236" r="58" class="shape-line" />
                  <g class="eyes">
                    <ellipse cx="398" cy="228" rx="14" :ry="circleEyeRy" class="eye-shell" />
                    <ellipse cx="433" cy="228" rx="14" :ry="circleEyeRy" class="eye-shell" />
                    <circle :cx="398 + pupilX" :cy="228 + pupilY" :r="circlePupilR" class="pupil" />
                    <circle :cx="433 + pupilX" :cy="228 + pupilY" :r="circlePupilR" class="pupil" />
                    <circle :cx="401 + pupilX * 0.5" :cy="225 + pupilY * 0.35" r="2.2" class="pupil-glint" />
                    <circle :cx="436 + pupilX * 0.5" :cy="225 + pupilY * 0.35" r="2.2" class="pupil-glint" />
                  </g>
                  <path class="mouth-line" d="M397 262C407 269 425 270 438 260" />
                </g>
              </g>

              <g class="shape-slot" :style="blobSlotStyle">
                <g class="shape shape--blob" :style="blobInnerStyle">
                  <path
                    class="shape-fill blob-fill"
                    d="M470 164C494 144 535 149 553 175C570 199 567 239 546 258C525 276 488 274 470 248C452 221 447 184 470 164Z"
                  />
                  <path
                    class="shape-line"
                    d="M470 164C494 144 535 149 553 175C570 199 567 239 546 258C525 276 488 274 470 248C452 221 447 184 470 164Z"
                  />
                  <g v-if="focusedField !== 'password'" class="eyes">
                    <ellipse cx="499" cy="203" rx="13" :ry="blobEyeRy" class="eye-shell" />
                    <ellipse cx="533" cy="203" rx="13" :ry="blobEyeRy" class="eye-shell" />
                    <circle :cx="499 + pupilX" :cy="203 + pupilY" :r="blobPupilR" class="pupil" />
                    <circle :cx="533 + pupilX" :cy="203 + pupilY" :r="blobPupilR" class="pupil" />
                  </g>
                  <g v-else class="eyes">
                    <ellipse cx="536" cy="205" rx="13" :ry="blobEyeRy" class="eye-shell" />
                    <circle :cx="536 + peekX" :cy="205 + peekY" :r="blobPupilR" class="pupil" />
                    <circle :cx="538 + peekX * 0.5" :cy="202 + peekY * 0.35" r="2.2" class="pupil-glint" />
                  </g>
                  <path class="mouth-line" d="M496 234C507 241 523 240 531 231" />
                </g>
              </g>

              <g class="shape-slot" :style="miniSlotStyle">
                <g class="shape shape--mini" :style="miniInnerStyle">
                  <path class="shape-fill mini-fill" d="M146 312C159 300 179 301 192 313C204 325 204 346 191 357C178 367 159 367 147 355C134 343 133 323 146 312Z" />
                  <path class="shape-line" d="M146 312C159 300 179 301 192 313C204 325 204 346 191 357C178 367 159 367 147 355C134 343 133 323 146 312Z" />
                  <ellipse cx="164" cy="330" rx="8.5" :ry="miniEyeRy" class="eye-shell" />
                  <ellipse cx="182" cy="330" rx="8.5" :ry="miniEyeRy" class="eye-shell" />
                  <circle :cx="164 + pupilX * 0.7" :cy="330 + pupilY * 0.7" :r="miniPupilR" class="pupil" />
                  <circle :cx="182 + pupilX * 0.7" :cy="330 + pupilY * 0.7" :r="miniPupilR" class="pupil" />
                </g>
              </g>
            </g>
          </svg>
        </div>
      </div>

      <div class="auth-panel panel-card">
        <div class="auth-panel__head">
          <p class="eyebrow">{{ mode === "login" ? "登录" : "注册" }}</p>
          <div class="tab-row auth-tabs">
            <button :class="{ active: mode === 'login' }" @click="switchMode('login')">登录</button>
            <button :class="{ active: mode === 'register' }" @click="switchMode('register')">注册</button>
          </div>
        </div>

        <form class="stack-form" @submit.prevent="submit">
          <label>
            用户名
            <input
              v-model.trim="form.username"
              class="field-input"
              placeholder="请输入用户名，至少 3 位"
              @focus="handleFocus('username', $event)"
              @input="handleFocus('username', $event)"
            />
          </label>

          <label>
            密码
            <input
              v-model.trim="form.password"
              class="field-input"
              type="password"
              placeholder="请输入密码，至少 6 位"
              @focus="handleFocus('password', $event)"
              @input="handleFocus('password', $event)"
            />
          </label>

          <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>

          <button class="solid-btn auth-submit" :disabled="userStore.loading || celebrating">
            {{
              userStore.loading
                ? "处理中..."
                : celebrating
                  ? "欢迎回来"
                  : mode === "login"
                    ? "进入我的站点"
                    : "注册并进入站点"
            }}
          </button>
        </form>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "../stores/user";

const router = useRouter();
const userStore = useUserStore();
const sceneRef = ref(null);

const mode = ref("login");
const errorMessage = ref("");
const focusedField = ref("username");
const introPhase = ref("scatter");
const switching = ref(false);
const celebrating = ref(false);
const motionReady = ref(false);
const blink = ref(false);
const look = reactive({ x: 0, y: 0 });
const sceneParallax = reactive({ x: 0, y: 0 });
const form = reactive({
  username: "",
  password: ""
});

let blinkTimer = 0;
let introTimer1 = 0;
let introTimer2 = 0;
let switchTimer1 = 0;
let switchTimer2 = 0;
let switchTimer3 = 0;
let successTimer = 0;

function clamp(value, min, max) {
  return Math.min(max, Math.max(min, value));
}

function clearSwitchTimers() {
  window.clearTimeout(switchTimer1);
  window.clearTimeout(switchTimer2);
  window.clearTimeout(switchTimer3);
}

function switchMode(nextMode) {
  if (mode.value === nextMode || switching.value) return;
  switching.value = true;
  motionReady.value = false;
  introPhase.value = "switch-scatter";
  errorMessage.value = "";
  focusedField.value = "username";
  clearSwitchTimers();

  switchTimer1 = window.setTimeout(() => {
    mode.value = nextMode;
  }, 220);

  switchTimer2 = window.setTimeout(() => {
    introPhase.value = "switch-forming";
  }, 520);

  switchTimer3 = window.setTimeout(() => {
    introPhase.value = "formed";
    switching.value = false;
    motionReady.value = true;
  }, 1280);
}

function handleFocus(field, event) {
  focusedField.value = field;
  const rect = event?.target?.getBoundingClientRect?.();
  if (!rect) return;

  const x = ((rect.left + rect.width / 2) / window.innerWidth - 0.5) * 2;
  const y = ((rect.top + rect.height / 2) / window.innerHeight - 0.5) * 2;
  look.x = clamp(x, -1, 1);
  look.y = clamp(y, -1, 1);
}

function handleSceneMove(event) {
  const rect = sceneRef.value?.getBoundingClientRect?.();
  if (!rect) return;
  const x = ((event.clientX - rect.left) / rect.width - 0.5) * 2;
  const y = ((event.clientY - rect.top) / rect.height - 0.5) * 2;
  sceneParallax.x = clamp(x, -1, 1);
  sceneParallax.y = clamp(y, -1, 1);
}

function resetSceneMove() {
  sceneParallax.x = 0;
  sceneParallax.y = 0;
}

const pupilX = computed(() => Math.round((look.x + sceneParallax.x * 0.45) * 6.5));
const pupilY = computed(() => Math.round((look.y + sceneParallax.y * 0.35) * 4.6));
const peekX = computed(() => Math.round((look.x + sceneParallax.x * 0.4) * 4));
const peekY = computed(() => Math.round((look.y + sceneParallax.y * 0.28) * 3));

const eyeScale = computed(() => (blink.value ? 0.18 : 1));
const triangleEyeRy = computed(() => 18 * eyeScale.value);
const mainEyeRy = computed(() => 22 * eyeScale.value);
const circleEyeRy = computed(() => 18 * eyeScale.value);
const blobEyeRy = computed(() => 17 * eyeScale.value);
const miniEyeRy = computed(() => 11 * eyeScale.value);

const trianglePupilR = computed(() => (blink.value ? 2 : 7.5));
const mainPupilR = computed(() => (blink.value ? 2.4 : 9));
const circlePupilR = computed(() => (blink.value ? 2 : 7));
const blobPupilR = computed(() => (blink.value ? 2 : 6.8));
const miniPupilR = computed(() => (blink.value ? 1.6 : 4.2));

const phasePreset = computed(() => {
  switch (introPhase.value) {
    case "scatter":
      return { factor: 0, deskScale: 0.7, deskOpacity: 0, scatterBoost: 1.2, opacity: 0.08 };
    case "forming":
      return { factor: 0.54, deskScale: 0.88, deskOpacity: 0.64, scatterBoost: 0.56, opacity: 0.92 };
    case "switch-scatter":
      return { factor: 0.18, deskScale: 0.82, deskOpacity: 0.48, scatterBoost: 1.35, opacity: 0.72 };
    case "switch-forming":
      return { factor: 0.8, deskScale: 0.95, deskOpacity: 0.88, scatterBoost: 0.22, opacity: 1 };
    default:
      return { factor: 1, deskScale: 1, deskOpacity: 1, scatterBoost: 0, opacity: 1 };
  }
});

function phaseTranslate(startX, startY) {
  return {
    x: startX * (1 - phasePreset.value.factor) * (1 + phasePreset.value.scatterBoost * 0.35),
    y: startY * (1 - phasePreset.value.factor) * (1 + phasePreset.value.scatterBoost * 0.35)
  };
}

const modeShift = computed(() => (mode.value === "login" ? -34 : 34));
const switchingShift = computed(() => {
  if (!switching.value) return 0;
  return mode.value === "login" ? 44 : -44;
});

const teamStyle = computed(() => ({
  transform: `translate(${modeShift.value + switchingShift.value + sceneParallax.x * 14}px, ${
    (focusedField.value === "username" ? look.y * 3 : 0) + sceneParallax.y * 8
  }px)`,
  opacity: phasePreset.value.opacity
}));

const deskStyle = computed(() => ({
  transform: `translate(${sceneParallax.x * 8}px, ${sceneParallax.y * 4}px) scaleX(${phasePreset.value.deskScale})`,
  opacity: phasePreset.value.deskOpacity
}));

function baseSlotTransform(startX, startY, extraX = 0, extraY = 0) {
  const phase = phaseTranslate(startX, startY);
  return `translate(${phase.x + extraX}px, ${phase.y + extraY}px)`;
}

const triangleSlotStyle = computed(() => ({
  transform: baseSlotTransform(-180, -130, focusedField.value === "password" ? -26 : 0, focusedField.value === "username" ? -4 : 2),
  opacity: focusedField.value === "password" ? 0.82 : phasePreset.value.opacity
}));

const squareSlotStyle = computed(() => ({
  transform: baseSlotTransform(0, -180, sceneParallax.x * 4, sceneParallax.y * 2),
  opacity: phasePreset.value.opacity
}));

const circleSlotStyle = computed(() => ({
  transform: baseSlotTransform(188, -84, focusedField.value === "password" ? -48 : 0, focusedField.value === "password" ? 10 : 0),
  opacity: focusedField.value === "password" ? 0.68 : phasePreset.value.opacity
}));

const blobSlotStyle = computed(() => ({
  transform: baseSlotTransform(
    236,
    -188,
    focusedField.value === "password" ? (mode.value === "login" ? -64 : 20) : sceneParallax.x * 5,
    focusedField.value === "password" ? 18 : sceneParallax.y * -4
  ),
  opacity: phasePreset.value.opacity
}));

const miniSlotStyle = computed(() => ({
  transform: baseSlotTransform(-126, 136, focusedField.value === "password" ? -22 : 0, focusedField.value === "password" ? 12 : 0),
  opacity: focusedField.value === "password" ? 0.74 : phasePreset.value.opacity
}));

const triangleInnerStyle = computed(() => ({
  transform: `translate(${sceneParallax.x * -3}px, ${sceneParallax.y * -2}px) rotate(${look.x * 4}deg)`
}));

const squareInnerStyle = computed(() => {
  const rotate = focusedField.value === "password" ? (mode.value === "login" ? -12 : 12) : look.x * 2;
  const introRotate =
    introPhase.value === "switch-scatter" ? 16 :
    introPhase.value === "scatter" ? 22 :
    introPhase.value === "forming" ? 8 : 0;
  const scale =
    introPhase.value === "scatter" ? 0.72 :
    introPhase.value === "forming" ? 0.9 :
    introPhase.value === "switch-scatter" ? 0.84 :
    introPhase.value === "switch-forming" ? 0.95 : 1;
  return {
    transform: `translate(${sceneParallax.x * 2}px, ${sceneParallax.y * 1.5}px) rotate(${rotate + introRotate}deg) scale(${scale})`
  };
});

const circleInnerStyle = computed(() => {
  const scale =
    introPhase.value === "scatter" ? 0.62 :
    introPhase.value === "forming" ? 0.88 :
    introPhase.value === "switch-scatter" ? 0.78 :
    introPhase.value === "switch-forming" ? 0.94 : 1;
  return {
    transform: `translate(${sceneParallax.x * 3}px, ${sceneParallax.y * 1.5}px) scale(${scale})`
  };
});

const blobInnerStyle = computed(() => {
  const rotate =
    focusedField.value === "password"
      ? mode.value === "login" ? -22 : 18
      : look.x * 3;
  const introRotate =
    introPhase.value === "switch-scatter" ? 32 :
    introPhase.value === "scatter" ? 28 :
    introPhase.value === "forming" ? 12 : 0;
  return {
    transform: `translate(${sceneParallax.x * 4}px, ${sceneParallax.y * -2}px) rotate(${rotate + introRotate}deg)`
  };
});

const miniInnerStyle = computed(() => {
  const scale =
    introPhase.value === "scatter" ? 0.48 :
    introPhase.value === "forming" ? 0.84 :
    introPhase.value === "switch-scatter" ? 0.58 :
    introPhase.value === "switch-forming" ? 0.92 : 1;
  return {
    transform: `translate(${sceneParallax.x * -2}px, ${sceneParallax.y * 2.5}px) scale(${scale})`
  };
});

async function submit() {
  errorMessage.value = "";
  celebrating.value = false;
  try {
    if (mode.value === "login") {
      await userStore.loginAction(form);
    } else {
      await userStore.registerAction(form);
    }
    celebrating.value = true;
    successTimer = window.setTimeout(() => {
      sessionStorage.setItem("blog_entry_animation", "auth-success");
      router.push("/");
    }, 980);
  } catch (error) {
    errorMessage.value = error?.response?.data?.message || "操作失败，请稍后再试";
    celebrating.value = false;
  }
}

function scheduleBlink() {
  const delay = 2300 + Math.random() * 2200;
  blinkTimer = window.setTimeout(() => {
    blink.value = true;
    window.setTimeout(() => {
      blink.value = false;
      scheduleBlink();
    }, 170);
  }, delay);
}

onMounted(() => {
  const usernameEl = document.querySelector("input");
  if (usernameEl) {
    const rect = usernameEl.getBoundingClientRect();
    const x = ((rect.left + rect.width / 2) / window.innerWidth - 0.5) * 2;
    const y = ((rect.top + rect.height / 2) / window.innerHeight - 0.5) * 2;
    look.x = clamp(x, -1, 1);
    look.y = clamp(y, -1, 1);
  }

  introTimer1 = window.setTimeout(() => {
    introPhase.value = "forming";
  }, 900);

  introTimer2 = window.setTimeout(() => {
    introPhase.value = "formed";
    motionReady.value = true;
  }, 2850);

  scheduleBlink();
});

onBeforeUnmount(() => {
  window.clearTimeout(blinkTimer);
  window.clearTimeout(introTimer1);
  window.clearTimeout(introTimer2);
  clearSwitchTimers();
  window.clearTimeout(successTimer);
});
</script>

<style scoped>
.auth-stage {
  min-height: calc(100vh - 152px);
  margin-top: 28px;
}

.auth-brand-flash {
  position: fixed;
  inset: 0;
  z-index: 30;
  display: grid;
  place-items: center;
  background: radial-gradient(circle, rgba(255, 209, 102, 0.14), rgba(8, 11, 17, 0.82));
  animation: flashFade 0.98s ease forwards;
  pointer-events: none;
}

.brand-flash__core {
  display: grid;
  gap: 14px;
  justify-items: center;
  color: #fff7ed;
  animation: flashRise 0.98s ease forwards;
}

.brand-badge--flash {
  width: 88px;
  height: 88px;
  border-radius: 28px;
  font-size: 2rem;
}

.brand-flash__core strong {
  font-size: 1.3rem;
  letter-spacing: 0.08em;
}

.auth-lane {
  position: relative;
  min-height: 700px;
}

.auth-scene,
.auth-panel {
  position: absolute;
  top: 0;
  bottom: 0;
  transition:
    left 0.88s cubic-bezier(0.2, 0.88, 0.22, 1),
    right 0.88s cubic-bezier(0.2, 0.88, 0.22, 1),
    transform 0.88s cubic-bezier(0.2, 0.88, 0.22, 1),
    opacity 0.88s ease;
}

.auth-scene {
  left: 0;
  width: calc(100% - 404px);
  overflow: hidden;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at 18% 18%, rgba(255, 138, 76, 0.16), transparent 24%),
    radial-gradient(circle at 82% 22%, rgba(255, 209, 102, 0.12), transparent 22%),
    rgba(255, 255, 255, 0.06);
}

.auth-panel {
  right: 0;
  width: 380px;
  display: grid;
  align-content: center;
  gap: 24px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.04)),
    rgba(8, 11, 17, 0.52);
}

.auth-stage--register .auth-scene {
  left: 404px;
}

.auth-stage--register .auth-panel {
  right: auto;
  left: 0;
}

.auth-stage--switching .auth-scene {
  transform: scale(0.984) translateY(6px);
}

.auth-stage--switching .auth-panel {
  transform: scale(0.984) translateY(-6px);
}

.auth-figure-wrap {
  position: relative;
  width: 100%;
  min-height: 520px;
  display: grid;
  place-items: center;
}

.auth-figure-glow {
  position: absolute;
  inset: auto 14% 14% 14%;
  height: 140px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(255, 138, 76, 0.34), transparent 68%);
  filter: blur(22px);
  transition: opacity 0.8s ease;
}

.auth-stage--celebrating .auth-figure-glow {
  opacity: 1;
  background: radial-gradient(circle, rgba(255, 209, 102, 0.46), transparent 68%);
}

.auth-figure {
  width: min(100%, 560px);
  overflow: visible;
}

.ambient path {
  stroke: rgba(255, 224, 198, 0.26);
  stroke-width: 5;
  stroke-linecap: round;
  fill: none;
}

.ambient {
  animation: floatLoop 5s ease-in-out infinite;
}

.ambient--two {
  animation-delay: 1s;
}

.ambient--three {
  animation-delay: 1.8s;
}

.spark circle {
  fill: rgba(255, 234, 214, 0.7);
}

.spark {
  animation: pulseSpark 3.2s ease-in-out infinite;
}

.spark--two {
  animation-delay: 0.9s;
}

.spark--three {
  animation-delay: 1.7s;
}

.fragment {
  opacity: 0;
  transform-origin: center;
  transition: opacity 0.6s ease;
}

.fragment path,
.fragment circle {
  fill: rgba(255, 236, 218, 0.4);
  stroke: rgba(255, 241, 228, 0.4);
  stroke-width: 2;
}

.auth-stage--scatter .fragment,
.auth-stage--forming .fragment {
  opacity: 1;
}

.fragment--one {
  animation: fragmentOne 2.8s cubic-bezier(0.2, 0.88, 0.22, 1) forwards;
}

.fragment--two {
  animation: fragmentTwo 2.8s cubic-bezier(0.2, 0.88, 0.22, 1) forwards;
}

.fragment--three {
  animation: fragmentThree 2.8s cubic-bezier(0.2, 0.88, 0.22, 1) forwards;
}

.fragment--four {
  animation: fragmentFour 2.8s cubic-bezier(0.2, 0.88, 0.22, 1) forwards;
}

.desk-line,
.desk-leg,
.shape-line,
.mouth-line {
  fill: none;
  stroke: rgba(255, 241, 228, 0.92);
  stroke-linecap: round;
  stroke-linejoin: round;
}

.desk-line,
.desk-leg {
  stroke-width: 8;
  stroke: rgba(255, 224, 206, 0.74);
}

.shape-line {
  stroke-width: 7;
}

.mouth-line {
  stroke-width: 4.5;
}

.shape-fill {
  opacity: 0.95;
}

.triangle-fill {
  fill: rgba(255, 183, 120, 0.26);
}

.square-fill {
  fill: rgba(255, 214, 102, 0.22);
}

.circle-fill {
  fill: rgba(255, 138, 76, 0.2);
}

.blob-fill {
  fill: rgba(255, 165, 132, 0.22);
}

.mini-fill {
  fill: rgba(255, 224, 160, 0.2);
}

.eye-shell {
  fill: #fff7f0;
  stroke: rgba(255, 243, 234, 0.5);
  stroke-width: 2;
  transition: ry 0.14s ease;
}

.pupil {
  fill: #2a1d17;
  transition: r 0.14s ease;
}

.pupil-glint {
  fill: #ffffff;
}

.shape-team,
.shape-slot,
.desk-group {
  transition:
    transform 1.45s cubic-bezier(0.18, 0.9, 0.22, 1),
    opacity 1.2s ease;
}

.shape {
  transform-origin: center;
  transition: transform 0.35s ease;
}

.auth-stage--motion-ready .shape--triangle {
  animation: wobbleTriangle 4.8s ease-in-out infinite;
}

.auth-stage--motion-ready .shape--square {
  animation: wobbleSquare 5.2s ease-in-out infinite;
}

.auth-stage--motion-ready .shape--circle {
  animation: wobbleCircle 4.4s ease-in-out infinite;
}

.auth-stage--motion-ready .shape--blob {
  animation: wobbleBlob 5.4s ease-in-out infinite;
}

.auth-stage--motion-ready .shape--mini {
  animation: wobbleMini 3.8s ease-in-out infinite;
}

.auth-stage--celebrating .shape-team {
  animation: teamCelebrate 0.92s ease-in-out;
}

.auth-panel__head {
  display: grid;
  gap: 14px;
}

.auth-tabs {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.auth-submit {
  width: 100%;
  justify-content: center;
  min-height: 52px;
}

@keyframes wobbleTriangle {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  35% { transform: translateY(-6px) rotate(-2deg); }
  70% { transform: translateY(4px) rotate(2deg); }
}

@keyframes wobbleSquare {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  30% { transform: translateY(4px) rotate(-1.4deg); }
  65% { transform: translateY(-7px) rotate(1.8deg); }
}

@keyframes wobbleCircle {
  0%, 100% { transform: translateY(0px) scale(1); }
  50% { transform: translateY(-9px) scale(1.03); }
}

@keyframes wobbleBlob {
  0%, 100% { transform: translateY(0px) rotate(0deg) scale(1); }
  40% { transform: translateY(-5px) rotate(2deg) scale(1.02); }
  72% { transform: translateY(4px) rotate(-2deg) scale(0.99); }
}

@keyframes wobbleMini {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-8px); }
}

@keyframes pulseSpark {
  0%, 100% { opacity: 0.3; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.2); }
}

@keyframes floatLoop {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-8px); }
}

@keyframes fragmentOne {
  0% { transform: translate(-54px, -48px) rotate(-32deg) scale(0.58); opacity: 0; }
  24% { opacity: 0.9; }
  70% { transform: translate(36px, 16px) rotate(18deg) scale(1); opacity: 0.72; }
  100% { transform: translate(118px, 84px) rotate(0deg) scale(0.18); opacity: 0; }
}

@keyframes fragmentTwo {
  0% { transform: translate(60px, -52px) rotate(0deg) scale(0.5); opacity: 0; }
  22% { opacity: 0.95; }
  68% { transform: translate(-28px, 54px) rotate(160deg) scale(1.1); opacity: 0.68; }
  100% { transform: translate(-146px, 118px) rotate(280deg) scale(0.18); opacity: 0; }
}

@keyframes fragmentThree {
  0% { transform: translate(-22px, 30px) rotate(-12deg) scale(0.4); opacity: 0; }
  28% { opacity: 0.88; }
  72% { transform: translate(64px, -58px) rotate(120deg) scale(1.05); opacity: 0.72; }
  100% { transform: translate(98px, -132px) rotate(220deg) scale(0.16); opacity: 0; }
}

@keyframes fragmentFour {
  0% { transform: translate(26px, 18px) rotate(0deg) scale(0.52); opacity: 0; }
  30% { opacity: 0.92; }
  70% { transform: translate(-48px, -42px) rotate(-140deg) scale(1.08); opacity: 0.64; }
  100% { transform: translate(-166px, -118px) rotate(-220deg) scale(0.18); opacity: 0; }
}

@keyframes teamCelebrate {
  0% { transform: translate(0, 0) scale(1); }
  24% { transform: translate(0, -16px) scale(1.04); }
  48% { transform: translate(0, 4px) scale(0.99); }
  72% { transform: translate(0, -8px) scale(1.02); }
  100% { transform: translate(0, 0) scale(1); }
}

@keyframes flashFade {
  0% { opacity: 0; }
  20% { opacity: 1; }
  100% { opacity: 0; }
}

@keyframes flashRise {
  0% { transform: translateY(22px) scale(0.92); opacity: 0; }
  22% { transform: translateY(0px) scale(1.02); opacity: 1; }
  100% { transform: translateY(-10px) scale(1); opacity: 0; }
}

@media (max-width: 1080px) {
  .auth-lane {
    min-height: auto;
    display: grid;
    gap: 24px;
  }

  .auth-scene,
  .auth-panel,
  .auth-stage--register .auth-scene,
  .auth-stage--register .auth-panel {
    position: static;
    width: auto;
    left: auto;
    right: auto;
    transform: none;
  }
}

@media (max-width: 768px) {
  .auth-stage {
    margin-top: 14px;
  }

  .auth-panel {
    padding: 20px;
    gap: 18px;
  }
}

@media (max-width: 480px) {
  .auth-panel {
    border-radius: 24px;
  }

  .auth-figure-wrap {
    min-height: 320px;
  }

  .auth-figure {
    width: 100%;
  }

  .auth-submit {
    min-height: 46px;
    font-size: 0.95rem;
  }
}
</style>

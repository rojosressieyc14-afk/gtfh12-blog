<template>
  <section id="users" class="manage-panel">
    <div class="review-head review-head--stack">
      <div>
        <p class="admin-label">用户管理</p>
        <h2>用户与封禁状态</h2>
      </div>
      <div class="toolbar-row">
        <input
          v-model.trim="userKeyword"
          class="filter-select filter-input"
          placeholder="搜索用户名 / 角色 / 状态"
          @keyup.enter="changeUserPage(1)"
        />
        <select v-model="userRoleFilter" class="filter-select" @change="changeUserPage(1)">
          <option value="">全部角色</option>
          <option value="admin">管理员</option>
          <option value="user">普通用户</option>
        </select>
        <button class="role-btn" @click="changeUserPage(1)">搜索</button>
      </div>
    </div>

    <div class="table-list">
      <article v-for="item in users" :key="item.id" class="table-card table-card--article">
        <div class="table-main">
          <h3>{{ item.username }}</h3>
          <p>{{ roleLabel(item.role) }} · {{ userStatusLabel(item.status) }}</p>
          <p v-if="item.banReason" class="reject-note">封禁原因：{{ item.banReason }}</p>
        </div>
        <div class="table-actions">
          <button class="role-btn" @click="toggleRole(item)">
            {{ item.role === "admin" ? "设为普通用户" : "设为管理员" }}
          </button>
          <button class="role-btn" @click="toggleStatus(item)">
            {{ item.status === "banned" ? "解除封禁" : "封禁用户" }}
          </button>
        </div>
      </article>
    </div>

    <div class="pager-row">
      <button class="role-btn" :disabled="userPage <= 1" @click="changeUserPage(userPage - 1)">上一页</button>
      <span>第 {{ userPage }} / {{ userTotalPages }} 页</span>
      <button class="role-btn" :disabled="userPage >= userTotalPages" @click="changeUserPage(userPage + 1)">下一页</button>
    </div>
  </section>
</template>

<script setup>
import { computed, ref } from "vue";
import { getUsers, updateUserRole, updateUserStatus } from "../../api/dashboard";

const emit = defineEmits(["flash"]);

const users = ref([]);
const userKeyword = ref("");
const userRoleFilter = ref("");
const userPage = ref(1);
const userTotal = ref(0);
const pageSize = 6;

function say(message) {
  emit("flash", message);
}

function roleLabel(role) {
  return role === "admin" ? "管理员" : "普通用户";
}

function userStatusLabel(status) {
  return status === "banned" ? "已封禁" : "正常";
}

const userTotalPages = computed(() => Math.max(1, Math.ceil(userTotal.value / pageSize)));

async function loadUsers() {
  const { data } = await getUsers({
    page: userPage.value,
    pageSize,
    keyword: userKeyword.value,
    role: userRoleFilter.value
  });
  users.value = data.items || [];
  userTotal.value = data.pagination?.total || 0;
}

async function changeUserPage(nextPage) {
  userPage.value = nextPage;
  await loadUsers();
}

async function toggleRole(item) {
  try {
    await updateUserRole(item.id, { role: item.role === "admin" ? "user" : "admin" });
    await loadUsers();
    say("用户角色已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新用户角色失败。");
  }
}

async function toggleStatus(item) {
  try {
    await updateUserStatus(item.id, { status: item.status === "banned" ? "active" : "banned" });
    await loadUsers();
    say("用户状态已更新。");
  } catch (error) {
    say(error?.response?.data?.message || error?.message || "更新用户状态失败。");
  }
}

loadUsers();
</script>

<style scoped>
.review-head--stack {
  align-items: flex-start;
}

.reject-note {
  color: #ffd2d2;
}

@media (max-width: 768px) {
  .toolbar-row {
    flex-direction: column;
  }
  .filter-input {
    min-width: 0;
    width: 100%;
  }
}

@media (max-width: 480px) {
  .table-card {
    flex-direction: column;
    align-items: flex-start;
  }
  .table-actions {
    flex-direction: column;
    width: 100%;
  }
  .table-actions button {
    width: 100%;
  }
  .table-main {
    width: 100%;
  }
}
</style>

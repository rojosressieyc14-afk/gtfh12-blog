# Login 记住账号密码（Remember Me + 用户名记忆）

Date: 2026-08-14
Status: Draft
Author: opencode

## Overview

登录页新增"记住账号密码"勾选：

- 勾选后：登录 cookie 有效期从 7 天延长到 30 天（关浏览器后仍保持登录），并把用户名存入 localStorage（下次自动填充，密码不落盘）。
- 不勾选：登录 cookie 为会话级（关浏览器即失效），不保存用户名。
- 复选框初始未勾选；若 localStorage 已有记忆的用户名，则预填充用户名并自动勾选。

不涉及 admin 端（YAGNI，admin 登录页维持现状）。

## 相关文件

- `server/internal/middleware/auth.go`
- `server/internal/handler/auth_handler.go`
- `server/internal/handler/handler_test.go`
- `web/src/views/AuthView.vue`

## Tasks

### T1: middleware/auth.go — cookie 有效期参数化

**改动**

- 常量 `authCookieMaxAge` 更名为导出的 `AuthCookieMaxAge`（7 天），新增 `RememberCookieMaxAge = 30 * 24 * 3600`。
- `SetSessionCookies(c *gin.Context, token string)` → `SetSessionCookies(c *gin.Context, token string, maxAge int)`，auth 与 csrf 两个 cookie 都使用传入的 maxAge。
- 注意：gin `c.SetCookie` 的 maxAge 为 0 时，Go http 包不写 Max-Age/Expires 属性 → 会话级 cookie（关浏览器失效），这正是"不勾选"的语义。

```go
const (
	AuthCookieName       = "blog_token"
	CSRFCookieName       = "blog_csrf"
	AuthCookieMaxAge     = 7 * 24 * 3600
	RememberCookieMaxAge = 30 * 24 * 3600
)
```

```go
func SetSessionCookies(c *gin.Context, token string, maxAge int) {
	secure := secureRequest(c)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AuthCookieName, token, maxAge, "/", "", secure, true)
	c.SetSameSite(http.SameSiteLaxMode)
	csrf := randomToken()
	c.SetCookie(CSRFCookieName, csrf, maxAge, "/", "", secure, false)
}
```

其余（`CSRF()` 中间件、`ClearSessionCookies`、`RequireAuth`）不变。

**验收**

- `SetSessionCookies` 所有调用点已同步新签名（全项目仅 `auth_handler.go` 两处）。
- maxAge=0 时两个 cookie 均无 Max-Age 属性。

### T2: auth_handler.go — Login 根据 remember 传 maxAge

```go
type authPayload struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Remember bool   `json:"remember"`
}
```

- `Register`：`middleware.SetSessionCookies(c, token, middleware.AuthCookieMaxAge)`（行为不变，保持 7 天）。
- `Login`：

```go
	maxAge := 0
	if payload.Remember {
		maxAge = middleware.RememberCookieMaxAge
	}
	middleware.SetSessionCookies(c, token, maxAge)
```

**验收**

- 登录不带 remember / remember=false → 会话级 cookie；remember=true → 30 天。
- 注册仍为 7 天。

### T3: handler_test.go — 修复过期断言并覆盖 remember

**背景**：Phase 3 起响应体不再含 `token`，但 `TestRegisterLoginMeFlow` 第 171 行仍 `result["token"].(string)`，运行时必 panic。本次一并修复为 cookie 断言，并新增 remember 分支测试。

修复 `TestRegisterLoginMeFlow` 的 register 段：

```go
	result := parseResponse(t, w)
	if _, ok := result["user"]; !ok {
		t.Fatalf("register: expected user in response, got %v", result)
	}
	cookies := w.Result().Cookies()
	var tokenCookie *http.Cookie
	for _, ck := range cookies {
		if ck.Name == middleware.AuthCookieName {
			tokenCookie = ck
			break
		}
	}
	if tokenCookie == nil || tokenCookie.Value == "" {
		t.Fatal("register: expected blog_token cookie")
	}
```

（需要新增 import `"blog/server/internal/middleware"`。）

新增测试：

```go
func TestLoginRememberCookie(t *testing.T) {
	r, _ := setupTest(t)

	seq := atomic.AddUint64(&testSeq, 1)
	username := fmt.Sprintf("remember_%d", seq)
	password := "TestPass123!"

	w := requestJSON(r, "POST", "/api/auth/register", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// 不勾选 -> 会话级 cookie（MaxAge == 0）
	w = requestJSON(r, "POST", "/api/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	if w.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	authCookie := findCookie(t, w, middleware.AuthCookieName)
	if authCookie.MaxAge != 0 {
		t.Fatalf("login without remember: expected session cookie (MaxAge 0), got %d", authCookie.MaxAge)
	}

	// 勾选 -> 30 天
	w = requestJSON(r, "POST", "/api/auth/login", fmt.Sprintf(`{"username":"%s","password":"%s","remember":true}`, username, password))
	if w.Code != http.StatusOK {
		t.Fatalf("login remember: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	authCookie = findCookie(t, w, middleware.AuthCookieName)
	if authCookie.MaxAge != middleware.RememberCookieMaxAge {
		t.Fatalf("login with remember: expected MaxAge %d, got %d", middleware.RememberCookieMaxAge, authCookie.MaxAge)
	}
}

func findCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, ck := range w.Result().Cookies() {
		if ck.Name == name {
			return ck
		}
	}
	t.Fatalf("expected cookie %q", name)
	return nil
}
```

**验收**

- `go test ./internal/handler/...` 通过（需要本地 MySQL `blog_test` 库；若数据库不可用，则记录说明并以 `go build` + `go vet` 兜底，不做其他绕过）。
- 旧断言不再引用 `result["token"]`。

### T4: AuthView.vue — 复选框 + 用户名记忆

**模板**（密码 label 与 `errorMessage` 之间，仅登录模式显示）：

```html
<label v-if="mode === 'login'" class="remember-row">
  <input v-model="form.remember" type="checkbox" class="remember-check" />
  <span>记住账号密码</span>
</label>
```

**script**：

- `form` 增加字段：`remember: false`
- 常量：`const REMEMBER_USERNAME_KEY = "blog_remembered_username";`

预填充（放入已有/新建的 `onMounted`）：

```js
const savedUsername = localStorage.getItem(REMEMBER_USERNAME_KEY);
if (savedUsername) {
  form.username = savedUsername;
  form.remember = true;
}
```

`submit()` 登录成功分支（`await userStore.loginAction(form);` 之后追加）：

```js
if (form.remember) {
  localStorage.setItem(REMEMBER_USERNAME_KEY, form.username);
} else {
  localStorage.removeItem(REMEMBER_USERNAME_KEY);
}
```

**scoped style 追加**：

```css
.remember-row {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--text-soft);
}
```

**验收**

- 勾选登录成功 → cookie 有效 30 天；刷新登录页用户名已填充且复选框勾选。
- 取消勾选登录成功 → 关浏览器重开后需重新登录；此前保存的用户名被清除。
- 注册流程不受影响（复选框仅登录模式出现）。

### T5: 验证

- `go build .\cmd\api\`、`go vet ./...`
- `go test ./internal/handler/...`（MySQL 可用时）
- `cd web; npx vite build`
- 手动冒烟：勾选/不勾选两条路径 + 刷新登录页预填充

### T6: 补充 —— remember 时 JWT 有效期延长到 30 天

**背景**：最终审查发现 `utils.GenerateJWT` 硬编码 24h 有效期（jwt.go:81），即使 cookie 30 天，登录态 24h 后即 401，30 天承诺不成立。用户已确认修复。

- `utils.GenerateJWT(userID uint, username, role string)` → `GenerateJWT(userID uint, username, role string, ttl time.Duration)`；新增常量 `TokenTTLDefault = 24 * time.Hour`、`TokenTTLRemember = 30 * 24 * time.Hour`。
- `AuthService.Login(username, password string, remember bool)`：remember → `TokenTTLRemember`，否则 `TokenTTLDefault`。
- `AuthService.Register`：`TokenTTLDefault`（行为不变）。
- `auth_handler.go` Login 传递 `payload.Remember`。
- `handler_test.go`：两处 `GenerateJWT` 补 `utils.TokenTTLDefault`；`TestLoginRememberCookie` 增加断言——remember 登录后解析 cookie token，剩余有效期 > 29 天。

**验收**

- remember 登录的 JWT `ExpiresAt` ≈ now+30d；不勾选仍 24h。
- 全部测试通过。

## Suggested commits

1. `feat(server): support remember-me cookie max-age in auth`（T1+T2）
2. `fix(server): update auth handler tests for cookie flow + remember coverage`（T3）
3. `feat(web): add remember-me checkbox with username autofill`（T4）

## Exit criteria

- 上述所有验收项通过；三端构建无回归；`blog_token` 在两种模式下有效期符合预期；用户名仅存 localStorage，密码不进任何本地存储。
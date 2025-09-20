document.addEventListener("DOMContentLoaded", () => {
  // ===== Валидация (как договорились: тихо, без алертов) =====
  const form = document.getElementById("login-form");
  const email = document.getElementById("email");
  const password = document.getElementById("password");
  const submitBtn = document.getElementById("submit-btn");
  const emailErr = document.getElementById("email-error");
  const passwordErr = document.getElementById("password-error");

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  const passwordRegex = /^.{6,}$/; // минимум 6 любых символов

  const touched = { email: false, password: false };

  const isEmailValid = v => emailRegex.test(v.trim());
  const isPasswordValid = v => passwordRegex.test(v.trim());

  function renderFieldState(input, errorEl, valid, message, showMessage) {
    if (valid) {
      input.classList.remove("border-red-500");
      errorEl.classList.add("hidden");
      errorEl.textContent = "";
    } else if (showMessage) {
      input.classList.add("border-red-500");
      errorEl.classList.remove("hidden");
      errorEl.textContent = message;
    } else {
      input.classList.remove("border-red-500");
      errorEl.classList.add("hidden");
      errorEl.textContent = "";
    }
  }
  function validateEmail(showMsg) {
    const ok = isEmailValid(email.value);
    renderFieldState(email, emailErr, ok, "Введите корректный e-mail.", showMsg);
    return ok;
  }
  function validatePassword(showMsg) {
    const ok = isPasswordValid(password.value);
    renderFieldState(password, passwordErr, ok, "Пароль должен быть минимум 6 символов.", showMsg);
    return ok;
  }
  function updateSubmitState() {
    submitBtn.disabled = !(isEmailValid(email.value) && isPasswordValid(password.value));
  }

  email.addEventListener("input", () => { if (touched.email) validateEmail(true); updateSubmitState(); });
  password.addEventListener("input", () => { if (touched.password) validatePassword(true); updateSubmitState(); });
  email.addEventListener("blur", () => { touched.email = true; validateEmail(true); updateSubmitState(); });
  password.addEventListener("blur", () => { touched.password = true; validatePassword(true); updateSubmitState(); });

  // ===== Универсальный POST JSON =====
  async function postJSON(url, data, opts = {}) {
    const { headers = {}, signal, credentials = "same-origin" } = opts;
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...headers },
      body: JSON.stringify(data ?? {}),
      credentials,
      signal
    });

    // Пытаемся распарсить JSON всегда (часто даже при ошибке приходит JSON)
    let payload = null;
    try { payload = await res.json(); } catch (_) { /* игнор */ }

    if (!res.ok) {
      const msg = payload?.message || res.statusText || "Ошибка запроса";
      const err = new Error(msg);
      err.status = res.status;
      err.payload = payload;
      throw err;
    }
    return payload; // успешный JSON
  }

  // ===== UI: message box в правом нижнем углу на 3 сек =====
  const errBox = document.getElementById("signin-error");
  const errMsg = document.getElementById("signin-error-text");
  let hideTimer = null;

  function showErrorBox(message) {
    if (hideTimer) clearTimeout(hideTimer);
    errMsg.textContent = message || "Неизвестная ошибка";
    errBox.classList.remove("hidden", "opacity-0");
    errBox.classList.add("opacity-100");
    hideTimer = setTimeout(hideErrorBox, 3000);
  }
  function hideErrorBox() {
    errBox.classList.add("opacity-0");
    // мягко скрываем (совпадает с duration-300)
    setTimeout(() => errBox.classList.add("hidden"), 300);
  }

  // ===== Вход через универсальный POST =====
  async function signIn({ email, password }) {
    return postJSON("/api/v1/auth/sign-in", { email, password });
  }

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    touched.email = true;
    touched.password = true;

    const okEmail = validateEmail(true);
    const okPass = validatePassword(true);
    updateSubmitState();
    if (!okEmail || !okPass) return;

    // блокируем кнопку на время запроса (минимальный UX)
    submitBtn.disabled = true;
    const originalText = submitBtn.textContent;
    submitBtn.textContent = "Входим…";

    try {
      await signIn({ email: email.value.trim(), password: password.value });
      // успех -> редирект
      window.location.href = "/dashboard";
    } catch (err) {
      // ошибка -> показать message box
      showErrorBox(err.message || "Ошибка входа");
      submitBtn.disabled = false;
      submitBtn.textContent = originalText;
    }
  });

  updateSubmitState();
});
function toggleMobileMenu() {
  const sidebar = document.getElementById('mobile-sidebar');
  const overlay = document.getElementById('mobile-overlay');
  if (sidebar && overlay) {
    sidebar.classList.toggle('-translate-x-full');
    overlay.classList.toggle('hidden');
    document.body.classList.toggle('overflow-hidden');
  }
}
function closeMobileMenu() {
  const sidebar = document.getElementById('mobile-sidebar');
  const overlay = document.getElementById('mobile-overlay');
  if (sidebar && overlay) {
    sidebar.classList.add('-translate-x-full');
    overlay.classList.add('hidden');
    document.body.classList.remove('overflow-hidden');
  }
}

document.addEventListener("DOMContentLoaded", async () => {
  const btn = document.getElementById("profileBtn");
  const dropdown = document.getElementById("profileDropdown");
  btn.addEventListener("click", () => dropdown.classList.toggle("hidden"));

  try {
    const res = await fetch("/api/v1/data/user/data", { credentials: "include" });
    if (!res.ok) throw new Error("HTTP " + res.status);
      const user = await res.json();
      window.__USER_DATA = user; // кэшируем, чтобы dashboard мог использовать без повторного запроса    document.getElementById("userName").textContent = user.name || "No Data";
      document.getElementById("userName").textContent = user.name || "No Data";
      document.getElementById("userRole").textContent = user.role?.name || "No Data";
    document.getElementById("dropdownFullName").textContent = (user.name || "No Data") + " " + (user.surname || "");
    document.getElementById("dropdownEmail").textContent = user.email || "No Data";
    document.getElementById("dropdownPhone").textContent = user.phone || "No Data";
    document.getElementById("dropdownRole").textContent = user.role?.desc || "No Data";
  } catch (e) {
    console.error("Failed to load user data:", e);
  }

  document.getElementById("logoutForm").addEventListener("submit", async (ev) => {
    ev.preventDefault();
    await fetch("/api/v1/auth/sign-out", { method: "POST", credentials: "include" });
    window.location.href = "/login";
  });
});
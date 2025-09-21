    document.addEventListener("DOMContentLoaded", async () => {
      const nameEl = document.getElementById("greetName");
      if (!nameEl) return;

      const applyName = (user) => {
        nameEl.textContent = (user?.name) ? user.name : "No Data";
      };

      // Use cached user data if available
      if (window.__USER_DATA) {
        applyName(window.__USER_DATA);
        return;
      }

      // Otherwise fetch user data
      try {
        const response = await fetch("/api/v1/data/user/data", { credentials: "include" });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const user = await response.json();
        applyName(user);
      } catch (error) {
        console.warn("Failed to load user for greeting:", error);
        applyName(null);
      }
    });
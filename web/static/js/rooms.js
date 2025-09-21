    const DEFAULT_BG = "https://images.unsplash.com/photo-1562774053-701939374585?w=1200&h=800&fit=crop&crop=center";

    const STATUS_MAP = {
      "active": "Активен",
      "disabled": "Выключен",
      "inactive": "Неактивен"
    };

    function cardHTML(room) {
      const id = room?.id || "";
      const number = room?.number ?? null;
      const name = room?.name || "Кабинет";
      const title = number ? `${number} Кабинет` : name;
      const status = STATUS_MAP[(room?.status || "").toLowerCase()] || "Неизвестно";

      // ПК-статистики в ответе нет — ставим прочерки
      const pcsOnline = "—";
      const pcsTotal = "—";

      return `
        <div class="hover-card room-card rounded-3xl shadow-xl overflow-hidden"
             style="background-image: url('${DEFAULT_BG}')"
             onclick="window.location.href='/rooms/${id}'">
          <div class="room-gradient absolute inset-0 flex flex-col justify-end p-6 lg:p-8">
            <div class="text-white">
              <h3 class="text-[26px] lg:text-[30px] font-bold leading-none mb-3">${title}</h3>
              <div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
                <div class="flex flex-col space-y-1 text-sm lg:text-base opacity-95">
                  <span class="font-medium">ПК: ${pcsOnline} онлайн / ${pcsTotal} всего</span>
                  <span>Статус: ${status}</span>
                </div>
                <button class="bg-white text-[#0085FF] font-semibold px-5 py-2.5 rounded-full text-sm lg:text-base hover:bg-opacity-95 transition-all flex items-center gap-2 w-fit shadow-lg">
                  Перейти
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>`;
    }

    function attachHoverEffects(scopeEl) {
      scopeEl.querySelectorAll('.room-card').forEach(card => {
        card.addEventListener('mouseenter', () => {
          card.style.transform = 'translateY(-6px) scale(1.02)';
          card.style.boxShadow = '0 25px 50px rgba(0, 133, 255, 0.2)';
        });
        card.addEventListener('mouseleave', () => {
          card.style.transform = 'translateY(-4px) scale(1)';
          card.style.boxShadow = '0 20px 40px rgba(0, 0, 0, 0.15)';
        });
      });
    }

    async function renderRooms() {
      const grid = document.getElementById("rooms-grid");
      try {
        const res = await fetch("/api/v1/data/room", { credentials: "include" });
        if (!res.ok) throw new Error("HTTP " + res.status);
        const rooms = await res.json();

        if (!Array.isArray(rooms) || rooms.length === 0) {
          grid.innerHTML = `<div class="col-span-full text-center text-gray-500 py-10">Данных нет</div>`;
          return;
        }

        grid.innerHTML = rooms.map(cardHTML).join("");
        attachHoverEffects(grid);
      } catch (e) {
        console.error("Failed to load rooms:", e);
        grid.innerHTML = `<div class="col-span-full text-center text-gray-500 py-10">Ошибка загрузки. Данных нет</div>`;
      }
    }

    document.addEventListener("DOMContentLoaded", renderRooms);

    // Пагинация остаётся декоративной (как просил)
    let currentPage = 1;
    const totalPages = 12;
    const prevBtn = document.querySelector('.pagination-btn:first-child');
    const nextBtn = document.querySelector('.pagination-btn:last-child');
    const pageInfo = document.querySelector('span');

    function updatePagination() {
      if (!pageInfo) return;
      pageInfo.textContent = `Страница ${currentPage} из ${totalPages}`;
      if (prevBtn) prevBtn.disabled = currentPage === 1;
      if (nextBtn) nextBtn.disabled = currentPage === totalPages;
    }
    updatePagination();
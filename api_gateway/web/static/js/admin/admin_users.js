document.addEventListener('DOMContentLoaded', () => {
    const tableCtn    = document.getElementById('usersTableContainer');
    const editModal   = document.getElementById('editUserModal');
    const addModal    = document.getElementById('addUserModal');
    const closeBtns   = document.querySelectorAll('.close-modal');
    const cancelEdit  = document.getElementById('cancelEdit');
    const cancelAdd   = document.getElementById('cancelAdd');
    const saveBtn     = document.getElementById('saveUser');
    const createBtn   = document.getElementById('createUser');
    const addBtn      = document.getElementById('addUserBtn');
    const applyBtn    = document.getElementById('applyFilters');
    const searchInput = document.getElementById('searchInput');

    let users = [];
    let currentUserId = null;

    const API = {
        list:   '/api/v1/admin/users',
        update: id => `/api/v1/admin/user/update/${id}`,
        create: '/api/v1/admin/user/create'
    };

    async function fetchUsers() {
        try {
            const res = await fetch(API.list);
            users = await res.json();
        } catch {
            users = [];
        }
    }

    function render(list = users) {
        if (!list.length) {
            tableCtn.innerHTML = `<div class="error-message"><i class="fas fa-exclamation-circle"></i><p>Пользователи не найдены</p></div>`;
            return;
        }
        const roles = { admin: 'Администратор', creator: 'Продавец', user: 'Пользователь' };
        let html = `<table>
      <thead><tr>
        <th>ID</th><th>Пользователь</th><th>Email</th><th>Роль</th><th>Дата регистрации</th><th>Действия</th>
      </tr></thead><tbody>`;
        list.forEach(u => {
            const date = new Date(u.created_at).toLocaleDateString('ru-RU');
            html += `<tr>
        <td>#${u.user_id}</td>
        <td>
          <div style="display:flex;align-items:center;gap:10px;">
            ${u.name}
          </div>
        </td>
        <td>${u.email}</td>
        <td><span class="role-badge role-${u.role}">${roles[u.role]||u.role}</span></td>
        <td>${date}</td>
        <td>
          <button class="action-btn edit" data-id="${u.user_id}"><i class="fas fa-edit"></i></button>
        </td>
      </tr>`;
        });
        html += `</tbody></table>`;
        tableCtn.innerHTML = html;
        bindButtons();
    }

    function bindButtons() {
        document.querySelectorAll('.action-btn.edit').forEach(btn => {
            btn.onclick = () => openEdit(+btn.dataset.id);
        });
    }

    function openEdit(id) {
        currentUserId = id;
        const u = users.find(x => x.user_id === id);
        if (!u) return;
        document.getElementById('userId').value   = u.user_id;
        document.getElementById('userName').value = u.name;
        document.getElementById('userEmail').value= u.email;
        document.getElementById('userRole').value = u.role;
        editModal.style.display = 'flex';
    }

    function closeModals() {
        editModal.style.display = 'none';
        addModal.style.display  = 'none';
        currentUserId = null;
    }

    function filterAndRender() {
        const role = document.getElementById('roleFilter').value;
        const term = searchInput.value.trim().toLowerCase();
        let list = [...users];
        if (role) list = list.filter(u => u.role === role);
        if (term) list = list.filter(u =>
            u.name.toLowerCase().includes(term) ||
            u.email.toLowerCase().includes(term)
        );
        render(list);
    }

    saveBtn.onclick = async () => {
        const newRole = document.getElementById('userRole').value;
        try {
            const res = await fetch(API.update(currentUserId), {
                method: 'PUT',
                headers: {'Content-Type':'application/json'},
                body: JSON.stringify({ role: newRole })
            });
            if (!res.ok) throw new Error();
            await reload();
            closeModals();
        } catch {
            alert('Ошибка при сохранении');
        }
    };

    addBtn.onclick = () => {
        document.getElementById('addUserForm').reset();
        addModal.style.display = 'flex';
    };

    createBtn.onclick = async () => {
        const f = document.getElementById('addUserForm');
        if (!f.checkValidity()) return f.reportValidity();
        const data = {
            name:     document.getElementById('newUserName').value,
            email:    document.getElementById('newUserEmail').value,
            role:     document.getElementById('newUserRole').value,
            password: document.getElementById('newUserPassword').value
        };
        try {
            const res = await fetch(API.create, {
                method: 'POST',
                headers: {'Content-Type':'application/json'},
                body: JSON.stringify(data)
            });
            if (!res.ok) throw new Error();
            await reload();
            closeModals();
        } catch {
            alert('Ошибка при создании');
        }
    };

    closeBtns.forEach(b => b.onclick = closeModals);
    cancelEdit.onclick = cancelAdd.onclick = closeModals;
    applyBtn.onclick   = filterAndRender;
    searchInput.oninput= filterAndRender;
    window.onclick     = e => { if (e.target===editModal||e.target===addModal) closeModals(); };

    async function reload() {
        await fetchUsers();
        render();
    }

    reload();
});

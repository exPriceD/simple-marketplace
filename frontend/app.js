const apiUrl = '';

function show(element) { element.style.display = ''; }
function hide(element) { element.style.display = 'none'; }

function setUser(username) {
  hide(document.getElementById('auth-section'));
  show(document.getElementById('user-section'));
  show(document.getElementById('filters-section'));
  document.getElementById('user-info').textContent = username;
}

function clearUser() {
  show(document.getElementById('auth-section'));
  hide(document.getElementById('user-section'));
  hide(document.getElementById('filters-section'));
  document.getElementById('user-info').textContent = '';
}

function getToken() {
  return localStorage.getItem('access_token');
}

function setToken(token) {
  localStorage.setItem('access_token', token);
}

function removeToken() {
  localStorage.removeItem('access_token');
}

async function apiRequest(path, options = {}) {
  options.headers = options.headers || {};
  if (getToken()) {
    options.headers['Authorization'] = 'Bearer ' + getToken();
  }
  const res = await fetch(apiUrl + path, options);
  if (!res.ok) throw await res.json();
  return res.json();
}

document.getElementById('login-form').onsubmit = async e => {
  e.preventDefault();
  const username = document.getElementById('login-username').value;
  const password = document.getElementById('login-password').value;
  try {
    const data = await apiRequest('/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: username, password })
    });
    setToken(data.access_token);
    setUser(username);
    loadListings();
  } catch (err) {
    alert('Ошибка входа');
  }
};

document.getElementById('register-form').onsubmit = async e => {
  e.preventDefault();
  const username = document.getElementById('register-username').value;
  const password = document.getElementById('register-password').value;
  try {
    await apiRequest('/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login: username, password })
    });
    alert('Регистрация успешна, войдите!');
  } catch (err) {
    alert('Ошибка регистрации');
  }
};

document.getElementById('logout-btn').onclick = () => {
  removeToken();
  clearUser();
  loadListings();
};

document.getElementById('listing-form').onsubmit = async e => {
  e.preventDefault();
  const title = document.getElementById('listing-title').value;
  const image_url = document.getElementById('listing-image').value;
  const price = Number(document.getElementById('listing-price').value);
  const description = document.getElementById('listing-description').value;
  try {
    await apiRequest('/listings', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, image_url, price, description })
    });
    document.getElementById('listing-form').reset();
    loadListings();
  } catch (err) {
    alert('Ошибка размещения объявления');
  }
};

document.getElementById('apply-filters').onclick = () => {
  loadListings(0);
};

let currentOffset = 0;
let limit = 10;
let total = 0;

async function loadListings(offset = 0) {
  currentOffset = offset;
  let url = '/listings?limit=' + limit + '&offset=' + offset;
  const min = document.getElementById('filter-min-price').value;
  const max = document.getElementById('filter-max-price').value;
  const sort = document.getElementById('sort-type').value;
  const order = document.getElementById('sort-order').value;
  if (min) url += '&price_min=' + min;
  if (max) url += '&price_max=' + max;
  if (sort) url += '&sort_by=' + sort;
  if (order) url += '&sort_dir=' + order;
  try {
    const data = await apiRequest(url);
    renderListings(data.items || []);
    total = data.total || 0;
    renderPagination();
  } catch (err) {
    renderListings([]);
  }
}

function renderListings(listings) {
  const el = document.getElementById('listings');
  el.innerHTML = '';
  listings.forEach(listing => {
    const div = document.createElement('div');
    div.className = 'listing';
    div.innerHTML =
      (listing.image_url ? `<img src="${listing.image_url}"><br>` : '') +
      `<b>${listing.title}</b> <span>${listing.price} ₽</span><br>` +
      `<span>${listing.description}</span><br>` +
      `<span>Автор: ${listing.author_login || ''}</span>` +
      (listing.is_owner ? ' <span style="color:green;">(Ваше)</span>' : '');
    el.appendChild(div);
  });
}

function renderPagination() {
  const el = document.getElementById('pagination');
  el.innerHTML = '';
  const pages = Math.ceil(total / limit);
  for (let i = 0; i < pages; i++) {
    const btn = document.createElement('button');
    btn.textContent = (i + 1);
    if (i * limit === currentOffset) btn.disabled = true;
    btn.onclick = () => loadListings(i * limit);
    el.appendChild(btn);
  }
}

function checkAuth() {
  const token = getToken();
  if (token) {
    let payload = {};
    try {
      payload = JSON.parse(atob(token.split('.')[1] || 'e30='));
    } catch {}
    setUser(payload.login || 'Пользователь');
    show(document.getElementById('filters-section'));
  } else {
    clearUser();
  }
}

function randomListing(i) {
  const titles = [
    'Телефон', 'Ноутбук', 'Книга', 'Велосипед', 'Кофеварка', 'Стол', 'Кресло', 'Гитара', 'Часы', 'Кроссовки',
    'Пылесос', 'Монитор', 'Рюкзак', 'Камера', 'Принтер', 'Планшет', 'Микрофон', 'Колонка', 'Мышь', 'Клавиатура'
  ];
  const descs = [
    'Почти новый', 'В отличном состоянии', 'Работает идеально', 'Без царапин', 'С гарантией', 'Редко использовался',
    'Оригинал', 'С документами', 'В комплекте зарядка', 'Быстрая доставка', 'Торг уместен', 'Самовывоз', 'Доставка возможна',
    'Отличный подарок', 'Для дома и офиса', 'Стильный дизайн', 'Компактный', 'Мощный', 'Лёгкий', 'Удобный'
  ];
  const images = [
    'https://placehold.co/300x200?text=Фото',
    'https://placehold.co/300x200?text=Image',
    'https://placehold.co/300x200?text=Pic',
    'https://placehold.co/300x200?text=Photo',
    'https://placehold.co/300x200?text=Goods'
  ];
  const t = titles[Math.floor(Math.random()*titles.length)] + ' ' + (Math.floor(Math.random()*1000)+i);
  const d = descs[Math.floor(Math.random()*descs.length)] + '. ' + descs[Math.floor(Math.random()*descs.length)] + '.';
  const img = images[Math.floor(Math.random()*images.length)];
  const price = Math.floor(Math.random()*90000)+1000;
  return { title: t, image_url: img, price, description: d };
}

const bulkBtn = document.getElementById('bulk-create-btn');
if (bulkBtn) {
  bulkBtn.onclick = async () => {
    bulkBtn.disabled = true;
    for (let i = 0; i < 10; i++) {
      const l = randomListing(i + Math.floor(Math.random()*10000));
      try {
        await apiRequest('/listings', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(l)
        });
      } catch {}
    }
    bulkBtn.disabled = false;
    loadListings();
  };
}

document.addEventListener('DOMContentLoaded', () => {
  checkAuth();
  loadListings();
}); 
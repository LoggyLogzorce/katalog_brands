// Элементы DOM
const registerForm = document.getElementById('registerForm');
const passwordInput = document.getElementById('password');
const confirmPasswordInput = document.getElementById('confirmPassword');
const passwordStrength = document.getElementById('passwordStrength');
const registerBtn = document.getElementById('registerBtn');
const passwordRequirements = {
    length: document.getElementById('reqLength'),
    number: document.getElementById('reqNumber'),
    letter: document.getElementById('reqLetter')
};

// Проверка сложности пароля
function checkPasswordStrength(password) {
    let strength = 0;
    let validRequirements = 0;

    // Сброс стилей требований
    Object.values(passwordRequirements).forEach(req => {
        req.classList.remove('valid');
    });

    // Проверка длины
    if (password.length >= 8) {
        strength += 40;
        passwordRequirements.length.classList.add('valid');
        validRequirements++;
    }

    // Проверка наличия цифр
    if (/\d/.test(password)) {
        strength += 30;
        passwordRequirements.number.classList.add('valid');
        validRequirements++;
    }

    // Проверка наличия букв
    if (/[a-zA-Z]/.test(password)) {
        strength += 30;
        passwordRequirements.letter.classList.add('valid');
        validRequirements++;
    }

    // Обновление индикатора силы пароля
    passwordStrength.style.width = `${strength}%`;

    // Цвет индикатора в зависимости от силы
    if (strength < 50) {
        passwordStrength.style.background = '#e74c3c';
    } else if (strength < 80) {
        passwordStrength.style.background = '#f39c12';
    } else {
        passwordStrength.style.background = '#2ecc71';
    }

    // Активация кнопки, если все требования выполнены
    registerBtn.disabled = validRequirements < 3;

    return strength;
}

// Проверка соответствия паролей
function checkPasswordMatch() {
    const password = passwordInput.value;
    const confirmPassword = confirmPasswordInput.value;

    if (password && confirmPassword && password !== confirmPassword) {
        confirmPasswordInput.style.borderColor = '#e74c3c';
        confirmPasswordInput.style.boxShadow = '0 0 0 3px rgba(231, 76, 60, 0.2)';
        return false;
    } else {
        confirmPasswordInput.style.borderColor = '#e0e0e0';
        confirmPasswordInput.style.boxShadow = 'none';
        return true;
    }
}

// Обработчики событий
passwordInput.addEventListener('input', function() {
    checkPasswordStrength(this.value);
});

confirmPasswordInput.addEventListener('input', checkPasswordMatch);

registerForm.addEventListener('submit', function(e) {
    e.preventDefault();

    // Проверка соответствия паролей
    if (!checkPasswordMatch()) {
        alert('Пароли не совпадают!');
        return;
    }

    // Сбор данных формы
    const formData = {
        name: document.getElementById('firstName').value,
        email: document.getElementById('email').value,
        password: document.getElementById('password').value,
    };

    // Эмуляция отправки данных на сервер
    console.log('Регистрация с данными:', formData);

    fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    }).then(r => r.json()).then(r => {
        if (r.error) {
            alert(r.error);
        } else {
            alert('Регистрация прошла успешно! Добро пожаловать!');
            window.location.href = '/';
        }
    });

    // Перенаправление на страницу входа
    setTimeout(() => {
        window.location.href = '/auth';
    }, 1000);
});

// Инициализация проверки пароля
passwordInput.dispatchEvent(new Event('input'));
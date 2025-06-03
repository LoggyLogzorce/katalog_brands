document.getElementById('loginForm').addEventListener('submit', function(e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    // Здесь будет логика входа
    console.log('Вход с данными:', {email, password});

    // В реальном приложении здесь будет AJAX-запрос к серверу
    fetch('/api/v1/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({email, password})
    }).then(r => r.json()).then(r => {
        if (r.error) {
            alert(r.error);
        } else {
            alert('Вход выполнен успешно!');
            window.location.href = '/';
        }
    });
});
(() => {
    const input_login = document.getElementById('input-login');
    const input_password = document.getElementById('input-password');
    const input_label = document.getElementById('input-data')

    document.addEventListener('keyup', event => {
        if (event.code === 'Enter')
            document.getElementById('login-btn').click();
    });

    function closeInvalid(event) {
        let classes_login = input_login.classList;
        let classes_password = input_password.classList;
        if (classes_login.contains('is-invalid'))
            classes_login.remove('is-invalid');
        if (classes_password.contains('is-invalid'))
            classes_password.remove('is-invalid');
        input_label.querySelector('.invalid-feedback').style = "";
    }
    input_login.addEventListener('input', closeInvalid);
    input_password.addEventListener('input', closeInvalid);

    document.getElementById('login-btn').addEventListener('click', () => {
        if (!input_login.value) {
            input_login.classList.add('is-invalid');
            return undefined;
        }
        if (!input_password.value) {
            input_password.classList.add('is-invalid');
            return undefined;
        }

        let xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/', true);

        xhr.setRequestHeader("Accept", "application/json");
        xhr.setRequestHeader("Content-Type", "application/json");

        xhr.onload = () => {
            var response = JSON.parse(xhr.responseText);
            if (xhr.status < 300) {
                if (response.error) {
                    console.log(response.error);
                    input_label.querySelector('.invalid-feedback')
                        .style.display = "block";
                    input_login.classList.add('is-invalid');
                    input_password.classList.add('is-invalid');
                } else {
                    window.location.href = response.redirect;
                }
            } else {
                alert(`ОШИБКА: "${response.error}"`)
            }
        };

        let send = new Object();
        let data = new Object();
        send['type'] = "login";
        data['login'] = input_login.value;
        data['password'] = input_password.value;
        send['data'] = JSON.stringify(data);

        xhr.send(JSON.stringify(send));
    });
})();
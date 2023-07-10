(() => {
    const full_name = document.getElementById('full-name');
    const login = document.getElementById('login');
    const email = document.getElementById('email');
    const password = document.getElementById('input-password');
    const password2 = document.getElementById('input-password2');
    const get_avatar = document.getElementById('get-avatar');
    const faculty = document.getElementById('faculty');
    const group = document.getElementById('group');
    const approval = document.getElementById('approval');
    const button = document.querySelector('button[type=submit]')

    password.addEventListener('input', (event) => {
        if (event.currentTarget.value.length) {
            password.removeEventListener('input', this);
            password2.removeAttribute('disabled');
        }
    });

    approval.addEventListener('change', (event) => {
        button.disabled = !Boolean(event.currentTarget.checked);
    });

    login.addEventListener('keyup', (event) => {
        let value = event.currentTarget.value.toLowerCase();
        event.currentTarget.value = value.replace(/[^a-z\d]/g, '');
        if (event.currentTarget.classList.contains('is-invalid'))
            if (event.currentTarget.value)
                event.currentTarget.classList.remove('is-invalid');
    });

    full_name.addEventListener('change', (event) => {
        if (event.currentTarget.classList.contains('is-invalid'))
            event.currentTarget.classList.remove('is-invalid');
    });

    var base64avatar = ""

    function avatarToBase64(input_element) {
        let file = input_element.files[0];
        var reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = function () {
            base64avatar = reader.result
        };
        reader.onerror = function (error) {
            console.log('Avatar convert error: ', error);
            base64avatar = "avatar conver arror: " + error
        };
    }

    function passwordVerification() {
        let label = password.parentElement.querySelector('.invalid-feedback');
        if (password.value.length < 5) {
            label.innerText = 'Пароль слишком короткий';
            password.classList.add('is-invalid');
        } else password.classList.remove('is-invalid');
        if (password.value && password2.value)
            if (password.value != password2.value) {
                label.innerText = '';
                password.classList.add('is-invalid');
                password2.classList.add('is-invalid');
            } else {
                password.classList.remove('is-invalid');
                password2.classList.remove('is-invalid');
            }
    }

    password.addEventListener('change', passwordVerification);
    password2.addEventListener('change', passwordVerification);

    button.addEventListener('click', () => {
        if (!full_name.value) {
            full_name.classList.add('is-invalid');
            return undefined;
        }
        if (!login.value) {
            login.classList.add('is-invalid');
            return undefined;
        }
        if (!password.value) {
            password.classList.add('is-invalid');
            return undefined;
        }
        if (get_avatar.value)
            avatarToBase64(get_avatar);
        else
            base64avatar = "none";

        let xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/', true);

        xhr.setRequestHeader("Accept", "application/json");
        xhr.setRequestHeader("Content-Type", "application/json");

        xhr.onload = () => {
            var response = JSON.parse(xhr.responseText);
            console.log(response);
            if (response.error === "")
                window.location.replace(response.redirect);
            else alert(response.error);
        };

        var send = (xhr) => {
            let send = new Object();
            let data = new Object();
            send['type'] = "registration";
            data['login'] = login.value;
            data['password'] = password.value;
            data['full_name'] = full_name.value;
            data['email'] = email.value;
            data['faculty'] = faculty.value;
            data['group'] = group.value;
            data['avatar_base64'] = base64avatar;
            send['data'] = JSON.stringify(data);
            xhr.send(JSON.stringify(send));
        }
        if (base64avatar) {
            send(xhr);
            return undefined;
        }
        // Если рендер не завершиться по истечении времени, останавливаем цикл
        // повторить с интервалом 0.5 секунды
        let timerId = setInterval((xhr) => { if (base64avatar) { send(xhr); base64avatar = ""; } }, 500, xhr);
        // остановить через 5 секунд
        setTimeout((xhr, timer) => { clearInterval(timer); send(xhr); }, 5000, xhr, timerId);
    });
})();
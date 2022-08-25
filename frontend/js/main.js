const openRegristerModal = () => {
    const loginPageId = document.querySelector('#login-page-id');
    const registerPageId = document.querySelector('#register-page-id');
    loginPageId.classList.remove('open');
    registerPageId.classList.remove('close');
    loginPageId.classList.add('close');
    registerPageId.classList.add('open');
};
const openRegisterBtn = document.querySelector('#open-register-btn-id');
openRegisterBtn.addEventListener('click', (e) => {
    e.preventDefault();
    openRegristerModal();
});
const showMessages = (msg = '') => {
    const messageContainer = document.querySelector('#message-container-id');
    messageContainer.textContent = msg;
    messageContainer.classList.add('show');
    setTimeout(() => {
        messageContainer.classList.remove('show');
    }, 3000);
};
const checkRegisterData = (userData) => {
    for (let key of Object.keys(userData)) {
        if (userData[key] === '') {
            showMessages(`Missing field: ${key}`);
            return [false, key];
        }
    }
    return [true, ''];
};

let socket;
const CreateWebSocket = () => {
    console.log('Attempting to open!');
    socket = new WebSocket('ws://localhost:8800/ws');
};
const validateUser = (resp) => {
    if (resp.Msg === 'Login successful') {
        //Create the cookie when logged in#
        CreateWebSocket();
        showMessages('Login successful');
        const loginPageId = document.querySelector('#login-page-id');
        const registerPageId = document.querySelector('#register-page-id');
        const mainPageId = document.querySelector('#main-page-id');
        loginPageId.classList.remove('open');
        registerPageId.classList.remove('open');
        mainPageId.classList.remove('close');
        loginPageId.classList.add('close');
        registerPageId.classList.add('close');
        mainPageId.style.display = 'grid';
    } else {
        showMessages(resp.Msg);
    }
};

const getRegisterData = () => {
    const firstName = document.querySelector('#first-name-id');
    const nickname = document.querySelector('#nickname-id');
    const lastName = document.querySelector('#last-name-id');
    const age = document.querySelector('#age-id');
    const gender = document.querySelector('#gender-id');
    const email = document.querySelector('#email-id');
    const password = document.querySelector('#password-id');
    const confirmPassword = document.querySelector('#confirm-password-id');
    const userData = {
        firstname: firstName.value,
        nickname: nickname.value,
        lastName: lastName.value,
        age: age.value,
        gender: gender.value,
        email: email.value,
        password: password.value,
        confirmPassword: confirmPassword.value,
    };
    return [checkRegisterData(userData)[0], userData];
};

const ClearRegistrationFields = () => {
    const firstName = document.querySelector('#first-name-id');
    const nickname = document.querySelector('#nickname-id');
    const lastName = document.querySelector('#last-name-id');
    const age = document.querySelector('#age-id');
    const gender = document.querySelector('#gender-id');
    const email = document.querySelector('#email-id');
    const password = document.querySelector('#password-id');
    const confirmPassword = document.querySelector('#confirm-password-id');
    firstName.value = '';
    nickname.value = '';
    lastName.value = '';
    age.value = '';
    gender.value = '';
    email.value = '';
    password.value = '';
    confirmPassword.value = '';
};
const registerBtn = document.querySelector('#register-btn-id');
registerBtn.addEventListener('click', (e) => {
    e.preventDefault();
    const userData = getRegisterData();
    if (userData[0]) {
        fetch('/register', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData[1]),
        })
            .then((response) => {
                return response.text();
            })
            .then((resp) => {
                showMessages(resp);
                if (resp === 'Register successful') {
                    ClearRegistrationFields();
                    setTimeout(() => {
                        const loginPageId =
                            document.querySelector('#login-page-id');
                        const registerPageId =
                            document.querySelector('#register-page-id');
                        registerPageId.classList.remove('open');
                        registerPageId.classList.add('close');
                        loginPageId.classList.remove('close');
                        loginPageId.classList.add('open');
                    }, 2500);
                }
                return resp;
            });
    }
});
const openLoginModal = () => {
    const loginPageId = document.querySelector('#login-page-id');
    const registerPageId = document.querySelector('#register-page-id');
    loginPageId.classList.remove('close');
    loginPageId.classList.add('open');
    registerPageId.classList.remove('open');
    registerPageId.classList.add('close');
};
const openLoginBtn = document.querySelector('#open-login-btn-id');
openLoginBtn.addEventListener('click', (e) => {
    e.preventDefault();
    openLoginModal();
});
const getloginData = () => {
    const emailOrUsername = document.querySelector(
        '#login-email-or-username-id'
    );
    const password = document.querySelector('#login-password-id');
    const userLoginData = {
        emailOrUsername: emailOrUsername.value,
        password: password.value,
    };
    Object.freeze();
    return userLoginData;
};
const loginBtn = document.querySelector('#login-btn-id');
loginBtn.addEventListener('click', (e) => {
    e.preventDefault();
    const userLoginData = getloginData();
    fetch('/login', {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(userLoginData),
    })
        .then((response) => {
            return response.json();
        })
        .then((resp) => {
            validateUser(resp);
            console.log(resp);
        });
});

function unSet(fields, revBtn) {
    setTimeout(function () {
        fields.forEach((field) => field.setAttribute('type', 'password'));
        revBtn.innerText = 'Reveal Password';
    }, 5000);
}

function revealPasswordBtn(id, className) {
    const revealBtn = document.querySelector(id);
    const inputFields = document.querySelectorAll(className);

    inputFields.forEach((eachField) => {
        if (eachField.getAttribute('type') === 'password') {
            eachField.setAttribute('type', 'text');
            revealBtn.innerText = 'Hide Password';
        } else if (eachField.getAttribute('type') === 'text') {
            eachField.setAttribute('type', 'password');
            revealBtn.innerText = 'Reveal Password';
        }
    });

    unSet(inputFields, revealBtn);
}
const openChatModal = (e) => {
    console.log(e);
    const chatModalContainer = document.querySelector(
        '#chat-modal-container-id'
    );
    const userID = e.getAttribute('data-user-id'); //data-user-id is the id of the user where we click on. This will be use to access the data on the database
    //when open a specific chat, we're going to get the chat data between the current user and the user tat they click
    chatModalContainer.style.display = 'flex';
    console.log(userID);
};
const closeChat = () => {
    const chatModalContainer = document.querySelector(
        '#chat-modal-container-id'
    );
    chatModalContainer.style.display = 'none';
};
const openPostModal = (e) => {
    console.log(e);
    const postModalContainer = document.querySelector(
        '#create-post-modal-container-id'
    );
    postModalContainer.style.display = 'flex';
};
const closeNewPost = () => {
    const postModalContainer = document.querySelector(
        '#create-post-modal-container-id'
    );
    postModalContainer.style.display = 'none';
};
const sendNewPost = () => {
    closeNewPost();
};

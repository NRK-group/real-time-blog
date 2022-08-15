const openRegristerModal = () => {
    const loginPageId = document.querySelector('#login-page-id');
    const registerPageId = document.querySelector('#register-page-id');
    loginPageId.classList.add('close');
    registerPageId.classList.add('open');
};
const openRegisterBtn = document.querySelector('#open-register-btn-id');
openRegisterBtn.addEventListener('click', (e) => {
    e.preventDefault();
    openRegristerModal();
});

const getRegisterData = () => {
    console.log('hello');
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
    Object.freeze;
    return userData;
};
const registerBtn = document.querySelector('#register-btn-id');
registerBtn.addEventListener('click', (e) => {
    e.preventDefault();
    const userData = getRegisterData();
    console.log(userData);
    fetch('/register', {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
    }).then((response) => {
        return response;
    });
});

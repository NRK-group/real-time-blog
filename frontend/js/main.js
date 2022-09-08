const SendResponsebtn = document.getElementById('send-response-btn');

let allPost, gUsers, gChatUsers;

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
function getCookie(name) {
    // Split cookie string and get all individual name=value pairs in an array
    var cookieArr = document.cookie.split(';');

    // Loop through the array elements
    for (var i = 0; i < cookieArr.length; i++) {
        var cookiePair = cookieArr[i].split('=');

        /* Removing whitespace at the beginning of the cookie name
        and compare it with the given string */
        if (name == cookiePair[0].trim()) {
            // Decode the cookie value and return
            return cookiePair[1];
        }
    }

    // Return null if not found
    return null;
}

const LoadMessages = (messageText, classType, sentTime) => {
    const MSG = AddMessages(messageText, classType, sentTime);
    const CHAT_CONTENT_CONTAINER = document.querySelector(
        '.chat-content-container'
    );
    CHAT_CONTENT_CONTAINER.prepend(MSG);
    // console.log('Prepending');
    // const FIRST_MSG = document.querySelector('.chat');
    // if (FIRST_MSG === null) {
    //     CHAT_CONTENT_CONTAINER.append(MSG);
    // } else {
    //     FIRST_MSG.prepend(MSG);
    // }
};

const AddMessages = (messageText, classType, sentTime) => {
    //Make the div that will hold everything
    const MSG_HOLDER = document.createElement('div');
    MSG_HOLDER.className = classType;
    //Create the inside divs
    const PROFILE_IMG = document.createElement('div');
    PROFILE_IMG.className = 'user-image';

    const MSG_CONTAINER = document.createElement('div');
    MSG_CONTAINER.className = 'chat-content';
    //Create a p for the message text
    const TEXT = document.createElement('p');
    TEXT.innerHTML = messageText;
    MSG_CONTAINER.append(TEXT);
    //Create a div for the date
    const DATE = document.createElement('div');
    DATE.innerHTML = `${sentTime[2]} ${sentTime[1]} ${
        sentTime[3]
    } ${sentTime[4].slice(0, -3)}`;
    DATE.classList.add('message-date');
    MSG_CONTAINER.append(DATE);

    MSG_HOLDER.append(PROFILE_IMG, MSG_CONTAINER);
    return MSG_HOLDER;
};

const DisplayMessage = (messageText, classType, sentTime) => {
    const MSG = AddMessages(messageText, classType, sentTime);
    const CHAT_CONTENT_CONTAINER = document.querySelector(
        '.chat-content-container'
    );
    CHAT_CONTENT_CONTAINER.append(MSG);
    CHAT_CONTENT_CONTAINER.scroll({
        top: CHAT_CONTENT_CONTAINER.scrollHeight,
    });
};

// let ;
let typing, debounce;

const StoppedTyping = () => {
    TYPING_MSG = document
        .querySelector('.fa-message')
        .classList.remove('animate-typing');
};

const Debounce = (callback, time) => {
    window.clearTimeout(debounce);
    debounce = window.setTimeout(callback, time);
};

let allUsers;

// ProccessMessage is a function that will display the message in the chat if the user has it open.
const ProcessMessage = (message) => {
    const TYPING_MSG = document.querySelector('.fa-message');
    if (message.message === ' ') {
        {
            let username;
            for (let i = 0; i < allUsers.length; i++) {
                if (allUsers[i].UserID === message.senderID) {
                    username = allUsers[i].Nickname;
                }
            }
            //User has started typing
            TYPING_MSG.innerHTML = `${username} Is Typing ...`;
            TYPING_MSG.classList.add('animate-typing');
            Debounce(StoppedTyping, 1750);
        }
    } else {
        DisplayMessage(message.message, 'chat', message.date.split(' '));
    }
};

const AddNotification = (i, senderID) => {
    console.log('Adding notification with:', senderID);
    const MESSAGE_BOX = document.getElementById(senderID);
    let notValue = parseInt(MESSAGE_BOX.innerHTML);
    notValue += parseInt(i);
    MESSAGE_BOX.innerHTML = notValue;
    MESSAGE_BOX.style.display = 'flex';
};

let socket;
const CreateWebSocket = () => {
    socket = new WebSocket('ws://localhost:8800/ws');
    socket.onopen = () => {
        //Access The cookie value
        let cookie = getCookie('session_token');
        if (cookie == null) {
            return;
        }

        // socket.send(cookie);
    };
    socket.onmessage = (text) => {
        let messageInfo = JSON.parse(text.data);
        if (messageInfo.message === 'e702c728-67f2-4ecd-9e79-4795010501ea') {
            const notificationContainer = document.querySelector(
                '#notification-container-id'
            );
            notificationContainer.classList.add('visible');
            return;
        }
        const CHAT_MODAL_CONTAINER = document.querySelector(
                '#chat-modal-container-id'
            ),
            SEND_BTN = document.querySelector('.send-chat-btn');
        // Check if the chat modal is open and that the reciever id is the ID of the user who sent the message
        messageInfo.recieverID = getCookie('session_token').split('&')[0];

        if (
            CHAT_MODAL_CONTAINER.style.display === 'flex' &&
            SEND_BTN.getAttribute('data-reciever-id') == messageInfo.senderID
        ) {
            console.log('Proccessing', messageInfo);
            ProcessMessage(messageInfo);
            return;
        }
        if (messageInfo.message != ' ') {
            //Show the notification animation
            console.log('You have a message from: ', messageInfo.senderID);

            AddNotification(1, messageInfo.senderID);
            //Return the notification to the db
            messageInfo.notification = true;
            messageInfo.userID = messageInfo.senderID;
            console.log('Sending back to the golang: ', messageInfo);
            socket.send(JSON.stringify(messageInfo));
        }
    };
};

const validateCoookie = () => {
    fetch('/vadidate')
        .then(async (response) => {
            resp = await response.json();
            if (resp.Msg === 'Login successful') {
                console.log('Valid cookie');
                validateUser(resp);
            }
        })
        .catch(() => {
            console.log('no valid cookie');
        });
};
const removeAllChildNodes = (parent) => {
    while (parent.firstChild) {
        parent.removeChild(parent.firstChild);
    }
};
const Logout = () => {
    fetch('/logout').then(async (response) => {
        resp = await response.text();
        showMessages(resp);
        const loginPageId = document.querySelector('#login-page-id');
        const registerPageId = document.querySelector('#register-page-id');
        const mainPageId = document.querySelector('#main-page-id');
        loginPageId.classList.remove('close');
        registerPageId.classList.remove('close');
        mainPageId.classList.remove('open');
        loginPageId.classList.add('open');
        registerPageId.classList.add('close');
        mainPageId.style.display = 'grid';
    });
};

const TypingMessage = (val) => {
    const USER_ID = getCookie('session_token').split('&')[0];
    const RECIEVER = document
        .querySelector('.send-chat-btn')
        .getAttribute('data-reciever-id');

    return JSON.stringify({
        message: val,
        userID: USER_ID,
        recieverID: RECIEVER,
    });
};
const NewPostCreated = (val) => {
    const USER_ID = getCookie('session_token').split('&')[0];
    return JSON.stringify({
        message: val,
        userID: USER_ID,
    });
};

const NewPostNotif = () => {
    socket.send(NewPostCreated('e702c728-67f2-4ecd-9e79-4795010501ea'));
};
const IsTyping = () => {
    //Send typing message when they start
    socket.send(TypingMessage(' '));
};

const validateUser = (resp) => {
    if (resp.Msg === 'Login successful') {
        //Create the cookie when logged in#
        CreateWebSocket();
        gUsers = [];
        gChatUsers = [];
        if (resp.Users) {
            gUsers = resp.Users;
        }
        if (resp.ChatUsers) {
            gChatUsers = resp.ChatUsers;
        }
        GetNotificationAmount();
        ShowUsers();
        allUsers = [...(gUsers || []), ...(gChatUsers || [])];
        allPost = resp.Posts;
        DisplayAllPost(resp.Posts);
        const loginPageId = document.querySelector('#login-page-id');
        const registerPageId = document.querySelector('#register-page-id');
        const mainPageId = document.querySelector('#main-page-id');
        loginPageId.classList.remove('open');
        registerPageId.classList.remove('open');
        mainPageId.classList.remove('close');
        loginPageId.classList.add('close');
        registerPageId.classList.add('close');
        mainPageId.style.display = 'grid';
        UpdateUserProfile(resp);
    } else {
        showMessages(resp.Msg);
    }
};

const UpdateUserProfile = (resp) => {
    document.getElementById(
        'profile-name-id'
    ).innerText = `${resp.User.Firstname}  ${resp.User.Lastname}`;
    document.getElementById(
        'profile-username-id'
    ).innerText = ` @${resp.User.Nickname}`;
    let yearCreated = resp.User.DateCreated.split(',')[1];
    document.querySelector(
        '#account-date-created-id'
    ).innerText = `since ${yearCreated}`;
    //User model
    document.getElementById('edit-first-name-id').value = resp.User.Firstname;
    document.getElementById('edit-last-name-id').value = resp.User.Lastname;
    document.getElementById('edit-nickname-id').value = resp.User.Nickname;
    document.getElementById('edit-age-id').value = resp.User.Age;
    document.getElementById('edit-email-id').value = resp.User.Email;
    document.getElementById('edit-gender-id').value = resp.User.Gender;

    console.log('RESP+++ ', resp.User);
};

const editBtn = document.getElementById('save-changes-btn');

editBtn.onclick = () => {
    console.log('dthdgh');
    EditUserProfile();
};

const EditUserProfile = () => {
    let lastName = document.getElementById('edit-last-name-id').value;
    let firstName = document.getElementById('edit-first-name-id').value;
    let nickname = document.getElementById('edit-nickname-id').value;
    let age = document.getElementById('edit-age-id').value;
    let gender = document.getElementById('edit-gender-id').value;
    let email = document.getElementById('edit-email-id').value;
    let password = document.getElementById('edit-password-id').value;
    let newPassword = document.getElementById('new-password-id').value;
    let confirmPassword = document.getElementById(
        'confirm-new-password-id'
    ).value;

    let UserInfo = {
        nickname: nickname,
        age: age,
        gender: gender,
        password: newPassword,
        confirmPassword: confirmPassword,
        oldPassword: password,
        email: email,
        firstName: firstName,
        lastName: lastName,
    };

    if (CheckRequirements(UserInfo) === '') {
        fetch('/updateuser', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(UserInfo),
        })
            .then(async (response) => {
                resp = await response.json();

                return resp;
            })
            .then((resp) => {
                showMessages(resp.Msg);
            });
    } else {
        showMessages(CheckRequirements(UserInfo));
    }
};

const GetNotificationAmount = () => {
    fetch('/Notify', {
        method: 'GET',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
    })
        .then(async (response) => {
            resp = await response.json();
            return resp;
        })
        .then((resp) => {
            resp.forEach((thisUser) => {
                console.log('This user: ', typeof thisUser.count);
                if (thisUser.count == 0 || thisUser.count == undefined) return;
                AddNotification(thisUser.count, thisUser.senderID);
            });
        });
};

//CheckNotificationDisplay will hide notification divs when their innerHTML is 0
const CheckNotificationDisplay = (arr) => {
    console.log(arr);
    arr.forEach((user) => {
        const NOTIF_BOX = document.getElementById(user.UserID);
        if (parseInt(NOTIF_BOX.innerHTML) < 1) {
            NOTIF_BOX.display = 'none';
        }
    });
};

const ShowUsers = (firstRun = true) => {
    console.log('showusers -> gUsers: ', gUsers);
    if (gUsers) {
        let usersDiv = document.getElementById('forum-users-container');
        let lastChat = document.getElementById('all-forum-users-container');
        let lastChatUsers = '';
        (gChatUsers || []).forEach((item, index) => {
            let username;
            item.Nickname.length < 8
                ? (username = item.Nickname)
                : (username = item.Nickname.slice(0, 6) + '...');
            let status;
            item.Status === 'Online'
                ? (status = 'online')
                : (status = 'offline');
            lastChatUsers =
                `<div
            key=${index}
            class="forum-user"
            data-user-id=${item.UserID}
            data-username=${item.Nickname}
            onclick="openChatModal(this)"
        >
        <span class="icon ${status}">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <path d="M256 512c141.4 0 256-114.6 256-256S397.4 0 256 0S0 114.6 0 256S114.6 512 256 512z"/>
        </svg>
        </span>
        <div class="user-image"></div>
        <div class="username">${username} <div class="notification" id="${item.UserID}">0</div></div>
        </div>` + lastChatUsers;
        });

        let usersDivTitle = document.getElementById('forum-users-title');
        let allUsersDivTitle = document.getElementById('forum-all-users-title');

        let users = '';
        (gUsers || []).forEach((item, index) => {
            let username;
            item.Nickname.length < 8
                ? (username = item.Nickname)
                : (username = item.Nickname.slice(0, 6) + '...');
            users =
                `<div
            key=${index}
            class="forum-user"
            data-user-id=${item.UserID}
            data-username=${item.Nickname}
            onclick="openChatModal(this)"
        >
        <span class="icon offline">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <path d="M256 512c141.4 0 256-114.6 256-256S397.4 0 256 0S0 114.6 0 256S114.6 512 256 512z"/>
        </svg>
        </span>
        <div class="user-image"></div>
        <div class="username">${username} <div class="notification" id="${item.UserID}">0</div></div>
        </div>` + users;

            // AddNotification(notifs, item.UserID)
        });

        lastChat.innerHTML = users;
        usersDiv.innerHTML = lastChatUsers;
        usersDivTitle.innerText = `Chats`;
        allUsersDivTitle.innerText = `Users`;
        CheckNotificationDisplay([...(gUsers || []), ...(gChatUsers || [])]);
    }
};

const CheckRequirements = (userInfo) => {
    //Check the username is between 3 and 11
    if (userInfo.nickname.length < 3 || userInfo.nickname.length > 11)
        return 'Nicknames must be 3-11 characters';
    if (userInfo.age > 115) return 'Max age is 115';
    if (userInfo.confirmPassword != userInfo.password)
        return "The passwords don't match";
    return '';
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
    const RESULT = [checkRegisterData(userData)[0], userData];
    if (!RESULT) return;
    if (CheckRequirements(userData) != '') {
        return [false, showMessages(CheckRequirements(userData))];
    }

    return RESULT;
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

const logoutBtn = document.getElementById('logout-btn');

logoutBtn.onclick = () => {
    Logout();
};

const registerBtn = document.querySelector('#register-btn-id');
registerBtn.addEventListener('click', (e) => {
    // CreateWebSocket()
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
            showMessages(resp.Msg);
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

const DisplayTenMessages = (messages) => {
    // console.log('DISlaying ten messages');
    const CURR_USER_ID = getCookie('session_token').split('&')[0];
    const CHAT_CONTENT_CONTAINER = document.querySelector(
        '.chat-content-container'
    );

    let prev = CHAT_CONTENT_CONTAINER.scrollHeight;
    messages.forEach((msg) => {
        let classNames = 'chat';
        if (msg.senderID === CURR_USER_ID) {
            classNames = 'chat sender';
        }

        LoadMessages(msg.message, classNames, msg.date.split(' '));
    });

    //Causing eventlistner to target
    CHAT_CONTENT_CONTAINER.scroll({
        top: CHAT_CONTENT_CONTAINER.scrollHeight - prev,
    });
};

//FetchMsgs loads the next 10 messages of the conversation onto the screen
const FetchMsgs = (chat, SEND_BTN) => {
    console.log('Getting Messages');
    fetch('/MessageInfo', {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(chat),
    })
        .then(async (response) => {
            resp = await response.json();
            SEND_BTN.setAttribute('data-chat-id', resp.chatID);
            //Add the first 10 messages
            if ((resp.Messages || []).length != 0)
                DisplayTenMessages(resp.Messages);
            return resp;
        })
        .catch((err) => {
            console.log('FetchMsgs() ERR:', err);
        });
};

function AllowMSG() {
    canRun = true;
    clearTimeout(timer);
}
let timer;
let canRun = true;
function GetMsg(users, SEND_BTN) {
    if (canRun) {
        FetchMsgs(users, SEND_BTN);
        canRun = false;
        timer = setTimeout(AllowMSG, 1000);
    }
}
const CHAT_CONTENT_CONTAINER = document.querySelector(
    '.chat-content-container'
);
let valid = false;

const DeleteChatNotifications = (usersID, recieversID) => {
    const DELETE = {
        userID: usersID,
        recieverID: recieversID,
    };
    fetch('/Notify', {
        method: 'PUT',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(DELETE),
    });

    //Remove the notifications from the clients page
    const MESSAGE_BOX = document.getElementById(usersID);
    MESSAGE_BOX.innerHTML = '0';
    MESSAGE_BOX.style.display = 'none';
};

const openChatModal = (e) => {
    console.log('Valid', valid);
    const RECIEVER_ID = e.getAttribute('data-user-id'); //data-user-id is the id of the user where we click on. This will be use to access the data on the database
    const RECIEVER_USERNAME = e.getAttribute('data-username');
    const CHAT_CONTAINER = document.querySelector('.chat-content-container');
    CHAT_CONTAINER.id = RECIEVER_ID;

    removeAllChildNodes(CHAT_CONTAINER);
    const chatModalContainer = document.querySelector(
        '#chat-modal-container-id'
    );
    const CHAT_USERNAME = document.querySelector('#chat-username-id');
    CHAT_USERNAME.innerHTML = RECIEVER_USERNAME;
    //when open a specific chat, we're going to get the chat data between the current user and the user tat they click
    chatModalContainer.style.display = 'flex';
    //Add the data to the send btn
    const SEND_BTN = document.querySelector('.send-chat-btn');
    SEND_BTN.setAttribute('data-reciever-id', RECIEVER_ID);
    //Now check golang for the chatID
    const USER_ID = getCookie('session_token').split('&')[0];
    let users = {
        userID: USER_ID,
        recieverID: RECIEVER_ID,
        X: 0,
    };

    DeleteChatNotifications(
        RECIEVER_ID,
        getCookie('session_token').split('&')[0]
    );
    GetMsg(users, SEND_BTN);

    CHAT_CONTENT_CONTAINER.addEventListener('scroll', CheckScrollTop);
    valid = true;
};

const ArrangeUsers = (userId) => {
    let user = '';
    gUsers.forEach((item, index) => {
        if (userId === item.UserID) {
            user = item;
            gUsers.splice(index, 1);
        }
    });

    if (user !== '') {
        document.getElementById(user.UserID).remove();
    }

    gChatUsers.forEach((item, inx) => {
        if (userId === item.UserID) {
            user = item;
            gChatUsers.splice(inx, 1);
        }
    });

    console.log('object ', user);
    gChatUsers = [...gChatUsers, user];
    ShowUsers(false);
};

const SendMessage = () => {
    //Get the message from the text box
    const TEXT_BOX = document.querySelector('.chat-input-box');
    const MSG = TEXT_BOX.value;
    //Get information using the send btns attributes
    const SEND_BTN = document.querySelector('.send-chat-btn');
    const SEND_TO = SEND_BTN.getAttribute('data-reciever-id');
    const SEND_FROM = getCookie('session_token').split('&')[0];
    const SENT_TIME = new Date();
    const SORTED = SENT_TIME.toString();
    const CHAT_ID = SEND_BTN.getAttribute('data-chat-id');
    const INFORMATION = {
        message: MSG.trim(),
        userID: SEND_FROM,
        recieverID: SEND_TO,
        date: SORTED,
        chatID: CHAT_ID,
    };

    if (MSG.trim().length !== 0) {
        socket.send(JSON.stringify(INFORMATION));
        DisplayMessage(MSG, 'chat sender', SORTED.split(' '));
        TEXT_BOX.value = '';
        console.log('gUsers in socket.send: ', gUsers);
        ArrangeUsers(SEND_TO);
    }
};

const closeChat = () => {
    valid = false;
    removeAllChildNodes(document.querySelector('.chat-content-container'));
    const chatModalContainer = document.querySelector(
        '#chat-modal-container-id'
    );
    chatModalContainer.style.display = 'none';
    //clear the text box
    document.querySelector('.chat-input-box').value = '';
    console.log('REmoving event listner');
    CHAT_CONTENT_CONTAINER.removeEventListener('scroll', CheckScrollTop);
};

const CheckScrollTop = () => {
    const CHAT_CONTENT_CONTAINER = document.querySelector(
        '.chat-content-container'
    );
    if (CHAT_CONTENT_CONTAINER.scrollTop !== 0 || !valid) return;
    const SEND_BTN = document.querySelector('.send-chat-btn');
    let chats = {
        userID: getCookie('session_token').split('&')[0],
        recieverID: SEND_BTN.getAttribute('data-reciever-id'),
        X: CHAT_CONTENT_CONTAINER.childElementCount,
    };
    GetMsg(chats, SEND_BTN);
};

const openPostModal = (e) => {
    const postModalContainer = document.querySelector(
        '#create-post-modal-container-id'
    );
    let postTitle = document.getElementById('new-post-title-id');
    let postCategory = document.getElementById('new-post-category-id');
    let postContent = document.getElementById('new-post-content-id');
    postTitle.value = '';
    postContent.value = '';
    postCategory.value = '';
    postModalContainer.style.display = 'flex';
};
const closeNewPost = () => {
    const postModalContainer = document.querySelector(
        '#create-post-modal-container-id'
    );
    postModalContainer.style.display = 'none';
};
const sendNewPost = () => {
    let postTitle = document.getElementById('new-post-title-id').value;
    let postCategory = document.getElementById('new-post-category-id').value;
    let postContent = document.getElementById('new-post-content-id').value;
    if (
        postTitle.length > 5 &&
        postCategory.length !== 0 &&
        postContent.length > 10
    ) {
        let newPost = {
            postTitle: postTitle,
            postCategory: postCategory,
            postContent: postContent,
        };
        fetch('/post', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newPost),
        })
            .then((response) => {
                return response.json();
            })
            .then((resp) => {
                showMessages(resp.Msg);
                allPost = resp.Posts;
                DisplayAllPost(resp.Posts);
            })
            .then(() => {
                NewPostNotif();
            });
        closeNewPost();
        return;
    }
    if (postTitle.length <= 5) {
        showMessages('Invalid length of title');
        return;
    }
    if (postCategory.length === 0) {
        showMessages('Choice category');
        return;
    }
    if (postContent.length <= 10) {
        showMessages('Invalid length of content');
        return;
    }
};

const SendResponse = (e) => {
    let content = document.getElementById('create-response-input-box').value;
    if (content !== '') {
        let newResponse = {
            postID: e.getAttribute('data-post-id'),
            responseContent: content,
        };

        fetch('/response', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newResponse),
        })
            .then((response) => {
                return response.json();
            })
            .then((resp) => {
                allPost = resp.Posts;
                DisplayAllPost(resp.Posts);
                document.getElementById('create-response-input-box').value = '';
                CreateResponses(allPost, e.getAttribute('data-post-id'));
            });

        return;
    }
    showMessages('Invalid length of content');
};

const openResponseModal = (postId) => {
    const responseModal = document.querySelector(
        '#response-modal-container-id'
    );
    SendResponsebtn.setAttribute('data-post-id', postId);

    let post;
    for (let item of allPost) {
        if (item.PostID === postId) {
            post = item;
            break;
        }
    }
    CreateResponses(allPost, postId);
    const responsePostContainer = document.querySelector(
        '#response-post-container'
    );
    let category = '';
    if (post.Category === 'golang') {
        category =
            '<div class="post-category golang golang-category">GoLang</div>';
    }
    if (post.Category === 'javascript') {
        category =
            '<div class="post-category javascript javascript-category">JavaScript</div>';
    }
    if (post.Category === 'rust') {
        category = '<div class="post-category rust rust-category">Rust</div>';
    }
    responsePostContainer.innerHTML = `
    <div class="post-title">${post.Title}</div>
    <div class="post-profile"> 
        <div class="post-user-profile">
            <div class="user-image"></div>
            <span>
                <div class="username">${post.UserID}</div>
                <div class="post-created">${post.Date}</div>
            </span>
        </div>
        ${category}
    </div>
    <div class="post-content overflow scrollbar-hidden">${post.Content}</div>`;
    responseModal.style.display = 'flex';
};

const CreateResponses = (allPost, postID) => {
    let comments = '';
    let allComments;
    for (let item of allPost) {
        if (item.PostID === postID) {
            allComments = item.Comments;
            break;
        }
    }
    let responseContainer = document.getElementById('all-reponse-container');
    if (allComments) {
        allComments.forEach((item) => {
            comments =
                `
                <div class="response-container">
                    <div class="response-user-profile">
                        <div class="user-image">
                        </div>
                        <span>
                            <div class="response-username">
                                @${item.UserID}
                                <span class="response-created">
                                ${item.Date}
                                </span>
                            </div>
                            <div class="response-content">
                                ${item.Content}
                            </div>
                        </span>
                    </div>
                </div>` + comments;
        });
    }
    responseContainer.innerHTML = comments;
};

const Favorite = (postId, react, node, iconNode) => {
    if (postId !== '') {
        if (react !== '0') {
            react = 0;
        } else {
            react = 1;
        }

        let newReaction = {
            postID: postId,
            react: react,
        };

        fetch('/favorite', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newReaction),
        })
            .then(async (response) => {
                return await response.json();
            })
            .then((resp) => {
                if (react === 1) {
                    iconNode.style.color = '#533de0';
                } else {
                    iconNode.style.color = '';
                }
                node.setAttribute('data-post-reaction', react);
                allPost = resp.Posts;
                if (
                    document
                        .querySelector('#favorite-post-id')
                        .classList.contains('all-post-active')
                ) {
                    let newAllPost = allPost.filter(
                        (post) => post.Favorite.react === 1
                    );
                    DisplayAllPost(newAllPost);
                }
            });
        return;
    }
};

const CreatePost = (
    postId,
    titleValue,
    userImageValue,
    usernameValue,
    postCreatedValue,
    postCategoryValue,
    postContentValue,
    react,
    UserID
) => {
    const postContainer = document.createElement('div');
    postContainer.className = 'post-container';
    postContainer.setAttribute('data-post-id', postId);
    const postTitle = document.createElement('div');
    postTitle.className = 'post-title';
    postTitle.textContent = titleValue;
    const postProfile = document.createElement('div');
    postProfile.className = 'post-profile';
    const postUserProfile = document.createElement('div');
    postUserProfile.className = 'post-user-profile';
    const userImage = document.createElement('div');
    userImage.className = 'user-image';
    //create an image to add the image here
    userImage.value = userImageValue; // this need to be converted to an image
    const span = document.createElement('span');
    const username = document.createElement('div');
    username.className = 'username';
    username.textContent = usernameValue;
    const postCreated = document.createElement('div');
    postCreated.className = 'post-created';
    postCreated.textContent = postCreatedValue;
    span.append(username, postCreated);
    postUserProfile.append(userImage, span);
    const postCategory = document.createElement('div');
    if (postCategoryValue === 'golang') {
        postCategory.classList.add(
            'post-category',
            'golang',
            'golang-category'
        );
    }
    if (postCategoryValue === 'javascript') {
        postCategory.classList.add(
            'post-category',
            'javascript',
            'javascript-category'
        );
    }
    if (postCategoryValue === 'rust') {
        postCategory.classList.add('post-category', 'rust', 'rust-category');
    }
    postCategory.textContent = postCategoryValue;
    postProfile.append(postUserProfile, postCategory);
    const postContent = document.createElement('div');
    postContent.className = 'post-content';
    postContent.innerHTML = postContentValue;
    const postButtons = document.createElement('div');
    postButtons.className = 'post-buttons';
    const favoriteBtn = document.createElement('div');
    favoriteBtn.className = 'favorite-btn';
    favoriteBtn.setAttribute('data-post-id', postId);
    favoriteBtn.tabIndex = '1';
    favoriteBtn.setAttribute('data-post-reaction', react);
    const favoriteIcon = document.createElement('span');
    favoriteIcon.className = 'favorite-icon';
    //--------
    favoriteBtn.onclick = () => {
        let nreact = favoriteBtn.getAttribute('data-post-reaction');
        Favorite(
            postId,
            nreact,
            (node = favoriteBtn),
            (iconNode = favoriteIcon)
        );
    };

    if (getCookie('session_token').split('&')[0] === UserID && react === 1) {
        favoriteIcon.style.color = '#533de0';
        favoriteBtn.setAttribute('data-post-reaction', react);
    }
    //-----
    favoriteIcon.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"> <path d="M48 0H336C362.5 0 384 21.49 384 48V487.7C384 501.1 373.1 512 359.7 512C354.7 512 349.8 510.5 345.7 507.6L192 400L38.28 507.6C34.19 510.5 29.32 512 24.33 512C10.89 512 0 501.1 0 487.7V48C0 21.49 21.49 0 48 0z"/></svg>`;
    favoriteBtn.append(favoriteIcon);
    const responseBtn = document.createElement('div');
    responseBtn.className = 'response-btn';
    const responseIcon = document.createElement('span');
    responseIcon.className = 'response-icon';
    responseIcon.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M256 32C114.6 32 .0272 125.1 .0272 240c0 49.63 21.35 94.98 56.97 130.7c-12.5 50.37-54.27 95.27-54.77 95.77c-2.25 2.25-2.875 5.734-1.5 8.734C1.979 478.2 4.75 480 8 480c66.25 0 115.1-31.76 140.6-51.39C181.2 440.9 217.6 448 256 448c141.4 0 255.1-93.13 255.1-208S397.4 32 256 32z"/></svg>`;
    responseBtn.onclick = () => {
        openResponseModal(postId);
    };
    responseBtn.append(responseIcon, 'Response');
    postButtons.append(favoriteBtn, responseBtn);
    postContainer.append(postTitle, postProfile, postContent, postButtons);
    return postContainer;
};
// this part need to be automated to all the post

const closeResponseModal = () => {
    const responseModal = document.querySelector(
        '#response-modal-container-id'
    );
    SendResponsebtn.setAttribute('data-post-id', '');
    responseModal.style.display = 'none';
};
const DisplayAllPost = (post) => {
    if (post) {
        const allPostContainer = document.querySelector(
            '#all-post-container-id'
        );
        removeAllChildNodes(allPostContainer);
        post.forEach(
            ({
                PostID,
                Title,
                ImgUrl,
                UserID,
                Date,
                Category,
                Content,
                Favorite,
            }) => {
                allPostContainer.append(
                    CreatePost(
                        PostID,
                        Title,
                        ImgUrl,
                        UserID,
                        Date,
                        Category,
                        Content,
                        Favorite.react,
                        Favorite.UserID
                    )
                );
            }
        );
    }
};
const openProfileModal = () => {
    const profileModal = document.querySelector('#profile-moadal-container-id');
    profileModal.style.display = 'flex';
};
const closeProfileModal = () => {
    const profileModal = document.querySelector('#profile-moadal-container-id');
    profileModal.style.display = 'none';
};
const selectFilter = (e) => {
    const allPost = document.querySelector('#all-post-id');
    allPost.classList.remove('all-post-active');
    const golangPost = document.querySelector('#golang-post-id');
    golangPost.classList.remove('golang-active');
    const javaScriptPost = document.querySelector('#javascript-post-id');
    javaScriptPost.classList.remove('javascript-active');
    const rustPost = document.querySelector('#rust-post-id');
    rustPost.classList.remove('rust-active');
    const yourPost = document.querySelector('#your-post-id');
    yourPost.classList.remove('all-post-active');
    const favoritePost = document.querySelector('#favorite-post-id');
    favoritePost.classList.remove('all-post-active');
    if (e.id === 'golang-post-id') {
        e.classList.add('golang-active');
    } else if (e.id === 'javascript-post-id') {
        e.classList.add('javascript-active');
    } else if (e.id === 'rust-post-id') {
        e.classList.add('rust-active');
    } else {
        e.classList.add('all-post-active');
    }
};
const openAllPost = (e) => {
    selectFilter(e);
    DisplayAllPost(allPost);
};
const openGoLangPost = (e) => {
    selectFilter(e);
    filterPost('golang');
};
const openJavaScriptPost = (e) => {
    selectFilter(e);
    filterPost('javascript');
};
const openRustPost = (e) => {
    selectFilter(e);
    filterPost('rust');
};
const openFavoritePost = (e) => {
    selectFilter(e);
    let newAllPost = allPost.filter((post) => post.Favorite.react === 1);
    DisplayAllPost(newAllPost);
};
const openYourPost = (e) => {
    selectFilter(e);
    let newAllPost = allPost.filter(
        (post) => post.UserID === getCookie('session_token').split('&')[1]
    );
    DisplayAllPost(newAllPost);
};

const filterPost = (tag) => {
    let newAllPost = allPost.filter((post) => post.Category === tag);
    DisplayAllPost(newAllPost);
};
const refreshThePost = () => {
    fetch('/post')
        .then((response) => {
            return response.json();
        })
        .then((resp) => {
            allPost = resp.Posts;
            DisplayAllPost(resp.Posts);
            return;
        })
        .then(() => {
            const middleContainer = document.querySelector('.middle-container');
            middleContainer.scrollTo({ top: 0, behavior: 'smooth' });
            const notificationContainer = document.querySelector(
                '#notification-container-id'
            );
            notificationContainer.classList.remove('visible');
            return;
        })
        .catch((err) => {
            console.log('refreshThePost()', err);
        });
    return;
};
let lastScrollTop = 0;
const scrollOnPost = (e) => {
    const allPost = document.querySelector('#all-post-id');
    let scrollTop = e.scrollTop;
    if (scrollTop <= lastScrollTop) {
        //if scrolling up
        if (scrollTop <= 30) {
            if (allPost.classList.contains('all-post-active')) refreshThePost();
        }
    }
    lastScrollTop = scrollTop;
};
const loadingPageTimer = (func, timeout = 300) => {
    let timer;
    return function (...args) {
        if (timer) {
            clearTimeout(timer);
        }

        timer = setTimeout(() => {
            func.apply(this, args);
        }, timeout);
    };
};
let closePage;
const closeLoadingPage = () => {
    console.log('Close');
    window.clearTimeout(closePage);
    closePage = setTimeout(function () {
        // Run the callback
        const loadingPage = document.querySelector('.loading-page-background');
        loadingPage.style.display = 'none';
        console.log('go');
    }, 2000);
};
closeLoadingPage();

const toggleSwitch = document.querySelector(
    '.theme-switch input[type="checkbox"]'
);
function switchTheme(e) {
    if (e.target.checked) {
        document.documentElement.setAttribute('data-theme', 'dark');
        localStorage.setItem('theme', 'dark');
    } else {
        document.documentElement.setAttribute('data-theme', 'light');
        localStorage.setItem('theme', 'light');
    }
}
toggleSwitch.addEventListener('change', switchTheme, false);
const currentTheme = localStorage.getItem('theme')
    ? localStorage.getItem('theme')
    : null;

if (currentTheme) {
    document.documentElement.setAttribute('data-theme', currentTheme);
    if (currentTheme === 'dark') {
        toggleSwitch.checked = true;
    }
}
const OpenNewPostCategory = (e) => {
    const postModalContainer = document.querySelector(
        '#create-post-modal-container-id'
    );
    let postTitle = document.getElementById('new-post-title-id');
    let postCategory = document.getElementById('new-post-category-id');
    let postContent = document.getElementById('new-post-content-id');
    postTitle.value = '';
    postContent.value = '';
    if (e.classList.contains('golang')) {
        postCategory.value = 'golang';
    }
    if (e.classList.contains('javascript')) {
        postCategory.value = 'javascript';
    }
    if (e.classList.contains('rust')) {
        postCategory.value = 'rust';
    }
    postModalContainer.style.display = 'flex';
};

const SendResponsebtn = document.getElementById('send-response-btn');

let allPost;

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

const DisplayMessage = (messageText, classType, sentTime) => {
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
    const CHAT_CONTENT_CONTAINER = document.querySelector(
        '.chat-content-container'
    );
    CHAT_CONTENT_CONTAINER.append(MSG_HOLDER);
    CHAT_CONTENT_CONTAINER.scroll({
        top: CHAT_CONTENT_CONTAINER.scrollHeight,
        behavior: 'smooth',
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
    const CHAT_MODAL_CONTAINER = document.querySelector(
            '#chat-modal-container-id'
        ),
        SEND_BTN = document.querySelector('.send-chat-btn');
    // Check if the chat modal is open and that the reciever id is the ID of the user who sent the message

    if (
        CHAT_MODAL_CONTAINER.style.display === 'flex' &&
        SEND_BTN.getAttribute('data-reciever-id') == message.senderID
    ) {
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
    }
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
        const MESSAGE_INFO = JSON.parse(text.data);
        ProcessMessage(MESSAGE_INFO);
    };
};

const validateCoookie = () => {
    fetch('/vadidate').then(async (response) => {
        resp = await response.json();
        if (resp.Msg === 'Login successful') {
            validateUser(resp);
        }
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

const IsTyping = () => {
    //Send typing message when they start
    socket.send(TypingMessage(' '));
};

const validateUser = (resp) => {
    if (resp.Msg === 'Login successful') {
        //Create the cookie when logged in#
        CreateWebSocket();
        ShowUsers(resp.Users);
        allUsers = resp.Users;
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
        console.log(resp);
        UpdateUserProfile(resp);
    } else {
        showMessages(resp.Msg);
    }
};

const UpdateUserProfile = (resp) => {
    document.getElementById(
        'profile-name'
    ).innerText = `${resp.User.Firstname}  ${resp.User.Lastname}`;
    document.getElementById(
        'profile-username'
    ).innerText = `@${resp.User.Nickname}`;

    //User model
    document.getElementById('edit-first-name-id').value = resp.User.Firstname;
    document.getElementById('edit-last-name-id').value = resp.User.Lastname;
    document.getElementById('edit-nickname-id').value = resp.User.Nickname;
    document.getElementById('edit-age-id').value = resp.User.Age;
    document.getElementById('edit-emial-id').value = resp.User.Email;
};

const ShowUsers = (Users) => {
    if (Users) {
        let usersDiv = document.getElementById('forum-users-container');
        let usersDivTitle = document.getElementById('forum-users-title');
        let users = '';
        Users.forEach((item, index) => {
            users =
                `<div
                    key=${index}
                    class="forum-user"
                    data-user-id=${item.UserID}
                    onclick="openChatModal(this)"
                >
                <div class="user-image"></div>
                <div class="username">@${item.Nickname}</div>
                </div>` + users;
        });
        usersDiv.innerHTML = users;
        usersDivTitle.innerText = `${Users.length} Active User`;
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
const openChatModal = (e) => {
    const chatModalContainer = document.querySelector(
        '#chat-modal-container-id'
    );
    const RECIEVER_ID = e.getAttribute('data-user-id'); //data-user-id is the id of the user where we click on. This will be use to access the data on the database
    //when open a specific chat, we're going to get the chat data between the current user and the user tat they click
    chatModalContainer.style.display = 'flex';
    //Add the data to the send btn
    const SEND_BTN = document.querySelector('.send-chat-btn');
    // const INFO_DIV = document.querySelector('.')
    SEND_BTN.setAttribute('data-reciever-id', RECIEVER_ID);
    //Now check golang for the chatID
    const USER_ID = getCookie('session_token').split('&')[0];
    let users = {
        userID: USER_ID,
        recieverID: RECIEVER_ID,
    };

    fetch('/MessageInfo', {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(users),
    }).then(async (response) => {
        resp = await response.json();
        console.log(resp.chatID);
        SEND_BTN.setAttribute('data-chat-id', resp.chatID);
        return resp;
    });
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
    }
};
const closeChat = () => {
    const chatModalContainer = document.querySelector(
        '#chat-modal-container-id'
    );
    chatModalContainer.style.display = 'none';
    //clear the text box
    document.querySelector('.chat-input-box').value = '';
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
        <div class="post-category golang golang-category">GoLang</div>
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

const CreatePost = (
    postId,
    titleValue,
    userImageValue,
    usernameValue,
    postCreatedValue,
    postCategoryValue,
    postContentValue,
    react
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
    favoriteBtn.tabIndex = '1';
    const favoriteIcon = document.createElement('span');
    favoriteIcon.className = 'favorite-icon';
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
                        Favorite.React
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
};
const openGoLangPost = (e) => {
    selectFilter(e);
};
const openJavaScriptPost = (e) => {
    selectFilter(e);
};
const openRustPost = (e) => {
    selectFilter(e);
};
const openFavoritePost = (e) => {
    selectFilter(e);
};
const openYourPost = (e) => {
    selectFilter(e);
};

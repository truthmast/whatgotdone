export default function loadUserKit(appId, initFn, signInFn) {
    const widgetJsUrl = "https://widget.userkit.io/widget.js";
    let scriptEls = document.getElementsByTagName("script");

    for (const el of scriptEls) {
        if (el.src == widgetJsUrl) {
            // widget.js already loaded, execute callback
            if (typeof initFn === 'function') {
                initFn(window.UserKit, window.UserKitWidget);
            }
            return;
        }
    }

    // If callback is provided, attach a listener for 'UserKitInit'
    if (typeof initFn === 'function') {
        document.addEventListener("UserKitInit", initFn);
    }

    // Attach listener for 'UserKitSignIn' event
    document.addEventListener("UserKitSignIn", () => {
        signInFn(window.UserKit, window.UserKitWidget);
    });

    // Load widget.js
    let userKitScript = document.createElement("script");
    userKitScript.setAttribute("src", widgetJsUrl);
    userKitScript.setAttribute(
        "data-app-id",
        appId
    );
    userKitScript.setAttribute("data-login-dismiss", "false");
    document.head.appendChild(userKitScript);

}

export function logout() {
    // eslint-disable-next-line no-unused-vars
    loadUserKit(process.env.VUE_APP_USERKIT_APP_ID, (userKit, userKitWidget) => {
        userKit.logout();
    });
}
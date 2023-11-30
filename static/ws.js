const hoso = document.querySelector("#hoso");
const playBtn = document.querySelector("#playBtn");
const play = "play";
const pause = "pause";
let action = play;

function toggleVideo() {
    if (action === "play") {
        // play video
        hoso.play();
        playBtn.innerHTML = pause;
        action = pause;
    } else {
        hoso.pause();
        playBtn.innerHTML = play;
        action = play;
    }
}

function onOpen() {
    document.querySelector("#status").innerHTML = "Connected";
}

function onClose() {
    document.querySelector("#status").innerHTML = "Disconnected";
}

function onConnecting() {
    document.querySelector("#status").innerHTML = "Connecting...";
}

document.body.addEventListener("htmx:wsConnecting", onConnecting);
document.body.addEventListener("htmx:wsOpen", onOpen);
document.body.addEventListener("htmx:wsClose", onClose);

document.body.addEventListener("htmx:wsConfigSend", (e) => {
    if (e.detail.parameters["videoStatus"]) {
        e.detail.parameters["action"] = action;
        toggleVideo();
    }
});

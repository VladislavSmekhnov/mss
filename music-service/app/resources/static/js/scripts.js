var audioPlayer;

window.addEventListener('DOMContentLoaded', event => {

    // Navbar shrink function
    var navbarShrink = function () {
        const navbarCollapsible = document.body.querySelector('#mainNav');
        if (!navbarCollapsible) {
            return;
        }
        if (window.scrollY === 0) {
            navbarCollapsible.classList.remove('navbar-shrink')
        } else {
            navbarCollapsible.classList.add('navbar-shrink')
        }

    };

    // Shrink the navbar 
    navbarShrink();

    // Shrink the navbar when page is scrolled
    document.addEventListener('scroll', navbarShrink);

    //  Activate Bootstrap scrollspy on the main nav element
    const mainNav = document.body.querySelector('#mainNav');
    if (mainNav) {
        new bootstrap.ScrollSpy(document.body, {
            target: '#mainNav',
            rootMargin: '0px 0px -40%',
        });
    };

    // Collapse responsive navbar when toggler is visible
    const navbarToggler = document.body.querySelector('.navbar-toggler');
    const responsiveNavItems = [].slice.call(
        document.querySelectorAll('#navbarResponsive .nav-link')
    );
    responsiveNavItems.map(function (responsiveNavItem) {
        responsiveNavItem.addEventListener('click', () => {
            if (window.getComputedStyle(navbarToggler).display !== 'none') {
                navbarToggler.click();
            }
        });
    });

    var hlsPlayer = function () {
        audioPlayer = document.getElementById("audio-player");
        const playPauseButton = document.getElementById("play-pause-button");
        const prevButton = document.getElementById("previous-btn");
        const nextButton = document.getElementById("next-btn");
        const seekSlider = document.getElementById("seek-slider");
        const currentTimeIndicator = document.getElementById("current-time");
        const durationDisplay = document.getElementById("duration");
        const muteButton = document.getElementById("mute-button");
        const volumeSlider = document.getElementById("volume-slider");
        const repeatButton = document.getElementById("repeat-btn");

        const songsPath = "../music_service/static/assets/audio/storage/"


        let currentStreamIndex = -1;
        let currentSelectionName = '';
        let isPlaying = false; // To track the player's playing state
        let repeatMode = 0; // To track the repeat button state: 0 = off, 1 = repeat one, 2 = repeat all

        const streamUrls = [];
        // Assuming audioFilesData contains the audio file data
        for (var i = 0; i < audioFilesData.length; i++) {
            // Assuming each audio file object has a "Path" property
            var currentTrack = audioFilesData[i];
            streamUrls.push(songsPath + currentTrack.Path);
        }

        function getTrackIndexById(trackId) {
            for (var i = 0; i < audioFilesData.length; i++) {
                if (String(audioFilesData[i].ID) == String(trackId)) {
                    return i;
                }
            }
            return -1;
        }

        // Function to load and play the current stream
        function loadStream(selectionName, trackId) {
            var oldStreamIndex = currentStreamIndex;
            currentStreamIndex = getTrackIndexById(trackId);
            currentSelectionName = selectionName;

            if (oldStreamIndex == currentStreamIndex) {
                togglePlayPause();
                return;
            }

            const currentStreamUrl = streamUrls[currentStreamIndex];
            const currentTrack = audioFilesData[currentStreamIndex];

            // Update the track title and artist
            document.querySelector(".track-title").textContent = currentTrack.TrackName;
            document.querySelector(".track-artist").textContent = currentTrack.Artist;

            // Update the HLS stream URL with the correct path
            if (Hls.isSupported()) {
                const hls = new Hls();
                hls.loadSource(currentStreamUrl);
                hls.attachMedia(audioPlayer);
                hls.on(Hls.Events.MANIFEST_PARSED, function () {
                    // HLS stream is ready, so enable player controls
                    playPauseButton.addEventListener("click", togglePlayPause);
                    prevButton.addEventListener("click", playPreviousStream);
                    nextButton.addEventListener("click", playNextStream);
                    seekSlider.addEventListener("input", updateSeek);
                    audioPlayer.addEventListener("timeupdate", updateTime);
                    audioPlayer.addEventListener("durationchange", updateDuration);
                    muteButton.addEventListener("click", toggleMute);
                    volumeSlider.addEventListener("input", updateVolume);
                    repeatButton.addEventListener("click", toggleRepeat);

                    updateRepeatButton();
                });
            } else if (audioPlayer.canPlayType("application/vnd.apple.mpegurl")) {
                // For Safari, use native HLS support if available
                audioPlayer.src = currentStreamUrl;
                audioPlayer.addEventListener("loadedmetadata", function () {
                    // HLS stream is ready, so enable player controls
                    playPauseButton.addEventListener("click", togglePlayPause);
                    prevButton.addEventListener("click", playPreviousStream);
                    nextButton.addEventListener("click", playNextStream);
                    seekSlider.addEventListener("input", updateSeek);
                    audioPlayer.addEventListener("timeupdate", updateTime);
                    audioPlayer.addEventListener("durationchange", updateDuration);
                    muteButton.addEventListener("click", toggleMute);
                    volumeSlider.addEventListener("input", updateVolume);
                    repeatButton.addEventListener("click", toggleRepeat);

                    updateRepeatButton();
                });
            } else {
                alert("HLS is not supported in this browser.");
            }
            togglePlayPause();
        }

        function togglePlayPause() {
            if (audioPlayer.paused) {
                audioPlayer.play();
                playPauseButton.innerHTML = "<span class='material-icons md-36'>pause</span>";
                isPlaying = true;
            } else {
                audioPlayer.pause();
                playPauseButton.innerHTML = "<span class='material-icons md-36'>play_arrow</span>";
                isPlaying = false;
            }
        }

        function playPreviousStream() {
            if (audioPlayer.currentTime > 5) {
                audioPlayer.currentTime = 0;
                // repeatCurrentStream();
            } else {
                var currentTrackId = audioFilesData[currentStreamIndex].ID;
                var trackNumber = getTrackNumberInSelection(currentTrackId, currentSelectionName);
                var newTrackNumber = (trackNumber - 1 + selectionDtosData[currentSelectionName].Tracks.length) % selectionDtosData[currentSelectionName].Tracks.length;
                loadStream(currentSelectionName, getTrackIdInSelection(newTrackNumber, currentSelectionName));
            }
        }

        function getTrackNumberInSelection(trackId, selectionName) {
            var selectionTracks = selectionDtosData[selectionName].Tracks;
            for (var i = 0; i < selectionTracks.length; i++) {
                if (String(selectionTracks[i].ID) == String(trackId)) {
                    return i;
                }
            }
        }

        function getTrackIdInSelection(trackNumber, selectionName) {
            return selectionDtosData[selectionName].Tracks[trackNumber].ID;
        }

        function playNextStream() {
            var currentTrackId = audioFilesData[currentStreamIndex].ID;
            var trackNumber = getTrackNumberInSelection(currentTrackId, currentSelectionName);
            var newTrackNumber = (trackNumber + 1) % selectionDtosData[currentSelectionName].Tracks.length;
            loadStream(currentSelectionName, getTrackIdInSelection(newTrackNumber, currentSelectionName));
        }

        function updateSeek() {
            const seekTime = (audioPlayer.duration * (seekSlider.value / 100)).toFixed(2);
            audioPlayer.currentTime = seekTime;
        }

        function updateTime() {
            const currentTime = audioPlayer.currentTime.toFixed(2);
            currentTimeIndicator.textContent = formatTime(currentTime);
            seekSlider.value = (audioPlayer.currentTime / audioPlayer.duration) * 100;
        }

        function updateDuration() {
            const duration = audioPlayer.duration.toFixed(2);
            durationDisplay.textContent = formatTime(duration);
        }

        function formatTime(time) {
            const minutes = Math.floor(time / 60);
            const seconds = Math.floor(time % 60);
            return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
        }

        function toggleMute() {
            audioPlayer.muted = !audioPlayer.muted;
            muteButton.innerHTML = audioPlayer.muted
                ? "<span class='material-icons'>volume_off</span>"
                : "<span class='material-icons'>volume_up</span>";
        }

        function updateVolume() {
            audioPlayer.volume = volumeSlider.value / 100;
        }

        function toggleRepeat() {
            repeatMode = (repeatMode + 1) % 3; // Cycle between 0, 1, and 2
            updateRepeatButton();
        }

        function updateRepeatButton() {
            switch (repeatMode) {
                case 2: // Repeat off
                    repeatButton.innerHTML = "<span class='material-icons'>repeat</span>";
                    repeatButton.style.color = "#909090";
                    break;

                case 0: // Repeat playlist
                    repeatButton.innerHTML = "<span class='material-icons'>repeat</span>";
                    repeatButton.style.color = "#fcc800";
                    break;

                case 1: // Repeat one track
                    repeatButton.innerHTML = "<span class='material-icons'>repeat_one</span>";
                    repeatButton.style.color = "#fcc800";
                    break;

            }
        }

        // Event listener for the "ended" event of the audio player when in "Repeat one" mode
        function repeatCurrentStream() {
            audioPlayer.currentTime = 0; // Repeat the current track by setting currentTime to 0
            if (isPlaying) {
                audioPlayer.play(); // Resume to play audio stream if it was playing before
            }
        }

        // Function to play the next track in "Repeat all" mode when current track ends
        audioPlayer.addEventListener("ended", function () {
            if (repeatMode === 0 && streamUrls.length > 1) {
                playNextStream();
            } else if (repeatMode === 1) {
                repeatCurrentStream();
            }
        });

        document.querySelectorAll(".play-button-popup").forEach(function (button) {
            button.addEventListener("click", function () {
                const streamUrl = songsPath + this.getAttribute("data-src");
                const selectionName = this.getAttribute('data-src-sel');
                const trackId = this.getAttribute('data-src-trackid');

                loadStream(selectionName, trackId);
            });
        });

    };

    hlsPlayer();
});

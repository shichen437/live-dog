<template>
    <div class="video-player-container">
        <el-dialog v-model="dialogVisible" :title="title" width="80%" :before-close="handleClose" destroy-on-close>
            <div class="video-container">
                <video ref="videoPlayer" class="video-js vjs-big-play-centered"></video>
            </div>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref, watch, computed, onBeforeUnmount, nextTick } from 'vue';
import videojs from 'video.js';
import 'video.js/dist/video-js.css';
import '@videojs/http-streaming';
import 'videojs-flvjs-es6';

const props = defineProps({
    modelValue: Boolean,
    title: String,
    videoUrl: String
});

const emit = defineEmits(['update:modelValue']);
const dialogVisible = ref(props.modelValue);
const videoPlayer = ref(null);
let player = null;

const isAudioFile = computed(() => {
    const fileType = props.videoUrl?.toLowerCase() || '';
    return fileType.endsWith('.aac') || fileType.endsWith('.mp3') || fileType.endsWith('.wav');
});

function initPlayer() {
    if (videoPlayer.value) {
        const options = {
            controls: true,
            fluid: true,
            preload: 'auto',
            controlBar: {
                children: [
                    'playToggle',
                    'volumePanel',
                    'currentTimeDisplay',
                    'timeDivider',
                    'durationDisplay',
                    'progressControl',
                    'playbackRateMenuButton',
                    'fullscreenToggle'
                ]
            }
        };

        const fileType = props.videoUrl.toLowerCase();
        if (fileType.endsWith('.flv')) {
            options.techOrder = ['flvjs', 'html5'];
            options.sources = [{
                src: props.videoUrl,
                type: 'video/x-flv'
            }];
        } else if (isAudioFile.value) {
            options.techOrder = ['html5'];
            options.sources = [{
                src: props.videoUrl,
                type: fileType.endsWith('.aac') ? 'audio/aac' : 
                      fileType.endsWith('.mp3') ? 'audio/mpeg' : 'audio/wav'
            }];
            options.height = '50px';
            delete options.fluid;
            delete options.controlBar.fullscreenToggle;
        } else {
            options.techOrder = ['html5'];
            options.sources = [{
                src: props.videoUrl,
                type: 'video/mp4'
            }];
        }

        player = videojs(videoPlayer.value, options);

        if (isAudioFile.value) {
            player.addClass('vjs-audio');
        }

        // 添加快退按钮
        const rewindButton = document.createElement('button');
        rewindButton.className = 'vjs-control vjs-button vjs-rewind-control';
        rewindButton.title = '快退 10 秒';
        rewindButton.onclick = () => {
            const time = player.currentTime();
            player.currentTime(Math.max(0, time - 10));
        };

        // 添加快进按钮
        const forwardButton = document.createElement('button');
        forwardButton.className = 'vjs-control vjs-button vjs-forward-control';
        forwardButton.title = '快进 5 秒';
        forwardButton.onclick = () => {
            const time = player.currentTime();
            player.currentTime(Math.min(player.duration(), time + 5));
        };

        // 插入按钮到控制栏
        const controlBar = player.getChild('ControlBar');
        const playToggle = controlBar.getChild('PlayToggle');
        controlBar.el().insertBefore(rewindButton, playToggle.el().nextSibling);
        controlBar.el().insertBefore(forwardButton, rewindButton.nextSibling);
    }
}

function destroyPlayer() {
    if (player) {
        player.dispose();
        player = null;
    }
}

watch(() => dialogVisible.value, (val) => {
    emit('update:modelValue', val);
});

watch(() => props.modelValue, (val) => {
    dialogVisible.value = val;
    if (val) {
        nextTick(() => {
            destroyPlayer();
            initPlayer();
        });
    }
});

const handleClose = () => {
    destroyPlayer();
    dialogVisible.value = false;
};

onBeforeUnmount(() => {
    destroyPlayer();
});
</script>

<style lang="scss" scoped>
.video-player-container {
    width: 100%;
}

.video-container {
    width: 100%;
    margin: 0 auto;
    background: #000;
    min-height: 50px;
}

:deep(.el-dialog__body) {
    padding: 10px;
}

:deep(.video-js) {
    width: 100%;
    height: 100%;

    &.vjs-audio {
        min-height: 50px !important;
        height: 50px !important;
        background-color: #000;

        .vjs-poster {
            display: none;
        }

        .vjs-big-play-button {
            display: none;
        }

        .vjs-control-bar {
            display: flex;
            position: relative;
            height: 100%;
            opacity: 1;
            background: rgba(0, 0, 0, 0.7);
        }

        .vjs-progress-control {
            position: absolute;
            left: 0;
            right: 0;
            width: 100%;
            height: 4px;
            top: -4px;
            transition: all 0.2s;

            &:hover {
                height: 8px;
                top: -8px;
            }
        }
    }

    .vjs-rewind-control {
        cursor: pointer;
        font-family: Arial, sans-serif;

        &::before {
            font-size: 14px;
            line-height: 2.2;
            content: "10s";
        }
    }

    .vjs-forward-control {
        cursor: pointer;
        font-family: Arial, sans-serif;

        &::before {
            font-size: 14px;
            line-height: 2.2;
            content: "5s";
        }
    }

    .vjs-big-play-button {
        background-color: rgba(0, 0, 0, 0.45);
        border-color: #fff;

        &:hover {
            background-color: #409eff;
        }
    }

    .vjs-control-bar {
        background-color: rgba(0, 0, 0, 0.7);
    }

    .vjs-slider-bar {
        background: #409eff;
    }

    .vjs-play-progress {
        background-color: #409eff;
    }

    .vjs-volume-level {
        background-color: #409eff;
    }
}
</style>
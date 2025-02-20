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
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';
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

function initPlayer() {
    if (videoPlayer.value) {
        const options = {
            controls: true,
            fluid: true,
            aspectRatio: '16:9',
            controlBar: {
                children: [
                    'playToggle',
                    'rewindButton',
                    'forwardButton',
                    'volumePanel',
                    'currentTimeDisplay',
                    'timeDivider',
                    'durationDisplay',
                    'progressControl',
                    'playbackRateMenuButton',
                    'fullscreenToggle'
                ]
            },
            techOrder: ['html5', 'flvjs'],
            flvjs: {
                mediaDataSource: {
                    isLive: false,
                    cors: true,
                    withCredentials: false
                }
            }
        };

        if (props.videoUrl.toLowerCase().endsWith('.flv')) {
            options.sources = [{
                src: props.videoUrl,
                type: 'video/x-flv'
            }];
        } else {
            options.sources = [{
                src: props.videoUrl,
                type: 'video/mp4'
            }];
        }

        player = videojs(videoPlayer.value, options);

        // 注册快进快退按钮
        const Button = videojs.getComponent('Button');
        
        // 创建快退按钮
        const RewindButton = videojs.extend(Button, {
            constructor: function() {
                Button.apply(this, arguments);
                this.controlText("快退 5 秒");
            },
            handleClick: function() {
                const time = this.player().currentTime();
                this.player().currentTime(Math.max(0, time - 5));
            },
            buildCSSClass: function() {
                return `vjs-rewind-control ${Button.prototype.buildCSSClass.call(this)}`;
            }
        });
        
        // 创建快进按钮
        const ForwardButton = videojs.extend(Button, {
            constructor: function() {
                Button.apply(this, arguments);
                this.controlText("快进 5 秒");
            },
            handleClick: function() {
                const time = this.player().currentTime();
                this.player().currentTime(Math.min(this.player().duration(), time + 5));
            },
            buildCSSClass: function() {
                return `vjs-forward-control ${Button.prototype.buildCSSClass.call(this)}`;
            }
        });

        // 注册组件
        videojs.registerComponent('RewindButton', RewindButton);
        videojs.registerComponent('ForwardButton', ForwardButton);
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
}

:deep(.el-dialog__body) {
    padding: 10px;
}

:deep(.video-js) {
    width: 100%;
    height: 100%;
    
.vjs-rewind-control {
        &::before {
            content: "⏪";
        }
    }
    
    .vjs-forward-control {
        &::before {
            content: "⏩";
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
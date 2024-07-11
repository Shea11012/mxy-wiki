---
date created: 2024-06-01 15:37
date modified: 2024-06-01 16:26
---
#ffmpeg

# 视频格式与编码

视频格式，常见的 mp4、mkv、mov、flv、avi、wmv 等，它们都是容器，这个容器里包含了视频流、音频流、字幕等
编码：对于每种流，本质都是文件的字节码，所以需要一个特定的编码方式

1. **MP4 (MPEG-4 Part 14)**
    
    - **视频编码**：H.264 (AVC), H.265 (HEVC), MPEG-4, VP9
    - **音频编码**：AAC, MP3, ALAC, AC-3
    - **字幕编码**：MPEG-4 Timed Text, WebVTT, SRT
2. **AVI (Audio Video Interleave)**
    
    - **视频编码**：DivX, Xvid, H.264, MJPEG
    - **音频编码**：MP3, AC-3, PCM, WMA
    - **字幕编码**：SRT (外部), SSA (外部)
3. **MKV (Matroska Video)**
    
    - **视频编码**：H.264, H.265, VP8, VP9, AV1
    - **音频编码**：AAC, MP3, Vorbis, FLAC, AC-3, DTS
    - **字幕编码**：SRT, SSA/ASS, VobSub, PGS
4. **MOV (QuickTime File Format)**
    
    - **视频编码**：H.264, H.265, MPEG-4, ProRes
    - **音频编码**：AAC, ALAC, MP3
    - **字幕编码**：Text, Closed Captions, SRT (外部)
5. **WMV (Windows Media Video)**
    
    - **视频编码**：WMV1, WMV2, WMV3, VC-1
    - **音频编码**：WMA, WMA Pro, WMA Lossless
    - **字幕编码**：SRT (外部)
6. **FLV (Flash Video)**
    
    - **视频编码**：Sorenson Spark, VP6, H.264
    - **音频编码**：MP3, AAC, Nellymoser
    - **字幕编码**：SRT (外部)
7. **WebM**
    
    - **视频编码**：VP8, VP9, AV1
    - **音频编码**：Vorbis, Opus
    - **字幕编码**：WebVTT
8. **3GP (3rd Generation Partnership Project)**
    
    - **视频编码**：H.263, H.264, MPEG-4
    - **音频编码**：AMR-NB, AMR-WB, AAC, HE-AAC
    - **字幕编码**：3GPP Timed Text
9. **MPEG (Moving Picture Experts Group)**
    
    - **视频编码**：MPEG-1, MPEG-2
    - **音频编码**：MP2, MP3, AAC
    - **字幕编码**：VobSub (MPEG-2), SRT (外部)
10. **OGG (Ogg Video)**
    
    - **视频编码**：Theora
    - **音频编码**：Vorbis, Opus
    - **字幕编码**：Kate, SRT (外部)

# ffmpeg 转码

将 mkv 文件转 mp4
```shell
ffmpeg -i in.mkv out.mp4
```

查看当前文件的信息
```shell
ffmpeg -i in.mkv
```

文件含有两个 stream：
- 视频流 1080P60FPS，码率 4059kb，编码 h264
- 音频流 48k 采样率，127kb 码率，aac 编码
![[../assets/工具/ffmpeg/IMG-20240601155423595.png]]

## 指定编码方式

mkv 文件视频流可能本来就是 H264 编码，mp4 也支持该编码，所以不需要对视频流重新编码
```shell
# -c:v 指定视频流编码器 copy是指不重新编码直接copy源视频流
ffmpeg -i in.mkv -c:v copy out.mp4
```

如果不指定编码，会用 mp4 的默认音视频编码器进行编码

指定编码，视频用 264，音频用 mp3，拼成 mp4
```shell
ffmpeg -i in.mkv -c:v libx264 -c:a libmp3lame out.mp4
```

通过 `ffmpeg -codecs` 查看当前机器支持的编码

## 分辨率 比特率 帧率调整

-s 指定目标输出分辨率，-r 帧率，-b 比特率
```shell
ffmpeg -i input.mp4 -s 720x480 -r 30 -b:v 2M -b:a 128k output.mp4
```
-s 分辨率是等比扩容，截断需要用 filter
```shell
# scale指定宽高，还原点；ow oh指定输出文件的宽高；iw ih输入文件宽高
ffmpeg -i input.mp4 -filter:v "scale=800:600;x=ow/2:y=0" output.mp4
```

## 提取音频

```shell
ffmpeg -i in.mkv out.mp4
```

## 视频截取和拼接

-ss 指定截取开始时间，-t 指定持续时间
```shell
ffmpeg -i input.mp4 -ss 00:00:20 -t 00:00:05 -c copy output.mp4
```
多个视频拼接，需要准备一个 input.txt
```shell
file 'input1.mp4'
file 'input2.mp4'
```
指定该文件为输入，-f 指定 concat 按顺序拼接，拼接视频文件的编码器和格式必须相同
```shell
ffmpeg -f concat -i input.txt -c copy output.mp4
```

## 滤镜

### 视频滤镜

1. **scale**
    
    - 用途：调整视频分辨率。
    - 示例：`ffmpeg -i input.mp4 -vf "scale=1280:720" output.mp4`
2. **crop**
    
    - 用途：裁剪视频。
    - 示例：`ffmpeg -i input.mp4 -vf "crop=640:480:10:10" output.mp4`
3. **transpose**
    
    - 用途：旋转或翻转视频。
    - 示例：`ffmpeg -i input.mp4 -vf "transpose=1" output.mp4` （顺时针旋转 90 度）
4. **hue**
    
    - 用途：调整视频的色调、饱和度、亮度。
    - 示例：`ffmpeg -i input.mp4 -vf "hue=s=0" output.mp4` （将视频变为黑白）
5. **overlay**
    
    - 用途：将一个视频叠加到另一个视频上。
    - 示例：`ffmpeg -i background.mp4 -i overlay.mp4 -filter_complex "overlay=10:10" output.mp4`
6. **drawtext**
    
    - 用途：在视频上添加文本。
    - 示例：`ffmpeg -i input.mp4 -vf "drawtext=text='Hello World':fontcolor=white:fontsize=24:x=10:y=10" output.mp4`
7. **fade**
    
    - 用途：添加淡入淡出效果。
    - 示例：`ffmpeg -i input.mp4 -vf "fade=t=in:st=0:d=5" output.mp4`
8. **fps**
    
    - 用途：更改视频的帧率。
    - 示例：`ffmpeg -i input.mp4 -vf "fps=24" output.mp4`
9. **split**
    
    - 用途：将输入视频分成多个输出流。
    - 示例：`ffmpeg -i input.mp4 -filter_complex "split[v1][v2];[v1]scale=640:360[v1out];[v2]scale=320:180[v2out]" -map "[v1out]" output1.mp4 -map "[v2out]" output2.mp4`
10. **vflip**
    
    - 用途：垂直翻转视频。
    - 示例：`ffmpeg -i input.mp4 -vf "vflip" output.mp4`

### 音频滤镜

1. **volume**
    
    - 用途：调整音量。
    - 示例：`ffmpeg -i input.mp4 -af "volume=0.5" output.mp4`
2. **aecho**
    
    - 用途：添加回声效果。
    - 示例：`ffmpeg -i input.mp3 -af "aecho=0.8:0.88:60:0.4" output.mp3`
3. **atempo**
    
    - 用途：调整音频速度（节奏）。
    - 示例：`ffmpeg -i input.mp3 -af "atempo=1.5" output.mp3`
4. **equalizer**
    
    - 用途：均衡器，调整特定频率的增益。
    - 示例：`ffmpeg -i input.mp3 -af "equalizer=f=1000:t=q:w=1:g=5" output.mp3`
5. **aresample**
    
    - 用途：重新采样音频。
    - 示例：`ffmpeg -i input.mp3 -af "aresample=44100" output.mp3`
6. **acompressor**
    
    - 用途：音频压缩器，减少动态范围。
    - 示例：`ffmpeg -i input.mp3 -af "acompressor" output.mp3`
7. **highpass**
    
    - 用途：高通滤波器，移除低频声音。
    - 示例：`ffmpeg -i input.mp3 -af "highpass=f=300" output.mp3`
8. **lowpass**
    
    - 用途：低通滤波器，移除高频声音。
    - 示例：`ffmpeg -i input.mp3 -af "lowpass=f=3000" output.mp3`
9. **adelay**
    
    - 用途：音频延迟。
    - 示例：`ffmpeg -i input.mp3 -af "adelay=1000|1000" output.mp3`
10. **pan**
    
    - 用途：调整声道平衡。
    - 示例：`ffmpeg -i input.mp3 -af "pan=stereo|c0=0.5*c0+0.5*c1|c1=0.5*c0+0.5*c1" output.mp3`

## ffmpeg 输入

-i 可以接本地视频文件，txt 格式的视频列表，还能接直播流
```shell
ffmpeg -i rtmp://example.com/live/stream -c copy output.mp4
```
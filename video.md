# Video

### Commands

```
~/Tools/ffmpeg -y -r 0.5 -i ~/Desktop/background/black_*.png -i ~/Desktop/voice/alien.mp3 -absf aac_adtstoasc -r 30 ~/Desktop/output.mp4

~/Tools/ffmpeg -f image2 -r 0.5 -i ~/Desktop/background/black_%04d.png ~/Desktop/output.mp4

~/Tools/ffmpeg -i ~/Desktop/test/galaxy.mp3 -loop 1 -f image2 -r 1 -i ~/Desktop/test/black_0%d.png -absf aac_adtstoasc  -t 5  ~/Desktop/test/output.mp4

~/Tools/ffmpeg -i ~/Desktop/test/galaxy.mp3 -loop 1 -r 1 -i ~/Desktop/test/black_0%d.png -absf aac_adtstoasc  -t 5  ~/Desktop/test/output.mp4

~/Tools/ffmpeg -i ~/Desktop/test/galaxy.mp3 -r 1 -i ~/Desktop/test/black_0%d.png -acodec aac -strict -2 -vcodec libx264 -ar 22050 -ab 128k -ac 2 -pix_fmt yuvj420p -y -t 5  ~/Desktop/test/output.mp4     可以

~/Tools/ffmpeg -i ~/Desktop/test/galaxy.mp3 -r 1 -i ~/Desktop/test/black_0%d.png -c:v libx264 -tune stillimage -c:a aac -b:a 192k -strict -2 -ar 22050 -ab 128k -ac 2 -pix_fmt yuvj420p -y  ~/Desktop/test/output.mp4  可以

~/Tools/ffmpeg -i ~/Desktop/test/galaxy.mp3 -r 1 -i ~/Desktop/test/black_0%d.png -c:v libx264 -tune stillimage -c:a aac -b:a 192k -ar 22050 -ac 2 -pix_fmt yuvj420p -y  ~/Desktop/test/output.mp4 可以

~/Tools/ffmpeg -r 1 -i ~/Desktop/test/black_0%d.png -i ~/Desktop/test/galaxy.mp3 -c:v libx264 -c:a aac -b:a 192k -ar 22050 -ac 2 -pix_fmt yuvj420p -shortest -y  ~/Desktop/test/output1.mp4
```
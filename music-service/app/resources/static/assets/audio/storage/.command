Example: ffmpeg -i input.wav -c:a aac -b:a 128k -map 0:0 -f segment -segment_time 10 -segment_list outputlist_mpegts.m3u8 -segment_format mpegts output_mpegts%03d.ts


ffmpeg -i 'DJ Snake, Justin Bieber – Let Me Love You.mp3' -c:a aac -b:a 192k -map 0:0 -f segment -segment_time 5 -segment_list djsnake_bieber-let_me_love_you.m3u8 -segment_format mpegts djsnake_bieber-let_me_love_you%03d.ts
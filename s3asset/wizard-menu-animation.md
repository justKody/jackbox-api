# Wizard menu animation

The looping wizard on the main menu (`LandingScreen`). Made from a green-screen video, with the green removed and exported as a transparent animated WebP.

## Output

| Setting | Value |
|---|---|
| File | `trivia-app/assets/wizard-menu.webp` |
| Format | Animated WebP, lossy, with transparency |
| Resolution | 480 × 480 px |
| Frame rate | 24 fps (42 ms per frame) |
| Frames | 240 |
| Length | about 10 s, loops forever |
| Size | 5.4 MB |

## Source

| Setting | Value |
|---|---|
| File | `backend/s3asset/Agent Video - SCENE 1 2D hand-drawn game-menu character animation_ single continuous take_ no cuts_.mp4` |
| Codec | 10-bit HEVC, no audio, 28 MB |
| Resolution | 1664 × 1248 |
| Frame rate | 24 fps, 241 frames (the last one repeats the first) |
| Background | green screen, measured at `#12E71F` |

## Processing steps

| Step | Setting | Why |
|---|---|---|
| Crop | `1248×1248` starting at x=242, y=0 | Full-height square centred on the wizard, including the spin |
| Scale | 480 × 480, lanczos | Sharp downscale |
| Green removal | `chromakey` with colour `0x12E71F`, similarity `0.16`, blend `0.06` | Removes the green and leaves soft edges |
| Spill fix | `despill` with type `green`, mix `0.15`, expand `0` | Light setting, because stronger values turned the cyan eyes blue |
| Frames kept | 240 | Drops the duplicate last frame so the loop is seamless |
| WebP encoding | `img2webp -lossy -q 60 -m 6 -d 42 -loop 0` | Quality 60, best compression, infinite loop |

## Why 480 px

It shows at `min(210, screenHeight × 0.22)` points, which is about 180 pt on a typical phone. That needs about 360 px on 2× screens and about 540 px on 3× screens. Anything bigger is downscaled by the phone anyway, and it costs file size, memory and battery.

## Recreate it

Run from the project root. It needs `ffmpeg` and `img2webp` (`brew install ffmpeg webp`) and about 100 MB of free disk space while it runs.

```bash
SRC="backend/s3asset/Agent Video - SCENE 1 2D hand-drawn game-menu character animation_ single continuous take_ no cuts_.mp4"
mkdir -p /tmp/wiz && ffmpeg -y -i "$SRC" -frames:v 240 -vf "crop=1248:1248:242:0,scale=480:480:flags=lanczos,format=yuv444p,chromakey=0x12E71F:0.16:0.06,despill=type=green:mix=0.15:expand=0,format=rgba" /tmp/wiz/f%03d.png && img2webp -loop 0 -lossy -q 60 -m 6 -d 42 /tmp/wiz/*.png -o trivia-app/assets/wizard-menu.webp && rm -rf /tmp/wiz
```

### Settings worth changing

- `scale=480:480`: change both numbers to `400:400` for a file of about 3.5–4 MB that looks slightly softer.
- `-q 60`: lower means smaller and blurrier; higher means sharper and bigger.
- Frame rate: for 12 fps (about half the size, but choppier), add `fps=12,` at the start of the `-vf` filters and change `-frames:v 240` to `-frames:v 120` and `-d 42` to `-d 83`.
- `chromakey` similarity: raise it if green edges show, lower it if the glow rings start disappearing.

## In the app

`trivia-app/src/screens/LandingScreen.tsx`, rendered with `expo-image`:

```tsx
<Image
  source={require('../../assets/wizard-menu.webp')}
  autoplay={!reduce}          // still frame when "reduce motion" is on
  contentFit="contain"
  style={{ width: hatSize, height: hatSize }}  // hatSize = min(210, screenHeight * 0.22)
/>
```

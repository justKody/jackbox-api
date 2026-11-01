"""Builds dataset.json: the OWID questions that have a voice-over clip.

A question's clip is s3asset/v2/<id>.mp3, the Hat reading that question aloud.
Questions without a clip are left out, so every question the game asks has a
voice. The game server embeds dataset.json as its question bank.

Run after fetch_owid.py or after adding clips: python3 build_dataset.py
"""

import json
from pathlib import Path

HERE = Path(__file__).parent
QUESTIONS = HERE / "owid_questions.json"
ASSETS = HERE.parents[1] / "s3asset"
CLIPS = "v2"
OUT = HERE / "dataset.json"


def main():
    questions = json.loads(QUESTIONS.read_text())
    out = []
    for q in questions:
        clip = f"{CLIPS}/{q['id']}.mp3"
        if not (ASSETS / clip).is_file():
            continue
        out.append({
            "id": q["id"],
            "text": q["text"],
            "answer": q["answer"],
            # path under s3asset, the server turns it into a url
            "voiceOver": clip,
            "year": q["year"],
            "source": q["source"],
            "url": q["url"],
        })
        print(f"ok   {q['id']:7} {clip}")

    OUT.write_text(json.dumps(out, indent=2, ensure_ascii=False) + "\n")
    print(f"\nwrote {len(out)} of {len(questions)} questions to {OUT}")


if __name__ == "__main__":
    main()

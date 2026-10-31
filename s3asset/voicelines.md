# Wizard Hat Host voicelines

Prompts for the host's reaction lines in the trivia app, one table per game moment. Each line has a few variants so the Hat doesn't repeat itself; the app picks one at random (skipping whichever played last).

All 69 lines are generated in `backend/s3asset/voicelines/<id>.mp3` (Wizard Hat Host, eleven_v4, one take each).

## Rules for every prompt

- **Voice:** Wizard Hat Host (`12z3QxL3NOpVt1ZCSQmA`), model `eleven_v4`, one take per line.
- **Tags go at the start or middle, never at the end.** No `[chuckles]` anywhere.
- **No names or numbers.** Audio is recorded once, so it can't say "Sam" or "42%". The on-screen bubble keeps the specific text (names, guesses, points); the voice gives the reaction.
- **Button reactions stay under ~2 seconds**, so tapping Higher or Lower never talks over the next thing.
- **Serious questions:** reaction lines are generic, so they work for every question, but avoid playing a jokey filler line straight after a sensitive question clip (poverty, child deaths and so on).
- **One line at a time.** A new line stops the one playing; the question clip always wins over a filler.

## Order within a round

1. New round line (`round_next` or `round_final`)
2. The question clip (`s3asset/questions/<id>.mp3` or `s3asset/v2/<id>.mp3`)
3. Guess phase lines, then vote phase lines, then reveal lines

Timers come from `src/game/useGame.ts`: guess 20s, vote 15s, and the "hurry" lines fire at 5 seconds left, the same point the Hat's mood switches to `shock`.

---

## 1. Game start

Plays once, when round 1 begins. The existing `welcome-back-to-the-show.mp3` can join this set for returning players.

| ID | Prompt |
|---|---|
| start_1 | [excitedly] Gather round, gather round! [playfully] The Hat is open for business. Six riddles, one champion! |
| start_2 | [warmly] Welcome, welcome, my clever friends! [mischievously] Let's see whose guesses are truly magical. |
| start_3 | [excitedly] Wands at the ready! [dramatic pause] The game... begins! |

## 2. New round

Before each question's clip, from round 2 on. The final round gets its own lines.

| ID | When | Prompt |
|---|---|---|
| round_next_1 | Rounds 2–5 | [playfully] Onward! Another riddle awaits. |
| round_next_2 | Rounds 2–5 | [excitedly] Next question, coming right up! |
| round_next_3 | Rounds 2–5 | [mischievously] Ooh, I've got a tricky one for you this time. |
| round_final_1 | Last round | [dramatic pause] The final round! [excitedly] Everything rides on this one! |
| round_final_2 | Last round | [suspensefully] Last riddle of the game... make it count. |

## 3. Guess phase: your turn

After the question clip, when you're the guesser (`GuessScreen.tsx`, `YourGuess`).

| ID | Prompt |
|---|---|
| guess_you_1 | [warmly] Your turn! Spin the dial to your best guess. |
| guess_you_2 | [playfully] Go on then, give that dial a whirl. |
| guess_you_3 | [encouragingly] Trust your gut... what's the number? |

## 4. Guess phase: someone else's turn

After the question clip, when another player is guessing.

| ID | Prompt |
|---|---|
| guess_other_1 | [curiously] They're turning the dial... get ready to call it! |
| guess_other_2 | [whispers] Watch closely now... higher or lower is coming. |

## 5. Guess phase: filler

The gap-filler. Plays once if nobody has locked in with about 10 seconds left.

| ID | Prompt |
|---|---|
| guess_filler_1 | [thoughtfully] Hmm... take your time. [playfully] But not too much time. |
| guess_filler_2 | [mischievously] Somewhere between zero and a hundred. That narrows it down! |
| guess_filler_3 | [warmly] No wrong answers here... [dramatic pause] well, actually, quite a few. |

## 6. Guess phase: running out of time

At 5 seconds left, if the guess isn't locked.

| ID | Prompt |
|---|---|
| guess_hurry_1 | [urgently] Quick, quick! Lock it in! |
| guess_hurry_2 | [gasps] The sands are running out! |
| guess_hurry_3 | [urgently] Five seconds! Decide, decide! |

## 7. Guess locked

When you tap "Lock in".

| ID | Prompt |
|---|---|
| guess_locked_1 | [proudly] Sealed with magic! |
| guess_locked_2 | [excitedly] Locked in! No take-backs. |
| guess_locked_3 | [playfully] Bold choice. I like it. |

## 8. Guess timed out

When the timer hits zero and the guess locks wherever the dial stopped.

| ID | Prompt |
|---|---|
| guess_timeout_1 | [sighs] Time's up! I've sealed it where the dial stopped. |
| guess_timeout_2 | [playfully] Out of time! The Hat has decided for you. |

## 9. Vote phase starts

When the guess is shown and everyone else must call it (`VoteScreen.tsx`).

| ID | Prompt |
|---|---|
| vote_start_1 | [excitedly] Now then... higher or lower? |
| vote_start_2 | [mischievously] Too high? Too low? [dramatic pause] Make your call! |
| vote_start_3 | [playfully] Place your bets, my friends. Higher or lower? |

## 10. You tapped Higher

| ID | Prompt |
|---|---|
| vote_higher_1 | [excitedly] Higher! Reaching for the stars! |
| vote_higher_2 | [playfully] Up, up and away! |
| vote_higher_3 | [impressed] Higher it is. Feeling optimistic! |

## 11. You tapped Lower

| ID | Prompt |
|---|---|
| vote_lower_1 | [mischievously] Lower! A cautious one, I see. |
| vote_lower_2 | [playfully] Down we go! |
| vote_lower_3 | [knowingly] Lower, eh? A bit of a sceptic. |

## 12. Vote phase: you're the guesser

While the others vote on your guess.

| ID | Prompt |
|---|---|
| vote_guesser_1 | [warmly] Everyone's betting on your guess now. Fingers crossed... |
| vote_guesser_2 | [whispers] Let's see if they trust you... |

## 13. Vote phase: running out of time

At 5 seconds left, if you haven't called.

| ID | Prompt |
|---|---|
| vote_hurry_1 | [urgently] Tick tock! Higher or lower? |
| vote_hurry_2 | [gasps] Hurry, hurry! Make your call! |

## 14. All calls in

| ID | Prompt |
|---|---|
| vote_allin_1 | [excitedly] All calls are in... let's reveal! |
| vote_allin_2 | [dramatic pause] Everyone's decided. [mischievously] Let's see who's right. |

## 15. Reveal build-up

Right before the answer appears (`RevealScreen.tsx`, "And the answer is…").

| ID | Prompt |
|---|---|
| reveal_build_1 | [suspensefully] And the answer is... |
| reveal_build_2 | [whispers] Drumroll, please... |
| reveal_build_3 | [dramatic pause] The moment of truth... |

## 16. Reveal result

Picked from what happened to you this round, using the same checks `RevealScreen.tsx` already uses for the bubble text.

| ID | When | Prompt |
|---|---|---|
| reveal_bullseye_1 | You guessed exactly | [gasps] A perfect bullseye! [excitedly] Pure wizardry! |
| reveal_bullseye_2 | You guessed exactly | [amazed] Spot on to the very number! I'm speechless! |
| reveal_close_1 | Your guess earned 80+ points | [proudly] So close! That's a fine guess. |
| reveal_close_2 | Your guess earned 80+ points | [warmly] Nearly there! The Hat is impressed. |
| reveal_ok_1 | Your guess earned 50–79 points | [warmly] Not bad at all! Points on the board. |
| reveal_far_1 | Your guess earned under 50 points | [sighs] Oh dear, not even close. [playfully] We'll pretend that didn't happen. |
| reveal_far_2 | Your guess earned under 50 points | [playfully] Way off! But I admire the confidence. |
| reveal_right_1 | Your call was right | [cheerfully] Spot on! You called it! |
| reveal_right_2 | Your call was right | [excitedly] Correct! The Hat salutes you. |
| reveal_wrong_1 | Your call was wrong | [sympathetically] Ooh, not this time. The next one's yours! |
| reveal_wrong_2 | Your call was wrong | [playfully] Wrong way! Happens to the best of us. |
| reveal_tie_1 | Someone else hit a bullseye, so no call wins | [gasps] A bullseye! Nobody wins the call this time. |
| reveal_nocall_1 | You didn't call | [sighs] You didn't make a call this time... |
| reveal_nocall_2 | You didn't call | [playfully] No call? Even a wild guess beats silence! |

## 17. Game over

On the final scores screen (`GameOverScreen.tsx`).

| ID | When | Prompt |
|---|---|---|
| gameover_win_1 | You came first | [excitedly] And our champion is... you! [proudly] Magnificent! |
| gameover_win_2 | You came first | [cheering] Victory! The Hat bows to you. |
| gameover_lose_1 | You didn't win | [warmly] Well played! The crown slipped away this time. |
| gameover_lose_2 | You didn't win | [playfully] Not a win, but a jolly good show. Another round? |
| gameover_end_1 | After the win or lose line | [warmly] That's the game! Thank you for playing, my friends. |

## 18. Filler gap lines

For any quiet stretch that isn't covered above: waiting on other players, between the reveal and "Next question", or a slow connection. Play at most one per gap.

| ID | Prompt |
|---|---|
| filler_1 | [playfully] Hmm, now where did I put my wand... |
| filler_2 | [mischievously] Did you know? I'm only wrong about one percent of the time. |
| filler_3 | [humming] Dum dee dum... |
| filler_4 | [thoughtfully] Numbers, numbers... they're everywhere, you know. |
| filler_5 | [warmly] Lovely to have you all here. |
| filler_6 | [whispers] Psst... the answer is never quite what you think. |

## 19. Matchmaking (optional)

| ID | When | Prompt |
|---|---|---|
| match_searching_1 | Finding rivals | [curiously] Searching the realm for worthy rivals... |
| match_found_1 | Match found | [excitedly] Rivals found! Let the game begin! |
